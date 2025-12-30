package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func HandleValidationErrors(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		

		for _ ,e := range validationError{
           

			root := strings.Split(e.Namespace(), ".")[0]
			rawPath := strings.TrimPrefix(e.Namespace(), root + ".")
			parts := strings.Split(rawPath, ".")

			 for i, part := range parts {
				if strings.Contains(part, "[") {
					idx := strings.Index(part, "[")
					base := camelToSnake(part[:idx])
					index := part[idx:]
					parts[i] = base + index
				} else {
					parts[i] = camelToSnake(part)
				}
				}

				fieldPath := strings.Join(parts, ".")
			
			switch e.Tag() {
			case "required":
				errors[fieldPath] = fmt.Sprintf("Field '%s' is required", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be a valid email address", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be at least %s characters long", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be at most %s characters long", fieldPath, e.Param())
			case "gt":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be greater than %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be less than %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be greater than or equal to %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be less than or equal to %s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be a valid UUID", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be a valid slug", fieldPath)
			case "oneof":
				allowedValues := strings.Split(e.Param(), " ")
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be one of the following: %s", fieldPath, strings.Join(allowedValues, ", "))
			case "search":
					errors[fieldPath] = fmt.Sprintf("Field '%s' must be a valid search query", fieldPath)
			case "datetime":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be a valid date and time", fieldPath)
			case "min_int":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be at least %s", fieldPath, e.Param())
			case "max_int":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be at most %s", fieldPath, e.Param())
			case "file_ext":
				errors[fieldPath] = fmt.Sprintf("Field '%s' must be a valid file extension", fieldPath)
			}
			


		}
		return gin.H{
			"error": errors,
		}
	}

	return gin.H{
		"error": "Wrong Request: " + err.Error(),
	}
}

func RegisterValidator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("validator engine not found")
	}

	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("min_int", func(fl validator.FieldLevel) bool {
		minStr := fl.Param()
		minInt, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return false
		}
	 return fl.Field().Int() >= minInt
	})

	v.RegisterValidation("max_int", func(fl validator.FieldLevel) bool {
		maxStr := fl.Param()
		maxInt, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			return false
		}
		return fl.Field().Int() <= maxInt
	})

	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {

		fileName := fl.Field().String()
		allowedStr := fl.Param()
		if allowedStr == "" {
			return false
		}

		allowedExt := strings.Fields(allowedStr)

		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")

		return slices.Contains(allowedExt, ext)
	})

	return nil
}


