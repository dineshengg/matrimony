package middleware

import (
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	log "github.com/sirupsen/logrus"
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
	log.Infof("time it took for http handler - %s, url path - %s, user agent - %s, %s", string(ctx.Method()), string(ctx.Path()), string(ctx.UserAgent()), time.Since(durationInfo.StartTime).String())
	return nil
}
