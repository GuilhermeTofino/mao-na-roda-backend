package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware loga informações sobre cada requisição
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Status: %d | Method: %s | Path: %s | Duration: %v", c.Writer.Status(), c.Request.Method, c.Request.URL.Path, duration)
	}
}

// RecoveryMiddleware recupera de pânicos e loga o erro
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recuperado: %v\nStack trace: %s", err, string(debug.Stack()))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
			}
		}()
		c.Next()
	}
}
