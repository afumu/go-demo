package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// 创建一个子命令
var versionCmd = &cobra.Command{
	Use: "version",
	//简短的描述
	Short: "version subcommand show git version info.",

	// 实际执行的代码
	Run: func(cmd *cobra.Command, args []string) {
		// 我这里实际调用了git的命令，获取返回结果
		output, err := ExecuteCommand("git", "version", args...)
		if err != nil {
			Error(cmd, args, err)
		}

		fmt.Fprint(os.Stdout, output)
	},
}

func init() {
	// 给根命令添加子命令
	rootCmd.AddCommand(versionCmd)
}
