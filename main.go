package main

import (
	"./generators/elevation"
)

func main() {
	elevation.Download(3)
	elevation.Process(2)
	elevation.Styleize()
	elevation.Import()
}
