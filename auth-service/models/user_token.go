package models

import (
	"gorm.io/gorm"
)

type UserToken struct {
	gorm.Model
	UserId uint   `gorm:"size:255;not null;" json:"userId"`
	Token  string `gorm:"size:255;not null;" json:"token"`
}

func (ut *UserToken) SaveUserToken() error {

	var err error
	err = GetDB().Create(&ut).Error
	if err != nil {
		return err
	}
	return nil
}
