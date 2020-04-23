package dialect

import (
	"net/url"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

const REDIS_DEFAULT_PORT uint16 = 6379

type Redis struct {
	options []redis.DialOption
	Values  url.Values
}

func (Redis) Name() string {
	return "redis"
}

func (Redis) QuoteIdent(ident string) string {
	return WrapWith(ident, "'", "'")
}

func (r *Redis) GetDSN(params ConnParams) string {
	r.options, r.Values = make([]redis.DialOption, 0), url.Values{}
	dsn := params.GetAddr("127.0.0.1", REDIS_DEFAULT_PORT)
	if params.Password != "" {
		r.options = append(r.options, redis.DialPassword(params.Password))
		r.Values.Set("auth", params.Password)
	}
	if dbno, err := strconv.Atoi(params.Database); err == nil {
		r.options = append(r.options, redis.DialDatabase(dbno))
		r.Values.Set("select", params.Database)
	}
	return dsn
}

func (r *Redis) Connect(params ConnParams) (redis.Conn, error) {
	dsn := r.GetDSN(params)
	return redis.Dial("tcp", dsn, r.options...)
}
