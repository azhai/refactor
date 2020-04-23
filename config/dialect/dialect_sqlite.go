package dialect

type Sqlite struct {
}

func (Sqlite) Name() string {
	return "sqlite"
}

func (Sqlite) QuoteIdent(ident string) string {
	return WrapWith(ident, "`", "`")
}

func (Sqlite) GetDSN(params ConnParams) string {
	user := ConcatWith(params.Username, params.Password)
	var dsn string
	if user != "" {
		dsn = user + "@"
	}
	dsn += params.Database + "?cache=shared&mode=rwc"
	return dsn
}
