package refactor_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/k0kubun/pp"

	"github.com/azhai/refactor"
	"github.com/azhai/refactor/config"
	"github.com/stretchr/testify/assert"
	"xorm.io/xorm"
)

var (
	verbose     = false
	configFile  = "./settings.yml"
	testSqlFile = "./mysql_test.sql"
)

func createTables(cfg config.IReverseSettings) (err error) {
	c, ok := cfg.GetConnConfig("default")
	if !ok {
		err = fmt.Errorf("the connection is not found")
		return
	}
	r, _ := config.NewReverseSource(c)
	var db *xorm.Engine
	if db, err = r.Connect(false); err != nil {
		return
	}
	var content []byte
	if content, err = ioutil.ReadFile(testSqlFile); err != nil {
		return
	}
	repl := strings.NewReplacer(
		"{{CURR_MONTH}}", time.Now().Format("200601"),
		"{{PREV_MONTH}}", time.Now().AddDate(0, -1, 0).Format("200601"),
		"{{EARLY_MONTH}}", time.Now().AddDate(0, -2, 0).Format("200601"),
	)
	sql := repl.Replace(string(content))
	_, err = db.Import(strings.NewReader(sql))
	return
}

func Test01CreateTables(t *testing.T) {
	cfg, err := config.ReadSettings(configFile)
	if verbose {
		pp.Println(cfg)
	}
	assert.NoError(t, err)
	err = createTables(cfg)
	assert.NoError(t, err)
	err = refactor.ExecReverseSettings(cfg, verbose)
	assert.NoError(t, err)
}

func Test02Reverse(t *testing.T) {
	fileName := "./models/default/models.go"
	_, err := ioutil.ReadFile(fileName)
	assert.NoError(t, err)
	// assert.EqualValues(t, "", string(bs))
}
