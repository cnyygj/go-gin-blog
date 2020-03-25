package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/Songkun007/go-gin-blog/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID int `gorm:"primary_key" json:"id"`
	CreatedOn int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn int `json:"deleted_on"`
}

func Setup() {

	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	// gorm默认表名是结构体名称的复数，如 type User struct {} // 默认表名是`users`
	// 可以通过定义DefaultTableNameHandler对默认表名应用任何规则
	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)			// 全局禁用表名复数，如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	db.LogMode(true)				// 启用Logger，显示详细日志
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// 替换替换默认的钩子
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)

	// Define callbacks for creating
	//func init() {
	//	DefaultCallback.Create().Register("gorm:begin_transaction", beginTransactionCallback)
	//	DefaultCallback.Create().Register("gorm:before_create", beforeCreateCallback)
	//	DefaultCallback.Create().Register("gorm:save_before_associations", saveBeforeAssociationsCallback)
	//	DefaultCallback.Create().Register("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//	DefaultCallback.Create().Register("gorm:create", createCallback)
	//	DefaultCallback.Create().Register("gorm:force_reload_after_create", forceReloadAfterCreateCallback)
	//	DefaultCallback.Create().Register("gorm:save_after_associations", saveAfterAssociationsCallback)
	//	DefaultCallback.Create().Register("gorm:after_create", afterCreateCallback)
	//	DefaultCallback.Create().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
	//}

	// 注册回调
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()

		// 通过 scope.Fields() 获取所有字段，判断当前是否包含所需字段
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {	// 可判断该字段的值是否为空
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifyTime` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {

	// scope.Get(...) 根据入参获取设置了字面值的参数，例如本文中是 gorm:update_column ，它会去查找含这个字面值的字段属性
	if _, ok := scope.Get("gorm:update_column"); !ok {
		// scope.SetColumn(...) 假设没有指定 update_column 的字段，我们默认在更新回调设置 ModifiedOn 的值
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 删除钩子函数
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string

		// 检查是否手动指定了 delete_option
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		// 获取我们约定的删除字段，若存在则 UPDATE 软删除，若不存在则 DELETE 硬删除
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				//  返回引用的表名
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				// scope.AddToVars 该方法可以添加值作为 SQL 的参数，也可用于防范 SQL 注入
				scope.AddToVars(time.Now().Unix()),
				// scope.QuotedTableName() 返回组合好的条件 SQL
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

func CloseDB() {
	defer db.Close()
}


