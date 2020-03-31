package jwt

import (
	"time"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Songkun007/go-gin-blog/pkg/util"
	"github.com/Songkun007/go-gin-blog/pkg/e"
	"github.com/Songkun007/go-gin-blog/pkg/logging"
)

// 中间件，请求方法引入中间件后，一个请求发过来，在函数调用链中，会先执行中间件的内容，真正请求的函数会被挂起（），
// 中间件的内容验证通过后，才会继续请求被挂起的函数
func JWT() gin.HandlerFunc {
	return func (c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token")

		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			logging.Info("token参数有误")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code" : code,
				"msg" : e.GetMsg(code),
				"data" : data,
			})

			// Abort 在被调用的函数中阻止挂起函数被继续调用
			c.Abort()
			return
		}

		// Next() 调用的函数中的链中执行挂起的函数
		c.Next()
	}
}
