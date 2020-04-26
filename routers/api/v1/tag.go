package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/astaxie/beego/validation"

	"github.com/Songkun007/go-gin-blog/pkg/e"
	"github.com/Songkun007/go-gin-blog/pkg/setting"
	"github.com/Songkun007/go-gin-blog/pkg/util"
	"github.com/Songkun007/go-gin-blog/pkg/app"
	"github.com/Songkun007/go-gin-blog/service/tag_service"
	"github.com/Songkun007/go-gin-blog/pkg/export"
	"github.com/Songkun007/go-gin-blog/pkg/logging"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}

	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:       name,
		State:      state,
		PageNum:    util.GetPage(c),
		PageSize:   setting.AppSetting.PageSize,
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	tags , err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = tags
	data["total"] = count

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type AddTagForm struct {
	Name		string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy 	string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     	int    `form:"state" valid:"Range(0,1)"`
}

// 新增文章标签
func AddTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddTagForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:       form.Name,
		CreatedBy:  form.CreatedBy,
		State:      form.State,
	}
	exist, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int		`form:"id" valid:"Required; Min(1)"`
	Name       string	`form:"name" valid:"Required; MaxSize(100)"`
	ModifiedBy string	`form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int		`form:"state" valid:"Range(0, 1)"`
}

// 修改文章标签
func EditTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	appG := app.Gin{C:c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()

	valid.Min(id, 1, "id").Message("ID必须大于0")
	if !valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 导出
func ExportTag(c *gin.Context) {
	appG := app.Gin{C : c}
	name := c.PostForm("name")
	state := -1

	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:       name,
		State:      state,
	}

	filename, err := tagService.Export()
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["export_url"] = export.GetExcelFullUrl(filename)
	data["export_save_url"] = export.GetExcelFullPath() + filename

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
