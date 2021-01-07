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
// sqlite3* prepareReading1Memory(char * filename, int n);
// void reading1Memory(sqlite3 *DB, int n);
import "C"
import (
	"testing"
	"unsafe"
)

func reading1MemoryNative(b *testing.B, filename string, n int) {
	cs := C.CString(filename)
	db := C.prepareReading1Memory(cs, C.int(n))
	C.free(unsafe.Pointer(cs))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		C.reading1Memory(db, C.int(n))
	}
	C.sqlite3_close(db)
}
