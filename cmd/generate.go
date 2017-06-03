// Copyright © 2017 Grofers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"log"
	"github.com/spf13/cobra"

	codon_generator "github.com/grofers/go-codon/generator"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates service code alongside codon service specification",
	Long: `After you have run init on your project folder you should edit
the specification and configuration files according to your project. Run
this command to finally generate code according to the specification. Your
working directory must be your project directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if codon_generator.Generate("golang") {
			log.Println("Generate successful")
		} else {
			log.Println("Generate failed")
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
