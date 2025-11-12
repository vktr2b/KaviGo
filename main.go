package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"kavigo/pkgs/cli"
	"kavigo/pkgs/dirhandler"
	"kavigo/pkgs/filehandler"
	"kavigo/pkgs/globvars"
	"kavigo/pkgs/parser"
	"kavigo/pkgs/sftp"
)

func init() {

	// initialize CLI
	cli.RunCli()

	// Read the yaml conf
	parser.ReadConf()

	// Create a Remote connection
	if globvars.CopyToRemote {
		err := sftp.CreateRemoteConn()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {

	d := globvars.D
	o := globvars.O
	v := globvars.V
	p := globvars.P
	copyToRemote := globvars.CopyToRemote

	data, err := parser.GetDataFromManga(d)
	checkFatalErr(err)

	for _, name := range data {
		origin := filepath.Clean(name.Directory + "/" + name.ChapterName)
		var toBe string

		if p && o == "" && !copyToRemote {
			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", name.Directory, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = (filepath.Clean(dirty))

			err := filehandler.CopyToPreserve(origin, toBe)
			checkFatalErr(err)

		} else if !p && len(o) >= 1 && !copyToRemote {
			dirhandler.CheckOutputDir(o)
			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", o, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = filepath.Clean(dirty)

			err := filehandler.CopyToPreserve(origin, toBe)
			checkFatalErr(err)

			err = os.Remove(origin)
			checkFatalErr(err)

		} else if p && len(o) >= 1 && !copyToRemote {
			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", o, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = filepath.Clean(dirty)

			err := filehandler.CopyToPreserve(origin, toBe)
			checkFatalErr(err)

		} else if copyToRemote && len(o) >= 1 {

			dirty := fmt.Sprintf("%v/%v/%v_v%v_chp%v%v", o, name.MangaName, name.MangaName, name.Volume, name.Chapter, name.Extention)

			//fmt.Println(o)

			toBe = filepath.Clean(dirty)

			//fmt.Println(o)

			toMk := filepath.Clean(fmt.Sprintf("%v/%v/", o, name.MangaName))

			//fmt.Println(toMk)

			sftp.CopyFilesToRemote(origin, toBe, toMk)

		} else if !copyToRemote {
			dirty := fmt.Sprintf("%v/%v_v%v_chp%v%v", name.Directory, name.MangaName, name.Volume, name.Chapter, name.Extention)

			toBe = filepath.Clean(dirty)

			os.Rename(origin, toBe)
		} else {
			fmt.Println("wrong config, this probably will not be triggered but if yes open a issue on github")
		}

		if v {
			fmt.Printf("original: %v \t Renamed to %v \n", origin, toBe)
		}

	}

	if copyToRemote {
		sftp.CloseRemoteConn()
	}

}

func checkFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
