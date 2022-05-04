package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/midoks/dagger/dagger-server/internal/conf"
)

var (
	db  *gorm.DB
	err error
)

func Init() error {
	dbType := conf.Database.Type
	dbUser := conf.Database.User
	dbPasswd := conf.Database.Password
	dbHost := conf.Database.Host
	dbPort := conf.Database.Port

	dbName := conf.Database.Name
	dbCharset := conf.Database.Charset
	dbPath := conf.Database.Path

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)

	switch dbType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True", dbUser, dbPasswd, dbHost, dbPort, dbName, dbCharset)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	case "sqlite3":
		// fmt.Println("sqlite3 path:", dbPath)
		os.MkdirAll("./data", os.ModePerm)
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{SkipDefaultTransaction: true})
		// &gorm.Config{SkipDefaultTransaction: true,}
		// // synchronous close
		db.Exec("PRAGMA synchronous = OFF;")
	default:
		fmt.Println("database type not found")
		return errors.New("database type not found")
	}

	if err != nil {
		fmt.Println("init db err,link error:", err)
		return err
	}

	sqlDB, sqlErr := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if sqlErr != nil {
		fmt.Println(sqlErr)
		return sqlErr
	}

	db.AutoMigrate(&User{})
	return nil
}
