package refactor_test

import (
	"fmt"

	"gitea.com/azhai/refactor"
	"gitea.com/azhai/refactor/config"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	connKey     = "default"
	configFile  = "./settings.yml"
	testSqlFile = "./mysql_test.sql"
)

func init() {
	var err error
	if err = createTables(); err != nil {
		panic(err)
	}
	if err = generateModels(); err != nil {
		panic(err)
	}
}

func getDataSource(cfg config.IReverseSettings, name string) (*config.DataSource, error) {
	conns := cfg.GetConnections(name)
	if c, ok := conns[name]; ok {
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

func getConnection(key string) (*xorm.Engine, error) {
	cfg, err := config.ReadSettings(configFile)
	if err != nil {
		return nil, err
	}
	var d *config.DataSource
	if d, err = getDataSource(cfg, key); err != nil {
		return nil, err
	}
	return d.Connect(false)
}

func createTables() error {
	db, err := getConnection(connKey)
	if err != nil {
		return err
	}
	_, err = db.ImportFile(testSqlFile)
	return err
}

func generateModels(names ...string) error {
	cfg, err := config.ReadSettings(configFile)
	if err != nil {
		return err
	}
	var d *config.DataSource
	conns := cfg.GetConnections(names...)
	for key := range conns {
		if d, err = getDataSource(cfg, key); err != nil {
			return err
		}
		for _, target := range cfg.GetReverseTargets() {
			target = target.MergeOptions(d.ConnKey, d.PartConfig)
			if err := refactor.Reverse(d, &target); err != nil {
				return err
			}
		}
	}
	return nil
}
