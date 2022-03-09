package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/hhong0326/goPostgresqlDocker.git/util"
)

var validCurreny validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// reflection value -> Interface
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is supported
		return util.IsSupportedCurrency(currency)
	}

	return false
}
