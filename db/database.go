package db

import (
	// "time"

	"golang.org/x/exp/slog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"valorant/models"
)

func Db() *gorm.DB {
  // loc, err := time.LoadLocation("Asia/Tokyo")
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err != nil {
    slog.Error("message", err)
  }

  return db
}

func AutoMigrate() {
  db := Db()
  db.AutoMigrate(&models.Stats{})
  db.AutoMigrate(&models.Matches{})
  db.AutoMigrate(&models.Ratings{})
}
