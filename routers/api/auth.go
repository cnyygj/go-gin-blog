package api

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"github.com/astaxie/beego/validation"

	"github.com/Songkun007/go-gin-blog/pkg/e"
	"github.com/Songkun007/go-gin-blog/pkg/util"
	"github.com/Songkun007/go-gin-blog/models"
	"github.com/Songkun007/go-gin-blog/pkg/logging"
	"github.com/Songkun007/go-gin-blog/pkg/gredis"
)

// 定义参数校验规则
type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}


// 获取权限
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username : username, Password : password}
	ok, _ := valid.Valid(&a)

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			// 生产Token
			token, err := util.GenerateToken(username, password)

			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code = e.SUCCESS
			}

			err = gredis.Set("token", data, 0)
			if err != nil {
				logging.Warn("redis set fail, ", err)
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg"  : e.GetMsg(code),
		"data" : data,
	})
}
