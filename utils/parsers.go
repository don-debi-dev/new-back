package utils

import "database/sql"

// ParseNullString is a utility function which parse sql null string to regular strings
func ParseNullString(p sql.NullString) string {
	if p.Valid == true {
		return p.String
	}
	return ""
}
