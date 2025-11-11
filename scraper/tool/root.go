package tool

import "github.com/spf13/cobra"

var Root = cobra.Command{
	Use: "",
}

func Run() {
	Root.Execute()
}

func init() {
	Root.AddCommand(&ScrapeLibrary)
}
