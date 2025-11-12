package globvars

var R string = "./volRanges"   // ranges file to load
var D string                   // directory where the manga is located
var O string                   // output directory if required
var PR string                  // pre loaded ranges from conf
var C string = "./kavigo.yaml" // yaml config file
var V bool                     // verbosity flag
var P bool                     // preserve original files flag
var Host string
var Port int
var User string
var Password string
var SshKey string
var CopyToRemote bool // wether remote options should be used or the local ones

type Range struct {
	Min, Max float64
	Volume   int
}

type Manga struct {
	Directory   string
	MangaName   string
	ChapterName string
	Chapter     float64
	Extention   string
	Volume      int
}

type Conf struct {
	Directories struct {
		Manga       string `yaml:"manga"`
		Destination string `yaml:"destination"`
	}
	Options struct {
		Verbosity    bool `yaml:"verbosity"`
		Preserve     bool `yaml:"preserve"`
		CopyToRemote bool `yaml:"copytoremote"`
	}
	Remote struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SshKey   string `yaml:"sshkey"`
	}
	Ranges string `yaml:"ranges"`
}
