package extractor

import (
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"kavigo/pkgs/filehandler"
	"kavigo/pkgs/globvars"
	"kavigo/pkgs/loader"
)

func GetDataFromManga(directory string) ([]globvars.Manga, error) {
	var manga []globvars.Manga

	//load the ranges file
	ranges, err := loader.LoadRanges(globvars.R)
	if err != nil {
		log.Fatal(err)
	}

	//extract the manga name from the top directory
	mn := filepath.Base(directory)
	//repalce every space with "_"
	mnf := strings.ReplaceAll(mn, " ", "_")

	// read the contents of the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	//loop through the contents of the directory and do some magic
	for _, file := range files {
		fileName := file.Name()
		extension := filehandler.GetFileExtension(fileName)

		chapter := extractOneChapterNumber(fileName, extension)

		volume := categorizeChapters(chapter, ranges)
		manga = append(manga, globvars.Manga{
			Directory:   directory,
			MangaName:   mnf,
			ChapterName: file.Name(),
			Chapter:     chapter,
			Volume:      volume,
			Extention:   extension,
		})

	}

	return manga, err

}

func categorizeChapters(num float64, ranges []globvars.Range) int {
	for _, r := range ranges {
		rounedDown := math.Floor(num)
		if rounedDown >= r.Min && rounedDown <= r.Max {
			return r.Volume
		}
	}
	return 0
}

func extractOneChapterNumber(fileName string, extention string) float64 {

	// strip filename of the extension
	s := strings.ReplaceAll(fileName, extention, "")

	// regex magic
	reg := regexp.MustCompile(`[0-9]*\.?[0-9]+`)

	// find all the numbers with the regex, in the striped filename
	findString := reg.FindAllString(s, -1)

	//Turn the findString slice to a String
	toString := strings.Join(findString, ", ")

	ret, err := strconv.ParseFloat(toString, 64)
	if err != nil {
		log.Fatal(err)
	}

	return ret

}
