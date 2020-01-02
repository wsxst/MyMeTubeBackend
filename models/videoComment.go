package models

import (
	db "MeTube/database"
	"github.com/jinzhu/gorm"
	"time"
)

type MVideoComment struct {
	Id 			   int64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	UserId         int64 `gorm:"column:user_id" json:"user_id"`
	VideoId        int64 `gorm:"column:video_id" json:"video_id"`
	Content        string `gorm:"column:content" json:"content"`
	CommentTime    time.Time `gorm:"column:comment_time" json:"comment_time"`
	LikeCount      int64 `gorm:"column:like_count" json:"like_count"`
	ReplyCount     int64 `gorm:"column:reply_count" json:"reply_count"`
}

func (p *MVideoComment) TableName() string {
	return "video_comment"
}

func GetCommentsByVideoIdSortedByTime(videoId int64, start int64, limit int64) []*MVideoComment {
	var comments []*MVideoComment

	db.ORM.Where("video_id = ?", videoId).Limit(limit).Offset((start - 1)*limit).Order("comment_time desc").Find(&comments)

	return comments
}

func GetCommentsByVideoIdSortedByHotCount(videoId int64, start int64, limit int64) []*MVideoComment {
	var comments []*MVideoComment

	// TODO:这里排序的规则还可以修改的更复杂一点
	db.ORM.Where("video_id = ?", videoId).Limit(limit).Offset((start - 1)*limit).Order("(like_count*0.7+reply_count*0.3) desc").Find(&comments)

	return comments
}

func AddAComment(userId int64, videoId int64, content string) (int64, error) {
	mVideoComment := &MVideoComment{
		UserId:      userId,
		VideoId:     videoId,
		Content:     content,
		CommentTime: time.Now(),
	}
	err := db.ORM.Model(&MVideoComment{}).Create(mVideoComment).Error
	return mVideoComment.Id, err
}

func RemoveAComment(commentId int64) error {
	return db.ORM.Delete(&MVideoComment{}, "id = ?", commentId).Error
}

func MinusLikeCountByCommentId(commentId int64) error {
	err := db.ORM.Model(&MVideoComment{}).
		Where("id = ?", commentId).
		Update("like_count", gorm.Expr("like_count - ?", 1)).Error

	if err != nil {
		return err
	}

	return nil
}

func AddLikeCountByCommentId(commentId int64) error {
	err := db.ORM.Model(&MVideoComment{}).
		Where("id = ?", commentId).
		Update("like_count", gorm.Expr("like_count + ?", 1)).Error

	if err != nil {
		return err
	}

	return nil
}

func MinusReplyCountByCommentId(commentId int64) error {
	err := db.ORM.Model(&MVideoComment{}).
		Where("id = ?", commentId).
		Update("reply_count", gorm.Expr("reply_count - ?", 1)).Error

	if err != nil {
		return err
	}

	return nil
}

func AddReplyCountByCommentId(commentId int64) error {
	err := db.ORM.Model(&MVideoComment{}).
		Where("id = ?", commentId).
		Update("reply_count", gorm.Expr("reply_count + ?", 1)).Error

	if err != nil {
		return err
	}

	return nil
}

func GetCommentByCommentId(id int64) *MVideoComment {
	comment := &MVideoComment{}

	//fmt.Println(id)
	err := db.ORM.Where("id = ?", id).First(comment).Error
	if err != nil  || comment.Id <= 0 {
		return nil
	}

	return comment
}
