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

func init() {
	var (
		err                                       error
		dbName, user, password, host, tablePrefix string
	)
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatalf("Fail to get section 'database': %v", err)
		// log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
	})

	db.Callback().Create().Before("gorm:create").Register("update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Before("gorm:update").Register("update_time_stamp", updateTimeStampForUpdateCallback)

	if err != nil {
		log.Println(err)
	}
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
