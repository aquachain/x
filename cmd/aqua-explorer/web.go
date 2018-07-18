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
	"fmt"
	"html/template"
	"net/http"

	"github.com/aquachain/x/internal/bindat"
	"gitlab.com/aquachain/aquachain/common/log"

	"gopkg.in/urfave/cli.v1"
)

func staticHandle(s string) http.HandlerFunc {
	img, _ := bindat.Asset(s)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(img)
	}
}

func Serve(ctx *cli.Context) error {
	addr := "localhost:8080"
	if a := ctx.String("addr"); a != "" {
		addr = a
	}
	diffchart, timingchart, distribchart := GenerateCharts(ctx)
	if diffchart == nil {
		return fmt.Errorf("Error")
	}
	http.HandleFunc("/", http.HandlerFunc(indexhandler))
	http.HandleFunc("/diff.png", diffchart.HandleChart)
	http.HandleFunc("/aquachain.png", staticHandle("aquachain.png"))
	http.HandleFunc("/distribution.png", distribchart.HandleChart)
	http.HandleFunc("/timing.png", timingchart.HandleChart)
	log.Info("Starting HTTP server", "addr", addr)
	return http.ListenAndServe(addr, http.DefaultServeMux)

}

var coolstyle = `
img {
    max-width: 100%;
    height: auto;
}
.glogo {
  float: right;
  clear: none;
}
.tlogo {
  float: left;
  clear: none;
}
#id {
  width: 100px;
  height: 100px;
}
#wr {
  clear: both;
  width: 100%;
  height: 100%;
}
`

const logo = `
                              _           _
  __ _  __ _ _   _  __ _  ___| |__   __ _(_)_ __
 / _ '|/ _' | | | |/ _' |/ __| '_ \ / _' | | '_ \
| (_| | (_| | |_| | (_| | (__| | | | (_| | | | | |
 \__,_|\__, |\__,_|\__,_|\___|_| |_|\__,_|_|_| |_|
          |_|

`

var IndexHTML = `
<!DOCTYPE html>
<head>
<title>{{.Title}}</title>
<style> ` + coolstyle + ` </style>
</head>
<body>
<div>
<div class="tlogo"><pre>{{.Logo}}</pre></div>
<div class="glogo"><img style="max-width: 100px; height: auto;" src="aquachain.png"></div>
</div>
<br>
<div id="wr">
{{range .Images}}
<img src="{{.}}"><br>
{{end}}
</div>
</body>
`

func indexhandler(w http.ResponseWriter, r *http.Request) {
	title := `Aquachain Explorer`
	data := map[string]interface{}{
		"Title":  title,
		"Logo":   logo,
		"Images": []string{"diff.png", "timing.png", "distribution.png"},
	}
	tmpl, err := template.New("").Parse(IndexHTML)
	if err != nil {
		fmt.Println(err)
		return
	}
	tmpl.Execute(w, data)
}
