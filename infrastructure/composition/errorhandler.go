package composition

import (
	"fmt"
	"net/http"

	"github.com/aboglioli/big-brother/pkg/errors"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err errors.Error) {
	switch e := err.(type) {
	case *errors.CustomError:
		switch e.Kind() {
		case errors.VALIDATION:
			fmt.Println("[Validation Error]")
			fmt.Println(e)

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    e.Code(),
					"message": e.Message(),
				},
			})
		case errors.INTERNAL:
			fmt.Println("[Internal Error]")
			fmt.Println(e)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "INTERNAL_ERROR",
			})
		}
	default:
		fmt.Println("[Unknown error]")
		fmt.Println(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "UNKNOWN_ERROR",
		})
	}
}
