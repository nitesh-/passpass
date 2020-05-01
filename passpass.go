package main

import(
	"fmt"
	openssl "gopkg.in/Luzifer/go-openssl.v3"
	"io/ioutil"
	"strings"
	"errors"
	"encoding/json"
	"github.com/howeyc/gopass"
	"github.com/pborman/getopt"
	"github.com/atotto/clipboard"
)

const EncryptedFilePath = "./.p"

func main() {

	optHelp := getopt.BoolLong("help", 0, "Help")
	optSet := getopt.StringLong("set", 's', "", "value must be {keyname:password}. Sets key & password.")
	optGet := getopt.StringLong("get", 'g', "", "value must be {keyname}. Get password for the corresponding key")
	optGetKeys := getopt.BoolLong("get-keys", 'p', "Print all keys", "")
	
	getopt.Parse()
	if *optHelp || (*optSet == "" && *optGet == "" && *optGetKeys == false) {
		getopt.Usage()
	} else {

		fmt.Printf("Please enter master password: ")

		masterPassword, _ := gopass.GetPasswd()

		if len(masterPassword) > 0 {
			if *optSet != "" {
				strSplit := strings.Split(*optSet, ":")
				SetPassword(string(masterPassword), strSplit[0], strSplit[1])
			} else if *optGet != "" {
				GetPassword(string(masterPassword), *optGet)
			} else if *optGetKeys {
				GetAllKeys(string(masterPassword))
			}
		} else {
			fmt.Println("Please provide master password")
		}
	}
}

func SetPassword(masterPassword string, key string, password string) {
	// Initiate the map which will store the json data
	encryptedDataMap := make(map[string]string)

	// Read the password file
	content, err := ioutil.ReadFile(EncryptedFilePath)
	if err == nil {
		json.Unmarshal(content, &encryptedDataMap)
	}
	
	// Initiate the openssl to encrypt the password
	o := openssl.New()
	enc, err := o.EncryptBytes(masterPassword + "" + key, []byte(password), openssl.DigestSHA256Sum)
	if err != nil {
		fmt.Printf("An error occurred while encrypting the string: %s\n", err)
	}

	// Make the encrypted string to the map
	encryptedDataMap[key] = string(enc)

	// Encode the data to JSON
	jsonBytes, err := json.Marshal(encryptedDataMap);

	if err == nil {
		// Write the JSON to file
		EncryptFile(masterPassword, EncryptedFilePath, string(jsonBytes))
	} else {
		fmt.Printf("An error occurred while encoding the string to JSON: %s\n", err)
	}
}

func GetPassword(masterPassword string, key string) {
	// Initiate the map which will store the json data
	encryptedDataMap := make(map[string]string)

	// Read the password file
	content, err := DecryptFile(masterPassword, EncryptedFilePath)
	if err == nil {
		// Decode the json string to data map
		json.Unmarshal([]byte(content), &encryptedDataMap)

		encPassword, encryptedKeyExists := encryptedDataMap[key]

		if encryptedKeyExists {
			// Initiate the openssl to decrypt the password
			o := openssl.New()
			// Decrypt the password
			dec, err := o.DecryptBytes(masterPassword + "" + key, []byte(encPassword), openssl.DigestSHA256Sum)

			if err == nil {
				// Try copying to clipboard
				err := clipboard.WriteAll(string(dec))
				// If unsuccessful, then print password
				if err != nil {
					fmt.Println(err)
					fmt.Printf("Password: %s\n", string(dec))
				} else {
					fmt.Println("Password copied to clipboard")
				}
			} else {
				if err.Error() == "invalid padding" {
					fmt.Println("The master password you provide may be incorrect.");
				} else {
					fmt.Printf("An error occurred while decrypting the string: %s\n", err)
				}
			}
		} else {
			fmt.Printf("The provided key %s does not exists.", key)
		}
	} else {
		fmt.Println(err)
	}
}

func GetAllKeys(masterPassword string) {
	content, err := DecryptFile(masterPassword, EncryptedFilePath)
	encryptedDataMap := make(map[string]string)
	if err == nil {
		json.Unmarshal([]byte(content), &encryptedDataMap)
		for key := range encryptedDataMap {
			fmt.Println(key)
		}
	} else {
		fmt.Println(err)
	}
}

func DecryptFile(masterPassword string, filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err == nil {
		o := openssl.New()
		dec, err := o.DecryptBytes(masterPassword, []byte(content), openssl.DigestSHA256Sum)
		if(err == nil) {
			return string(dec), nil
		}
		return "", errors.New("The master password may be incorrect")
	}

	return "", errors.New("Invalid Encrypted file path")
}

func EncryptFile(masterPassword string, filePath string, fileData string) bool {
	o := openssl.New()
	encFileData, err := o.EncryptBytes(masterPassword, []byte(fileData), openssl.DigestSHA256Sum)
	if(err == nil) {
		ioutil.WriteFile(EncryptedFilePath, encFileData, 0644)
		return true
	}

	return false
}
