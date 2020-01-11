package models

import (
	db "MeTube/database"
	"github.com/jinzhu/gorm"
	"time"
)

type MUserReply struct {
	Id 			int64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	SendUserId  int64 `gorm:"column:send_user_id" json:"send_user_id"`
	RecvUserId  int64 `gorm:"column:recv_user_id" json:"recv_user_id"`
	CommentId   int64 `gorm:"column:comment_id" json:"comment_id"`
	Content     string `gorm:"column:content" json:"content"`
	ReplyTime   time.Time `gorm:"column:reply_time" json:"reply_time"`
	Level       int64 `gorm:"column:level" json:"level"`
	LikeCount   int64 `gorm:"column:like_count" json:"like_count"`
}

func (p *MUserReply) TableName() string {
	return "user_reply"
}

func GetOnePageRepliesByCommentId(commentId int64, replyStart int64, replyFoldLimit int64) []*MUserReply {
	var replies []*MUserReply

	err := db.ORM.Where("comment_id = ?", commentId).Limit(replyFoldLimit).Offset((replyStart - 1)*replyFoldLimit).Order("like_count desc").Order("reply_time desc").Find(&replies).Error
	if err != nil {
		return nil
	}

	return replies
}

func AddAReply(commentId int64, sendUserId int64, recvUserId int64, content string, level int64) (int64, error) {
	mUserReply := &MUserReply{
		CommentId:   commentId,
		SendUserId:  sendUserId,
		RecvUserId:  recvUserId,
		Content:     content,
		ReplyTime:   time.Now(),
		Level:       level,
	}
	err := db.ORM.Model(&MUserReply{}).Create(mUserReply).Error
	return mUserReply.Id, err
}

func RemoveAReply(replyId int64) error {
	return db.ORM.Delete(&MUserReply{}, "id = ?", replyId).Error
}

func MinusLikeCountByReplyId(replyId int64) error {
	err := db.ORM.Model(&MUserReply{}).
		Where("id = ?", replyId).
		Update("like_count", gorm.Expr("like_count - ?", 1)).Error

	if err != nil {
		return err
	}

	return nil
}

func AddLikeCountByReplyId(replyId int64) error {
	err := db.ORM.Model(&MUserReply{}).
		Where("id = ?", replyId).
		Update("like_count", gorm.Expr("like_count + ?", 1)).Error

	if err != nil {
		return err
	}

	return nil
}

func GetReplyByReplyId(id int64) *MUserReply {
	reply := &MUserReply{}

	//fmt.Println(id)
	err := db.ORM.Where("id = ?", id).First(reply).Error
	if err != nil  || reply.Id <= 0 {
		return nil
	}

	return reply
}
