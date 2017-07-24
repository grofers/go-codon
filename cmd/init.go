package cmd

import (
	"os"
	"log"
	"github.com/spf13/cobra"

	codon_bootstrap "github.com/grofers/go-codon/bootstrap"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstraps an empty project folder for use with codon.",
	Long: `Sets up your project folder with proper folder structure and
config/specification files to be used to generate service code. Run this
command to begin your project. Your working directory must be your project
directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if codon_bootstrap.Bootstrap("golang") {
			log.Println("Bootstrap successful")
		} else {
			log.Println("Bootstrap failed")
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
