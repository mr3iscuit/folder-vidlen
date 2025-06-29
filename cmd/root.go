package cmd

import (
	"fmt"
	"os"

	"github.com/mr3iscuit/folder-vidlen/course"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   getCommandName(),
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		pth := args[0]

		crs := course.NewCourseLength(pth)
		len, err := crs.GetCourseLength()
		if err != nil {
			return err
		}
		fmt.Printf("%s", len)

		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func getCommandName() string {
	return os.Args[0]
}
