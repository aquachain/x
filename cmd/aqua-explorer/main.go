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
	"os"
	"path/filepath"

	"github.com/aquachain/x/internal/debug"
	"github.com/aquachain/x/utils"
	"github.com/aquanetwork/aquachain/common/log"
	"gopkg.in/urfave/cli.v1"
)

var gitCommit = ""

var (
	app = utils.NewApp(gitCommit, "usage")
)

func init() {
	app.Name = "aqua-explorer"
	app.Action = switcher
	app.Flags = append(debug.Flags, []cli.Flag{
		cli.StringFlag{
			Value: filepath.Join(utils.DataDirFlag.Value.String(), "aquachain.ipc"),
			Name:  "rpc",
			Usage: "path or url to rpc",
		},
		cli.StringFlag{
			Value: "latest.png",
			Name:  "o",
			Usage: "output",
		},
		cli.StringFlag{
			Value: "localhost:8080",
			Name:  "addr",
			Usage: "Interface to serve HTTP. Use :8080 for all interfaces",
		},
		cli.Uint64Flag{
			Value: 25000,
			Name:  "from",
			Usage: "begin at block",
		},
	}...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Crit(err.Error())
	}
}

func switcher(ctx *cli.Context) error {
	if err := debug.Setup(ctx); err != nil {
		return err
	}
	return Serve(ctx)
	// if argc := ctx.NArg(); argc == 0 {
	// 	return fmt.Errorf("command not found")
	// }
	// arg1 := ctx.Args()[0]
	// switch arg1 {
	// case "serve":
	// 	return Serve(ctx)
	// default:
	// 	return fmt.Errorf("command not found")
	// }

}
