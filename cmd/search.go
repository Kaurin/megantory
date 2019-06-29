package cmd

import (
	"strings"

	"github.com/Kaurin/megantory/lib/search"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches all the things",
	Long: `
* Groks your AWS credentials file for all the profiles defined therein
* Performs a free-text search against all profiles / regions / supported services
* Does all the searches mentioned above and processing with a high degree of concurrency
* Returns a list of found resources with breadcrumbs where to find them

Example: megantory search "foobar"

`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Starting the Search CMD")
		search.Search(strings.Join(args, " "))
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
