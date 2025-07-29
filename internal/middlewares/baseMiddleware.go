package middlewares

import "github.com/gin-gonic/gin"

type IBaseMiddleware interface {
	Handle() gin.HandlerFunc
}
