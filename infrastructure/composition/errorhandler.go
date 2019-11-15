package composition

import (
	"fmt"
	"net/http"

	"github.com/aboglioli/big-brother/errors"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err errors.Error) {
	switch e := err.(type) {
	case *errors.ValidationError:
		fmt.Println("[Validation Error]")
		fmt.Println(e)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    e.Code(),
				"message": e.Message(),
			},
		})
	case *errors.InternalError:
		fmt.Println("[Internal Error]")
		fmt.Println(e)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL_ERROR",
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "UNKNOWN_ERROR",
		})
	}
}
