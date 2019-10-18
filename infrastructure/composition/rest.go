package composition

import (
	"fmt"
	"net/http"

	"github.com/aboglioli/big-brother/config"
	"github.com/gin-gonic/gin"
)

func StartREST() {
	c := config.Get()
	r := gin.Default()

	r.GET("/compositions", generic)
	r.GET("/compositions/:compositionId", generic)
	r.POST("/compositions", generic)
	r.PUT("/compositions/:compositionId", generic)
	r.DELETE("/compositions/:compositionId", generic)

	r.POST("/compositions/:compositionId/dependencies", generic)
	r.PUT("/compositions/:compositionId/dependencies/:dependencyId", generic)
	r.DELETE("/compositions/:compositionId/dependencies/:dependencyId", generic)

	r.Run(fmt.Sprintf(":%d", c.Port))
}

func generic(c *gin.Context) {
	compID := c.Param("compositionId")
	depID := c.Param("dependencyId")

	c.JSON(http.StatusOK, gin.H{
		"compositionId": compID,
		"dependencyId":  depID,
	})
}
