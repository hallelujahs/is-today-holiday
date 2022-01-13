package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hallelujahs/is-today-holiday/internal/helper"
	"github.com/hallelujahs/is-today-holiday/internal/schema"
	"github.com/winking324/rzap"
	"github.com/winking324/rzap-gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"path/filepath"
	"time"
)

func init() {
	logFile := filepath.Join(schema.GlobalEnv.LogPath, fmt.Sprintf("%s.log", schema.GlobalEnv.AppName))
	rzap.NewGlobalLogger([]zapcore.Core{
		rzap.NewCore(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    100,
			MaxAge:     10,
			MaxBackups: 10,
			Compress:   true,
		}, zap.InfoLevel),
	})
}

func main() {
	holidays := helper.Holidays{}
	holidays.Load(schema.GlobalEnv.ConfigsPath)

	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for range ticker.C {
			holidays.Load(schema.GlobalEnv.ConfigsPath)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(rzap_gin.Logger(nil), rzap_gin.Recovery(nil, true))
	router.GET("/holiday", func(c *gin.Context) {
		if holidays.IsTodayHoliday() {
			c.String(http.StatusOK, "1")
			return
		}
		c.String(http.StatusOK, "0")
	})

	if err := router.Run(fmt.Sprintf(":%d", schema.GlobalEnv.AppPort)); err != nil {
		panic(fmt.Sprintf("Setup HTTP Server failed, error: %s", err))
	}
}
