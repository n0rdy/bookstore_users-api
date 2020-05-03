package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/n0rdy/bookstore_utils-go/rest_errors"
	"strings"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return rest_errors.NewNotFoundError("no records matching given id")
		}
		return rest_errors.NewInternalServerError("error parsing db response", err)
	}

	switch sqlErr.Number {
	case 1062:
		return rest_errors.NewBadRequestError("invalid data")
	}
	return rest_errors.NewInternalServerError("error processing request", err)
}
