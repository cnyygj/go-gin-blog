package routers

import (
	"github.com/Songkun007/go-gin-blog/routers/api"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/Songkun007/go-gin-blog/middleware/jwt"
	"github.com/Songkun007/go-gin-blog/pkg/setting"
	"github.com/Songkun007/go-gin-blog/routers/api/v1"
	"github.com/Songkun007/go-gin-blog/pkg/upload"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	// 获取静态文件
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	// 权限获取
	r.GET("/auth", api.GetAuth)
	// 上传图片
	r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		// 获取标签列表
		apiv1.GET("/tags", v1.GetTags)

		// 新建标签
		apiv1.POST("/tags", v1.AddTag)

		// 更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)

		// 删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		// 获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)

		// 获取文章列表
		apiv1.GET("/articles", v1.GetArticles)

		// 新建文章
		apiv1.POST("/articles", v1.AddArticle)

		// 更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)

		// 删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
