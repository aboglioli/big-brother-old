package errors

import (
	"fmt"
	"net/http"

	"github.com/aboglioli/big-brother/pkg/errors"

	"github.com/gin-gonic/gin"
)

func Handle(c *gin.Context, err error) {
	printStack(err)

	switch e := err.(type) {
	case *errors.Unauthorized:
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code": e.Code(),
			},
		})
	case *errors.Internal:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": "INTERNAL",
			},
		})
	case *errors.Status:
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    e.Code(),
				"message": e.Message(),
				"fields":  e.Fields(),
			},
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "UNKNOWN",
		})
	}
}

func printStack(err error) string {
	str := ""
	if err != nil {
		switch e := err.(type) {
		case *errors.Unauthorized:
			str += fmt.Sprintf("[Authorization]: %s\n", e)
		case *errors.Internal:
			str += fmt.Sprintf("[Internal]: %s\n", e)
		case *errors.Status:
			str += fmt.Sprintf("[Status]: %s\n", e)
		case *errors.Validation:
			str += fmt.Sprintf("[Validation]: %s\n", e)
		default:
			str += fmt.Sprintf("[Unknown]: %s\n", e)
		}

		ref, ok := err.(errors.Reference)
		if ok {
			str += printStack(ref.Reference())
		}
	}

	return str
}
