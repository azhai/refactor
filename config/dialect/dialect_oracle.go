package dialect

const ORACLE_DEFAULT_PORT uint16 = 1521

type Oracle struct {
}

func (Oracle) Name() string {
	return "oracle"
}

func (Oracle) QuoteIdent(ident string) string {
	return WrapWith(ident, "{", "}")
}

func (Oracle) GetDSN(params ConnParams) string {
	user := ConcatWith(params.Username, params.Password)
	dsn := "oracle://"
	if user != "" {
		dsn += user + "@"
	}
	dsn += params.GetAddr("127.0.0.1", ORACLE_DEFAULT_PORT)
	if params.Database != "" {
		dsn += "?database" + params.Database
	}
	return dsn
}
