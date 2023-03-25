package service

import "strings"

// check the table name and database name is correct
func invalidnameset(tablename, dbname string) (res bool) {
	if len(tablename) > 0 && len(dbname) > 0 && strings.ContainsRune(tablename, ' ') && strings.ContainsRune(dbname, ' ') {
		res = true
	} else {
		res = false
	}

	return
}
