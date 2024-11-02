package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("short.db"))
	if err != nil {
		panic(err)
	}
	DB = db
	return DB

}
func Migarete(db *gorm.DB) {

	db.AutoMigrate(&Url{})

}

type Url struct {
	gorm.Model
	ID       int `gorm:"primaryKey"`
	Url      string
	ShortUrl string
}
