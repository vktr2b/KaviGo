package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

type Manga struct {
	Directory   string
	MangaName   string
	ChapterName string
	Chapter     float64
	Extention   string
	Volume      int
}

type Range struct {
	Min, Max float64
	Volume   int
}

var d string
var o string
var r string
var v bool
var p bool

func main() {

	runCli()

	//fmt.Println(getDataFromManga(d))

	data, err := getDataFromManga(d)
	checkFatalErr(err)

	for _, name := range data {
		origin := filepath.Clean(name.Directory + "/" + name.ChapterName)
		var toBe string

		if p == true && o == "" {
			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", name.Directory, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = (filepath.Clean(dirty))

			err := customCopy(origin, toBe)
			checkFatalErr(err)

		} else if p == false && len(o) >= 1 {
			checkOutputDir(o)
			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", o, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = filepath.Clean(dirty)

			err := customCopy(origin, toBe)
			checkFatalErr(err)

			err = os.Remove(origin)
			checkFatalErr(err)

		} else if p == true && len(o) >= 1 {
			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", o, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = filepath.Clean(dirty)

			err := customCopy(origin, toBe)
			checkFatalErr(err)

		} else {
			dirty := fmt.Sprintf("%v/%v_v%v_chp%v%v", name.Directory, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = filepath.Clean(dirty)

			os.Rename(origin, toBe)
		}

		if v == true {
			fmt.Printf("original: %v \t Renamed to %v \n", origin, toBe)
		}

	}

}

func getDataFromManga(directory string) ([]Manga, error) {
	var manga []Manga

	//load the ranges file
	r, err := loadRanges(r)
	checkFatalErr(err)

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
		extension := getFileExtension(fileName)

		chapter := extractOneChapterNumber(fileName, extension)

		volume := categorigeChapters(chapter, r)
		manga = append(manga, Manga{
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

func extractOneChapterNumber(fileName string, extention string) float64 {

	// strip filename of the extension
	s := strings.ReplaceAll(fileName, extention, "")

	// regex
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

func loadRanges(rangeFile string) ([]Range, error) {
	var ranges []Range

	file, err := os.Open(rangeFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			log.Fatal("File format is wrong")
		}
		min, _ := strconv.ParseFloat(parts[0], 64)
		max, _ := strconv.ParseFloat(parts[1], 64)
		volume, _ := strconv.Atoi(parts[2])
		ranges = append(ranges, Range{Min: min, Max: max, Volume: volume})
	}
	return ranges, scanner.Err()
}

func categorigeChapters(num float64, ranges []Range) int {
	for _, r := range ranges {
		rounedDown := math.Floor(num)
		if rounedDown >= r.Min && rounedDown <= r.Max {
			return r.Volume
		}
	}
	return 0
}

func getFileExtension(f string) string {
	extension := f[len(f)-4:]

	return extension
}

func runCli() {

	app := &cli.App{
		Name:  "KaviGo",
		Usage: "fight the loneliness!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "d",
				Value:       "",
				Usage:       "input directory",
				Destination: &d,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "o",
				Value:       "",
				Usage:       "Output directory",
				Destination: &o,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "r",
				Value:       "./volRanges",
				Usage:       "Path to the Volume ranges file (comma-delimited)",
				Destination: &r,
				Required:    false,
			},
			&cli.BoolFlag{
				Name:        "v",
				Value:       false,
				Usage:       "Verbose output",
				Destination: &v,
				Required:    false,
			},
			&cli.BoolFlag{
				Name:        "p",
				Value:       false,
				Usage:       "Preserve original files",
				Destination: &p,
				Required:    false,
			},
		},
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func checkFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func checkOutputDir(o string) {

	if _, err := os.Stat(o); os.IsNotExist(err) {
		var answer string

		// Ask user if he wants to create the directory if it does not exist
		for {
			fmt.Println("Directory " + o + "does not exist \n Would you like it to be created ?(yes/no):")
			fmt.Scanln(&answer)
			if answer == "yes" || answer == "no" {
				break
			}
			fmt.Println("Invalid input. Please try again.")

		}

		// if answer yes make the directory otherwise exit
		if answer == "yes" {
			mkdirErr := os.MkdirAll(o, 0755)
			if mkdirErr != nil {
				log.Fatalln("There was a problem with creating directories:", mkdirErr)
			}

		} else {
			os.Exit(1)
		}

	}

}

func customCopy(original string, destination string) error {

	// read the original file
	origin, err := os.ReadFile(original)
	if err != nil {
		return err
	}

	// get the directory of the path, if it does not exist create it
	destinationPath := filepath.Dir(destination)
	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		mkdirErro := os.MkdirAll(destinationPath, 0755)
		if mkdirErro != nil {
			//log.Fatalln("There was a problem with creating directories:", mkdirErro)
			return mkdirErro
		}
	}

	// write a new file with the new name to the new directory
	erro := os.WriteFile(destination, []byte(origin), 0644)
	if erro != nil {
		//log.Fatalln("There was a problem with writing the file:", erro)
		return erro
	}

	return nil
}
