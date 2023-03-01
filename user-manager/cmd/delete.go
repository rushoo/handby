package cmd

import (
	"eagle/internal/model"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Delete() *cobra.Command {
	user := &model.User{}

	cmd := &cobra.Command{
		GroupID: "CRS",
		Use:     "delete",
		Short:   "delete resource",
		Long:    `delete resource by name`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := user.DeleteUser(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		},
		SuggestionsMinimumDistance: 1,
		SuggestFor:                 []string{"remove", "truncate"},
	}

	cmd.Flags().StringVarP(&user.Name, "name", "n", "", "delete by name")
	cmd.MarkFlagRequired("name")

	return cmd
}
