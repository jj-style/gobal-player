package server

import (
	"sync"

	"github.com/gin-gonic/gin"
)

func loadTemplates(r *gin.Engine) gin.HandlerFunc {
	var once sync.Once
	return func(c *gin.Context) {
		once.Do(func() {
			r.LoadHTMLGlob("templates/*")
		})
		c.Next()
	}
}
