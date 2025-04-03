package main

import (
	"flag"
	CustomCommand "gorm-gen-skeleton/app/command"
	"gorm-gen-skeleton/internal/bootstrap"
	"gorm-gen-skeleton/internal/command"
	"gorm-gen-skeleton/internal/variable"
)

func main() {
	flag.Parse()
	variable.Init()
	bootstrap.Init()
	cmd := command.New()
	cmd.AddCommand(CustomCommand.NewCommand(cmd.Root())).Execute()
}
