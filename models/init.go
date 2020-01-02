package models

import (
	"MeTube/database"
)

func AutoMigrate() {

	// 自动迁移数据库格式
	database.ORM.AutoMigrate(&MUser{})
	database.ORM.AutoMigrate(&MVideo{})
	database.ORM.AutoMigrate(&MUserReply{})
	database.ORM.AutoMigrate(&MVideoComment{})
	database.ORM.AutoMigrate(&MReplyLike{})
	database.ORM.AutoMigrate(&MCommentLike{})
	database.ORM.AutoMigrate(&MView{})
}
