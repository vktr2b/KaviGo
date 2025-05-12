package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"kavigo/pkgs/cli"
	"kavigo/pkgs/dirhandler"
	"kavigo/pkgs/extractor"
	"kavigo/pkgs/filehandler"
	"kavigo/pkgs/globvars"
)

func main() {

	d := globvars.D
	o := globvars.O
	v := globvars.V
	p := globvars.P

	cli.RunCli()

	//fmt.Println(getDataFromManga(d))

	data, err := extractor.GetDataFromManga(d)
	checkFatalErr(err)

	for _, name := range data {
		origin := filepath.Clean(name.Directory + "/" + name.ChapterName)
		var toBe string

		if p == true && o == "" {
			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", name.Directory, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = (filepath.Clean(dirty))

			err := filehandler.CopyToPreserve(origin, toBe)
			checkFatalErr(err)

		} else if p == false && len(o) >= 1 {
			dirhandler.CheckOutputDir(o)
			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", o, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = filepath.Clean(dirty)

			err := filehandler.CopyToPreserve(origin, toBe)
			checkFatalErr(err)

			err = os.Remove(origin)
			checkFatalErr(err)

		} else if p == true && len(o) >= 1 {
			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", o, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = filepath.Clean(dirty)

			err := filehandler.CopyToPreserve(origin, toBe)
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

func checkFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
