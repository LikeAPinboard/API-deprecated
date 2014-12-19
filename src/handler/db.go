package handler

var DSN string

func SetDSN(dsn string) {
	DSN = dsn
}

func GetDSN() string {
	return DSN
}
