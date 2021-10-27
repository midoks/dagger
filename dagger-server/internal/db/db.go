package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	// "gorm.io/driver/sqlite"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"github.com/midoks/dagger/dagger-server/internal/conf"
)

var (
	db  *gorm.DB
	err error
)

func Init() error {
	dbType := conf.GetString("db.type", "sqlite3")
	dbUser := conf.GetString("db.user", "root")
	dbPasswd := conf.GetString("db.password", "root")
	dbHost := conf.GetString("db.host", "127.0.0.1")
	dbPort, _ := conf.GetInt64("db.port", 3306)

	dbName := conf.GetString("db.name", "dagger")
	dbCharset := conf.GetString("db.charset", "utf8mb4")
	dbPath := conf.GetString("db.path", "data/imail.db3")

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
		fmt.Println("sqlite3 path:", dbPath)
		// os.MkdirAll("./data", os.ModePerm)

		// db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{SkipDefaultTransaction: true})
		// //&gorm.Config{SkipDefaultTransaction: true,}
		// // synchronous close
		// db.Exec("PRAGMA synchronous = OFF;")
	default:
		fmt.Println("database type not found")
		return errors.New("database type not found")
	}

	if err != nil {
		fmt.Println("init db err,link error:", err)
		return err
	}

	// fmt.Println("init db success!")

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
