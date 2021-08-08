package txt

import (
	"bufio"
	"github.com/txt2kindle/pkg/util"
	"golang.org/x/text/transform"
	"log"
	"os"
	"regexp"
	"strings"
)

type ChapterTitle struct {
	MaxLength int
	Patterns []*regexp.Regexp
}

func NewChapterTitle(patternFile string) *ChapterTitle {
	chapterTitle :=  &ChapterTitle{
		MaxLength: 60,
		Patterns: []*regexp.Regexp{},
	}
	chapterTitle.parsePatterns(patternFile)
	return chapterTitle
}

func (chapterTitle *ChapterTitle)parsePatterns(patternFile string) {
	fileEnc := util.DetectFileEnc(patternFile)

	file, err := os.Open(patternFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := transform.NewReader(file, fileEnc.NewDecoder())

	scanner := bufio.NewScanner(reader)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			log.Printf("skip comment line %s\n", line)
		} else {
			pattern, err := regexp.Compile(line)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("loaded title pattern %s\n", line)
			chapterTitle.Patterns = append(chapterTitle.Patterns, pattern)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (chapterTitle *ChapterTitle)IsTitle(line string) bool {
	if line == "" {
		return false
	}
	for _, regex := range chapterTitle.Patterns {
		if len(line) < chapterTitle.MaxLength && regex.MatchString(line) {
			return true
		}
	}
	return false
}