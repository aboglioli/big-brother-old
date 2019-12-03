package errors

import (
	rawErrors "errors"
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/tests/assert"
)

func stackErrors() (error, *errors.Internal, *errors.Validation, *errors.Status) {
	err := rawErrors.New("Inner error")
	internalErr := errors.NewInternal("INTERNAL").SetPath("internal.go").SetMessage("Internal").SetRef(err)
	validationErr := errors.NewValidation("VALIDATION").SetPath("validation.go").SetMessage("Validation").SetRef(internalErr)
	validationErr.AddWithMessage("field1", "FIELD_1", "Invalid field 1")
	validationErr.AddWithMessage("field2", "FIELD_2", "Invalid field %d", 2)
	statusErr := errors.NewStatus("STATUS").SetPath("status.go").SetStatus(404).SetMessage("Status").SetRef(validationErr)

	return err, internalErr, validationErr, statusErr
}

func TestErrors(t *testing.T) {
	_, internalErr, validationErr, statusErr := stackErrors()

	assert.Err(t, internalErr)
	assert.ErrCode(t, internalErr, "INTERNAL")
	assert.ErrMessage(t, internalErr, "Internal")
	assert.NotNil(t, internalErr.Reference())

	assert.Err(t, validationErr)
	assert.ErrCode(t, validationErr, "VALIDATION")
	assert.ErrMessage(t, validationErr, "Validation")
	assert.NotNil(t, validationErr.Reference())
	assert.Equal(t, validationErr.Size(), 2)
	assert.ErrValidation(t, validationErr, "field1", "FIELD_1")
	assert.ErrValidation(t, validationErr, "field2", "FIELD_2")

	assert.Err(t, statusErr)
	assert.ErrCode(t, statusErr, "STATUS")
	assert.ErrMessage(t, statusErr, "Status")
	assert.NotNil(t, statusErr.Reference())

	log := printStack(statusErr)
	t.Log(log)
}
