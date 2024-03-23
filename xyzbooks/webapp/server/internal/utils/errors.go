package utils

func IsSQLNoRowsErr(err error) bool {
	return err.Error() == "sql: no rows in result set"
}
