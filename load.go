package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
)

type ScriptOutput struct {
	Result [][]string
}

func executeReadScript(script string, opts []string) (*ScriptOutput, error) {
	output, err := exec.Command(
		"python",
		append([]string{script}, opts...)...,
	).Output()

	if err != nil {
		return nil, err
	}

	var o ScriptOutput
	d := json.NewDecoder(bytes.NewReader(output))
	if err := d.Decode(&o); err != nil {
		return nil, err
	}

	return &o, nil
}

func importToDB(res *ScriptOutput, name string, dbm *gorp.DbMap) error {
	count, err := IsImportedEquipmentInfo(name, dbm)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("skipping import: %s", name)
		return nil
	}

	for _, row := range res.Result {
		if row[0] == "工事設計認証を受けた者の氏名又は名称" {
			continue
		}

		t, err := time.Parse("2006-01-02", row[6])
		if err != nil {
			continue
		}

		obj := EquipmentInfo{
			CertifiedName: row[0],
			EquipmentType: row[1],
			Model:         row[2],
			AuthNumber:    row[3],
			RadioType:     row[4],
			IsApplied1421: row[5],
			AuthDate:      t,
			Note:          row[7],
			File:          name,
		}

		log.Printf("insert: %v", obj)

		if err := dbm.Insert(&obj); err != nil {
			if !strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
				return err
			}
			log.Printf("skipping insert: %s", err.Error())
		}
	}
	return nil
}

func Load(src string, dst string) error {
	dbm, err := InitDb(dst)
	if err != nil {
		return err
	}
	defer dbm.Db.Close()

	script, err := GetReadScript()
	if err != nil {
		return err
	}

	list, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, fi := range list {
		if fi.IsDir() {
			continue
		}

		name := fi.Name()
		switch {
		case strings.HasSuffix(name, ".xls"):
			log.Printf("open %s", name)

			opts := []string{"--src", filepath.Join(src, name)}
			res, err := executeReadScript(script, opts)
			if err != nil {
				return err
			}
			if err := importToDB(res, name, dbm); err != nil {
				return err
			}
		}
	}

	return nil
}
