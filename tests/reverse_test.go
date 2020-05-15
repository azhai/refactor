package refactor_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"gitea.com/azhai/refactor"
	"gitea.com/azhai/refactor/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
	"xorm.io/xorm"
)

var (
	configFile  = "./settings.yml"
	testSqlFile = "./mysql_test.sql"
)

func getDataSource(cfg config.IReverseSettings, name string) (*config.DataSource, error) {
	if c, ok := cfg.GetConnConfig(name); ok {
		d := config.NewDataSource(c, name)
		if d.ReverseSource == nil {
			err := fmt.Errorf("the driver %s is not exists", c.DriverName)
			return d, err
		}
		return d, nil
	}
	err := fmt.Errorf("the connection named %s is not found", name)
	return nil, err
}

func createTables(cfg config.IReverseSettings) error {
	d, err := getDataSource(cfg, "default")
	if err != nil {
		return err
	}
	var db *xorm.Engine
	if db, err = d.Connect(false); err == nil {
		_, err = db.ImportFile(testSqlFile)
	}
	return err
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
