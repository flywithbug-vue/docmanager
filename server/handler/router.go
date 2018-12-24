package handler

import (
	"doc-manager/server/middleware"
	"strings"

	"github.com/gin-gonic/gin"
)

type stateType int

const (
	routerTypeNormal stateType = iota
	routerTypeNeedAuth
)

type ginHandleFunc struct {
	handler    gin.HandlerFunc
	routerType stateType
	method     string
	route      string
}




//host:port/auth_prefix/prefix/path
func RegisterRouters(r *gin.Engine, prefix string, authPrefix string) {
	jwtR := r.Group(prefix + authPrefix)
	jwtR.Use(middleware.JWTAuthMiddleware())
	for _, v := range routers {
		route := strings.ToLower(v.route)
		switch v.routerType {
		case routerTypeNeedAuth:
			funcDoRouteNeedAuthRegister(strings.ToUpper(v.method), route, v.handler, jwtR)
		case routerTypeNormal:
			route = strings.ToLower(prefix + route)
			funcDoRouteRegister(strings.ToUpper(v.method), route, v.handler, r)
		}
	}
}

//使用jwt过滤的routerType==routerTypeNeedAuth
func funcDoRouteNeedAuthRegister(method, route string, handler gin.HandlerFunc, jwtR *gin.RouterGroup) {
	switch method {
	case "POST":
		jwtR.POST(route, handler)
	case "PUT":
		jwtR.PUT(route, handler)
	case "HEAD":
		jwtR.HEAD(route, handler)
	case "DELETE":
		jwtR.DELETE(route, handler)
	case "OPTIONS":
		jwtR.OPTIONS(route, handler)
	default:
		jwtR.GET(route, handler)
	}
}

//普通routerType==routerTypeNormal
func funcDoRouteRegister(method, route string, handler gin.HandlerFunc, r *gin.Engine) {
	switch method {
	case "POST":
		r.POST(route, handler)
	case "PUT":
		r.PUT(route, handler)
	case "HEAD":
		r.HEAD(route, handler)
	case "DELETE":
		r.DELETE(route, handler)
	case "OPTIONS":
		r.OPTIONS(route, handler)
	default:
		r.GET(route, handler)
	}
}

var routers = []ginHandleFunc{
	{
		handler:htmlHandler,
		routerType:routerTypeNormal,
		method:"GET",
		route:"/doc/html",

	},
	{
		handler:    RegisterHandler,
		routerType: routerTypeNormal,
		method:     "POST",
		route:      "/register",
	},
	{
		handler:    LoginHandler,
		routerType: routerTypeNormal,
		method:     "POST",
		route:      "/login",
	},
	{
		handler:    LogoutHandler,
		routerType: routerTypeNeedAuth,
		route:      "/logout",
		method:     "POST",
	},
	{
		handler:GetUserInfoHandler,
		routerType:routerTypeNeedAuth,
		method:"GET",
		route:"/user/info",
	},
}
