package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Keeper-Security/aws-credential-process/internal/vault"
)

func main() {
	/*
		This is the entry point for the AWS CLI Credential Process. It is 
		invoked by the AWS CLI with the following config in ~/.aws/config:

			credential_process = path/to/aws-auth --uid <UID>

		The AWS CLI will call this executable with passed arguments. It 
		credential process expects the UID of the credential to be fetched from
		the Vault, which is then marshalled into JSON, and returned to the AWS 
		CLI via Stdout.
	*/

	var uid string

	flag.StringVar(&uid, "uid", "", "Credential UID")
	flag.Parse()

	if uid == "" {
		fmt.Println("UID is required")
		os.Exit(1)
	}

	// Returns non-flag arguments. This should be an empty slice.
	if len(flag.Args()) != 0 {
		fmt.Println("Unknown arguments provided")
		os.Exit(1)
	}

	// Fetch the credential from the Vault
	credential, err := vault.FetchCredential(uid)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Marshal the credential to JSON
	json, err := json.Marshal(credential)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Return the JSON to Stdout for the AWS CLI to consume
	fmt.Println(string(json))
	os.Exit(0)
}