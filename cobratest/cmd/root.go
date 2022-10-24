package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 创建一个根命令
var rootCmd = &cobra.Command{
	// use 指定使用信息，命令是怎么呗调用的
	Use: "git",
	// 简短的描述
	Short: "git 短描述",
	// 长描述
	Long: `git 长描述，长描述，长描述`,
	// 实际执行操作的函数
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("发生错误")
	},
}

func Execute() {
	rootCmd.Execute()
}
