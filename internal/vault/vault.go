package vault

import (
	"fmt"
	"os"
	"path/filepath"

	ksm "github.com/keeper-security/secrets-manager-go/core"
)

type ConfigOptions struct {
	ConfigFile       string
	ConfigFileBackup string
}

type Credential struct {
	Version 		int
	AccessKeyId     string
	SecretAccessKey string
}

// Build the config options based on the given options.
func buildConfigOptions(h string) ConfigOptions {
	return ConfigOptions{
		ConfigFile:       filepath.Join(h, ".config", "keeper", "aws-credential-process.json"),
		ConfigFileBackup: filepath.Join(h, "aws-credential-process.json"),
	}
}

// Find the config.json file
func getConfig(options ConfigOptions) (string, error) {
	// If the ConfigFile exists, use it, else check ConfigFileBackup. If
	// neither exist, returns an error.
	if _, err := os.Stat(options.ConfigFile); err == nil {
		return options.ConfigFile, nil
	} else if _, err := os.Stat(options.ConfigFileBackup); err == nil {
		return options.ConfigFileBackup, nil
	} else {
		return "", fmt.Errorf("config file not found")
	}
}

// Fetch the AWS credential from the Vault via the Keeper Secrets Manager based
// on the UID provided. 
//
// It is expected that the record has two fields named: "Access Key ID" and 
// "Secret Access Key".
func FetchCredential(uid string) (*Credential, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Get config file path to be used for the KSM
	config, err := getConfig(buildConfigOptions(homeDir))
	if err != nil || config == "" {
		fmt.Println(err)
		os.Exit(1)
	}

	sm := ksm.NewSecretsManager(
		&ksm.ClientOptions{Config: ksm.NewFileKeyValueStorage(config)})

	records, err := sm.GetSecrets([]string{uid})
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("no records found for UID: %s", uid)
	}

	return &Credential{
		Version: 		 1,
		AccessKeyId:     records[0].GetFieldValueByLabel("Access Key ID"),
		SecretAccessKey: records[0].GetFieldValueByLabel("Secret Access Key"),
	}, nil
}