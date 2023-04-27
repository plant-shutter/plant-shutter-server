package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vincent-vinf/go-jsend"

	"github.com/plant-shutter/plant-shutter-server/pkg/db"
	"github.com/plant-shutter/plant-shutter-server/pkg/storage"
	"github.com/plant-shutter/plant-shutter-server/pkg/utils"
	"github.com/plant-shutter/plant-shutter-server/pkg/utils/config"
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

	router := r.Group("/api/device")
	router.Use(DeviceAuthMiddleware())
	router.POST("/:device/:project/upload", upload)

	utils.ListenAndServe(r, *port)
}

func upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		internalErr(c, "image not found", err)
		return
	}
	src, err := file.Open()
	if err != nil {
		internalErr(c, "open file err", err)
		return
	}

	if err = imageStorage.Save("1", "2", "3", src); err != nil {
		internalErr(c, "sava file err", err)
		return
	}
	defer src.Close()

	c.JSON(http.StatusOK, jsend.Success(fmt.Sprintf("'%s' uploaded!", file.Filename)))
}

func internalErr(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusInternalServerError, jsend.SimpleErr(msg))
	logger.Errorf("%s, err: %s", msg, err)
}

func DeviceAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Abort()
		logger.Info(c.Request.Header)
	}
}
