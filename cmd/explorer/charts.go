// Copyright 2018 The aquachain Authors
// This file is part of the aquachain/x project.
//
// aquachain is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// aquachain is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with aquachain. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"io"
	"strconv"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func drawChart(c *chart.Chart, w io.Writer) error {
	err := c.Render(chart.PNG, w)
	if err != nil {
		return err
	}
	return nil
}

// see switch and comments below for variadic 'labels' usage
func newChart(xseries, yseries []float64, labels ...string) *chart.Chart {

	var (
		title = ""
		xname = ""
		yname = ""
	)

	switch len(labels) {
	case 1: // title
		title = labels[0]
	case 2: // title, y (assumes x is block number)
		title = labels[0]
		yname = labels[1]
	case 3: //  title, x and y
		title = labels[0]
		xname = labels[1]
		yname = labels[2]
	default:
		panic("too many labels") // programming error
	}

	return &chart.Chart{
		Title: title,
		TitleStyle: chart.Style{
			Show:      true,
			FontColor: chart.GetDefaultColor(1),
			DotColor:  chart.GetDefaultColor(1),
		},
		Canvas: chart.Style{
			FillColor: drawing.ColorFromHex("000000"),
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:    50,
				Left:   25,
				Right:  25,
				Bottom: 10,
			},
			FillColor:   drawing.ColorFromHex("000000"),
			StrokeColor: drawing.ColorGreen,
			FontColor:   chart.GetDefaultColor(3),
			DotColor:    chart.GetDefaultColor(3),
		},
		Width:  1000,
		Height: 500,
		XAxis: chart.XAxis{
			Name: xname,
			NameStyle: chart.Style{
				Show:      true,
				FontColor: chart.GetDefaultColor(3),
				DotColor:  chart.GetDefaultColor(3),
			},
			Style: chart.Style{
				Show:      true,
				FontColor: chart.GetDefaultColor(3),
				DotColor:  chart.GetDefaultColor(3),
			},
			TickPosition: chart.TickPositionUnderTick,
			ValueFormatter: func(v interface{}) string {
				f := v.(float64)
				return strconv.Itoa(int(f))
			},
		},
		YAxis: chart.YAxis{
			Name: yname,
			NameStyle: chart.Style{
				Show:      true,
				FontColor: chart.GetDefaultColor(3),
				DotColor:  chart.GetDefaultColor(3),
			},
			Style: chart.Style{
				Show:      true,
				FontColor: chart.GetDefaultColor(3),
				DotColor:  chart.GetDefaultColor(3),
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:                true,
					StrokeColor:         chart.GetDefaultColor(3),
					FillColor:           chart.GetDefaultColor(3),
					TextRotationDegrees: 90,
				},
				XValues: xseries,
				YValues: yseries,
			},
		},
	}
}
