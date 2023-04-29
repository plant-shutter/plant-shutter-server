package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/plant-shutter/plant-shutter-server/pkg/orm"
	"github.com/plant-shutter/plant-shutter-server/pkg/utils/config"
)

var (
	db   *sql.DB
	once sync.Once
	cfg  config.Mysql
)

func Init(config config.Mysql) {
	cfg = config
}

func getInstance() *sql.DB {
	if db == nil {
		once.Do(func() {
			source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Passwd, cfg.Host, cfg.Port, cfg.Database)
			var err error
			db, err = sql.Open("mysql", source)
			if err != nil {
				panic(err)
			}
			db.SetConnMaxLifetime(time.Minute * 3)
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(10)
		})
	}
	return db
}

func Close() {
	if db != nil {
		_ = db.Close()
	}
}

func GetDeviceByID(id int) (*orm.Device, error) {
	devices, err := selectDevices("where id = ?", id)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, nil
	}

	return &devices[0], nil
}

func GetDeviceByName(name string) (*orm.Device, error) {
	devices, err := selectDevices("where name_ = ?", name)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, nil
	}

	return &devices[0], nil
}

func selectDevices(where string, args ...any) ([]orm.Device, error) {
	db := getInstance()
	rows, err := db.Query("select * from device "+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []orm.Device
	for rows.Next() {
		var v orm.Device
		if err = rows.Scan(&v.ID, &v.Name, &v.Info, &v.UserID, &v.LastActivity, &v.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, v)
	}

	return res, nil
}

func GetProjectByID(id int) (*orm.Project, error) {
	projects, err := selectProjects("where id = ?", id)
	if err != nil {
		return nil, err
	}
	if len(projects) == 0 {
		return nil, nil
	}

	return &projects[0], nil
}

func GetProjectByName(name string) (*orm.Project, error) {
	projects, err := selectProjects("where name_ = ?", name)
	if err != nil {
		return nil, err
	}
	if len(projects) == 0 {
		return nil, nil
	}

	return &projects[0], nil
}

func selectProjects(where string, args ...any) ([]orm.Project, error) {
	db := getInstance()
	rows, err := db.Query("select * from project "+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []orm.Project
	for rows.Next() {
		var v orm.Project
		if err = rows.Scan(&v.ID, &v.UserID, &v.DeviceID, &v.Name, &v.Info, &v.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, v)
	}

	return res, nil
}

func GetProjectLatestImage(projectID int) (*orm.Image, error) {
	images, err := selectImages("where project_id = ? order by created_at desc limit 1", projectID)
	if err != nil {
		return nil, err
	}
	if len(images) == 0 {
		return nil, nil
	}

	return &images[0], nil
}

func selectImages(where string, args ...any) ([]orm.Image, error) {
	db := getInstance()
	rows, err := db.Query("select * from image "+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []orm.Image
	for rows.Next() {
		var v orm.Image
		if err = rows.Scan(&v.ID, &v.ProjectID, &v.Name, &v.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, v)
	}

	return res, nil
}

func AddImage(image *orm.Image) (err error) {
	db := getInstance()
	stmt, err := db.Prepare("insert into image (project_id , name_, created_at) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(image.ProjectID, image.Name, image.CreatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	image.ID = int(id)

	return
}
