package routes

import (
	"producer/logger"

	"github.com/gin-gonic/gin"
)

// Route temp
func Route(listenAddr, postAddr string, postData gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.POST(postAddr, postData)
	for _, routeInfo := range router.Routes() {
		logger.SugarLogger.Debug("path: ", routeInfo.Path, "\thandler: ", routeInfo.Handler, "\tmethod: ", routeInfo.Method, ",\tregistered routes")
	}
	return router
}
