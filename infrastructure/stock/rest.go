package stock

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aboglioli/big-brother/config"
	"github.com/aboglioli/big-brother/quantity"
	"github.com/aboglioli/big-brother/stock"
	"github.com/gin-gonic/gin"
)

func StartREST() {
	c := config.Get()

	r := gin.Default()

	r.GET("/ping", ping)
	r.GET("/product/:id/stock", getStockByProductID)
	r.GET("/stock/:id", getStockByID)
	r.POST("/stock", createStock)
	r.PUT("/stock/:id", updateStock)

	r.Run(fmt.Sprintf(":%d", c.StockPort))
}

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

type UpdateRequest struct {
	Action   string            `json:"action" binding:"required"`
	Quantity quantity.Quantity `json:"quantity" binding:"required"`
}

func updateStock(c *gin.Context) {
	stockID := c.Param("id")

	var body UpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println(err)
		return
	}

	stockService, err := stock.NewService()
	if err != nil {
		log.Println(err)
		return
	}

	var stock *stock.Stock
	if body.Action == "add" {
		stock, err = stockService.AddQuantity(stockID, body.Quantity)
	} else if body.Action == "substract" {
		stock, err = stockService.SubstractQuantity(stockID, body.Quantity)
	}

	if err != nil {
		log.Println(err)
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
