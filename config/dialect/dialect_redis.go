package dialect

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

const REDIS_DEFAULT_PORT uint16 = 6379

type Redis struct {
	addr    string
	options []redis.DialOption
	Values  url.Values
}

func (Redis) Name() string {
	return "redis"
}

func (Redis) ImporterPath() string {
	return "github.com/gomodule/redigo/redis"
}

func (Redis) QuoteIdent(ident string) string {
	return WrapWith(ident, "'", "'")
}

func (r *Redis) ParseDSN(params ConnParams) string {
	r.Values = url.Values{}
	r.addr = params.GetAddr("127.0.0.1", REDIS_DEFAULT_PORT)
	if params.Password != "" {
		r.options = append(r.options, redis.DialPassword(params.Password))
		r.Values.Set("auth", params.Password)
	}
	if dbno, err := strconv.Atoi(params.Database); err == nil {
		r.options = append(r.options, redis.DialDatabase(dbno))
		r.Values.Set("select", params.Database)
	}
	return r.addr
}

func (r *Redis) Connect() (redis.Conn, error) {
	if r.addr == "" {
		return nil, fmt.Errorf("the address of redis server is empty")
	}
	return redis.Dial("tcp", r.addr, r.options...)
}
