package composition

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/gin-gonic/gin"
)

func validateAuthAndPermission(c *gin.Context, perm string) errors.Error {
	return nil
}
