package storage

// import (
// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgconn"
// )

// // func ErrorAlreadyExists(err error) bool {
// // 	if pgErr, ok := err.(*pgconn.PgError); ok {
// // 		return pgErr.Code == "23505"
// // 	}
// // 	return false
// // }

// // func CheckErr(err error) bool {
// // 	if pgErr, ok := err.(*pgconn.PgError); ok {
// // 		return pgErr.Code == "23514"
// // 	}
// // 	return false
// // }

// // func ErrNotFound() error {
// // 	return pgx.ErrNoRows
// // }

// func ErrorAlreadyExists(err error) bool {
// 	if pgErr, ok := err.(*pgconn.PgError); ok {
// 		return pgErr.Code == "23505"
// 	}
// 	return false
// }

// func CheckErr(err error) bool {
// 	if pgErr, ok := err.(*pgconn.PgError); ok {
// 		return pgErr.Code == "23514"
// 	}
// 	return false
// }

// func ErrNotFound() error {
// 	return pgx.ErrNoRows
// }
import (
	"database/sql"
	"strings"
)

func ErrorAlreadyExists(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func ErrNotFound() error {
	return sql.ErrNoRows
}
