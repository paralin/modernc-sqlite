// Copyright 2021 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func createDB(b *testing.B, inMemory bool, driverName string) *sql.DB {
	var dsn string
	if inMemory {
		dsn = ":memory:"
	} else {
		dsn = path.Join(b.TempDir(), "test.db")
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		b.Fatal(err)
	}
	return db
}

func createTestTable(b *testing.B, db *sql.DB, nRows int) {
	if _, err := db.Exec("drop table if exists t"); err != nil {
		b.Fatal(err)
	}

	if _, err := db.Exec("create table t(i int)"); err != nil {
		b.Fatal(err)
	}

	if nRows > 0 {
		s, err := db.Prepare("insert into t values(?)")
		if err != nil {
			b.Fatal(err)
		}
		defer s.Close()

		if _, err := db.Exec("begin"); err != nil {
			b.Fatal(err)
		}

		for i := 0; i < nRows; i++ {
			if _, err := s.Exec(int64(i)); err != nil {
				b.Fatal(err)
			}
		}
		if _, err := db.Exec("commit"); err != nil {
			b.Fatal(err)
		}
	}
}

func getFuncName(i interface{}) string {
	// get function name as "package.function"
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	// return last component
	comps := strings.Split(fn, ".")
	return comps[len(comps)-1]
}
