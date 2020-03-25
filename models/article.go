package models

import (
	_ "time"
	_ "github.com/jinzhu/gorm"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`	// 实际是一个嵌套的struct，它利用TagID与Tag模型相互关联，在执行查询的时候，能够达到Article、Tag关联查询的功能

	Title string `json:"title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}


// 根据ID判断文章是否存在
func ExistArticleById(id int) bool {
	var article Article
	db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&article)

	if article.ID > 0 {
		return true
	}

	return false
}

// 获取文章总数
func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)

	return
}

// 批量获取文章
func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	// Preload 方法的参数应该是主结构体的字段名，
	// 使用预加载实际上是执行类两步操作，以上面的为例：
	// 首先是SELECT * FROM blog_articles; 和SELECT * FROM blog_tag WHERE id IN (1,2,3,4)
	// 那么在查询出结构后，gorm内部处理对应的映射逻辑，将其填充到Article的Tag中
	// 参考：http://gorm.io/zh_CN/docs/preload.html

	// 类似的关联查询有Join 和 Releated

	return
}

// 根据ID获取指定文章
func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Related(&article.Tag)

	// 两张表是如何关联起来的？
	// 首先，Article有一个结构体成员是TagID，就是外键。gorm会通过类名+ID的方式去找到这两个类之间的关联关系
	// 另外，Article有一个结构体成员是Tag，就是我们嵌套在Article里的Tag结构体，我们可以通过Related进行关联查询

	return
}

// 编辑文章
func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

// 新增文章
func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID		:      data["tag_id"].(int),
		Title		:      data["title"].(string),
		Desc		:      data["desc"].(string),
		Content		:	   data["content"].(string),
		CreatedBy	:	   data["created_by"].(string),
		State		:      data["state"].(int),
	})

	// v.(I) 为类型断言
	// v表示一个接口值，I表示接口类型，用于判断一个接口值的实际类型是否为某个类型，或一个非接口值的类型是否实现了某个接口类型

	return true
}

// 删除单篇文章
func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(&Article{})

	return true
}

// 硬删除，硬删除要使用 Unscoped()
func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{})

	return true
}

// 回调，自动更新添加时间和更新时间
//func (article *Article) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//
//	return nil
//}

//func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//
//	return nil
//}
