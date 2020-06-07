package main

import(
	"fmt"
	"github.com/nitesh-/passpass/lib"
	"crypto/sha1"
	"io"
)

func main() {
	extractor := lib.Extractor("!  &uhHkd&jd8U42m24emk")
	stringGroups := extractor.GetStringGroups()
	combinations := extractor.GenerateCombinations(stringGroups)

	for _, value := range combinations {
		h := sha1.New()
		io.WriteString(h, value)
		fmt.Println(1, string(h.Sum(nil)))
	}
}
