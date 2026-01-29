package logger

import (
	"log/slog"
)

// LogError 记录错误信息到error.log文件
func LogError(message string, err error, args ...any) {

	if err != nil {

		if args == nil {
			args = []any{"error", err}
			slog.Error(message, args...)
		} else {
			newArgs := make([]any, len(args))
			newArgs = append(newArgs, "error")
			newArgs = append(newArgs, message)
			for _, arg := range args {
				newArgs = append(newArgs, arg)
			}
			slog.Error(message, newArgs...)
		}
	} else {
		slog.Error(message, args...)
	}

}

// LogInfo 记录信息到控制台
func LogInfo(message string, args ...any) {
	slog.Info(message, args...)
}
