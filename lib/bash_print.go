package lib 

import(
	"fmt"
)

// Structure which contains error codes.
type BashPrinterStruct struct {
	SUCCESS string
	ERROR string
	WARN string
	NORMAL string
	bashTextColorMap map[string]string
}

// This is a constructor
func BashPrinter() BashPrinterStruct {
	
	bashTextColorMap := make(map[string]string, 4)
	bashTextColorMap["success"] = "\033[0;32m"
	bashTextColorMap["error"] = "\033[0;31m"
	bashTextColorMap["normal"] = "\033[0;0m"
	bashTextColorMap["warn"] = "\033[0;33m"

	return BashPrinterStruct {
		SUCCESS: "success",
		ERROR: "error",
		WARN: "warn",
		NORMAL: "normal",
		bashTextColorMap: bashTextColorMap,
	}
}

// Returns the message after wrapping it with color codes.
func (bp *BashPrinterStruct) GetMessage(message string, messageType string) string {
	color, colorExists := bp.bashTextColorMap[messageType]

	if !colorExists {
		color = bp.bashTextColorMap["normal"]
	}

	var msg string = color + message + bp.bashTextColorMap["normal"]
	return msg
}

// Prints the message after wrapping it with color codes.
func (bp *BashPrinterStruct) PrintMessage(message string, messageType string) {
	var msg string = bp.GetMessage(message, messageType)
	fmt.Println(msg)
}