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

	"gitlab.com/aquachain/x/internal/debug"
	"gitlab.com/aquachain/x/utils"
	"gitlab.com/aquachain/aquachain/common/log"
	"gopkg.in/urfave/cli.v1"
)

var gitCommit = ""

var (
	app = utils.NewApp(gitCommit, "usage")
)

func init() {
	app.Name = "aqua-supplyproxy"
	app.Action = run
	app.Flags = append(debug.Flags, []cli.Flag{
		cli.StringFlag{
			Value: filepath.Join(utils.DataDirFlag.Value.String(), "aquachain.ipc"),
			Name:  "rpc",
			Usage: "path or url to rpc",
		},
		cli.StringFlag{
			Value: "localhost:8080",
			Name:  "addr",
			Usage: "Interface to serve HTTP. Use :8080 for all interfaces",
		},
	}...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Crit(err.Error())
	}
}
