package models

import db "MeTube/database"

type MView struct {
	Id 			  int64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	VideoId       int64 `gorm:"column:video_id" json:"video_id"`
	UserId        int64 `gorm:"column:user_id" json:"user_id"`
}

func (p *MView) TableName() string {
	return "view"
}

func JudgeView(videoId int64, userId int64) bool {
	videoView :=&MView{}
	db.ORM.Where("user_id = ? and video_id = ?", userId, videoId).First(videoView)
	if videoView.Id > 0 {
		return true
	} else {
		return false
	}
}

func CountView(videoId int64) int64 {
	var viewCount int64
	db.ORM.Where("video_id = ?", videoId).Count(&viewCount)
	return viewCount
}

func AddVideoView(videoId int64, userId int64){
	if !JudgeView(videoId, userId) {
		videoView := &MView{
			VideoId: videoId,
			UserId:  userId,
		}
		db.ORM.Create(videoView)
	}
}
