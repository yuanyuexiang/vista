package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator 创建新的验证器
func NewValidator() *CustomValidator {
	validate := validator.New()

	// 使用JSON标签名作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validator: validate}
}

// ValidateStruct 验证结构体
func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
	return cv.validator.Struct(obj)
}

// Engine 返回底层验证器
func (cv *CustomValidator) Engine() interface{} {
	return cv.validator
}
