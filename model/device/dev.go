package device

import (
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
