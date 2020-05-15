# passpass
PassPass is a simple command-line Password Manager built in Golang. It prompts the user for a master password every time and hence does not store the master password.

Every password is encrypted using the master password and stored in a file `p.psdb` in JSON format. The JSON is also encrypted using the master password. The default directory is `$HOME/.passpass/p.psdb`

The password file is portable and user just needs to remember the master password as it is not stored anywhere.


#### Usage:
------
	go run passpass.go --help

#### Generating the build:
	go build passpass.go

#### Options:
	Usage: passpass [-p] [-g value] [--help] [-s value] [parameters ...]
	 -d, --delete=value value must be {keyname}. Deletes the key.
	 -f, --password-file=value Provide the path of password file
	 -g, --get=value    value must be {keyname}. Get password for the corresponding key
	 --help             Help
	 -p, --get-keys     Print all keys
	 -s, --set=value    value must be {keyname:password}. Sets key & password.

#### Examples:
###### Set a password

	go run passpass.go -s 'keyname:password'

###### Retrieve password

	go run passpass.go -g keyname

###### Delete password

	go run passpass.go -d keyname

###### Retrieve all keys

	go run passpass.go --get-keys

###### Copy to Clipboard Dependencies:
	OSX - No Dependencies
	Windows 7 (probably work on other Windows) - No Dependencies
	Linux, Unix (requires 'xclip' or 'xsel' command to be installed)
	
### Change Log:
1. Added strong password constraint.
2. Users can change password
3. Password for a key shall be received using prompt so that it does not get logged into history

