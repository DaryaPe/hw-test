package dbstorage

import (
	"database/sql"
)

// StringToNull преобразует string в sql.NullString.
func StringToNull(value string) sql.NullString {
	return sql.NullString{String: value, Valid: len(value) > 0}
}
