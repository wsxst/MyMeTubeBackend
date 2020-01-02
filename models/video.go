package models

import (
	db "MeTube/database"
	"github.com/jinzhu/gorm"
	"time"
)

type MVideo struct {
	Id           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Title        string    `gorm:"column:title" json:"title"`
	Info         string    `gorm:"column:info" json:"info"`
	Url          string    `gorm:"column:url" json:"url"`
	Avatar       string    `gorm:"column:avatar" json:"avatar"`
	Owner        int64     `gorm:"column:owner" json:"owner"`
	TypeName     string    `gorm:"column:typename" json:"typename"`
	ViewCount    int64     `gorm:"column:view_count" json:"view_count"`
	CommentCount int64     `gorm:"column:comment_count" json:"comment_count"`
}

func (p *MVideo) TableName() string {
	return "videos"
}

func GetVideoByID(id int64) *MVideo {
	video := &MVideo{}

	//fmt.Println(id)
	err := db.ORM.Where("id = ?", id).First(video).Error
	if err != nil  || video.Id <= 0 {
		return nil
	}

	return video
}

func AddCommentCountByVideoId(videoId int64) error {
	err := db.ORM.Model(&MVideo{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func AddViewCountByVideoId(videoId int64) error {
	err := db.ORM.Model(&MVideo{}).Where("id = ?", videoId).Update("view_count", gorm.Expr("view_count + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func MinusCommentCountByVideoId(videoId int64) error {
	err := db.ORM.Model(&MVideo{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}
