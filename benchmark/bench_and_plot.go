package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"testing"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

var debug = false

func main() {
	for _, benchFunc := range allBenchmarksOfNRows {
		for _, isMemoryDB := range inMemory {

			// create graph handle
			graph := newGraph(fmt.Sprintf("%s | In-Memory: %v", getFuncName(benchFunc), isMemoryDB))

			var driverMetrics [][]float64

			// drivers
			for _, driver := range drivers {
				// this slices accumulate values as float64, for later plotting
				var (
					xValues    []float64
					yValues    []float64
					rowsPerSec float64
				)

				// number of rows in table
				for _, e := range rowCountsE {
					// generate benchmark name
					var benchName = fmt.Sprintf("%s_%s", getFuncName(benchFunc), makeName(isMemoryDB, driver, e))

					// if debug {
					// 	rowsPerSec = rand.Float64() * 200000
					// } else {

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

					// }

					// print result to console
					fmt.Println(benchName, "\t", fmt.Sprintf("%10.0f", rowsPerSec), "rows/sec")

					// adjust axes values
					xValues = append(xValues, float64(e))
					yValues = append(yValues, rowsPerSec)

					// adjust max for Y axis
					yMax := (int(rowsPerSec/150000) + 1) * 150000 // a special case of round()
					graph.YAxis.Range = &chart.ContinuousRange{
						Min: 0,
						Max: math.Max(float64(yMax), graph.YAxis.Range.GetMax()),
					}
				}

				// add series to graph
				graph.Series = append(graph.Series, newSeries(xValues, yValues, driver))

				// save driver metric for latter annotation
				driverMetrics = append(driverMetrics, yValues)

			}

			// add annotations: the ratio of second driver over first
			annotations := &chart.AnnotationSeries{}
			for i := range driverMetrics[0] {
				annotations.Annotations = append(annotations.Annotations, newRatioAnnotation(
					float64(rowCountsE[i]),
					driverMetrics[1][i],
					driverMetrics[1][i]/driverMetrics[0][i],
				))
			}
			graph.Series = append(graph.Series, annotations)

			// add legend
			//note we have to do this as a separate step because we need a reference to graph
			graph.Elements = []chart.Renderable{
				chart.Legend(graph, chart.Style{
					FontSize:    12,
					StrokeColor: drawing.ColorFromHex("1e1e1e1").WithAlpha(0),
					FillColor:   drawing.ColorFromHex("1e1e1e1"),
					FontColor:   drawing.ColorFromHex("d4d4d4").WithAlpha(192),
				}),
			}

			// render graph into file
			if err := os.MkdirAll("out", 0775); err != nil {
				log.Fatal(err)
			}
			f, err := os.Create(path.Join("out", fmt.Sprintf("%s_memory:%v.png", getFuncName(benchFunc), isMemoryDB)))
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			if err := graph.Render(chart.PNG, f); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func newGraph(title string) *chart.Chart {
	return &chart.Chart{
		Title: title,
		TitleStyle: chart.Style{
			Show:      true,
			FontColor: drawing.ColorFromHex("d4d4d4").WithAlpha(128),
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
			FillColor: drawing.ColorFromHex("252526"),
		},
		Canvas: chart.Style{
			FillColor: drawing.ColorFromHex("1e1e1e1"),
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show:        true,
				FontColor:   drawing.ColorFromHex("d4d4d4").WithAlpha(128),
				StrokeColor: drawing.ColorFromHex("d4d4d4").WithAlpha(128),
			},
			NameStyle: chart.Style{
				Show:      true,
				FontColor: drawing.ColorFromHex("d4d4d4").WithAlpha(128),
			},
			Name:  "rows",
			Ticks: genXticks(),
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{},
			Style: chart.Style{
				Show:        true,
				FontColor:   drawing.ColorFromHex("d4d4d4").WithAlpha(128),
				StrokeColor: drawing.ColorFromHex("d4d4d4").WithAlpha(128),
			},
			Name: "rows/sec",
			NameStyle: chart.Style{
				Show:      true,
				FontColor: drawing.ColorFromHex("d4d4d4").WithAlpha(128),
			},
		},
	}
}

func newSeries(xValues, yValues []float64, driverName string) chart.Series {
	// create series
	var (
		seriesName  string
		seriesColor drawing.Color
	)

	if driverName == "sqlite3" {
		seriesName = "CGo"
		seriesColor = drawing.ColorFromHex("d5d5a5")
	} else {
		seriesName = "Go"
		seriesColor = drawing.ColorFromHex("569cd5")
	}

	series := &chart.ContinuousSeries{
		Name: seriesName,
		Style: chart.Style{
			DotColor:    seriesColor,
			DotWidth:    2.5,
			Show:        true,
			StrokeColor: seriesColor,
			StrokeWidth: 1.5,
		},
		XValues: xValues,
		YValues: yValues,
	}

	return series
}

func newRatioAnnotation(x, y, ratio float64) chart.Value2 {
	return chart.Value2{
		XValue: x,
		YValue: y,
		Label:  fmt.Sprintf("%.2fx", ratio),
		Style: chart.Style{
			FontSize:            8,
			TextHorizontalAlign: chart.TextHorizontalAlignLeft,
			FillColor:           drawing.ColorFromHex("574255").WithAlpha(0),
			FontColor:           drawing.ColorFromHex("d4d4d4"),
			StrokeColor:         drawing.ColorFromHex("d4d4d4").WithAlpha(0),
			TextRotationDegrees: 45,
		},
	}
}

func genXticks() []chart.Tick {
	var ticks []chart.Tick
	for i, e := range rowCountsE {
		ticks = append(ticks, chart.Tick{
			Value: float64(e),
			Label: fmt.Sprintf("1e%d", i+1),
		})
	}
	return ticks
}
