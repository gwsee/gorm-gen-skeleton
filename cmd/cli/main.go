package main

import (
	CustomCommand "gorm-gen-skeleton/app/command"
	_ "gorm-gen-skeleton/internal/bootstrap"
	"gorm-gen-skeleton/internal/command"
)

func main() {
	cmd := command.New()
	cmd.AddCommand(CustomCommand.NewCommand(cmd.Root())).Execute()
}
