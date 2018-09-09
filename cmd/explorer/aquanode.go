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
	"gitlab.com/aquachain/aquachain/aqua"
	"gitlab.com/aquachain/aquachain/node"
	"gitlab.com/aquachain/aquachain/params"
	"gitlab.com/aquachain/x/utils"
	"gopkg.in/urfave/cli.v1"
)

type aquaConfig struct {
	Aqua aqua.Config
	Node node.Config
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = "aquachain-x"
	cfg.Version = params.VersionWithCommit(gitCommit)
	cfg.HTTPModules = append(cfg.HTTPModules, "aqua")
	cfg.WSModules = append(cfg.WSModules, "aqua")
	cfg.IPCPath = "aquachain.ipc"
	return cfg
}

func makeConfigNode(ctx *cli.Context) (*node.Node, aquaConfig) {
	// Load defaults.
	cfg := aquaConfig{
		Aqua: aqua.DefaultConfig,
		Node: defaultNodeConfig(),
	}

	// Apply flags.
	utils.SetNodeConfig(ctx, &cfg.Node)
	stack, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}
	utils.SetAquaConfig(ctx, stack, &cfg.Aqua)
	return stack, cfg
}

func makeFullNode(ctx *cli.Context) *node.Node {
	stack, cfg := makeConfigNode(ctx)
	utils.RegisterAquaService(stack, &cfg.Aqua)
	return stack
}
