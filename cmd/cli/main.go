package cmd

import (
	"fmt"
	"os"

	sourcecollector "github.com/hitesh22rana/sourcecollector/pkg"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sourcecollector",
	Short: "A simple tool to consolidate multiple files into a single .txt file",
	Long: `A simple tool to consolidate multiple files into a single .txt file.
Perfect for feeding your files to AI tools without any fuss.`,
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")
		fast, _ := cmd.Flags().GetBool("fast")

		sc, err := sourcecollector.NewSourceCollector(input, output, fast)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		sourceTree, err := sc.GenerateSourceTree()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		sourcetreeStructure, err := sc.GenerateSourceTreeStructure(sourceTree)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := sc.SaveSourceCode(sourceTree, sourcetreeStructure); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringP("input", "i", "", "Input directory path")
	rootCmd.Flags().StringP("output", "o", "output.txt", "Output file path")
	rootCmd.Flags().Bool("fast", false, "Faster result but may result in unordered data, default(false)")
	rootCmd.MarkFlagRequired("input")
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
