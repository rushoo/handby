package cmd

import (
	"eagle/internal/model"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var Male = "Male"
var FeMale = "Female"

type UpdateOptions struct {
	Name string
	Age  int
	Sex  string
}

func Update() *cobra.Command {
	user := &model.User{}
	cmd := &cobra.Command{
		GroupID: "CRS",
		Use:     "update",
		Short:   "update resource",
		Long:    `update resource by name ....`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := user.UpdateUser(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		},
		SuggestionsMinimumDistance: 1,
	}

	cmd.Flags().StringVarP(&user.Name, "name", "n", "", "name")
	cmd.Flags().IntVarP(&user.Age, "age", "a", 0, "age")
	cmd.Flags().BoolVarP(&user.Sex, "sex", "s", false, "sex")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagsRequiredTogether("name", "age")
	return cmd
}
