package globvars

var R string // ranges file to load
var D string // directory where the manga is located
var O string // output directory if required
var V bool   // verbosity flag
var P bool   // preserve original files flag

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
