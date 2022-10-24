package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 创建一个子命令
var versionCmd = &cobra.Command{
	Use: "version",
	//简短的描述
	Short: "version描述版本信息",
	Long:  "version详细的版本信息，详细的版本信息，详细的版本信息",
	// 实际执行的代码
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("这是一个版本信息")
	},
}

func init() {
	// 给根命令添加子命令
	rootCmd.AddCommand(versionCmd)
}
