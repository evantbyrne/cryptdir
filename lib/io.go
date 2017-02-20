package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Files map[string]string `json:"files"`
}

const dirMode = 0700
const fileMode = 0600

func ConfigRead(configDir string, password []byte) (config Config, nonce, key, salt []byte) {
	var (
		err       error
		cleartext []byte
	)

	cleartext, nonce, key, salt = DecryptFile(configDir+"/.cryptdir", password)

	if string(cleartext) != "" {

		if err = json.Unmarshal(cleartext, &config); err != nil {
			fmt.Println("Could not parse ctyptdir config file JSON.")
			os.Exit(1)
		}

		return config, nonce, key, salt
	}

	return Config{Files: make(map[string]string)}, nonce, key, salt
}

func ConfigWrite(configDir string, config Config, nonce, key, salt []byte) {
	var (
		configMarshalled []byte
		err              error
	)

	if configMarshalled, err = json.MarshalIndent(config, "", "\t"); err != nil {
		fmt.Printf("Could not convert config to JSON format.")
		os.Exit(1)
	}

	EncryptFile(configDir+"/.cryptdir", configMarshalled, nonce, key, salt)
}

func DecryptFile(path string, password []byte) (cleartext, nonce, key, salt []byte) {

	if _, err := os.Stat(path); err == nil {

		var (
			ciphertext []byte
			configData []byte
		)

		if configData, err = ioutil.ReadFile(path); err != nil {
			fmt.Printf("Could not read file at '%s'.\n", path)
			os.Exit(1)
		}

		if _, ciphertext, nonce, salt, err = VersionedSplit(configData); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		key = NewKey(salt, password)

		if cleartext, err = Decrypt(ciphertext, nonce, key); err != nil {
			fmt.Printf("Could not decrypt '%s' with given password.\n", path)
			os.Exit(1)
		}

	} else {
		nonce = NewNonce()
		salt = NewNonce()
		key = NewKey(salt, password)
	}

	return cleartext, nonce, key, salt
}

func EncryptFile(path string, cleartext, nonce, key, salt []byte) {
	var (
		ciphertext []byte
		err        error
	)

	if ciphertext, err = Encrypt(cleartext, nonce, key); err != nil {
		fmt.Println("Unexpected error: could not encrypt password vault.")
		os.Exit(1)
	}

	if err := ioutil.WriteFile(path, VersionedJoin(salt, nonce, ciphertext), fileMode); err != nil {
		fmt.Printf("Could not write file at '%s'.", path)
		os.Exit(1)
	}
}

func getRawFileName(config Config, fileName string) string {
	if _, ok := config.Files[fileName]; ok {
		return config.Files[fileName]
	}

	allRawNames := make([]string, 0, len(config.Files))

	for _, value := range config.Files {
		allRawNames = append(allRawNames, value)
	}

	return getRawFileNameUnique(allRawNames)
}

func getRawFileNameUnique(allRawNames []string) string {
	var (
		rawPath string
	)

	rawPath = RandomString(rawFileNameLength, false, true, true)

	for _, value := range allRawNames {

		if value == rawPath {
			return getRawFileNameUnique(allRawNames)
		}
	}

	return rawPath
}

func mustGetEnvPassword() (password string) {
	password = os.Getenv(envPassword)
	if password == "" {
		log.Fatal(messageLocked)
	}

	return password
}

func mustGetWorkingDir() (workingDir string) {
	var (
		err error
	)

	workingDir, err = os.Getwd()
	if err != nil {
		log.Fatal(messageWorkingDir)
	}

	return workingDir
}
