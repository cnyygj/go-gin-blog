package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/astaxie/beego/validation"

	"github.com/Songkun007/go-gin-blog/pkg/app"
	"github.com/Songkun007/go-gin-blog/pkg/e"
	"github.com/Songkun007/go-gin-blog/pkg/util"
	"github.com/Songkun007/go-gin-blog/service/auth_service"
)

// 定义参数校验规则
type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}


// 获取权限
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	username := c.Query("username")
	password := c.Query("password")

	a := auth{Username : username, Password : password}
	ok, _ := valid.Valid(&a)
	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}

	authService := auth_service.Auth{
		UserName: username,
		Password: password,
	}
	exists, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_AUTH, nil)
	}

	// 生产Token
	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	data := make(map[string]interface{})
	data["token"] = token

	appG.Response(http.StatusOK, e.SUCCESS, data)
}
