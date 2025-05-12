package loader

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"kavigo/pkgs/globvars"
)

func LoadRanges(rangeFile string) ([]globvars.Range, error) {

	var ranges []globvars.Range

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
		ranges = append(ranges, globvars.Range{Min: min, Max: max, Volume: volume})
	}
	return ranges, scanner.Err()
}
