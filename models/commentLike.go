package models

import (
	db "MeTube/database"
)

type MCommentLike struct {
	Id 			  int64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CommentId     int64 `gorm:"column:comment_id" json:"comment_id"`
	UserId        int64 `gorm:"column:user_id" json:"user_id"`
	Flag          bool `gorm:"column:flag" json:"flag"`
}

func (p *MCommentLike) TableName() string {
	return "comment_like"
}

func AddLikeByCommentId(commentId int64, userId int64) error {
	commentLike :=&MCommentLike{}
	db.ORM.Where("user_id = ? and comment_id = ?", userId, commentId).First(commentLike)
	if commentLike.Id > 0 {
		err := db.ORM.Model(&MCommentLike{}).Where("user_id = ? and comment_id = ?", userId, commentId).Update("flag", true).Error
		if err != nil {
			return err
		}
	} else {
		commentLike = &MCommentLike{CommentId:commentId, UserId:userId, Flag:true}
		db.ORM.Create(commentLike)
	}

	return nil
}

func CancelLikeByCommentId(commentId int64, userId int64) error {
	err := db.ORM.Model(&MCommentLike{}).Where("user_id = ? and comment_id = ?", userId, commentId).Update("flag", false).Error
	if err != nil {
		return err
	}

	return nil
}

func JudgeLikeComment(userId int64, commentId int64) bool {
	commentLike :=&MCommentLike{}
	db.ORM.Where("user_id = ? and comment_id = ?", userId, commentId).First(commentLike)
	if commentLike.Id > 0 && commentLike.Flag {
		return true
	} else {
		return false
	}
}

func CountLikeComment(commentId int64) int64 {
	var likeCount int64
	db.ORM.Where("comment_id = ? and flag = TRUE", commentId).Count(&likeCount)
	return likeCount
}
