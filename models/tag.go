package models

import (
	_ "time"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	Model

	Name string `json:"name"`
	CreatedBy string `json:"create_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

// 获取符合条件的所有记录
func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var (
		tags []Tag
		err error
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

// 根据条件查询总条数
func GetTagTotal(maps interface{}) (int, error) {
	var count int

	err := db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// 判断标签名是否已存在
func ExitsTagByName(name string) (bool, error) {
	var tag Tag

	err := db.Select("id").Where("name = ? AND deleted_on = ?", name, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

// 根据ID判断标签是否存在
func ExistTagById(id int) (bool, error) {
	var tag Tag

	err := db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

// 新增标签
func AddTag(name string, state int, createBy string) error {
	tag := Tag{
		Name:       name,
		CreatedBy:  createBy,
		State:      state,
	}
	err := db.Create(&tag).Error
	if err != nil {
		return err
	}

	return nil
}

// 编辑标签信息
func EditTag(id int, data interface{}) error {
	err := db.Model(&Tag{}).Where("id = ? AND deleted_on = ?", id, 0).Update(data).Error
	if err != nil {
		return err
	}

	return nil
}

// 删除单个标签
func DeleteTag(id int) error {
	if err := db.Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}

	return nil
}

// 硬删除，硬删除要使用 Unscoped()
func CleanAllTag() (bool, error) {
	if err := db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{}).Error; err != nil {
		return false, err
	}

	return true, nil
}


/*
gorm的Callbacks

可以将回调方法定义为模型结构的指针

在创建、更新、查询、删除时将被调用
如果任何回调返回错误，gorm将停止未来操作并回滚所有更改

gorm所支持的回调方法：
创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
删除：BeforeDelete、AfterDelete
查询：AfterFind
 */
//func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//
//	return nil
//}
//
//func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//
//	return nil
//}


