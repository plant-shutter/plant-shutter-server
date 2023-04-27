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
	return selectDevice("id = ?", id)
}

func selectDevice(where string, args ...any) (*orm.Device, error) {
	db := getInstance()
	rows, err := db.Query("select * from device where "+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	v := &orm.Device{}
	if rows.Next() {
		if err = rows.Scan(&v.ID, &v.Name, &v.Info, &v.LastActivity, &v.CreatedAt); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("not found")
	}

	return v, nil
}
