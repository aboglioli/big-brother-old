package cost

import (
	"fmt"
	"net/http"

	"github.com/aboglioli/big-brother/config"
	"github.com/gin-gonic/gin"
)

func StartREST() {
	c := config.Get()

	r := gin.Default()

	r.GET("/ping", ping)

	r.Run(fmt.Sprintf(":%d", c.CostPort))
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
