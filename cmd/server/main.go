package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vincent-vinf/go-jsend"

	"github.com/plant-shutter/plant-shutter-server/pkg/db"
	"github.com/plant-shutter/plant-shutter-server/pkg/orm"
	"github.com/plant-shutter/plant-shutter-server/pkg/storage"
	"github.com/plant-shutter/plant-shutter-server/pkg/utils"
	"github.com/plant-shutter/plant-shutter-server/pkg/utils/config"
)

const (
	DeviceNameHeader = "Device-Name"

	GinDeviceKey = "device"
)

var (
	configPath = flag.String("config-path", "config.yaml", "")
	port       = flag.Int("port", 8000, "")
	logger     = logrus.New()
)

var (
	imageStorage *storage.Storage
)

func init() {
	flag.Parse()
}

func main() {
	var err error
	if err = config.Load(*configPath); err != nil {
		logger.Fatal(err)
	}
	imageStorage, err = storage.New(config.Global.Image)
	if err != nil {
		logger.Fatal(err)
	}
	db.Init(config.Global.Mysql)

	r := gin.New()
	//gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(utils.Cors())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, jsend.SimpleErr("page not found"))
	})

	clientRouter := r.Group("/api/device")
	clientRouter.Use(DeviceAuthMiddleware())
	clientRouter.POST("/project/:project/upload", uploadImage)

	customerRouter := r.Group("/api/customer")
	customerRouter.GET("/project/:project/image", getImage)

	utils.ListenAndServe(r, *port)
}

func getImage(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("project"))
	project, err := db.GetProjectByID(pid)
	if err != nil {
		internalErr(c, "get project failed", err)
		return
	}
	if project == nil {
		c.JSON(http.StatusBadRequest, jsend.SimpleErr("project not found"))
		return
	}

	image, err := db.GetProjectLatestImage(project.ID)
	if err != nil {
		internalErr(c, "get latest image failed", err)
		return
	}
	c.File(imageStorage.GetPath(project.ID, image.Name))
}

func uploadImage(c *gin.Context) {
	//device := getDeviceFormReq(c)
	// todo: check device with project
	pid, _ := strconv.Atoi(c.Param("project"))
	project, err := db.GetProjectByID(pid)
	if err != nil {
		internalErr(c, "get project failed", err)
		return
	}
	if project == nil {
		c.JSON(http.StatusBadRequest, jsend.SimpleErr("project not found"))
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, jsend.SimpleErr("image not found"))
		return
	}

	src, err := file.Open()
	if err != nil {
		internalErr(c, "open file err", err)
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		internalErr(c, "read file err", err)
		return
	}

	fileName := utils.GetFileName(data)

	image := &orm.Image{
		ProjectID: project.ID,
		Name:      fileName,
		CreatedAt: time.Now(),
	}
	if err = db.AddImage(image); err != nil {
		internalErr(c, "insert image to db failed", err)
		return
	}

	if err = imageStorage.Save(project.ID, fileName, bytes.NewReader(data)); err != nil {
		internalErr(c, "sava file err", err)
		return
	}

	c.JSON(http.StatusOK, jsend.Success(fmt.Sprintf("'%s' uploaded!", file.Filename)))
}

func internalErr(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusInternalServerError, jsend.SimpleErr(msg))
	logger.Errorf("%s, err: %s", msg, err)
}

func DeviceAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var deviceName string
		if d, ok := c.Request.Header[DeviceNameHeader]; ok && len(d) > 0 && d[0] != "" {
			deviceName = d[0]
		} else {
			c.Abort()
			c.JSON(http.StatusUnauthorized, jsend.SimpleErr("device name not found in header"))
			return
		}

		device, err := db.GetDeviceByName(deviceName)
		if err != nil {
			c.Abort()
			internalErr(c, "get device by name err", err)
			return
		}
		if device == nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, jsend.SimpleErr("device name not found"))
			return
		}
		c.Set(GinDeviceKey, *device)
	}
}

func getDeviceFormReq(c *gin.Context) *orm.Device {
	d, ok := c.Get(GinDeviceKey)
	if !ok {
		return nil
	}
	device, ok := d.(orm.Device)
	if !ok {
		return nil
	}

	return &device
}
