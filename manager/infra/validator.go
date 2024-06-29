package infra

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type ValidatorSupport struct {
	logger *Logger
}

func NewValidatorSupport(logger *Logger) *ValidatorSupport {
	instance := ValidatorSupport{
		logger: logger,
	}
	instance.initialize()
	return &instance
}

func (v *ValidatorSupport) initialize() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			v, ok := fld.Tag.Lookup("form")
			if !ok {
				v, ok = fld.Tag.Lookup("uri")
			}
			if ok {
				name := strings.SplitN(v, ",", 2)[0]
				if name == "-" {
					return ""
				}
				return name
			}
			return fld.Name
		})
	}
}

func (v *ValidatorSupport) Response(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}
	return errs
}
