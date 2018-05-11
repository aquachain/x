package main

import (
	"fmt"

	chart "github.com/wcharczuk/go-chart"
)

//
// func drawchart_diff(ctx *cli.Context) (err error) {
// 	start := time.Now()
// 	w, err := os.Create(ctx.String("o"))
// 	if err != nil {
// 		return err
// 	}
// 	graph, err := getchart_diff(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	err = graph.chart.Render(chart.PNG, w)
// 	if err != nil {
// 		utils.Fatalf("Render error: %v\n", err)
// 	}
// 	fmt.Printf("Image render done in %v\n", time.Since(start))
// 	return err
// }

func float2Tick(f []float64) (ticks []chart.Tick) {
	for i, v := range f {
		if i%4 == 0 {
			ticks = append(ticks, chart.Tick{Label: fmt.Sprintf("%v", int(v)), Value: v})
		}
	}
	return ticks
}
