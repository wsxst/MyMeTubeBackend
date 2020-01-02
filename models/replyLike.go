package models

import db "MeTube/database"

type MReplyLike struct {
	Id 			  int64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	ReplyId       int64 `gorm:"column:reply_id" json:"reply_id"`
	UserId        int64 `gorm:"column:user_id" json:"user_id"`
	Flag          bool `gorm:"column:flag" json:"flag"`
}

func (p *MReplyLike) TableName() string {
	return "reply_like"
}

func AddLikeByReplyId(replyId int64, userId int64) error {
	replyLike :=&MReplyLike{}
	db.ORM.Where("user_id = ? and reply_id = ?", userId, replyId).First(replyLike)
	if replyLike.Id > 0 {
		err := db.ORM.Model(&MReplyLike{}).Where("user_id = ? and reply_id = ?", userId, replyId).Update("flag", true).Error
		if err != nil {
			return err
		}
	} else {
		replyLike = &MReplyLike{ReplyId:replyId, UserId:userId, Flag:true}
		db.ORM.Create(replyLike)
	}

	return nil
}

func CancelLikeByReplyId(replyId int64, userId int64) error {
	err := db.ORM.Model(&MReplyLike{}).Where("user_id = ? and reply_id = ?", userId, replyId).Update("flag", false).Error
	if err != nil {
		return err
	}

	return nil
}

func JudgeLikeReply(userId int64, replyId int64) bool {
	replyLike :=&MReplyLike{}
	db.ORM.Where("user_id = ? and reply_id = ?", userId, replyId).First(replyLike)
	if replyLike.Id > 0 && replyLike.Flag {
		return true
	} else {
		return false
	}
}

func CountLikeReply(replyId int64) int64 {
	var likeCount int64
	db.ORM.Where("reply_id = ? and flag = TRUE", replyId).Count(&likeCount)
	return likeCount
}

