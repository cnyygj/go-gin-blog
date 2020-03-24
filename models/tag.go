package models

import (
	_ "time"
	_ "github.com/jinzhu/gorm"
)

type Tag struct {
	Model

	Name string `json:"name"`
	CreatedBy string `json:"create_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

// 获取符合条件的所有记录
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

// 根据条件查询总条数
func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

// 判断标签名是否已存在
func ExitsTagByName(name string) bool {
	var tag Tag

	db.Select("id").Where("name = ?", name).First(&tag)

	if tag.ID > 0 {
		return true
	}

	return false
}

// 根据ID判断标签是否存在
func ExistTagById(id int) bool {
	var tag Tag

	db.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

// 新增标签
func AddTag(name string, state int, createBy string) bool {
	db.Create(&Tag {
		Name:       name,
		CreatedBy:  createBy,
		State:      state,
	})

	return true
}

// 编辑标签信息
func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Update(data)

	return true
}

// 删除单个标签
func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})

	return true
}

// 硬删除，硬删除要使用 Unscoped()
func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})

	return true
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


