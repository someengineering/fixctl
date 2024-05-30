package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/someengineering/fixctl/auth"
	"github.com/someengineering/fixctl/format"
	"github.com/someengineering/fixctl/search"
	"github.com/someengineering/fixctl/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "fixctl",
		Short: "fixctl is the Fix Security CLI tool",
		Long:  `fixctl allows you to search the Fix Security Graph and export cloud inventory data for further processing.`,
		Run:   executeSearch,
	}

	apiEndpoint string
	fixToken    string
	fixJWT      string
	workspace   string
	username    string
	password    string
	formatType  string
	searchStr   string
	csvHeaders  string
	withEdges   bool
	verbose     bool
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&apiEndpoint, "endpoint", "https://app.fix.security", "API endpoint URL (env FIX_ENDPOINT)")
	rootCmd.PersistentFlags().StringVar(&fixToken, "token", "", "Auth token (env FIX_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&workspace, "workspace", "", "Workspace ID (env FIX_WORKSPACE)")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Username (env FIX_USERNAME)")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password (env FIX_PASSWORD)")
	rootCmd.PersistentFlags().StringVar(&formatType, "format", "json", "Output format: json, yaml or csv")
	rootCmd.PersistentFlags().StringVar(&searchStr, "search", "", "Search string")
	rootCmd.PersistentFlags().StringVar(&csvHeaders, "csv-headers", "id,name,kind,/ancestors.cloud.reported.id,/ancestors.account.reported.id,/ancestors.region.reported.id", "CSV headers")
	rootCmd.PersistentFlags().BoolVar(&withEdges, "with-edges", false, "Include edges in search results")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose output")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.SetEnvPrefix("FIX")
	viper.AutomaticEnv()
}

func initConfig() {
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func executeSearch(cmd *cobra.Command, args []string) {
	invalidArgs := false
	username, password, err := utils.SanitizeCredentials(viper.GetString("username"), viper.GetString("password"))
	if err != nil {
		logrus.Errorln("Invalid username or password:", err)
		invalidArgs = true
	}
	searchStr, err := utils.SanitizeSearchString(viper.GetString("search"))
	if err != nil {
		logrus.Errorln("Invalid search string:", err)
		invalidArgs = true
	}
	apiEndpoint, err := utils.SanitizeAPIEndpoint(viper.GetString("endpoint"))
	if err != nil {
		logrus.Errorln("Invalid API endpoint:", err)
		invalidArgs = true
	}
	fixToken, err := utils.SanitizeToken(viper.GetString("token"))
	if err != nil {
		logrus.Errorln("Invalid token:", err)
		invalidArgs = true
	}
	workspace, err := utils.SanitizeWorkspaceId(viper.GetString("workspace"))
	if err != nil {
		logrus.Errorln("Invalid workspace ID:", err)
		invalidArgs = true
	}
	csvHeaders, err := utils.SanitizeCSVHeaders(viper.GetString("csv-headers"))
	if err != nil {
		logrus.Errorln("Invalid CSV headers:", err)
		invalidArgs = true
	}
	formatType, err := utils.SanitizeOutputFormat(viper.GetString("format"))
	if err != nil {
		logrus.Errorln("Invalid output format:", err)
		invalidArgs = true
	}
	if invalidArgs {
		os.Exit(1)
	}

	if fixToken == "" && username != "" && password != "" {
		var err error
		fixJWT, err = auth.LoginAndGetJWT(apiEndpoint, username, password)
		if err != nil {
			logrus.Errorln("Login error:", err)
			return
		}
	}
	if fixToken != "" {
		var err error
		fixJWT, err = auth.GetJWTFromToken(apiEndpoint, fixToken)
		if err != nil {
			logrus.Errorln("Login error:", err)
			return
		}
	}
	if fixJWT == "" {
		logrus.Errorln("Either token or username and password are required")
		os.Exit(1)
	}

	results, errs := search.SearchGraph(apiEndpoint, fixJWT, workspace, searchStr, withEdges)
	firstResult := true
	for result := range results {
		var output string
		var err error
		switch formatType {
		case "yaml":
			if !firstResult {
				fmt.Println("---")
			} else {
				firstResult = false
			}
			output, err = format.ToYAML(result)
		case "csv":
			output, err = format.ToCSV(result, csvHeaders)
		default:
			output, err = format.ToJSON(result)
		}
		if err != nil {
			fmt.Printf("Error formatting output: %v\n", err)
			os.Exit(2)
		}
		fmt.Print(output)
	}

	if err, ok := <-errs; ok {
		logrus.Errorln("Search error:", err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func SetVersion(version string) {
	rootCmd.Version = version
}
