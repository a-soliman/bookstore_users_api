package mysql_utils

import (
	"strings"

	"github.com/a-soliman/bookstore_utils-go/rest_errors"
	"github.com/go-sql-driver/mysql"
)

const (
	errorNoRows = "no rows in result set"
)

// ParseError given an error it trys to convert it to mySqlError, and returns the appropriate restErr
func ParseError(err error) rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return rest_errors.NewNotFoundError("no record matching given id")
		}
		return rest_errors.NewInternalServerError("error parsing database response", nil)
	}
	switch sqlErr.Number {
	case 1062:
		return rest_errors.NewBadRequestError("duplicated key")
	}
	return rest_errors.NewInternalServerError("error processing request", nil)
}
