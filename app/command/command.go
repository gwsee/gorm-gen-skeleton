package command

import (
	"flag"
	AppCommand "gorm-gen-skeleton/internal/command"
	"gorm-gen-skeleton/internal/variable"

	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var DaoPath = flag.String("d", "./dao/query", "dao文件生成的地址")

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
	c.root.PersistentFlags().StringP("cfg", "c", "", "动态配置文件 没必要每个项目都生成对应的文件")
	c.root.PersistentFlags().StringP("dao", "d", "", "动态dao/query生成的目录,dao model就在相同的目录下;此命令需要注册 不然就不生效")
}

func (c *Command) RegisterCmds() []AppCommand.Interface {
	return []AppCommand.Interface{
		&FooCommand{},
		newGenCommand(),
	}
}

func newGenCommand() AppCommand.Interface {
	cfg := gen.Config{
		OutPath:           *DaoPath,
		OutFile:           "",
		ModelPkgPath:      "",
		Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	}

	cfg.WithImportPkgPath("gorm.io/plugin/soft_delete") //目前不是所有的表都加了这个 就暂时不自动加了
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
				"int": func(detailType gorm.ColumnType) (dataType string) {
					if detailType.Name() == "deleted_at" {
						return "soft_delete.DeletedAt"
					}
					return "int"
				},
				"bigint": func(detailType gorm.ColumnType) (dataType string) {
					if detailType.Name() == "deleted_at" {
						return "soft_delete.DeletedAt"
					}
					return "int"
				},
				// "decimal": func(detailType gorm.ColumnType) (dataType string) { return "float32" },
			},
		),
	)
}
