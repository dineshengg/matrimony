package middleware

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

type MiddleWare struct {
	middle_router  *routing.RouteGroup
	logging        bool
	authentication bool
}

func NewMiddleWare(router *routing.RouteGroup, centralized bool, logRequest bool, authenticate bool) *MiddleWare {
	md := &MiddleWare{
		middle_router:  router,
		logging:        true,
		authentication: true,
	}

	if authenticate == true {
		md.middle_router.Use(NewAuthentication(router).Authenticate)
	}

	if logRequest == true {
		md.middle_router.Use(NewLogger(router).TimeDuration)
	}

	if centralized == true {
		md.middle_router.Use(NewLogger(router).CentralizedLogging)
	}

	if logRequest == true {
		md.middle_router.Use(NewLogger(router).LogRequest)
	}

	return md
}
