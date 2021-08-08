package txt

import (
	"bufio"
	"fmt"
	"github.com/leotaku/mobi"
	"github.com/txt2kindle/pkg/util"
	"golang.org/x/text/language"
	"golang.org/x/text/transform"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Txt struct {
	Title string
	FilePath string
	Authors []string
	Chapters []*TxtChapter
}

func NewTxt(filePath string, authors []string) *Txt {
	title := filepath.Base(filePath)
	title = strings.TrimSuffix(title, ".txt")
	return &Txt{
		Title: title,
		FilePath: filePath,
		Authors: authors,
		Chapters: []*TxtChapter{},
	}
}

func (txt *Txt)ConvertAsKindle(saveDir string)  {
	kindleFilePath := path.Join(saveDir, fmt.Sprintf("%s.azw3", txt.Title))
	fmt.Printf("Saving kindle file to %s", kindleFilePath)
	chapters := make([]mobi.Chapter, len(txt.Chapters))

	for i := 0; i < len(txt.Chapters); i++ {
		txtCh := txt.Chapters[i]
		builder := strings.Builder{}
		builder.WriteString("<h1>")
		builder.WriteString(txtCh.Title)
		builder.WriteString("</h1>")
		for line := range txtCh.Lines {
			builder.WriteString("<p>")
			builder.WriteString(txtCh.Lines[line])
			builder.WriteString("</p>")
		}

		ch := mobi.Chapter{
			Title:  txtCh.Title,
			Chunks: mobi.Chunks(builder.String()),
		}
		chapters = append(chapters, ch)
	}

	mb := mobi.Book{
		Title:       txt.Title,
		Authors:     txt.Authors,
		CreatedDate: time.Now(),
		Language:    language.Italian,
		Chapters:    chapters,
		UniqueID:    rand.Uint32(),
	}

	// Convert book to PalmDB database
	db := mb.Realize()

	// Write database to file
	f, _ := os.Create(kindleFilePath)
	err := db.Write(f)
	if err != nil {
		panic(err)
	}
}

func (txt *Txt) Parse(titlePatternFile string) {
	chapterTitle := NewChapterTitle(titlePatternFile)
	fileEnc := util.DetectFileEnc(txt.FilePath)
	file, err := os.Open(txt.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := transform.NewReader(file, fileEnc.NewDecoder())

	scanner := bufio.NewScanner(reader)

	txtChapter := newTxtChapter(txt.Title)
	txt.Chapters = append(txt.Chapters, txtChapter)
	for scanner.Scan() {
		line := scanner.Text()
		if chapterTitle.IsTitle(line) {
			txtChapter = newTxtChapter(line)
			txt.Chapters = append(txt.Chapters, txtChapter)
		} else {
			txtChapter.Lines = append(txtChapter.Lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type TxtChapter struct {
	Title string
	Lines []string
}

func newTxtChapter(title string) *TxtChapter {
	return &TxtChapter {
		Title: title,
		Lines: []string{},
	}
}

