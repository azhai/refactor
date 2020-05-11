package refactor_test

import (
	"io/ioutil"
	"testing"

	"gitea.com/azhai/refactor/config"
	"github.com/stretchr/testify/assert"
)

var (
	configFile  = "./settings.yml"
	testSqlFile = "./mysql_test.sql"
)

func init() {
	cfg, err := config.ReadSettings(configFile)
	if err != nil {
		panic(err)
	}
	if err = createTables(cfg); err != nil {
		panic(err)
	}
	if err = generateModels(cfg); err != nil {
		panic(err)
	}
}

func TestReverse(t *testing.T) {
	fileName := "./models/default/models.go"
	_, err := ioutil.ReadFile(fileName)
	assert.NoError(t, err)
	// assert.EqualValues(t, "", string(bs))
}
