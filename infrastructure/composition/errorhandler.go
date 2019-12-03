package composition

import (
	"fmt"
	"net/http"

	"github.com/aboglioli/big-brother/pkg/errors"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errors.Internal:
		fmt.Println("[Internal Error]")
		fmt.Println(e)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL_ERROR",
		})
	case *errors.Validation:
		fmt.Println("[Validation Error]")
		fmt.Println(e)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    e.Code(),
				"message": e.Message(),
			},
		})
	default:
		fmt.Println("[Unknown error]")
		fmt.Println(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "UNKNOWN_ERROR",
		})
	}
}
