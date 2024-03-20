package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/someengineering/fixctl/auth"
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
	flag.Parse()

	apiEndpoint := utils.GetEnvOrDefault("FIX_API_ENDPOINT", *apiEndpointPtr)
	username := utils.GetEnvOrDefault("FIX_USERNAME", *usernamePtr)
	password := utils.GetEnvOrDefault("FIX_PASSWORD", *passwordPtr)
	fixToken := utils.GetEnvOrDefault("FIX_TOKEN", *fixTokenPtr)
	workspaceID := utils.GetEnvOrDefault("FIX_WORKSPACE", *workspacePtr)
	searchStr := *searchStrPtr

	if apiEndpoint == "" {
		apiEndpoint = "https://app.fix.security"
	}
	if workspaceID == "" {
		fmt.Println("Workspace ID is required")
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

	response, err := search.SearchTable(apiEndpoint, fixToken, workspaceID, searchStr)
	if err != nil {
		fmt.Println("Search error:", err)
		return
	}

	fmt.Println(response)
}
