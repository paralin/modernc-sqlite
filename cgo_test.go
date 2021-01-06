// Copyright 2017 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite // import "modernc.org/sqlite"

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const gcoDriver = "sqlite3"

//BenchmarkReading1MemoryNative-12    	 1000000	      7386 ns/op	      88 B/op	       3 allocs/op
//BenchmarkReading1MemoryCGO-12    	  896272	      1121 ns/op	     112 B/op	       6 allocs/op

func BenchmarkReading1MemoryNative(b *testing.B) {
	db, err := sql.Open(driverName, "file::memory:")
	if err != nil {
		b.Fatal(err)
	}

	defer func() {
		db.Close()
	}()

	if _, err := db.Exec(`
	create table t(i int);
	begin;
	`); err != nil {
		b.Fatal(err)
	}

	s, err := db.Prepare("insert into t values(?)")
	if err != nil {
		b.Fatal(err)
	}

	defer s.Close()

	for i := 0; i < b.N; i++ {
		if _, err := s.Exec(int64(i)); err != nil {
			b.Fatal(err)
		}
	}
	if _, err := db.Exec("commit"); err != nil {
		b.Fatal(err)
	}

	r, err := db.Query("select * from t")
	if err != nil {
		b.Fatal(err)
	}

	defer r.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !r.Next() {
			b.Fatal(err)
		}
		r.Scan()
	}
	b.StopTimer()
	if *oRecsPerSec {
		b.SetBytes(1e6)
	}
}

func BenchmarkReading1MemoryCGO(b *testing.B) {
	db, err := sql.Open(gcoDriver, "file::memory:")
	if err != nil {
		b.Fatal(err)
	}

	defer func() {
		db.Close()
	}()

	if _, err := db.Exec(`
	create table t(i int);
	begin;
	`); err != nil {
		b.Fatal(err)
	}

	s, err := db.Prepare("insert into t values(?)")
	if err != nil {
		b.Fatal(err)
	}

	defer s.Close()

	for i := 0; i < b.N; i++ {
		if _, err := s.Exec(int64(i)); err != nil {
			b.Fatal(err)
		}
	}
	if _, err := db.Exec("commit"); err != nil {
		b.Fatal(err)
	}

	r, err := db.Query("select * from t")
	if err != nil {
		b.Fatal(err)
	}

	defer r.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !r.Next() {
			b.Fatal(err)
		}
		r.Scan()
	}
	b.StopTimer()
	if *oRecsPerSec {
		b.SetBytes(1e6)
	}
}
