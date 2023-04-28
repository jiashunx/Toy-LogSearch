package log

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "os"
    "path"
    "runtime"
)

var logger *zap.Logger

func init() {
    ec := zap.NewProductionEncoderConfig()
    // 设置日志记录时间格式
    ec.EncodeTime = zapcore.ISO8601TimeEncoder
    // 日志输出为json格式
    encoder := zapcore.NewJSONEncoder(ec)
    _ = os.MkdirAll("./logs/Toy-LogSearch", 0777)
    // 日志文件
    file, _ := os.OpenFile("./logs/Toy-LogSearch/Toy-LogSearch.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    syncer := zapcore.AddSync(file)
    core := zapcore.NewTee(
        // 同时向控制台和文件写日志
        // zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
        zapcore.NewCore(encoder, syncer, zapcore.DebugLevel),
    )
    logger = zap.New(core)
}

func Info(message string, fields ...zap.Field) {
    callerFields := getCallerInfoForLog()
    fields = append(fields, callerFields...)
    logger.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
    callerFields := getCallerInfoForLog()
    fields = append(fields, callerFields...)
    logger.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
    callerFields := getCallerInfoForLog()
    fields = append(fields, callerFields...)
    logger.Error(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
    callerFields := getCallerInfoForLog()
    fields = append(fields, callerFields...)
    logger.Warn(message, fields...)
}

func getCallerInfoForLog() (callerFields []zap.Field) {
    // 回溯两层，拿到写日志的调用方的函数信息
    pc, file, line, ok := runtime.Caller(2)
    if !ok {
        return
    }
    funcName := runtime.FuncForPC(pc).Name()
    //Base函数返回路径的最后一个元素，只保留函数名
    funcName = path.Base(funcName)
    callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
    return
}
