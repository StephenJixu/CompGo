package selog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

type SeLogger struct {
	selogger *zap.Logger
}

// create standard logger
func NewStdLogger(level ...string) *SeLogger {
	encoderConfig := GeneralConfig()

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	var core zapcore.Core
	if level == nil {
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	} else {
		switch level[0] {
		case "debug":
			core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
		case "info":
			core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
		default:
			core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
		}
	}

	logger = zap.New(core)

	return &SeLogger{selogger: logger}

}

func (sl *SeLogger) Debug(msg string) {
	sl.selogger.Debug(msg)
}

func (sl *SeLogger) Info(msg string) {
	sl.selogger.Info(msg)
}

func (sl *SeLogger) Warn(msg string) {
	sl.selogger.Warn(msg)
}

func (sl *SeLogger) Error(msg string) {
	sl.selogger.Error(msg)
}

func (sl *SeLogger) Fatal(msg string) {
	sl.selogger.Fatal(msg)
}

func (sl *SeLogger) Panic(msg string) {
	sl.selogger.Panic(msg)
}

func (sl *SeLogger) Debugf(format string, args ...interface{}) {
	sl.selogger.Sugar().Debugf(format, args...)
}

func (sl *SeLogger) Infof(format string, args ...interface{}) {
	sl.selogger.Sugar().Infof(format, args...)
}

func (sl *SeLogger) Warnf(format string, args ...interface{}) {
	sl.selogger.Sugar().Warnf(format, args...)
}

func (sl *SeLogger) Errorf(format string, args ...interface{}) {
	sl.selogger.Sugar().Errorf(format, args...)
}

func (sl *SeLogger) Fatalf(format string, args ...interface{}) {
	sl.selogger.Sugar().Fatalf(format, args...)
}

func (sl *SeLogger) Panicf(format string, args ...interface{}) {
	sl.selogger.Sugar().Panicf(format, args...)
}

func (sl *SeLogger) Debugt(msg string) {
	callerFields := getCallerInfoForLog()
	logger.Debug(msg, callerFields...)
}

func (sl *SeLogger) Infot(msg string) {
	callerFields := getCallerInfoForLog()
	sl.selogger.Info(msg, callerFields...)
}

func (sl *SeLogger) Warnt(msg string) {
	callerFields := getCallerInfoForLog()
	sl.selogger.Warn(msg, callerFields...)
}

func (sl *SeLogger) Errort(msg string) {
	callerFields := getCallerInfoForLog()
	sl.selogger.Error(msg, callerFields...)
}

func (sl *SeLogger) Fatalt(msg string) {
	callerFields := getCallerInfoForLog()
	sl.selogger.Fatal(msg, callerFields...)
}

func (sl *SeLogger) Panict(msg string) {
	callerFields := getCallerInfoForLog()
	sl.selogger.Panic(msg, callerFields...)
}

// ???????????????Logger
func (sl *SeLogger) Named(name string) Logger {
	format_name := formatName(name)
	sl.selogger = sl.selogger.Named(format_name)
	return sl
}

func formatName(name string) string {
	format_str := ""
	space := " "
	name = fmt.Sprintf("[%s]", name)
	origin_length := len(name)
	if 20-origin_length < 0 {
		cut_str := name[:18]
		format_str = fmt.Sprintf("%s.]", cut_str)
	} else {
		space_count := 20 - origin_length
		end_spaces := strings.Repeat(space, space_count)
		format_str = fmt.Sprintf("%s%s", name, end_spaces)
	}
	return format_str
}

func (sl *SeLogger) With(fields ...zap.Field) Logger {
	return &SeLogger{selogger: sl.selogger.With(fields...)}
}

func Infoc(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Info(message, fields...)
}
func Debugc(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Debug(message, fields...)
}
func Errorc(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Error(message, fields...)
}
func Warnc(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Warn(message, fields...)
}

func getCallerInfoForLog() (callerFields []zap.Field) {
	pc, file, line, ok := runtime.Caller(2) // ?????????????????????????????????????????????????????????
	if !ok {
		return
	}

	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) //Base????????????????????????????????????????????????????????????
	fileName := path.Base(file)
	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", fileName), zap.Int("line", line))

	return
}

func GeneralConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	// ????????????????????????????????????
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	// ??????????????????????????????????????????????????????
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return encoderConfig
}

func init() {
	// encoderConfig := zap.NewProductionEncoderConfig()
	// // ????????????????????????????????????
	// encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	// // ??????????????????????????????????????????????????????
	// encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// // ??????Encoder ??????JSONEncoder???????????????????????????JSON?????????
	// encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// file, err := os.OpenFile("../logs/test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// fileWriteSyncer := zapcore.AddSync(file)
	// core := zapcore.NewTee(
	// 	// ??????????????????????????????????????? ?????????????????????????????????????????????????????????????????????Debug ????????????????????????????????????Info
	// 	zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	// 	zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	// )
	// logger = zap.New(core)

	logger = NewStdLogger().selogger

}
