module modernc.org/sqlite

go 1.20

require (
	github.com/google/pprof v0.0.0-20221118152302-e6195bd50e26
	github.com/klauspost/cpuid/v2 v2.2.3
	github.com/mattn/go-sqlite3 v1.14.22
	golang.org/x/sys v0.16.0
	modernc.org/fileutil v1.3.0
	modernc.org/gc/v3 v3.0.0-20240107210532-573471604cb6
	modernc.org/libc v1.41.0
	modernc.org/mathutil v1.6.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/tools v0.17.0 // indirect
	modernc.org/memory v1.7.2 // indirect
	modernc.org/strutil v1.2.0 // indirect
	modernc.org/token v1.1.0 // indirect
)

retract [v1.16.0, v1.17.2] // https://gitlab.com/cznic/sqlite/-/issues/100

retract v1.19.0 // module source tree too large (max size is 524288000 bytes)

retract v1.20.1 // https://gitlab.com/cznic/sqlite/-/issues/123

retract v1.29.4 // tagged accidentally w/o builders checking the commit
