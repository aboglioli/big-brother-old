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
			"error": gin.H{
				"code": "INTERNAL",
			},
		})
	case *errors.Status:
		fmt.Println("[Status Error]")
		fmt.Println(e)

		status := http.StatusBadRequest
		if e.Status() > 0 {
			status = e.Status()
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"status":  status,
				"code":    e.Code(),
				"message": e.Message(),
			},
		})
	case *errors.Validation:
		fmt.Println("[Validation Error]")
		fmt.Println(e)

		fields := make([]errors.Field, len(e.Fields()))
		for _, f := range e.Fields() {
			fields = append(fields, f)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    e.Code(),
				"message": e.Message(),
				"fields":  fields,
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
