package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/MrHanson/gin-blog/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func Setup() {
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", setting.DatabaseSetting.User, setting.DatabaseSetting.Password, setting.DatabaseSetting.Host, setting.DatabaseSetting.Name)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   setting.DatabaseSetting.TablePrefix,
			SingularTable: true,
		},
	})

	if err != nil {
		log.Println(err)
	}

	db.Callback().Create().Before("gorm:create").Register("update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Before("gorm:update").Register("update_time_stamp", updateTimeStampForUpdateCallback)
}

func CloseDB() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}

func updateTimeStampForCreateCallback(db *gorm.DB) {
	nowTime := time.Now().Unix()
	updateFieldValue("CreatedOn", nowTime, db)
	updateFieldValue("ModifiedOn", nowTime, db)
}

func updateTimeStampForUpdateCallback(db *gorm.DB) {
	nowTime := time.Now().Unix()
	updateFieldValue("ModifiedOn", nowTime, db)
}

func updateFieldValue(name string, value interface{}, db *gorm.DB) {
	schema := db.Statement.Schema
	if schema == nil {
		return
	}
	if modifyTimeField, ok := schema.FieldsByName[name]; ok {
		modifyTimeField.Set(db.Statement.ReflectValue, value)
	}
}
