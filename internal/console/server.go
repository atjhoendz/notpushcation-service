package console

import "github.com/spf13/cobra"

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {

}
