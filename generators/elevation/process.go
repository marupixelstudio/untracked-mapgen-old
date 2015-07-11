package elevation

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var profiles []int = []int{10, 50, 100} // contour line elevations in meters

func process(filename string) {
	src := sourceDir + filename
	destTmpl := processedDir + strings.Replace(filename, ".hgt", "", 1)
	var wg sync.WaitGroup
	wg.Add(len(profiles))
	for _, profile := range profiles {
		go func(profile int) {
			defer wg.Done()

			// Build shapefiles
			dest := fmt.Sprintf("%vc%v.shp", destTmpl, profile)
			args := []string{"-i", "10", "-snodata", "32767", "-a", "height", src, dest}
			fmt.Println("gdal_contour", strings.Join(args, " "))
			cmd := exec.Command("gdal_contour", args...)
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}

			// Index the shapefiles
			cmd = exec.Command("shapeindex", dest)
			cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}(profile)
	}
	wg.Wait()
}

func Process(parallelism int) {
	pool := NewPool(parallelism)
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		c := pool.Borrow()
		go func(f os.FileInfo) {
			defer pool.Return(c)
			process(f.Name())
			log.Println(f.Name())
		}(f)
	}
}
