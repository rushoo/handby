package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "eagle",
	Short:   "eagle is an example",                            // 简短介绍
	Long:    "eagle is an example to show how to use cobra  ", // 完整介绍
	Version: "0.0.1",                                          // 设置版本号，如果添加了可以
	Run:     runHelp,
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func init() {
	//命令分组
	basicCommandQ := cobra.Group{
		Title: "Basic Command(Q)",
		ID:    "Q",
	}
	rootCmd.AddGroup(&basicCommandQ)
	rootCmd.AddCommand(Get(), List())

	basicCommandCRS := cobra.Group{
		Title: "Basic Command(CRS)",
		ID:    "CRS",
	}

	rootCmd.AddGroup(&basicCommandCRS)
	rootCmd.AddCommand(Create(), Delete(), Update())

	// 设置使用介绍模版
	// rootCmd.SetUsageTemplate(rootUsageTemplate)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
