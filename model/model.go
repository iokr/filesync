package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/dzhenquan/filesync/web/config"
)

var DB *gorm.DB

// Base Model
type BaseModel struct {
	ID        uint64 	`gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", config.DBConfig.URL)
	if err != nil {
		return nil, err
	}

	DB = db
	db.SingularTable(true)
	db.AutoMigrate(&TaskFileInfo{}, &TaskFileLog{}, &TaskFileCopy{})
	return db, err
}