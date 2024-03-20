# fixctl
Fix CLI tool

## Usage
```
Usage of fixctl:
  --endpoint: API endpoint URL (default "https://app.fix.security")
  --format: Output format: json or yaml (default "json")
  --help: Display help information (default "false")
  --password: Password (default "")
  --search: Search string (default "")
  --token: Auth token (default "")
  --username: Username (default "")
  --with-edges: Include edges in search results (default "false")
  --workspace: Workspace ID (default "")
```

Alternatively the following environment variables can be used:
```
FIX_ENDPOINT
FIX_USERNAME
FIX_PASSWORD
FIX_TOKEN
FIX_WORKSPACE
```

If no token is provided, the username and password will be used to authenticate and obtain a token. Does not support MFA.
