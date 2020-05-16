package refactor_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"xorm.io/xorm"

	"gitea.com/azhai/refactor"
	"gitea.com/azhai/refactor/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

var (
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
	if db, err = r.Connect(false); err == nil {
		_, err = db.ImportFile(testSqlFile)
	}
	return
}

func Test01CreateTables(t *testing.T) {
	cfg, err := config.ReadSettings(configFile)
	pp.Println(cfg)
	assert.NoError(t, err)
	err = createTables(cfg)
	assert.NoError(t, err)
	err = refactor.ExecReverseSettings(cfg)
	assert.NoError(t, err)
}

func Test02Reverse(t *testing.T) {
	fileName := "./models/default/models.go"
	_, err := ioutil.ReadFile(fileName)
	assert.NoError(t, err)
	// assert.EqualValues(t, "", string(bs))
}
