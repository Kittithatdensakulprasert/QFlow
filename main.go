package main

import (
	"qflow/config"
	"qflow/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	r := gin.Default()
	router.Setup(r)

	r.Run(":" + cfg.Port)
}
