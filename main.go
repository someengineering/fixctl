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

func main() {
	apiEndpointPtr := flag.String("api-endpoint", "", "API endpoint URL")
	fixTokenPtr := flag.String("token", "", "Auth token")
	workspacePtr := flag.String("workspace", "", "Workspace ID")
	searchStrPtr := flag.String("search", "", "Search string")
	usernamePtr := flag.String("username", "", "Username")
	passwordPtr := flag.String("password", "", "Password")
	formatPtr := flag.String("format", "json", "Output format: json or yaml")
	withEdgesPtr := flag.Bool("with-edges", false, "Include edges in search results")
	flag.Parse()

	apiEndpoint := utils.GetEnvOrDefault("FIX_API_ENDPOINT", *apiEndpointPtr)
	username := utils.GetEnvOrDefault("FIX_USERNAME", *usernamePtr)
	password := utils.GetEnvOrDefault("FIX_PASSWORD", *passwordPtr)
	fixToken := utils.GetEnvOrDefault("FIX_TOKEN", *fixTokenPtr)
	workspaceID := utils.GetEnvOrDefault("FIX_WORKSPACE", *workspacePtr)
	searchStr := *searchStrPtr
	withEdges := *withEdgesPtr
	outputFormat := *formatPtr

	if apiEndpoint == "" {
		apiEndpoint = "https://app.fix.security"
	}
	if workspaceID == "" {
		fmt.Println("Workspace ID is required")
		os.Exit(1)
	}
	if outputFormat != "json" && outputFormat != "yaml" {
		fmt.Println("Invalid output format")
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
