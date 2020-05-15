package main

import(
	"strings"
	"os"
	"runtime"
	"github.com/pborman/getopt"
	"github.com/howeyc/gopass"
	"github.com/atotto/clipboard"
	"github.com/nitesh-/passpass/lib"
	"fmt"
)

func main() {

	optHelp := getopt.BoolLong("help", 0, "Help")
	optSet := getopt.StringLong("set", 's', "", "value must be {keyname}. It will prompt for the password.")
	optGet := getopt.StringLong("get", 'g', "", "value must be {keyname}. Get password for the corresponding key.")
	optDelete := getopt.StringLong("delete", 'd', "", "value must be {keyname}. Deletes the key.")
	optPasswordFile := getopt.StringLong("password-file", 'f', "", "Provide the path of password file")
	optGetKeys := getopt.BoolLong("get-keys", 'p', "Print all keys", "")
	optChangePassword := getopt.BoolLong("change-password", 'c', "Change password", "")
	
	getopt.Parse()

	if *optHelp || (*optSet == "" && *optGet == "" && *optDelete == "" && *optGetKeys == false && *optChangePassword == false) {
		getopt.Usage()
	} else {

		fmt.Printf("Please enter master password: ")

		masterPassword, _ := gopass.GetPasswd()

		bashPrinter := lib.BashPrinter()
		passwordError, passwordInvalid := lib.ValidatePassword(string(masterPassword))
		if !passwordInvalid {
			homeDir := userHomeDir()
			EncryptedFilePath := ""
			// Set the password db file path
			if *optPasswordFile != "" {
				EncryptedFilePath = *optPasswordFile
				_, err := os.Stat(EncryptedFilePath)
				if err != nil {
					// Trying to save empty json encrypted file
					openSsl := lib.OpenSsl()
					err := openSsl.EncryptFile(string(masterPassword), EncryptedFilePath, "{}")
					if err != nil {
						bashPrinter.PrintMessage(err.Error(), bashPrinter.ERROR)
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
				var key string = *optSet
				if key != "" {
					fmt.Printf("Please enter password for %s: ", key)
					password, err := gopass.GetPasswd()

					if err == nil {
						err := lib.SetPassword(EncryptedFilePath, string(masterPassword), key, string(password))
						if err == nil {
							bashPrinter.PrintMessage("Password set successfully in " + EncryptedFilePath, bashPrinter.SUCCESS)
						} else {
							bashPrinter.PrintMessage(err.Error(), bashPrinter.ERROR)
						}
					} else {
						bashPrinter.PrintMessage(err.Error(), bashPrinter.ERROR)
					}
				} else {
					bashPrinter.PrintMessage("Please provide key", bashPrinter.ERROR)
				}
			} else if *optGet != "" {
				password, err := lib.GetPassword(EncryptedFilePath, string(masterPassword), *optGet)
				if err == nil {
					// Try copying to clipboard
					err := clipboard.WriteAll(password)
					// If unsuccessful, then print password
					if err != nil {
						bashPrinter.PrintMessage("Password: " + password, bashPrinter.NORMAL)
					} else {
						bashPrinter.PrintMessage("Password copied to clipboard", bashPrinter.SUCCESS)
					}
				} else {
					bashPrinter.PrintMessage(err.Error(), bashPrinter.ERROR)
				}
			} else if *optDelete != "" {
				err := lib.DeletePassword(EncryptedFilePath, string(masterPassword), *optDelete)
				if err == nil {
					bashPrinter.PrintMessage("Password deleted successfully.", bashPrinter.SUCCESS)
				} else {
					bashPrinter.PrintMessage(err.Error(), bashPrinter.ERROR)
				}
			} else if *optGetKeys {
				keySlice, err := lib.GetAllKeys(EncryptedFilePath, string(masterPassword))

				if err == nil {
					for key := range keySlice {
						if keySlice[key] != "" {
							bashPrinter.PrintMessage(keySlice[key], bashPrinter.NORMAL)
						}
					}
				} else {
					bashPrinter.PrintMessage(err.Error(), bashPrinter.ERROR)
				}
			} else if *optChangePassword {
				fmt.Printf("Please enter New master password : ")
				newMasterPassword, err := gopass.GetPasswd()

				if err == nil {
					err := lib.ChangePassword(EncryptedFilePath, string(masterPassword), string(newMasterPassword))
					if err == nil {
						bashPrinter.PrintMessage("Password changed successfully.", bashPrinter.SUCCESS)
					} else {
						bashPrinter.PrintMessage(err.Error(), bashPrinter.ERROR)
					}
				} else {
					bashPrinter.PrintMessage(err.Error(), bashPrinter.ERROR)
				}
			}
		} else {
			bashPrinter.PrintMessage("Password must have atleast:", bashPrinter.NORMAL)
			bashPrinter.PrintMessage(strings.Join(passwordError[:], "\n"), bashPrinter.NORMAL)
		}
	}
}


// https://stackoverflow.com/questions/7922270/obtain-users-home-directory#answer-41786440
// Returns the Home directory for Linux, Mac and Windows 
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
