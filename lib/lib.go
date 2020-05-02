package lib 

import(
	openssl "gopkg.in/Luzifer/go-openssl.v3"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strings"
)

const EncryptedFilePath = "./.p"

func SetPassword(encryptedFilePath string, masterPassword string, key string, password string) error {
	// Initiate the map which will store the json data
	encryptedDataMap := make(map[string]string)

	// Read the password file
	content, err := DecryptFile(masterPassword, encryptedFilePath)
	if err != nil && !strings.Contains(err.Error(), "Invalid Encrypted file path") {
		return err
	} else {
		json.Unmarshal([]byte(content), &encryptedDataMap)
		
		// Initiate the openssl to encrypt the password
		o := openssl.New()
		enc, err := o.EncryptBytes(masterPassword + "" + key, []byte(password), openssl.DigestSHA256Sum)
		if err != nil {
			return errors.New("An error occurred while encrypting the password: " + err.Error())
		}

		// Make the encrypted string to the map
		encryptedDataMap[key] = string(enc)

		// Encode the data to JSON
		jsonBytes, err := json.Marshal(encryptedDataMap)

		if err == nil {
			// Write the JSON to file
			err := EncryptFile(masterPassword, encryptedFilePath, string(jsonBytes))
			if(err != nil) {
				return errors.New("An error occurred while encrypting the file: " + err.Error())
			}
		} else {
			return errors.New("An error occurred while encoding the string to JSON: " + err.Error())
		}

		return nil
	}
}

func GetPassword(encryptedFilePath string, masterPassword string, key string) (string, error) {
	// Initiate the map which will store the json data
	encryptedDataMap := make(map[string]string)

	// Read the password file
	content, err := DecryptFile(masterPassword, encryptedFilePath)
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
				return string(dec), nil
			} else {
				if err.Error() == "invalid padding" {
					return "", errors.New("The master password you provide may be incorrect.")
				} else {
					return "", errors.New("An error occurred while decrypting the string: " + err.Error())
				}
			}
		} else {
			return "", errors.New("The provided key " + key + " does not exists.")
		}
	} else {
		return "", err
	}
}

func DeletePassword(encryptedFilePath string, masterPassword string, key string) error {
	// Initiate the map which will store the json data
	encryptedDataMap := make(map[string]string)

	// Read the password file
	content, err := DecryptFile(masterPassword, encryptedFilePath)
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
				err := EncryptFile(masterPassword, encryptedFilePath, string(jsonBytes))
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

func GetAllKeys(encryptedFilePath string, masterPassword string) ([]string, error) {
	content, err := DecryptFile(masterPassword, encryptedFilePath)
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

func EncryptFile(masterPassword string, filePath string, fileData string) error {
	o := openssl.New()
	encFileData, err := o.EncryptBytes(masterPassword, []byte(fileData), openssl.DigestSHA256Sum)
	if(err == nil) {
		err := ioutil.WriteFile(filePath, encFileData, 0644)
		if err != nil {
			return err
		}
	}
	return err
}