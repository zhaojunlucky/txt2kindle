package util

import (
	"bufio"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os"
	"unicode/utf8"
)

func DetectFileEnc(txtFile string) encoding.Encoding {
	file, err := os.Open(txtFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		b := []byte(line)

		if !utf8.FullRune(b) {
			return simplifiedchinese.GB18030
		}

	}
	return encoding.Nop
}
