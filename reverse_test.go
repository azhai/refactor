// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package refactor

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gitea.com/azhai/refactor/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"xorm.io/xorm"
)

var result = fmt.Sprintf(`package models

type A struct {
	Id int %sjson:"Id" xorm:"integer"%s
}

type B struct {
	Id int %sjson:"Id" xorm:"INTEGER"%s
}
`, "`", "`", "`", "`")

func reverse(rFile string) error {
	f, err := os.Open(rFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return reverseFromReader(f)
}

func reverseFromReader(rd io.Reader) error {
	var cfg config.ReverseConfig
	err := yaml.NewDecoder(rd).Decode(&cfg)
	if err != nil {
		return err
	}
	for _, target := range cfg.Targets {
		if err := RunReverse(&cfg.Source, &target); err != nil {
			return err
		}
	}

	return nil
}

func TestReverse(t *testing.T) {
	err := reverse("./testdata/goxorm.yml")
	assert.NoError(t, err)

	bs, err := ioutil.ReadFile("./models/models.go")
	assert.NoError(t, err)
	assert.EqualValues(t, result, string(bs))
}

func TestReverse2(t *testing.T) {
	type Outfw struct {
		Id       int    `xorm:"not null pk autoincr"`
		Sql      string `xorm:"default '' TEXT"`
		Template string `xorm:"default '' TEXT"`
		Filename string `xorm:"VARCHAR(50)"`
	}

	dir, err := ioutil.TempDir(os.TempDir(), "reverse")
	assert.NoError(t, err)

	e, err := xorm.NewEngine("sqlite3", filepath.Join(dir, "db.db"))
	assert.NoError(t, err)

	assert.NoError(t, e.Sync2(new(Outfw)))

	fp, err := os.Open("./testdata/goxorm.yml")
	err = reverseFromReader(fp)
	assert.NoError(t, err)
}
