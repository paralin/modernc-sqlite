// Copyright 2021 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"path"
	"testing"
)

// debug is only for development, to test plotting
var debug = false

func main() {
	for _, benchFunc := range allBenchmarksOfNRows {
		for _, isMemoryDB := range inMemory {
			// create graph
			graph := &GraphCompareOfNRows{
				title:      fmt.Sprintf("%s | In-Memory: %v", getFuncName(benchFunc), isMemoryDB),
				rowCountsE: rowCountsE,
				palette:    LightPalette,
			}

			// drivers
			for _, driver := range drivers {
				// this slice accumulates values as float64, for later plotting
				var (
					seriesValues []float64
					rowsPerSec   float64
				)

				// number of rows in table
				for _, e := range rowCountsE {
					if debug {
						// in debug mode we just generate random value to quickly see how information is plotted
						rowsPerSec = rand.Float64() * 200000
					} else {
						// run benchmark
						result := testing.Benchmark(func(b *testing.B) {
							// create DB
							db := createDB(b, isMemoryDB, driver)
							defer db.Close()

							// run bench
							benchFunc(b, db, int(math.Pow10(e)))
						})

						// calculate rows/sec
						rowsPerSec = math.Pow10(e) * float64(result.N) / result.T.Seconds()
					}

					// print result to console (FYI)
					benchName := fmt.Sprintf("%s_%s", getFuncName(benchFunc), makeName(isMemoryDB, driver, e))
					fmt.Println(benchName, "\t", fmt.Sprintf("%10.0f", rowsPerSec), "rows/sec")

					// add corresponding value to series
					seriesValues = append(seriesValues, rowsPerSec)
				}

				// add series to graph
				graph.AddSeries(driver, seriesValues)
			}

			// render graph into file
			outputFilename := path.Join("out", fmt.Sprintf("%s_memory:%v.png", getFuncName(benchFunc), isMemoryDB))
			if err := graph.Render(outputFilename); err != nil {
				log.Fatal(err)
			}
			log.Printf("plot written into %s\n", outputFilename)
		}
	}
}
