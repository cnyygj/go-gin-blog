package tag_service

import (
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"

	"github.com/Songkun007/go-gin-blog/models"
	"github.com/Songkun007/go-gin-blog/pkg/logging"
	"github.com/Songkun007/go-gin-blog/pkg/gredis"
	"github.com/Songkun007/go-gin-blog/service/cache_service"
	"github.com/Songkun007/go-gin-blog/pkg/export"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

// 根据name来检测是否已存在相关记录
func (t *Tag) ExistByName() (bool, error) {
	return models.ExitsTagByName(t.Name)
}

// 根据ID来检测是否已存在相关记录
func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagById(t.ID)
}

// 新增tag
func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

// 编辑tag
func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State > 0 {
		data["state"] = t.State
	}

	return models.EditTag(t.ID, data)
}

// 删除tag
func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

// 统计总数
func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

// 获取tag列表
func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil ,err
	}

	err = gredis.Set(key, tags, 3600)
	if err != nil {
		logging.Warn("redis: multiple tags informations sets fail, ", err)
	}

	return tags, nil
}

// 导出功能封装
func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	// 根据标题创建列
	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		// 给每列赋上对应的值
		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + time + ".xlsx"

	fullPath := export.GetExcelFullPath() + filename
	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// 导入
func (t *Tag) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := xlsx.GetRows("标签信息")
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}

			models.AddTag(data[1], 1, data[2])
		}
	}

	return nil
}

// 条件封装
func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
