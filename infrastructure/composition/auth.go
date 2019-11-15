package composition

import (
	"strings"

	"github.com/aboglioli/big-brother/auth"
	"github.com/aboglioli/big-brother/errors"
	"github.com/gin-gonic/gin"
)

var unauthorizedErr = errors.NewValidation().SetPath("infrastructure/composition/auth.validateAuthentication").SetCode("UNAUTHORIZED").SetMessage("Unauthorized")

func validateAuthentication(c *gin.Context) (*auth.User, errors.Error) {
	token := c.GetHeader("Authorization")

	if !strings.Contains(token, "Bearer ") {
		return nil, unauthorizedErr
	}

	token = token[7:]
	user, err := auth.Validate(token)
	if err != nil {
		return nil, unauthorizedErr
	}

	return user, nil
}

func validateAuthAndPermission(c *gin.Context, perm string) errors.Error {
	user, err := validateAuthentication(c)
	if err != nil {
		return err
	}

	if !user.HasPermission(perm) {
		return unauthorizedErr
	}

	return nil
}
