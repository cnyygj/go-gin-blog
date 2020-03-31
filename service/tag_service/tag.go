package tag_service

import (
	"encoding/json"

	"github.com/Songkun007/go-gin-blog/models"
	"github.com/Songkun007/go-gin-blog/pkg/logging"
	"github.com/Songkun007/go-gin-blog/pkg/gredis"
	"github.com/Songkun007/go-gin-blog/service/cache_service"
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
