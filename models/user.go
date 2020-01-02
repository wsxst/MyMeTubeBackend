package models

import db "MeTube/database"

type MUser struct {
	Id 			   int64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	NickName       string `gorm:"column:nickname" json:"nickname"`
	Avatar         string `gorm:"column:avatar" json:"avatar"`
}

func (p *MUser) TableName() string {
	return "users"
}

func GetUserByID(id int64) *MUser {
	user := &MUser{}

	err := db.ORM.Where("id = ?", id).First(user).Error
	if err != nil || user.Id <= 0 {
		return nil
	}

	return user
}