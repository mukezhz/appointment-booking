package common

import (
	"errors"

	"clean-architecture/pkg/errorz"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func HandleDBError(err error, domainErrs map[error]bool) error {
	if err == nil {
		return nil
	}
	var customErr error
	var mysqlErr *mysql.MySQLError
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		customErr = errorz.ErrNotFound
	case errors.Is(err, gorm.ErrInvalidData):
		customErr = errorz.ErrBadRequest
	case errors.Is(err, gorm.ErrInvalidTransaction):
		customErr = errorz.ErrInvalidTransactionState
	case errors.Is(err, gorm.ErrDuplicatedKey):
		customErr = errorz.ErrAlreadyExists
	case errors.As(err, &mysqlErr) && mysqlErr.Number == 1062:
		customErr = errorz.ErrAlreadyExists
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		customErr = errorz.ErrAlreadyExists
	default:
		customErr = errorz.ErrSomethingWentWrong
	}

	for domainErr := range domainErrs {
		if errors.Is(domainErr, customErr) {
			customErr = domainErr
			break
		}
	}
	if customErr != nil {
		return customErr
	}

	return err
}
