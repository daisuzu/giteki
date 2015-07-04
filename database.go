package main

import (
	"database/sql"
	"strings"

	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
)

func InitDb(dst string) (*gorp.DbMap, error) {
	db, err := sql.Open("sqlite3", dst)
	if err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	var t *gorp.TableMap
	t = dbmap.AddTable(EquipmentInfo{}).SetKeys(true, "ID")
	t.SetUniqueTogether(
		"CertifiedName",
		"EquipmentType",
		"Model",
		"AuthNumber",
		"RadioType",
		"AuthDate",
	)
	t = dbmap.AddTable(RadioAccessTechnology{}).SetKeys(true, "ID")
	t.ColMap("EquipmentType").SetUnique(true)

	if err := dbmap.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	for k, v := range ratDef {
		obj := RadioAccessTechnology{EquipmentType: k, Description: v}
		if err := dbmap.Insert(&obj); err != nil {
			if !strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
				return nil, err
			}
		}
	}

	return dbmap, nil
}
