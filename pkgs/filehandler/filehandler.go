package filehandler

import (
	"os"
	"path/filepath"
)

func GetFileExtension(f string) string {
	extension := f[len(f)-4:]

	return extension
}

func CopyToPreserve(original string, destination string) error {

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
