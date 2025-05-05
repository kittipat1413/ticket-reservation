package customvalidator

import (
	"time"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	v "github.com/kittipat1413/go-common/framework/validator"
)

// Ensure ThaiTimezoneValidator implements the CustomValidator interface.
var _ v.CustomValidator = (*ThaiTimezoneValidator)(nil)

const (
	// ThaiTimezoneValidatorTag is the tag identifier for the Thai timezone validation.
	ThaiTimezoneValidatorTag = "thaitimezone"
)

// DateValidator implements the CustomValidator interface for date validation.
type ThaiTimezoneValidator struct{}

// Tag returns the tag identifier for the date validator.
func (*ThaiTimezoneValidator) Tag() string {
	return ThaiTimezoneValidatorTag
}

// Func returns the validation function for date validation.
func (*ThaiTimezoneValidator) Func() validator.Func {
	return validateThaiTimezone
}

// Translation returns the translation text and custom translation function for the date validator.
func (*ThaiTimezoneValidator) Translation() (string, validator.TranslationFunc) {
	translationText := "{0} must be in Thai timezone (UTC+7)."

	// Custom translation function to handle parameters
	customTransFunc := func(ut ut.Translator, fe validator.FieldError) string {
		// {0} will be replaced with fe.Field()
		t, _ := ut.T(fe.Tag(), fe.Field())
		return t
	}

	return translationText, customTransFunc
}

func validateThaiTimezone(fl validator.FieldLevel) bool {
	// Check if the input is of type time.Time
	if value, ok := fl.Field().Interface().(time.Time); ok {
		// Check if the time is in the Thai time zone
		_, offset := value.Zone()
		// Check if the offset is 7 hours (25200 seconds)
		return offset == 7*60*60
	}
	return false
}
