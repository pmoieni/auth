package log

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var StdLogger *log.Logger // std logger returns log.Logger type which is used for http logger

func InitLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, err := os.OpenFile("./run.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	StdLogger = zap.NewStdLog(Logger)
	defer Logger.Sync()
}
