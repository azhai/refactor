package common

import (
	"fmt"
	"strings"

	"gitea.com/azhai/refactor/config"
	"gitea.com/azhai/refactor/config/dialect"
	"github.com/azhai/gozzo-utils/redisw"
	"github.com/gomodule/redigo/redis"
	"github.com/k0kubun/pp"
)

const (
	MAX_TIMEOUT     = 86400 * 30 // 接近无限时间
	SESS_ONLINE_KEY = "onlines"  // 在线用户
	SESS_TOKEN_KEY  = "_token_"
	SESS_PREFIX     = "sess" // 会话缓存前缀
	SESS_TIMEOUT    = 7200   // 会话缓存时间
	SESS_LIST_SEP   = ";"    // 角色名之间的分隔符
)

func InitCache(cfg *config.Settings, name string, verbose bool) (*SessionRegistry, error) {
	var sessreg *SessionRegistry
	c, ok := cfg.Connections[name]
	if !ok || c.DriverName != "redis" {
		return nil, nil
	}
	if verbose {
		d := dialect.GetDialectByName(c.DriverName)
		drv, dsn := d.Name(), d.GetDSN(c.Params)
		args := d.(*dialect.Redis).Values.Encode()
		pp.Println(drv, dsn, args)
	}
	sessreg = NewRegistry(redisw.ConnParams(c.Params))
	return sessreg, nil
}

func SessListJoin(data []string) string {
	return strings.Join(data, SESS_LIST_SEP)
}

func SessListSplit(data string) []string {
	return strings.Split(data, SESS_LIST_SEP)
}

type SessionRegistry struct {
	sessions map[string]*Session
	Onlines  *redisw.RedisHash
	*redisw.RedisWrapper
}

func NewRegistry(params redisw.ConnParams) *SessionRegistry {
	r := redisw.NewRedisPool(params, -1)
	return &SessionRegistry{
		sessions:     make(map[string]*Session),
		Onlines:      redisw.NewRedisHash(r, SESS_ONLINE_KEY, MAX_TIMEOUT),
		RedisWrapper: r,
	}
}

func (sr SessionRegistry) GetKey(token string) string {
	return fmt.Sprintf("%s:%s", SESS_PREFIX, token)
}

func (sr *SessionRegistry) GetSession(token string) *Session {
	key := sr.GetKey(token)
	if sess, ok := sr.sessions[key]; ok && sess != nil {
		return sess
	}
	sess := NewSession(sr, key)
	if _, err := sess.SetVal(SESS_TOKEN_KEY, token); err == nil {
		sr.sessions[key] = sess
	}
	return sess
}

func (sr *SessionRegistry) DelSession(token string) bool {
	key := sr.GetKey(token)
	if sess, ok := sr.sessions[key]; ok {
		succ, err := sess.DeleteAll()
		if succ && err == nil {
			delete(sr.sessions, key)
			return true
		}
	}
	return false
}

type Session struct {
	reg *SessionRegistry
	*redisw.RedisHash
}

func NewSession(reg *SessionRegistry, key string) *Session {
	hash := redisw.NewRedisHash(reg.RedisWrapper, key, SESS_TIMEOUT)
	return &Session{reg: reg, RedisHash: hash}
}

func (sess *Session) GetKey() string {
	token, err := sess.GetString(SESS_TOKEN_KEY)
	if err == nil && token != "" {
		return sess.reg.GetKey(token)
	}
	return ""
}

func (sess *Session) AddFlash(messages ...string) (int, error) {
	key := fmt.Sprintf("flash:%s", sess.GetKey())
	args := append([]interface{}{key}, redisw.StrToList(messages)...)
	return redis.Int(sess.Exec("RPUSH", args...))
}

// 数量n为最大取出多少条消息，-1表示所有消息
func (sess *Session) GetFlashes(n int) ([]string, error) {
	key := fmt.Sprintf("flash:%s", sess.GetKey())
	return redis.Strings(sess.Exec("LRANGE", key, 0, n))
}

// 绑定用户角色，返回旧的sid
func (sess *Session) BindRoles(uid string, roles []string, kick bool) (string, error) {
	newSid := sess.GetKey()
	oldSid, _ := sess.reg.Onlines.GetString(uid) // 用于踢掉重复登录
	if oldSid == newSid {                        // 同一个token
		oldSid = ""
	}
	_, err := sess.reg.Onlines.SetVal(uid, newSid)
	_, err = sess.SetVal("uid", uid)
	_, err = sess.SetVal("roles", SessListJoin(roles))
	if kick && oldSid != "" { // 清空旧的session
		sess.reg.Delete(oldSid)
	}
	return oldSid, err
}
