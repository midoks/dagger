package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"

	"github.com/midoks/dagger/dagger-server/internal/conf"
)

var (
	db  *gorm.DB
	err error
)

func Init() {
	dbType := config.GetString("db.type", "sqlite3")
	dbUser := config.GetString("db.user", "root")
	dbPasswd := config.GetString("db.password", "root")
	dbHost := config.GetString("db.host", "127.0.0.1")
	dbPort, _ := config.GetInt64("db.port", 3306)

	dbName := config.GetString("db.name", "imail")
	dbCharset := config.GetString("db.charset", "utf8mb4")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True", dbUser, dbPasswd, dbHost, dbPort, dbName, dbCharset)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("init db err,link error:", err)
		return
	}

	fmt.Println("init db success!")

	sqlDB, sqlErr := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	if sqlErr != nil {
		fmt.Println(sqlErr)
		return
	}

	db.AutoMigrate(&User{})

}
