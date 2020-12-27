package mysql_utils

import (
	"strings"

	"github.com/a-soliman/bookstore_users_api/utils/errors"
	"github.com/go-sql-driver/mysql"
)

const (
	errorNoRows = "no rows in result set"
)

// ParseError given an error it trys to convert it to mySqlError, and returns the appropriate restErr
func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerError("error parsing database response")
	}
	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError("duplicated key")
	}
	return errors.NewInternalServerError("error processing request")
}
