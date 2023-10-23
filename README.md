# Keeper AWS CLI Credential Process

By default, the AWS CLI uses credentials stored in plaintext in `~/.aws/credentials`. With this credential process, you can now use the Keeper Vault to store your AWS credentials, removing the need to have them on disk on the endpoint.

Instead, AWS will use this executable to fetch your AWS credential from your Vault securely using the Keeper Secrets Manager(KSM).

## Requirements

Usage requires:

- AWS credentials (see [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-services-iam-create-creds.html) on how to generate Access Keys)
- Keeper Secrets Manager (KSM) enabled
- AWS CLI v2

Development requires the above plus:

- Go > v1.21

## Setup

### Vault

The first step in the setup of the integration is to add you AWS `Access Key ID` and your `Secret Access Key` to a record in your Vault. There is no built in record type for this kind of secret; however, we can [create a custom record](https://docs.keeper.io/user-guides/record-types#custom-record-types) for this purpose alone. 

What you name this custom record type is up to you. However, this credential process looks for fields named `Access Key ID` and `Secret Access Key`, specifically. These fields must be present for successful authentication.

> Note: Field names are case sensitive. 

#TODO: Add image

Once you have created your custom field, you can now use it to create a record for your AWS Access Key. This record should be stored in a shared folder that your KSM application has permission to access.

Once safely stored, you are now able to delete the Access Key from your AWS credential file.

### KSM

The integration expects a KSM Application Configuration file at either `.config/keeper/aws-credential-process.json` or `aws-credential-process.json` relative to the user's home directory. It must have access to a Shared Folder that contains the AWS Access key required.

> For help setting up the KSM and obtaining a config file, head to the [official docs](https://docs.keeper.io/secrets-manager/secrets-manager/quick-start-guide)

### AWS Config

In your AWS config, which is usually located at `~/.aws/config`, add the following line to any profile you are using via the CLI. 

```ini
# Add the UID for your AWS Access Key
credential_process = /path/to/keeper-aws-credential-process --uid <Record UID>
```

## Usage 

Once configured as above, the AWS CLI will now automatically fetch your authentcation credential from the Keeper Vault. You can test that it works by using any CLI command in which you have an appropriate IAM role for; such as:

```shell
# List all s3 buckets
aws s3 ls
```

If the command completes without error, congratulations, you are now fully set up.

## Contributing

This module uses the built-in Golang tooling for building and testing:

```shell
# Run unit tests
go test ./...

# Build a local binary
go build -o keeper-aws-credential-process ./cmd/aws-credential-process/main.go
```

For bugs, changes, etc., please [submit an issue](https://github.com/Keeper-Security/aws-credential-process/issues).