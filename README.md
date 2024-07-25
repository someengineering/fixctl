# fixctl
Fix CLI tool

## Installation

### Binaries
Download the latest release from the [releases page](https://github.com/someengineering/fixctl/releases).
On macOS and Linux make sure to make the binary executable.

Linux Example:
```bash
curl -Lo fixctl https://github.com/someengineering/fixctl/releases/download/0.0.7/fixctl-linux-amd64-0.0.7
chmod +x fixctl
```

### Homebrew
```bash
brew install someengineering/tap/fixctl
```

## Usage
```
Usage of fixctl:
  --csv-headers: CSV headers (comma-separated, relative to /reported by default) (default "id,name,kind,/ancestors.cloud.reported.id,/ancestors.account.reported.id,/ancestors.region.reported.id")
  --endpoint: API endpoint URL (env FIX_ENDPOINT) (default "https://app.fix.security")
  --format: Output format: json, yaml or csv (default "json")
  --help: Display help information (default "false")
  --search: Search string (default "")
  --token: Auth token (env FIX_TOKEN) (default "")
  --with-edges: Include edges in search results (default "false")
  --workspace: Workspace ID (env FIX_WORKSPACE) (default "")
```

If an environment variable is set, it will be used and the command line flag ignored.

Go to your [user settings](https://app.fix.security/user-settings) and create an API token. Set the `FIX_TOKEN` environment variable to the token value.
Then go to your [workspace settings](https://app.fix.security/workspace-settings) and export `FIX_WORKSPACE` to the workspace ID you want to query.

### Example
Search for available AWS EBS volumes that have not been accessed in the last 7 days and output in CSV format.
```bash
$ fixctl --format csv --search "is(aws_ec2_volume) and volume_status = available and last_access > 7d"
vol-0adeedfc71dcbe9d5,ResotoEKS-dynamic-pvc-e575191f-d4f3-4253-96e4-399ded05bf14,aws_ec2_volume,aws,752466027617,eu-central-1
vol-0ae5f3fad85b7b3c6,vol-0ae5f3fad85b7b3c6,aws_ec2_volume,aws,625596817853,eu-central-1
vol-0fe068d91a8aaaced,ResotoEKS-dynamic-pvc-08ded29a-70c9-4d36-9d28-727140850d96,aws_ec2_volume,aws,752466027617,eu-central-1
```

The default output format for `fixctl` is JSON. Here we search for the same orphaned volumes and use `jq` to format the output as `aws ec2 delete-volume` commands.
```bash
$ fixctl --search "is(aws_ec2_volume) and volume_status = available and last_access > 30d" | jq -r '. | "aws ec2 delete-volume --volume-id \(.reported.id) --region \(.ancestors.region.reported.id) --profile \(.ancestors.account.reported.id)"'
aws ec2 delete-volume --volume-id vol-0adeedfc71dcbe9d5 --region eu-central-1 --profile 752466027617
aws ec2 delete-volume --volume-id vol-0ae5f3fad85b7b3c6 --region eu-central-1 --profile 625596817853
aws ec2 delete-volume --volume-id vol-0fe068d91a8aaaced --region eu-central-1 --profile 752466027617
```
