package main

import(
	"fmt"
	"os"
	"github.com/nitesh-/passpass/lib"
)

const EncryptedFilePath = "./testcase.psdb"

func main() {
	
	masterPassword := "Hellow1%"
	newMasterPassword := "Hellow1%%"

	keyPasswordMap := make(map[string]string)

	keyPasswordMap["gmail"] = "Hello1234"
	keyPasswordMap["yahoo"] = "Hello12345"
	keyPasswordMap["outlook"] = "Hello12346"
	keyPasswordMap["yandex"] = "Hello12347"

	testCaseOutputMap := make(map[string]string)

	// Validate Password test case
	_ , passwordInvalid := lib.ValidatePassword(masterPassword)
	if !passwordInvalid {
		testCaseOutputMap["ValidatePassword"] = "passed"
	} else {
		testCaseOutputMap["ValidatePassword"] = "failed"
	}

	// Setting Password
	testCaseOutputMap["SetPassword"] = "passed"
	testCaseOutputMap["GetPassword"] = "passed"
	for key := range keyPasswordMap {
		fmt.Printf("Setting Password for key %s\n", key)
		err := lib.SetPassword(EncryptedFilePath, masterPassword, key, keyPasswordMap[key])
		if err == nil {
			fmt.Printf("Verifying Password for key %s\n", key)
			password, err := lib.GetPassword(EncryptedFilePath, masterPassword, key)
			if err == nil && password == keyPasswordMap[key] {
				fmt.Printf("Password for key %s verified\n", key)
			} else {
				testCaseOutputMap["GetPassword"] = "failed"
			}
		} else {
			testCaseOutputMap["SetPassword"] = "failed"
			fmt.Println(err)
		}
	}

	// Get all keys
	testCaseOutputMap["GetAllKeys"] = "passed"
	keySlice, err := lib.GetAllKeys(EncryptedFilePath, masterPassword)
	if err == nil {
		for key := range keyPasswordMap {
			isKeyFound := false
			for key1 := range keySlice {
				if keySlice[key1] == key {
					isKeyFound = true
					break
				}
			}
			if !isKeyFound {
				fmt.Printf("Key %s is not found\n", key)
				testCaseOutputMap["GetAllKeys"] = "failed"
			} else {
				fmt.Printf("Key %s found\n", key)
			}
		}
	} else {
		testCaseOutputMap["GetAllKeys"] = "failed"
		fmt.Println(err)
	}
	
	// Change the password
	testCaseOutputMap["ChangePassword"] = "passed"
	if lib.ChangePassword(EncryptedFilePath, masterPassword, newMasterPassword) != nil {
		testCaseOutputMap["ChangePassword"] = "failed"
		fmt.Println(err)
	}

	// Delete keys
	testCaseOutputMap["DeletePassword"] = "passed"
	for key := range keyPasswordMap {
		fmt.Printf("Deleting Password for key %s\n", key)
		// Delete using new master password
		err := lib.DeletePassword(EncryptedFilePath, newMasterPassword, key)
		if(err == nil) {
			fmt.Printf("key %s deleted.\n", key)
		} else {
			testCaseOutputMap["DeletePassword"] = "failed"
		}
	}

	fmt.Println("================")
	fmt.Println("Test Case Output")
	fmt.Println("================")

	bashPrinter := lib.BashPrinter()
	
	for key, _ := range testCaseOutputMap {
		printColor := bashPrinter.SUCCESS
		if testCaseOutputMap[key] == "failed" {
			printColor = bashPrinter.ERROR
		}

		fmt.Println(key + " -> " + bashPrinter.GetMessage(testCaseOutputMap[key], printColor))
	}

	os.Remove(EncryptedFilePath)
}
