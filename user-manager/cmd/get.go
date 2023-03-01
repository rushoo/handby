package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"

	"eagle/internal/model"
)

func Get() *cobra.Command {
	user := &model.User{}
	cmd := &cobra.Command{
		GroupID: "Q",
		Use:     "get",
		Short:   "get resource",
		Long:    `get resource from local storage`,
		Run: func(cmd *cobra.Command, args []string) {
			u, err := user.Get()
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("User: {name: %s age: %d sex: %v} \n", u.Name, u.Age, u.Sex)
		},
		SuggestionsMinimumDistance: 1,
		SuggestFor:                 []string{"find"},
	}

	cmd.Flags().StringVarP(&user.Name, "name", "n", "", "get resource by name")
	return cmd
}
