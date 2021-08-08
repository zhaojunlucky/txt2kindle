package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/txt2kindle/pkg/txt"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	txtFileFlag := flag.String("t", "", "txt file to convert")
	chapterPatternFlag := flag.String("c", "chapters-setting.txt", "txt chapter pattern file path")
	flag.Parse()

	chapterFilePath := *chapterPatternFlag
	txtFilePath := *txtFileFlag

	if len(chapterFilePath) == 0 {
		panic("invalid input chapter pattern file path")
	} else if _, err := os.Stat(chapterFilePath); err != nil {
		panic("input chapter pattern file path doesn't exist, " + chapterFilePath)
	}

	if len(txtFilePath) == 0 {
		panic("invalid input txt file to convert")
	} else if _, err := os.Stat(txtFilePath); err != nil {
		panic("input txt file path doesn't exist, " + txtFilePath)
	}


	txtObj := txt.NewTxt(txtFilePath, []string{"MagicWorldZ"})
	txtObj.Parse(chapterFilePath)

	if confirmConvert(txtObj) {
		txtObj.ConvertAsKindle(filepath.Dir(txtFilePath))
	}
}

func confirmConvert(txtObj *txt.Txt) bool {
	fmt.Printf("Novel: %s\n", txtObj.Title)
	chapterCnt := len(txtObj.Chapters)
	fmt.Printf("Total Parsed Chapters: %d\n", chapterCnt)
	fmt.Printf("Chapters Review:\n")
	var begin = int(math.Min(8, float64(chapterCnt)))
	for i := 0; i < begin; i++ {
		fmt.Printf("\t%s\n", txtObj.Chapters[i].Title)
	}
	fmt.Printf("\t...\n")

	end := int(math.Max(float64(begin), float64(chapterCnt - 8)))
	for i := end; i < chapterCnt; i++ {
		fmt.Printf("\t%s\n", txtObj.Chapters[i].Title)
	}

	HandleCtrlC()
	fmt.Println("Press enter to continue the convert, Ctrl+C to exit...")

	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Unknown error", err)
		return false
	}

	return true
}

func HandleCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Stop convert...")
		os.Exit(1)
	}()
}