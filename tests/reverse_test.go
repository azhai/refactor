package refactor_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	fileName := "./models/default/models.go"
	_, err := ioutil.ReadFile(fileName)
	assert.NoError(t, err)
	// assert.EqualValues(t, "", string(bs))
}
