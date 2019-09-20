package rest

import (
	"log"
	"net/http"

	"github.com/aboglioli/big-brother/stock"
	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func createStock(c *gin.Context) {
	var body stock.CreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println(err)
		return
	}

	stockService, err := stock.NewService()
	if err != nil {
		log.Println(err)
		return
	}

	stock, err := stockService.Create(body)
	if err != nil {
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, stockToJSON(stock))
}

func getStockByID(c *gin.Context) {
	id := c.Param("id")

	stockService, err := stock.NewService()
	if err != nil {
		log.Println(err)
		return
	}

	stock, err := stockService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(), // TODO: errors from infrastructure!!!
		})
		return
	}

	c.JSON(http.StatusOK, stockToJSON(stock))
}

func getStockByProductID(c *gin.Context) {
	id := c.Param("id")

	stockService, err := stock.NewService()
	if err != nil {
		log.Println(err)
		return
	}

	stock, err := stockService.GetByProductID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(), // TODO: errors from infrastructure!!!
		})
		return
	}

	c.JSON(http.StatusOK, stockToJSON(stock))
}

func stockToJSON(s *stock.Stock) gin.H {
	return gin.H{
		"id":          s.ID.String(),
		"productId":   s.ProductID.String(),
		"productName": s.ProductName,
		"quantity":    s.Quantity,
	}
}
