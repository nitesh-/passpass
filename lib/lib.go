package lib 

import(
	"errors"
	"encoding/json"
	"strings"
	"unicode"
)

// Sets the password for a particular key.
func SetPassword(encryptedFilePath string, masterPassword string, key string, password string) error {
	// Initiate the map which will store the json data
	encryptedDataMap := make(map[string]string)

	// Read the password file
	openSsl := OpenSsl()
	content, err := openSsl.DecryptFile(masterPassword, encryptedFilePath)
	if err != nil && !strings.Contains(err.Error(), "Invalid Encrypted file path") {
		return err
	} else {
		json.Unmarshal([]byte(content), &encryptedDataMap)
		
		// Initiate the openssl to encrypt the password
		enc, err := openSsl.EncryptString(masterPassword + "" + key, password)
		if err != nil {
			return err
		}

		// Make the encrypted string to the map
		encryptedDataMap[key] = enc

		// Encode the data to JSON
		jsonBytes, err := json.Marshal(encryptedDataMap)

		if err == nil {
			// Write the JSON to file
			err := openSsl.EncryptFile(masterPassword, encryptedFilePath, string(jsonBytes))
			if(err != nil) {
				return errors.New("An error occurred while encrypting the file: " + err.Error())
			}
		} else {
			return errors.New("An error occurred while encoding the string to JSON: " + err.Error())
		}

		return nil
	}
}

// Returns the password for a particular key
func GetPassword(encryptedFilePath string, masterPassword string, key string) (string, error) {
	// Initiate the map which will store the json data
	encryptedDataMap := make(map[string]string)

	// Read the password file
	openSsl := OpenSsl()
	content, err := openSsl.DecryptFile(masterPassword, encryptedFilePath)
	if err == nil {
		// Decode the json string to data map
		json.Unmarshal([]byte(content), &encryptedDataMap)

		encPassword, encryptedKeyExists := encryptedDataMap[key]

		if encryptedKeyExists {
			// Decrypt the password
			dec, err := openSsl.DecryptString(masterPassword + "" + key, encPassword)

			if err == nil {
				return dec, nil
			} else {
				return "", err
			}
		} else {
			return "", errors.New("The provided key " + key + " does not exists.")
		}
	} else {
		return "", err
	}
}

// Deletes the key and password
func DeletePassword(encryptedFilePath string, masterPassword string, key string) error {
	// Initiate the map which will store the json data
	encryptedDataMap := make(map[string]string)

	// Read the password file
	openSsl := OpenSsl()
	content, err := openSsl.DecryptFile(masterPassword, encryptedFilePath)
	if err == nil {
		// Decode the json string to data map
		json.Unmarshal([]byte(content), &encryptedDataMap)

		_, encryptedKeyExists := encryptedDataMap[key]
		if encryptedKeyExists {
			delete(encryptedDataMap, key)
			// Encode the data to JSON
			jsonBytes, err := json.Marshal(encryptedDataMap)

			if err == nil {
				// Write the JSON to file
				err := openSsl.EncryptFile(masterPassword, encryptedFilePath, string(jsonBytes))
				if(err != nil) {
					errors.New("An error occurred while encrypting the file: " + err.Error())
				}
			} else {
				errors.New("An error occurred while encoding the string to JSON: " + err.Error())
			}
			return nil
		} else {
			return errors.New("The provided key " + key + " does not exists.")
		}
	}
	return err
}

// Returns all keys
func GetAllKeys(encryptedFilePath string, masterPassword string) ([]string, error) {
	openSsl := OpenSsl()
	content, err := openSsl.DecryptFile(masterPassword, encryptedFilePath)
	encryptedDataMap := make(map[string]string)
	if err == nil {
		json.Unmarshal([]byte(content), &encryptedDataMap)
		encryptedDataMapLen := len(encryptedDataMap)
		if(encryptedDataMapLen > 0) {
			keySlice := make([]string, encryptedDataMapLen)
			for key := range encryptedDataMap {
				keySlice = append(keySlice, key)
			}
			return keySlice, nil
		} else {
			return make([]string, 0), errors.New("There are no keys available.")
		}
		
	} else {
		return make([]string, 0), err
	}
}

// Changes the master password. Old master password is required
func ChangePassword(encryptedFilePath string, oldMasterPassword string, newMasterPassword string) error {
	
	// Initiate the map which will store the json data
	encryptedDataMap := make(map[string]string)

	// Read the password file
	openSsl := OpenSsl()
	content, err := openSsl.DecryptFile(oldMasterPassword, encryptedFilePath)
	if err == nil {
		// Decode the json string to data map
		json.Unmarshal([]byte(content), &encryptedDataMap)
		encryptedDataMapLen := len(encryptedDataMap)
		if(encryptedDataMapLen > 0) {
			for key := range encryptedDataMap {
				dec, err := openSsl.DecryptString(oldMasterPassword + "" + key, encryptedDataMap[key])
				if err == nil {
					enc, err := openSsl.EncryptString(newMasterPassword + "" + key, dec)
					if err != nil {
						return errors.New("An error occurred while encrypting the password using New master password: " + err.Error())
					}
					encryptedDataMap[key] = enc
				} else {
					return errors.New("An error occurred while decrypting the password using Old master password: " + err.Error())
				}
			}

			// Encode the data to JSON
			jsonBytes, err := json.Marshal(encryptedDataMap)

			if err == nil {
				// Write the JSON to file
				err := openSsl.EncryptFile(newMasterPassword, encryptedFilePath, string(jsonBytes))
				if(err != nil) {
					return errors.New("An error occurred while encrypting the file: " + err.Error())
				}
			} else {
				return errors.New("An error occurred while encoding the data to JSON: " + err.Error())
			}
		} else {
			return errors.New("There are no keys available.")
		}
	} else {
		return errors.New("An error occurred while decrypting the password file using Old Password: " + err.Error())
	}
	return nil
}

// Validates the password
func ValidatePassword(password string) ([]string, bool) {
	errorMessageArray := map[string]string{
		"upper case": "1 upper case character",
		"lower case": "1 lower case character",
		"numeric": "1 number",
		"special": "1 special character",
		"length": "6 characters",
	}
	hasError := false
	ret := make([]string, 0)
	if(len(password) <= 6) {
		ret = append(ret, errorMessageArray["length"])
		hasError = true
	} else {
		next: for name, classes := range map[string][]*unicode.RangeTable {
			"upper case": {unicode.Upper, unicode.Title},
			"lower case": {unicode.Lower},
			"numeric":    {unicode.Number, unicode.Digit},
			"special":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
		} {
				for _, r := range password {
						if unicode.IsOneOf(classes, r) {
								continue next
						}
				}
				ret = append(ret, errorMessageArray[name])
				hasError = true
		}
	}
	return ret, hasError
}