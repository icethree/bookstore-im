package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/cast"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// defaultLogger logger实例
var fileWriter *FileWriter
var defaultLogger *zap.SugaredLogger
var _, timeZone = time.Now().Zone()

const (
	loggerKey = "ContextLoggerInstance"
)

// fileWriter 日志实体
type FileWriter struct {
	Path        string
	Prefix      string
	File        *os.File
	CurrentDate int64
}

var FileWriterHub sync.Map
var SysLogHub sync.Map

func NewSysLog(logPath, prefix string) *log.Logger {
	w := NewFileWriter(logPath, prefix)
	le, ok := SysLogHub.Load(w)
	if ok {
		return le.(*log.Logger)
	}

	l := log.New(w, "", log.Ldate|log.Lmicroseconds)

	SysLogHub.Store(w, l)
	return l
}

func NewFileWriter(logPath, prefix string) *FileWriter {
	w := &FileWriter{
		Path:   logPath,
		Prefix: prefix,
	}
	writer, ok := FileWriterHub.Load(w.GetKey())
	if ok {
		return writer.(*FileWriter)
	}

	FileWriterHub.Store(w.GetKey(), w)
	return w
}

func (fw *FileWriter) GetKey() string {
	return fw.Path + fw.Prefix
}

func (fw *FileWriter) Close() {
	fw.File.Close()
	FileWriterHub.Delete(fw.GetKey())
}

// Write 实现write接口
func (fw *FileWriter) Write(p []byte) (int, error) {
	fw.CheckDate()
	return fw.File.Write(p)
}

// CheckDate 检测日期，隔天就要新建log文件
func (fw *FileWriter) CheckDate() {
	var err error
	now := time.Now()
	date := (now.Unix() + int64(timeZone)) / 86400

	if date == fw.CurrentDate {
		return
	}

	fw.File.Close()

	fileName := fmt.Sprintf("%v%v%v%v", fw.Path, fw.Prefix, now.Format("20060102"), ".logs")
	fw.File, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(fmt.Sprintf("open logs file failed. err:%s, file:%s\n", err, fileName))
	}
	fw.CurrentDate = date
	fmt.Printf("logger file create. file:%v, time:%v\n", fileName, time.Now())
}

// Init 初始化Logger配置
func Init(logPath string, prefix string, level zapcore.Level) {
	fileWriter = new(FileWriter)
	fileWriter.Path = logPath
	fileWriter.Prefix = prefix

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format(time.RFC3339Nano))
	}
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = customTimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(fileWriter),
		level,
	)
	defaultLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0)).Sugar()
}

// Logger zap suger logs 接口
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
}

func Debug(v ...interface{}) {
	defaultLogger.Debug(v)
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	defaultLogger.Info(v)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	defaultLogger.Info(v)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v)
}

func Error(v ...interface{}) {
	defaultLogger.Error(v)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v)
}

func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues)
}
func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues)
}

func WithContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return defaultLogger
	}

	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return ctxLogger
	} else {
		return defaultLogger
	}
}

func NewContext(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).Desugar().With(fields...).Sugar())
}

func NewContextWithLoggerFromParent(parent context.Context) context.Context {
	ctx := context.Background()

	traceID := cast.ToString(parent.Value("id"))
	ctx = context.WithValue(ctx, "id", traceID)

	return context.WithValue(ctx, loggerKey, WithContext(ctx).Desugar().Sugar())
}
