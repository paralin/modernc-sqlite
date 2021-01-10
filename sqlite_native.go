// Copyright 2017 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build cgo,cgobench

package sqlite

// #cgo LDFLAGS: -lsqlite3
// #include <stdlib.h>
// #include <sqlite3.h>
//
// sqlite3* prepareReading(char * filename, int n);
// void reading(sqlite3 *DB, int n);
// sqlite3* prepareInsertComparative(char * filename, int n);
// void insertComparative(sqlite3 *DB, int n);
import "C"
import (
	"testing"
	"unsafe"

	"modernc.org/libc"
	sqlite3 "modernc.org/sqlite/lib"
)

func benchmarkReadNativeC(b *testing.B, filename string, n int) {
	cs := C.CString(filename)
	db := C.prepareReading(cs, C.int(n))
	C.free(unsafe.Pointer(cs))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		C.reading(db, C.int(n))
	}
	C.sqlite3_close(db)
}

func benchmarkReadNativeGO(b *testing.B, filename string, n int) {
	tls := libc.NewTLS()
	cs := C.CString(filename)
	db := prepareReading(tls, uintptr(unsafe.Pointer(cs)), int32(n))
	C.free(unsafe.Pointer(cs))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reading(tls, db, int32(n))
	}
	sqlite3.Xsqlite3_close(tls, db)
}

func benchmarkInsertComparativeNativeC(b *testing.B, filename string, n int) {
	cs := C.CString(filename)
	db := C.prepareInsertComparative(cs, C.int(n))
	C.free(unsafe.Pointer(cs))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		C.insertComparative(db, C.int(n))
	}
	C.sqlite3_close(db)
}

func benchmarkInsertComparativeNativeGO(b *testing.B, filename string, n int) {
	tls := libc.NewTLS()
	cs := C.CString(filename)
	db := prepareInsertComparative(tls, uintptr(unsafe.Pointer(cs)), int32(n))
	C.free(unsafe.Pointer(cs))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insertComparative(tls, db, int32(n))
	}
	sqlite3.Xsqlite3_close(tls, db)
}
