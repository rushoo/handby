package cmd

import (
	"eagle/internal/model"
	"fmt"

	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	user := &model.User{}
	cmd := &cobra.Command{
		GroupID: "Q",
		Use:     "list",
		Short:   "list resource",
		Long:    `list resource from local storage`,
		Run: func(cmd *cobra.Command, args []string) {
			users, err := user.List()
			if err != nil {
				fmt.Println(err)
			}

			for _, v := range users {
				fmt.Printf("User: {name: %s age: %d sex: %v} \n", v.Name, v.Age, v.Sex)
			}
		},
	}
	return cmd
}
