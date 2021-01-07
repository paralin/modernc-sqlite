// Copyright 2017 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build cgo,cgobench

package sqlite

// #cgo LDFLAGS: -lsqlite3
// #include <stdlib.h>
// #include <stdio.h>
// #include <sqlite3.h>
// #include <string.h>
//
// sqlite3* prepareReading1(char * filename, int n);
// void reading1native(sqlite3 *DB, int n);
import "C"
import (
	"testing"
	"unsafe"

	"modernc.org/libc"
	sqlite3 "modernc.org/sqlite/lib"
)

func reading1NativeC(b *testing.B, filename string, n int) {
	cs := C.CString(filename)
	db := C.prepareReading1(cs, C.int(n))
	C.free(unsafe.Pointer(cs))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		C.reading1native(db, C.int(n))
	}
	C.sqlite3_close(db)
}

func reading1NativeGO(b *testing.B, filename string, n int) {
	tls := libc.NewTLS()
	cs := C.CString(filename)
	db := prepareReading1(tls, uintptr(unsafe.Pointer(cs)), int32(n))
	C.free(unsafe.Pointer(cs))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reading1native(tls, db, int32(n))
	}
	sqlite3.Xsqlite3_close(tls, db)
}
