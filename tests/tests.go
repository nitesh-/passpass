package main

import(
	"fmt"
	"os"
	"../lib"
)

const EncryptedFilePath = "./testcase.psdb"

func main() {
	
	masterPassword := "HelloWorld1234"

	keyPasswordMap := make(map[string]string)

	bashTextColorMap := make(map[string]string)
	bashTextColorMap["Green"] = "\033[0;32m"
	bashTextColorMap["Red"] = "\033[0;31m"
	bashTextColorMap["Normal"] = "\033[0;0m"

	keyPasswordMap["gmail"] = "Hello1234"
	keyPasswordMap["yahoo"] = "Hello12345"
	keyPasswordMap["outlook"] = "Hello12346"
	keyPasswordMap["yandex"] = "Hello12347"

	testCaseOutputMap := make(map[string]string)

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

	// Delete keys
	testCaseOutputMap["DeletePassword"] = "passed"
	for key := range keyPasswordMap {
		fmt.Printf("Deleting Password for key %s\n", key)
		err := lib.DeletePassword(EncryptedFilePath, masterPassword, key)
		if(err == nil) {
			fmt.Printf("key %s deleted.\n", key)
		} else {
			testCaseOutputMap["DeletePassword"] = "failed"
		}
	}

	fmt.Println("================")
	fmt.Println("Test Case Output")
	fmt.Println("================")
	for key, _ := range testCaseOutputMap {
		printColor := "Green"
		if testCaseOutputMap[key] == "failed" {
			printColor = "Red"
		}
		fmt.Println(key + " -> " + bashTextColorMap[printColor] + testCaseOutputMap[key] + bashTextColorMap["Normal"])
	}

	os.Remove(EncryptedFilePath)
}
