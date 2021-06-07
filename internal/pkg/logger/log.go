package logger

import (
	"fmt"
	"github.com/YeHeng/go-web-api/internal/pkg/plugin"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/config"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	plugin.AddPlugin(&logPlugin{})
}

type logPlugin struct {
}

func (m *logPlugin) Destroy() {
}

var log *zap.SugaredLogger

func (m *logPlugin) Init() {

	fmt.Println(color.Green("* [logging init]"))

	cfg := config.Get().Logger

	if err := os.MkdirAll(cfg.Folder, 0777); err != nil {
		fmt.Println(err.Error())
	}

	encoder := getEncoder()
	level := zapcore.DebugLevel
	_ = level.Set(cfg.Level)

	core := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook())),
		level)

	var logger *zap.Logger

	if gin.Mode() == gin.ReleaseMode {
		logger = zap.New(core)
	} else {
		logger = zap.New(core, zap.AddCaller(), zap.Development())
	}

	defer logger.Sync()
	log = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func hook() *lumberjack.Logger {
	cfg := config.Get().Logger
	return &lumberjack.Logger{
		Filename:   cfg.Folder + cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  cfg.LocalTime,
	}
}

func Get() *zap.SugaredLogger {
	return log
}