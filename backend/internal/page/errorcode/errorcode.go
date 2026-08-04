package errorcode

const (
	AccountExists      = "a"
	UnmatchedPasswords = "b"
	IncorrectUsername  = "c"
	IncorrectPassword  = "d"
	InternalSqlError   = "e"
	InternalError      = "f"
)

func CodeToMessage(errID string) string {
	errMap := map[string]string{
		AccountExists:      "That account already exists!",
		UnmatchedPasswords: "The passwords do not match!",
		IncorrectUsername:  "Incorrect username!",
		IncorrectPassword:  "Incorrect password!",
		InternalSqlError:   "Internal Server Error: Running SQL query failed!",
		InternalError:      "Internal Server Error",
	}
	if _, ok := errMap[errID]; ok {
		return errMap[errID]
	}
	return "Unknown error"
}
