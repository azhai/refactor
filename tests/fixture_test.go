package refactor_test

import (
	"fmt"

	"gitea.com/azhai/refactor"
	"gitea.com/azhai/refactor/config"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

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

func generateModels(cfg config.IReverseSettings, names ...string) error {
	var d *config.DataSource
	conns := cfg.GetConnections(names...)
	for key := range conns {
		var err error
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
