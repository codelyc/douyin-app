package handler

import (
	"douyin-app/domain"
	"douyin-app/service"
	"douyin-app/util"
	"github.com/gin-gonic/gin"
	"log"
)

/*
	点赞操作：
		1： 点赞，直接在favorite表创建 userid与favoriteId进行关联的数据，同时更新video表中 FavoriteCount 属性值
		2： 取消点赞，直接删除favorite表中这条记录，同时更新video表中 FavoriteCount 属性值

	actionTYpe == 1 点赞
		1) 先查询有无favorite record,无则创建，有则delete_time = 0
	actionType == 2 取消点赞
		1) delete_time = now()
*/
func VideoFavoriteHandler(c *gin.Context) {
	//token
	token := c.Query("token")
	if token == "" {
		log.Println("VideoFavoriteHandler Token <nil>")
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
		return
	}

	//userId
	userId, err := util.ParseToken(token)
	if err != nil {
		log.Println("VideoFavoriteHandler ParseToken Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.WrongAuth,
		})
		return
	}

	//videoId
	vid, err := util.Str2Int64(c.Query("video_id"))
	if err != nil {
		log.Println("VideoFavoriteHandler Str2Int64 Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}

	//actionType
	actionType := c.Query("action_type")

	//favoriteAction
	err = service.FavoriteAction(c, userId, vid, actionType)
	if err != nil {
		log.Println("VideoFavoriteHandler FavoriteAction Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.InternalServerError,
		})
		return
	}

	//response
	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
	})
}

/*
	返回点赞list
 	1. 从favorite表中找出 userid= ?? 符合所有的favorite记录存储在favoriteList中
	2. 遍历favoriteList得到videoList
*/
func VideoFavoriteListHandler(c *gin.Context) {
	var (
		userId   int64
		err      error
		videoDOs []*domain.VideoDO
	)

	//token
	token := c.Query("token")
	if token != "" {
		userId, err = util.ParseToken(token)
		if err != nil {
			log.Println("VedioFavoriteListHandler ParseToken Failed", err)
			util.MakeResponse(c, &util.HttpResponse{
				StatusCode: util.WrongAuth,
			})
			return
		}
	}

	//favoriteList
	videoDOs, err = service.GetFavoriteList(c, userId)
	if err != nil {
		log.Println("VedioFavoriteListHandler GetFavoriteList Failed", err)
		util.MakeResponse(c, &util.HttpResponse{
			StatusCode: util.ParamError,
		})
		return
	}

	//response
	util.MakeResponse(c, &util.HttpResponse{
		StatusCode: util.Success,
		ReturnVal: map[string]interface{}{
			"video_list": videoDOs,
		},
	})
}
