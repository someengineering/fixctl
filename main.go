package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/someengineering/fixctl/auth"
	"github.com/someengineering/fixctl/format"
	"github.com/someengineering/fixctl/search"
	"github.com/someengineering/fixctl/utils"
)

func customUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(flag.CommandLine.Output(), "  --%s: %s (default %q)\n", f.Name, f.Usage, f.DefValue)
	})
}

func main() {
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)
	flag.Usage = customUsage
	apiEndpointPtr := flag.String("endpoint", "https://app.fix.security", "API endpoint URL (env FIX_ENDPOINT)")
	fixTokenPtr := flag.String("token", "", "Auth token (env FIX_TOKEN)")
	workspacePtr := flag.String("workspace", "", "Workspace ID (env FIX_WORKSPACE)")
	searchStrPtr := flag.String("search", "", "Search string")
	usernamePtr := flag.String("username", "", "Username (env FIX_USERNAME)")
	passwordPtr := flag.String("password", "", "Password (env FIX_PASSWORD)")
	formatPtr := flag.String("format", "json", "Output format: json, yaml or csv")
	csvHeadersPtr := flag.String("csv-headers", "id,name,kind,/ancestors.cloud.reported.id,/ancestors.account.reported.id,/ancestors.region.reported.id", "CSV headers (comma-separated, relative to /reported by default)")
	withEdgesPtr := flag.Bool("with-edges", false, "Include edges in search results")
	help := flag.Bool("help", false, "Display help information")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	withEdges := *withEdgesPtr

	invalidArgs := false
	username, password, err := utils.SanitizeCredentials(utils.GetEnvOrDefault("FIX_USERNAME", *usernamePtr), utils.GetEnvOrDefault("FIX_PASSWORD", *passwordPtr))
	if err != nil {
		fmt.Println("Invalid username or password:", err)
		invalidArgs = true
	}
	searchStr, err := utils.SanitizeSearchString(*searchStrPtr)
	if err != nil {
		fmt.Println("Invalid search string:", err)
		invalidArgs = true
	}
	apiEndpoint, err := utils.SanitizeAPIEndpoint(utils.GetEnvOrDefault("FIX_ENDPOINT", *apiEndpointPtr))
	if err != nil {
		fmt.Println("Invalid API endpoint:", err)
		invalidArgs = true
	}
	fixToken, err := utils.SanitizeToken(utils.GetEnvOrDefault("FIX_TOKEN", *fixTokenPtr))
	if err != nil {
		fmt.Println("Invalid token:", err)
		invalidArgs = true
	}
	workspaceID, err := utils.SanitizeWorkspaceId(utils.GetEnvOrDefault("FIX_WORKSPACE", *workspacePtr))
	if err != nil {
		fmt.Println("Invalid workspace ID:", err)
		invalidArgs = true
	}
	csvHeaders, err := utils.SanitizeCSVHeaders(*csvHeadersPtr)
	if err != nil {
		fmt.Println("Invalid CSV headers:", err)
		invalidArgs = true
	}
	outputFormat, err := utils.SanitizeOutputFormat(*formatPtr)
	if err != nil {
		fmt.Println("Invalid output format:", err)
		invalidArgs = true
	}
	if invalidArgs {
		os.Exit(1)
	}

	if fixToken == "" && username != "" && password != "" {
		var err error
		fixToken, err = auth.LoginAndGetJWT(apiEndpoint, username, password)
		if err != nil {
			fmt.Println("Login error:", err)
			return
		}
	}
	if fixToken == "" {
		fmt.Println("Either token or username and password are required")
		os.Exit(1)
	}

	results, errs := search.SearchGraph(apiEndpoint, fixToken, workspaceID, searchStr, withEdges)
	firstResult := true
	for result := range results {
		var output string
		var err error
		switch outputFormat {
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
		fmt.Println("Search error:", err)
	}
}
