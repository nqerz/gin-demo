package main

import (
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/znqerz/gin-demo/global"
	"github.com/znqerz/gin-demo/pkg/logger"
)

var (
	configFile = kingpin.Flag("config", "Path to config file.").String()
)

func main() {
	kingpin.Version("gin-demo")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	secret := os.Getenv("SECRET_KEY_VALUE")
	version := os.Getenv("TARGET_VERSION")
	configPath := os.Getenv("CONFIG_FILE_PATH")

	if len(configPath) == 0 {
		configPath = *configFile
	}
	if len(configPath) == 0 {
		logger.Fatalf("Missing config file configuration, you can use `--config.file` or set environment variable `CONFIG_FILE_PATH`")
	}

	logger.Infof("config path: %s", configPath)
	if err := global.InitSetting(configPath); err != nil {
		logger.Fatalf("global.InitSetting err: %v", err)
	}

	callback, _ := initLogger()
	defer callback()

	gin.DefaultWriter = logger.GetWrite()
	gin.DefaultErrorWriter = logger.GetWrite()

	r := gin.Default()
	

	message := "pong; " + secret + ";" + version
	logger.Infof(message)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": message,
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func initLogger() (func(), error) {
	c := global.LogSetting
	logger.SetVersion("")
	logger.SetLevel(c.Level)
	logger.SetFormatter(c.Format)
	logger.SetReportCaller(true)

	var lumber *lumberjack.Logger
	if c.Output != "" {
		switch c.Output {
		case "stdout":
			logger.SetOutput(os.Stdout)
		case "stderr":
			logger.SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
				lumber := &lumberjack.Logger{
					Filename:   name,
					MaxSize:    10, // megabytes
					MaxBackups: 3,
					MaxAge:     28,
					Compress:   true,
				}
				logger.SetOutput(lumber)
			}
		}
	}

	return func() {
		if lumber != nil {
			lumber.Close()
		}
	}, nil
}
