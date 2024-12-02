package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Range struct {
	Min, Max float64
	Volume   int
}

var ranges []Range
var fileTitle string

func main() {
	flagDirPath := flag.String("d", "", "Manga directory path")
	flagRangesPath := flag.String("r", "", "Ranges file location")
	flagDestination := flag.String("o", "", "Set the destination of where to create the manga folder")
	flagVerbose := flag.Bool("v", false, "Set output verbosity")
	flagPreserve := flag.Bool("p", false, "Preserver original files")
	flagSpecial := flag.Bool("s", false, "Create special chapters")

	flag.Parse()

	dirPath := *flagDirPath
	mangaNameDir := filepath.Base(dirPath)
	mangaName := strings.ReplaceAll(mangaNameDir, " ", "_")
	var destination string

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	numbers := make([]interface{}, 0)
	fileNames := make([]string, 0)
	for _, each := range files {

		chapterNumber := extractChapterNumber(each.Name())

		numbers = append(numbers, chapterNumber)

		fileTitle = extractFileTitle(each.Name())

		fileNames = append(fileNames, fileTitle)
	}

	foundFileNames := deduplicateSlice(fileNames)

	// load the volRanges file
	if len(*flagRangesPath) == 0 {
		ranges, err = loadRanges("volRanges")
		if err != nil {
			log.Fatal("Failed to load default volRanges file. Please ensure it exists or provide a custom file path with -r flag. \n", err)
		}
	} else {
		ranges, err = loadRanges(*flagRangesPath)
		if err != nil {
			log.Fatalf("Failed to load ranges from file: %v, \n %v", *flagRangesPath, err)
		}
	}

	if *flagDestination != "" {
		destination = *flagDestination
	} else {
		destination = dirPath
	}

	// map to hold categories and their corresponding numbers
	volumeMap := make(map[int][]float64)

	for _, item := range numbers {
		switch v := item.(type) {
		case int:
			volume := categorigeChapters(float64(v), ranges)
			volumeMap[volume] = append(volumeMap[volume], float64(v))
		case float64:
			volume := categorigeChapters(v, ranges)
			volumeMap[volume] = append(volumeMap[volume], v)
		default:
			fmt.Printf("unexpected type: %T", v)
		}
	}

	// loop through the volumeMap
	for volume, nums := range volumeMap {

		//loop through the chapter numbers
		for _, inNums := range nums {
			var oldPath string
			stringChapter := strconv.FormatFloat(inNums, 'f', -1, 64)
			stringVolume := strconv.Itoa(volume)

			for _, oldFilename := range foundFileNames {
				q := dirPath + "/" + oldFilename + stringChapter + ".cbz"
				if fileExists(q) {
					//fmt.Println("this file exists", oldFilename)
					oldPath = dirPath + "/" + oldFilename + stringChapter + ".cbz"

				}
			}

			newPath := destination + "/" + mangaName + "_v" + stringVolume + "_chp" + stringChapter + ".cbz"
			newPathPreserved := destination + "/" + mangaName + "/" + mangaName + "_v" + stringVolume + "_chp" + stringChapter + ".cbz"
			var newPathSpecial string
			var newPathSpecialPreserved string
			var message string
			var newerPath string

			// path based on if special flag was used or not
			switch *flagSpecial {
			case true:
				newPathSpecial = destination + "/" + mangaName + "_v" + stringVolume + "_chp" + stringChapter + "_SP" + stringChapter + ".cbz"
				newPathSpecialPreserved = destination + "/" + mangaName + "/" + mangaName + "_v" + stringVolume + "_chp" + stringChapter + "_SP" + stringChapter + ".cbz"
			case false:
				newPathSpecial = destination + "/" + mangaName + "_v" + stringVolume + "_chp" + stringChapter + ".cbz"
				newPathSpecialPreserved = destination + "/" + mangaName + "/" + mangaName + "_v" + stringVolume + "_chp" + stringChapter + ".cbz"
			}

			// check if chapter is normal or a special
			if inNums == math.Trunc(inNums) {
				message = fmt.Sprintf("Chapter: %s \t\t Renamed to: %s", oldPath, newPath)

				switch *flagPreserve {
				case true:
					newerPath = newPathPreserved
				case false:
					newerPath = newPath
				}

			} else {
				message = fmt.Sprintf("SPECIAL chapter: %s \t Renamed to: %s", oldPath, newPathSpecial)

				switch *flagPreserve {
				case true:
					newerPath = newPathSpecialPreserved
				case false:
					newerPath = newPathSpecial
				}
			}

			// verbose switch
			if *flagVerbose {
				fmt.Println(message)
			}

			// preserver switch
			if *flagPreserve {
				preserverOriginal(oldPath, newerPath)
			} else {
				os.Rename(oldPath, newerPath)
			}

		}
	}
}

func extractChapterNumber(filenName string) interface{} {

	stripedOfType := strings.ReplaceAll(filenName, ".cbz", "")

	reg := regexp.MustCompile(`[0-9]*\.?[0-9]+`)

	findString := reg.FindAllString(stripedOfType, -1)

	stringFromSlice := strings.Join(findString, ", ")

	if len(findString) > 0 {
		return convertStringToNumber(stringFromSlice)
	}

	return nil
}

func convertStringToNumber(s string) interface{} {
	if strings.Contains(s, ".") {
		f, _ := strconv.ParseFloat(s, 64)
		return f
	} else {
		i, _ := strconv.Atoi(s)
		return i
	}
}

func loadRanges(rangeFile string) ([]Range, error) {

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

func preserverOriginal(source string, destination string) {
	origin, err := os.ReadFile(source)

	if err != nil {
		fmt.Println("there was an error with", source, err)
	}

	destinationPath := filepath.Dir(destination)

	mkdirErro := os.MkdirAll(destinationPath, 0755)
	if mkdirErro != nil {
		log.Fatalln("There was a problem with creating directories:", mkdirErro)
	}

	erro := os.WriteFile(destination, []byte(origin), 0644)
	if erro != nil {
		log.Fatalln("There was a problem with writing the file:", erro)
	}

}

func extractFileTitle(originalFileName string) string {

	removeType := strings.ReplaceAll(originalFileName, ".cbz", "")

	reg := regexp.MustCompile(`[0-9]*\.?[0-9]+`)

	result := reg.ReplaceAllString(removeType, "")

	return result

}

func deduplicateSlice(dupSlice []string) []string {
	// map to track unique file names
	unique := make(map[string]struct{})

	var result []string

	// loop over the input slice
	for _, str := range dupSlice {
		// check if string already exists in unique map, if doesn't set append it to the result slice and set the unique[str] to empty
		if _, exists := unique[str]; !exists {
			unique[str] = struct{}{}
			result = append(result, str)
		}
	}
	return result
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}
