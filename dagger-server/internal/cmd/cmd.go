package cmd

import (
	"time"

	"github.com/midoks/dagger/dagger-server/internal/conf"
	"github.com/midoks/dagger/dagger-server/internal/db"
	go_logger "github.com/phachon/go-logger"
	"github.com/urfave/cli"
)

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

//nolint:deadcode,unused
func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

//nolint:deadcode,unused
func durationFlag(name string, value time.Duration, usage string) cli.DurationFlag {
	return cli.DurationFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func initLogger() {

	logger = go_logger.NewLogger()

	// 文件输出配置
	fileConfig := &go_logger.FileConfig{
		Filename: "./logs/info.log", // 日志输出文件名，不自动存在
		// 如果要将单独的日志分离为文件，请配置LealFrimeNem参数。
		LevelFileName: map[int]string{
			logger.LoggerLevel("error"): "./logs/error.log", // Error 级别日志被写入 error .log 文件
			logger.LoggerLevel("debug"): "./logs/debug.log", // Debug 级别日志被写入到 debug.log 文件中
		},
		MaxSize:    1024 * 1024, // 文件最大值（KB），默认值0不限
		MaxLine:    100000,      // 文件最大行数，默认 0 不限制
		DateSlice:  "d",         // 文件根据日期切分， 支持 "Y" (年), "m" (月), "d" (日), "H" (时), 默认 "no"， 不切分
		JsonFormat: false,       // 写入文件的数据是否 json 格式化
		Format:     "",          // 如果写入文件的数据不 json 格式化，自定义日志格式
	}
	// 添加 file 为 logger 的一个输出
	logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)

	// logger.Infof("hello,world,now:%s", time.Now().Format("2006/1/2 15:04:05"))
}

func Init() {
	// fmt.Println("cmd init")
	conf.Load("conf/app.conf")
	initLogger()
	db.Init()
}
