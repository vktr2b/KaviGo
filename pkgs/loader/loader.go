package loader

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"kavigo/pkgs/globvars"
)

func LoadRanges() ([]globvars.Range, error) {

	var ranges []globvars.Range
	var source io.Reader

	if globvars.PR != "" {
		source = strings.NewReader(globvars.PR)
	} else {
		source, err := os.Open(globvars.R)
		if err != nil {
			return nil, err
		}
		defer source.Close()

	}

	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			log.Fatal("Error while parsing the Ranges")
		}
		min, _ := strconv.ParseFloat(parts[0], 64)
		max, _ := strconv.ParseFloat(parts[1], 64)
		volume, _ := strconv.Atoi(parts[2])
		ranges = append(ranges, globvars.Range{Min: min, Max: max, Volume: volume})
	}
	return ranges, scanner.Err()
}
