package main

import (
	"fmt"
	"html/template"
	"net/http"

	"gopkg.in/urfave/cli.v1"
)

func Serve(ctx *cli.Context) error {
	addr := "localhost:8080"
	if a := ctx.String("addr"); a != "" {
		addr = a
	}
	diffchart, timingchart, distribchart := GenerateCharts(ctx)

	http.HandleFunc("/", http.HandlerFunc(indexhandler))
	http.HandleFunc("/diff.png", diffchart.HandleChart)
	// http.HandleFunc("/txblock.png", mustgetchart_txblock(ctx).HandleChart)
	http.HandleFunc("/distribution.png", distribchart.HandleChart)
	http.HandleFunc("/timing.png", timingchart.HandleChart)
	return http.ListenAndServe(addr, http.DefaultServeMux)

}

var coolstyle = `
img {
    max-width: 100%;
    height: auto;
}
#wr {
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
<pre>{{.Logo}}</pre>
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
