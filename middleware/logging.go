package middleware

import (
	"log"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
)

type DurationInfo struct {
	StartTime time.Time
}

type Logger struct {
	log_router *routing.RouteGroup
}

func NewLogger(router *routing.RouteGroup) *Logger {
	log := &Logger{
		log_router: router,
	}
	return log
}

func (logger *Logger) CentralizedLogging(ctx *routing.Context) error {
	//log all the data to centralized logging system like elastic search
	return nil
}

func (logger *Logger) TimeDuration(ctx *routing.Context) error {
	t1 := time.Now()
	ctx.Set("DurationInfo", DurationInfo{StartTime: t1})

	return nil
}

func (logger *Logger) LogRequest(ctx *routing.Context) error {
	durationInfo := ctx.Get("DurationInfo").(DurationInfo)
	log.Println("time it took for http handler - %s, url path - %s, user agent - %s, %v", ctx.Method(), ctx.Path(), ctx.UserAgent(), time.Since(durationInfo.StartTime))
	return nil
}
