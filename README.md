# fixctl
Fix CLI tool

## Installation

### Binaries
Download the latest release from the [releases page](https://github.com/someengineering/fixctl/releases).
On macOS and Linux make sure to make the binary executable.

Example:
```bash
curl -Lo fixctl https://github.com/someengineering/fixctl/releases/download/0.0.3/fixctl-linux-amd64-0.0.3
chmod +x fixctl
```

### Homebrew
```bash
brew install someengineering/tap/fixctl
```

## Usage
```
Usage of fixctl:
  --endpoint: API endpoint URL (env FIX_ENDPOINT) (default "https://app.fix.security")
  --format: Output format: json or yaml (default "json")
  --help: Display help information (default "false")
  --password: Password (env FIX_PASSWORD) (default "")
  --search: Search string (default "")
  --token: Auth token (env FIX_TOKEN) (default "")
  --username: Username (env FIX_USERNAME) (default "")
  --with-edges: Include edges in search results (default "false")
  --workspace: Workspace ID (env FIX_WORKSPACE) (default "")
```

If no token is provided, the username and password will be used to authenticate and obtain a token. Does not support MFA.
If an environment variable is set, it will be used and the command line flag ignored.
