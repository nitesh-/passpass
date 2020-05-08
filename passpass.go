package main

import(
	"fmt"
	"strings"
	"os"
	"runtime"
	"github.com/pborman/getopt"
	"github.com/howeyc/gopass"
	"github.com/atotto/clipboard"
	"./lib"
)

func main() {

	optHelp := getopt.BoolLong("help", 0, "Help")
	optSet := getopt.StringLong("set", 's', "", "value must be {keyname:password}. Sets key & password.")
	optGet := getopt.StringLong("get", 'g', "", "value must be {keyname}. Get password for the corresponding key.")
	optDelete := getopt.StringLong("delete", 'd', "", "value must be {keyname}. Deletes the key.")
	optPasswordFile := getopt.StringLong("password-file", 'f', "", "Provide the path of password file")
	optGetKeys := getopt.BoolLong("get-keys", 'p', "Print all keys", "")
	
	getopt.Parse()

	if *optHelp || (*optSet == "" && *optGet == "" && *optDelete == "" && *optGetKeys == false) {
		getopt.Usage()
	} else {

		fmt.Printf("Please enter master password: ")

		masterPassword, _ := gopass.GetPasswd()

		if len(masterPassword) > 0 {
			homeDir := userHomeDir()
			EncryptedFilePath := ""
			// Set the password db file path
			if *optPasswordFile != "" {
				EncryptedFilePath = *optPasswordFile
				_, err := os.Stat(EncryptedFilePath)
				if err != nil {
					// Trying to save empty json encrypted file
					err := lib.EncryptFile(string(masterPassword), EncryptedFilePath, "{}")
					if err != nil {
						fmt.Println(err)
						os.Exit(0)
					}
				}
			} else {
				EncryptedFilePath = homeDir + "/.passpass/p.psdb"
				_, err := os.Stat(EncryptedFilePath)
				if err != nil {
					os.MkdirAll(strings.Replace(EncryptedFilePath, "/p.psdb", "", 1), 0700)
				}
			}

			if *optSet != "" {
				strSplit := strings.Split(*optSet, ":")
				err := lib.SetPassword(EncryptedFilePath, string(masterPassword), strSplit[0], strSplit[1])
				if err == nil {
					fmt.Println("Password set successfully in " + EncryptedFilePath)
				} else {
					fmt.Println(err)
				}
			} else if *optGet != "" {
				password, err := lib.GetPassword(EncryptedFilePath, string(masterPassword), *optGet)
				if err == nil {
					// Try copying to clipboard
					err := clipboard.WriteAll(password)
					// If unsuccessful, then print password
					if err != nil {
						fmt.Printf("Password: %s\n", password)
					} else {
						fmt.Println("Password copied to clipboard")
					}
				} else {
					fmt.Println(err)
				}
			} else if *optDelete != "" {
				err := lib.DeletePassword(EncryptedFilePath, string(masterPassword), *optDelete)
				if err == nil {
					fmt.Println("Password deleted successfully.")
				} else {
					fmt.Println(err)
				}
			} else if *optGetKeys {
				keySlice, err := lib.GetAllKeys(EncryptedFilePath, string(masterPassword))

				if err == nil {
					for key := range keySlice {
						if keySlice[key] != "" {
							fmt.Println(keySlice[key])
						}
					}
				} else {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println("Please provide master password")
		}
	}
}

/**
 * https://stackoverflow.com/questions/7922270/obtain-users-home-directory#answer-41786440
 *
 *
 */
func userHomeDir() string {
    if runtime.GOOS == "windows" {
        home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
        if home == "" {
            home = os.Getenv("USERPROFILE")
        }
        return home
    } else if runtime.GOOS == "linux" {
        home := os.Getenv("XDG_CONFIG_HOME")
        if home != "" {
            return home
        }
    }
    return os.Getenv("HOME")
}

/** build command
GOOS=windows GOARCH=amd64 go build -o bin/windows-amd/passpass.exe passpass.go
GOOS=windows GOARCH=386 go build -o bin/windows-386/passpass.exe passpass.go
GOOS=darwin GOARCH=amd64 go build -o bin/osx-amd/passpass passpass.go
GOOS=darwin GOARCH=386 go build -o bin/osx-386/passpass passpass.go
go build -o bin/linux/passpass passpass.go
*/
