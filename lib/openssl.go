package lib 

import(
	openssl "github.com/Luzifer/go-openssl/v3"
	"io/ioutil"
	"errors"
)

type openSsl struct {}

// This is a constructor
func OpenSsl() openSsl {
	return openSsl{}
}

// Encrypts the input string
func (op *openSsl) EncryptString(encKey string, s string) (string, error) {
	o := openssl.New()
	enc, err := o.EncryptBytes(encKey, []byte(s), openssl.DigestSHA256Sum)
	if err != nil {
		return "", errors.New("An error occurred while encrypting the string: " + err.Error())
	}

	// Make the encrypted string to the map
	return string(enc), nil
}

// Decrypts the input string
func (op *openSsl) DecryptString(encKey string, s string) (string, error) {
	o := openssl.New()
	dec, err := o.DecryptBytes(encKey, []byte(s), openssl.DigestSHA256Sum)
	if err == nil {
		return string(dec), nil
	} else if err.Error() == "invalid padding" {
		return "", errors.New("The encryption key you provide may be incorrect.")
	} else {
		return "", errors.New("An error occurred while decrypting the string: " + err.Error())
	}
}

// Reads the encrypted file and returns the decrypted content
func (op *openSsl) DecryptFile(masterPassword string, filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err == nil {
		dec, err := op.DecryptString(masterPassword, string(content))
		if(err == nil) {
			return string(dec), nil
		}
		return "", errors.New("The master password may be incorrect")
	}
	return "", errors.New("Invalid Encrypted file path")
}

// Encrypts the content and writes to the file
func (op *openSsl) EncryptFile(masterPassword string, filePath string, fileData string) error {
	encFileData, err := op.EncryptString(masterPassword, fileData)
	if(err == nil) {
		err := ioutil.WriteFile(filePath, []byte(encFileData), 0644)
		if err != nil {
			return err
		}
	}
	return err
}