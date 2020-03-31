package article_service

import (
	"encoding/json"
	"github.com/Songkun007/go-gin-blog/models"
	"github.com/Songkun007/go-gin-blog/pkg/gredis"
	"github.com/Songkun007/go-gin-blog/pkg/logging"
	"github.com/Songkun007/go-gin-blog/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

// 添加文章
func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}

	return nil
}

// 编辑文章
func (a *Article) Edit() error {
	data := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	}
	err := models.EditArticle(a.ID, data)
	if err == nil {
		// 数据更新，清理redis中的缓存
		if delErr := gredis.LikeDeletes("ARTICLE_"); delErr != nil {
			logging.Info("redis: keys delete fail, ", delErr)
		}
	}

	return err
}

// 获取单篇文章
func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	err = gredis.Set(key, article, 3600)
	if err != nil {
		logging.Warn("redis: single article informations sets fail, article_id:", a.ID, err)
	}

	return article, nil
}

// 获取所有文章
func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)

	cache := cache_service.Article{
		TagID:    a.TagID,
		State:    a.State,
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	err = gredis.Set(key, articles, 3600)
	if err != nil {
		logging.Warn("redis: multiple article informations sets fail, ", err)
	}

	return articles, nil
}

// 删除一片文章
func (a *Article) Delete() error {
	err := models.DeleteArticle(a.ID)
	if err == nil {
		// 删除redis缓存
		delErr := gredis.LikeDeletes("ARTICLE_")
		if delErr != nil {
			logging.Warn("redis: single article informations delete fail, article_id:", a.ID, delErr)
		}
	}

	return err
}

// 根据id判断文章是否存在
func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleById(a.ID)
}

// 统计符合条件的文章的总条数
func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

// 条件封装
func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}