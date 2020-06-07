package lib 

import(
	openssl "github.com/Luzifer/go-openssl/v3"
	"io/ioutil"
	"crypto/sha1"
	"io"
	"errors"
)

type openSsl struct {}

// This is a constructor
func OpenSsl() openSsl {
	return openSsl{}
}

// Encrypts the input string
func (op *openSsl) EncryptString(encKey string, s string) (string, error) {
	var encKeySha1 string = op.encodeSha1(encKey)
	o := openssl.New()
	enc, err := o.EncryptBytes(encKeySha1, []byte(s), openssl.DigestSHA256Sum)
	if err != nil {
		return "", errors.New("An error occurred while encrypting the string: " + err.Error())
	}

	// Make the encrypted string to the map
	return string(enc), nil
}

// Decrypts the input string
func (op *openSsl) DecryptString(encKey string, s string) (string, error) {
	var encKeySha1 string = op.encodeSha1(encKey)
	o := openssl.New()
	dec, err := o.DecryptBytes(encKeySha1, []byte(s), openssl.DigestSHA256Sum)
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
		// Here since the final encryption is done using masterPassword, we will decrypt in reverse manner
		dec, err := op.DecryptString(masterPassword, string(content))
		if(err == nil) {
			extractor := Extractor(masterPassword)
			stringGroups := extractor.GetStringGroups()
			combinations := extractor.GenerateCombinations(stringGroups)
			// Loop the combination in reverse manner
			for i := len(combinations)-1; i >= 0; i-- {
				value := combinations[i]
				dec, _ = op.DecryptString(value, dec)
			}
			return string(dec), nil
		}
		return "", errors.New("The master password may be incorrect")
	}
	return "", errors.New("Invalid Encrypted file path")
}

// Encrypts the content and writes to the file
func (op *openSsl) EncryptFile(masterPassword string, filePath string, fileData string) error {
	// The final encryption must be done using the master password. 
	// The initial layers of encryption must be done using combination of password
	extractor := Extractor(masterPassword)
	stringGroups := extractor.GetStringGroups()
	combinations := extractor.GenerateCombinations(stringGroups)
	encFileData := fileData
	for _, value := range combinations {
		encFileData, _ = op.EncryptString(value, encFileData)
	}

	encFileData, err := op.EncryptString(masterPassword, encFileData)

	if(err == nil) {

		err := ioutil.WriteFile(filePath, []byte(encFileData), 0644)
		if err != nil {
			return err
		}
	}
	return err
}

// Encode string to SHA1
func (op *openSsl) encodeSha1(encKey string) string {
	h := sha1.New()
	io.WriteString(h, encKey)
	encKeySha1 := string(h.Sum(nil))
	return encKeySha1
}