package command

import (
	"github.com/hinss/go-custom/framework/cobra"
)

// AddKernelCommands will add all command/* to root command
func AddKernelCommands(root *cobra.Command) {

	// app 启动命令
	root.AddCommand(initAppCommand())
	// build 启动命令
	root.AddCommand(initBuildCommand())
	// dev启动命令
	root.AddCommand(initDevCommand())
	// new 脚手架命令
	root.AddCommand(initNewCommand())
	// command 生成相关命令
	root.AddCommand(initCmdCommand())
	// provider 代码生成相关命令
	root.AddCommand(initProviderCommand())
	// middleware 迁移命令
	root.AddCommand(initMiddlewareCommand())
	// swagger 命令
	root.AddCommand(initSwaggerCommand())


}
