package dialect

const PGSQL_DEFAULT_PORT uint16 = 5432

type Postgres struct {
}

func (Postgres) Name() string {
	return "postgres"
}

func (Postgres) QuoteIdent(ident string) string {
	return WrapWith(ident, `"`, `"`)
}

func (Postgres) GetDSN(params ConnParams) string {
	dsn := "user=" + params.Username
	if params.Password != "" {
		dsn += " password=" + params.Password
	}
	if params.Host != "" {
		dsn += " host=" + params.Host
	}
	if port := params.StrPort(PGSQL_DEFAULT_PORT); port != "" {
		dsn += " port=" + port
	}
	if params.Database != "" {
		dsn += " dbname=" + params.Database
	}
	return dsn
}
