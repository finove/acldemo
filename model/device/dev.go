package device

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	MAC    string `json:"MAC"`
	ID     string `json:"ID"`
	Expire string `json:"Expire"`
	Remark string `json:"Remark"`
	Active bool   `json:"Active"`
}

func UserDB() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("failed to connect database")
		return
	}
	db.AutoMigrate(&Device{})
	db.Create(&Device{
		MAC:    "000EA93D209C",
		Remark: "测试1",
	})
	var dev Device
	db.First(&dev, "mac = ?", "000EA93D209C")
	log.Info().Msgf("find device %+v", dev)
}
