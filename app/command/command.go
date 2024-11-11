package command

import (
	AppCommand "gorm-gen-skeleton/internal/command"
	"gorm-gen-skeleton/internal/variable"

	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Command struct {
	root *cobra.Command
}

var _ (AppCommand.CommandInterface) = (*Command)(nil)

func NewCommand(root *cobra.Command) *Command {
	return &Command{
		root: root,
	}
}

func (c *Command) GlobalFlags() {
	c.root.PersistentFlags().StringP("foo", "f", "", "foo flag.")
}

func (c *Command) RegisterCmds() []AppCommand.Interface {
	return []AppCommand.Interface{
		&FooCommand{},
		newGenCommand(),
	}
}

func newGenCommand() AppCommand.Interface {
	cfg := gen.Config{
		OutPath:           "./dao/query",
		OutFile:           "",
		ModelPkgPath:      "gox",
		Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
		FieldNullable:     false,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: true, //这样才能添加对应的联合索引
		FieldWithTypeTag:  true,
	}
	cfg.WithImportPkgPath("gorm.io/plugin/soft_delete")
	return AppCommand.NewGenCommand(
		AppCommand.WithConfig(cfg),
		AppCommand.WithDB(variable.DB),
		// AppCommand.WithTables([]string{"user"}),
		AppCommand.WithIgnoreFileds([]string{}),
		AppCommand.WithMethods(
			map[string][]any{},
		),
		AppCommand.WithDataMap(
			map[string]func(detailType gorm.ColumnType) (dataType string){
				"tinyint": func(detailType gorm.ColumnType) (dataType string) {
					if detailType.Name() == "is_del" {
						return "soft_delete.DeletedAt"
					}
					return "int8"
				},
				// "decimal": func(detailType gorm.ColumnType) (dataType string) { return "float32" },
			},
		),
	)
}
