package cmd

import (
	"eagle/internal/model"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// 命令行提示中的example
var createExample = `
# Create a User

create -n Jack -a 10 -s 

`

func Create() *cobra.Command {
	user := &model.User{}

	cmd := &cobra.Command{
		Use:                   "create",
		Short:                 "create resource",
		Long:                  `crate resource to local storage`,
		Example:               createExample,
		DisableFlagsInUseLine: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			//检查用户名是否存在，存在就啥也不做返回
			if _, ok := user.CheckExist(); ok {
				fmt.Fprintln(os.Stderr, errors.New("duplicate name"))
				os.Exit(1)
			} else {
				log.Println("开始新建用户")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			//执行新建用户逻辑，出错就啥也不做返回
			if err := user.CreateUser(); err != nil {
				//fmt.Fprintln(os.Stderr, err.Error())
				//os.Exit(1)
				log.Println(err)
			}
		},
		GroupID:                    "CRS",                     //命令所属分组
		SuggestionsMinimumDistance: 1,                         // 开启建议
		SuggestFor:                 []string{"save", "store"}, // 开启建议命令

	}

	cmd.Flags().StringVarP(&user.Name, "name", "n", "", "name")
	cmd.Flags().IntVarP(&user.Age, "age", "a", 0, "age")
	cmd.Flags().BoolVarP(&user.Sex, "sex", "s", false, "sex")

	// 设置必选项和必须同时出现
	err := cmd.MarkFlagRequired("name")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	cmd.MarkFlagsRequiredTogether("name", "age")

	return cmd
}
