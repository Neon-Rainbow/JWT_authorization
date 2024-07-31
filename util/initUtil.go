package util

import (
	"JWT_authorization/util/MySQL"
	"JWT_authorization/util/Redis"
	"JWT_authorization/util/logger"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"time"
)

func Init() error {
	// This is a placeholder for the init function

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		err := MySQL.InitMySQL()
		if err != nil {
			ctx = context.WithValue(ctx, "error", err)
			return
		}
		return
	}()

	go func() {
		defer wg.Done()
		err := Redis.InitRedis()
		if err != nil {
			ctx = context.WithValue(ctx, "error", err)
			cancel()
			return
		}
		return
	}()

	go func() {
		defer wg.Done()
		// 获取当前时间
		currentTime := time.Now()
		// 根据需求生成日志文件名
		logFileName := fmt.Sprintf("./logFile/%s.log", currentTime.Format("2006-01-02 15:04:05"))
		logger := logger.InitLogger(logFileName)
		defer logger.Sync()
		zap.ReplaceGlobals(logger)
	}()

	go func() {
		wg.Wait()
		ctx = context.WithValue(ctx, "result", true)
		cancel()
	}()

	select {
	case <-ctx.Done():
		if ctx.Err().Error() == context.DeadlineExceeded.Error() {
			return fmt.Errorf("context deadline exceeded")
		}
		if ctx.Value("error") != nil {
			return ctx.Value("error").(error)
		}
		if ctx.Value("result") != nil {
			return nil
		}
		return errors.New("unknown error")
	}
}
