package composition

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aboglioli/big-brother/composition"
	"github.com/aboglioli/big-brother/config"
	"github.com/aboglioli/big-brother/infrastructure/events"
	"github.com/gin-gonic/gin"
)

// StartREST starts API server
func StartREST() {
	// Dendencies resolution
	eventManager, err := events.GetManager()
	if err != nil {
		log.Fatal(err)
	}

	compositionRepository, rawErr := composition.NewRepository()
	if rawErr != nil {
		log.Fatal(err)
	}

	compositionService := composition.NewService(compositionRepository, eventManager)

	// Create context
	rest := &RESTContext{
		compositionService: compositionService,
	}

	// Start Gin server
	c := config.Get()
	r := gin.Default()

	// Routes
	r.GET("/composition/:compositionId", rest.GetByID)
	r.POST("/composition", rest.Post)
	r.PUT("/composition/:compositionId", rest.Put)
	r.DELETE("/composition/:compositionId", rest.Delete)

	r.Run(fmt.Sprintf(":%d", c.Port))
}

type RESTContext struct {
	compositionService composition.Service
}

func (r *RESTContext) GetByID(c *gin.Context) {
	compID := c.Param("compositionId")

	comp, err := r.compositionService.GetByID(compID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"composition": comp,
	})
}

func (r *RESTContext) Post(c *gin.Context) {
	var body composition.CreateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	comp, err := r.compositionService.Create(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "CREATED",
		"composition": comp,
	})
}

func (r *RESTContext) Put(c *gin.Context) {
	compID := c.Param("compositionId")
	var body composition.UpdateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	comp, err := r.compositionService.Update(compID, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "UPDATED",
		"composition": comp,
	})
}

func (r *RESTContext) Delete(c *gin.Context) {
	compID := c.Param("compositionId")

	if err := r.compositionService.Delete(compID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "DELETED",
	})
}
