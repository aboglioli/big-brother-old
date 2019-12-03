package errors

import (
	"fmt"
	"strings"
)

type Field struct {
	Path    string `json:"path"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Validation struct {
	Generic
	fields []Field
}

func NewValidation(code string) *Validation {
	return &Validation{
		Generic: *New(code),
		fields:  []Field{},
	}
}

func (v *Validation) Add(path string, code string) *Validation {
	v.fields = append(v.fields, Field{
		Path: path,
		Code: code,
	})
	return v
}

func (v *Validation) AddWithMessage(path string, code string, msg string, args ...interface{}) *Validation {
	v.fields = append(v.fields, Field{path, code, fmt.Sprintf(msg, args...)})
	return v
}

func (v *Validation) Fields() []Field {
	return v.fields
}

func (v *Validation) Size() int {
	return len(v.fields)
}

func (v *Validation) SetPath(path string) *Validation {
	v.Generic.SetPath(path)
	return v
}

func (v *Validation) SetMessage(msg string, args ...interface{}) *Validation {
	v.Generic.SetMessage(msg, args...)
	return v
}

func (v *Validation) SetRef(err error) *Validation {
	v.Generic.SetRef(err)
	return v
}

func (v *Validation) Error() string {
	str := fmt.Sprintf("%s\n", v.Generic.Error())
	fieldsStr := []string{}
	for _, f := range v.fields {
		fieldsStr = append(fieldsStr, fmt.Sprintf("{%s, %s}", f.Path, f.Code))
	}
	str += fmt.Sprintf("\t- fields: %s", strings.Join(fieldsStr, ", "))
	return str
}
