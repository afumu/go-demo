package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// 创建一个根命令
var rootCmd = &cobra.Command{
	// use 指定使用信息，命令是怎么呗调用的
	Use: "git",
	// 简短的描述
	Short: "Git is a distributed version control system.",
	// 长描述
	Long: `Git is a free and open source distributed version control system
designed to handle everything from small to very large projects 
with speed and efficiency.`,
	// 实际执行操作的函数
	Run: func(cmd *cobra.Command, args []string) {
		Error(cmd, args, errors.New("unrecognized command"))
	},
}

func Execute() {
	rootCmd.Execute()
}
