package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/astaxie/beego/validation"

	"github.com/Songkun007/go-gin-blog/pkg/e"
	"github.com/Songkun007/go-gin-blog/pkg/setting"
	"github.com/Songkun007/go-gin-blog/pkg/util"
	"github.com/Songkun007/go-gin-blog/models"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")
	maps := make(map[string] interface{})
	data := make(map[string] interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})
}

// 新增文章标签
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	// 表单验证
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS

	if ! valid.HasErrors() {
		if ! models.ExitsTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}

// 修改文章标签
func EditTag(c *gin.Context) {

}

// 删除文章标签
func DeleteTag(c *gin.Context) {

}