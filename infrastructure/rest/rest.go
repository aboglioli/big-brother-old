package rest

import (
	"fmt"

	"github.com/aboglioli/big-brother/config"
	"github.com/gin-gonic/gin"
)

func Start() {
	c := config.Get()

	r := gin.Default()

	r.GET("/ping", ping)
	r.GET("/product/:id/stock", getStockByProductID)
	r.GET("/stock/:id", getStockByID)
	r.POST("/stock", createStock)

	r.Run(fmt.Sprintf(":%d", c.Port))
}
