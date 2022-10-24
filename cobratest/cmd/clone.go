package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 创建一个子命令
var cloneCmd = &cobra.Command{
	Use: "clone",
	//简短的描述
	Short: "clone克隆项目到本地",
	Long:  "clone克隆项目到本地,clone克隆项目到本地,clone克隆项目到本地,",
	// 实际执行的代码
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clone 信息")
	},
}

func init() {
	// 给根命令添加子命令
	rootCmd.AddCommand(cloneCmd)
}
