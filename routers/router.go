package routers

import (
	"net/http"

	"MeTube/controllers"

	"github.com/gin-gonic/gin"
)

// LoadRouters 初始化router
func LoadRouters(router *gin.Engine) {
	loadRouters(router)
}

func loadRouters(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Status": 0,
		})
	})

	router.POST("/viewvideo", controllers.ViewVideo)
	router.POST("/getvideocomments", controllers.GetVideoComments)
	router.POST("/comment", controllers.CommentAVideo)
	router.POST("/reply", controllers.ReplyAUser)
	router.POST("/replylike", controllers.ReplyLike)
	router.POST("/commentlike", controllers.CommentLike)
	router.GET("/getvideo/:videoid", controllers.GetVideo)
	router.POST("/removecomment", controllers.RemoveComment)
	router.POST("/removereply", controllers.RemoveReply)
	router.POST("/getreplies", controllers.GetReplies)
}
