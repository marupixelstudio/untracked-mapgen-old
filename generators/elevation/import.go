package elevation

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func pipe(shp2sql *exec.Cmd, psql *exec.Cmd) {
	psql.Stdin, _ = shp2sql.StdoutPipe()
	psql.Start()
	shp2sql.Run()
	psql.Wait()
}

func imp(filename string) {
	src := processedDir + filename
	shp2sql := exec.Command("shp2pgsql", "-a", "-g", "geometry", src, "contours")
	psql := exec.Command("psql", "-q", "-d", "untracked")
	pipe(shp2sql, psql)
}

func Import() {
	files, err := ioutil.ReadDir(processedDir)
	if err != nil {
		log.Fatal(err)
	}
	schemaCreated := false
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".shp") {
			log.Println(f.Name())
			if schemaCreated == false {
				src := processedDir + f.Name()
				shp2sql := exec.Command("shp2pgsql", "-p", "-I", "-g", "geometry", src, "contours")
				psql := exec.Command("psql", "-q", "-d", "untracked")
				pipe(shp2sql, psql)
				schemaCreated = true
			}
			imp(f.Name())
		}
	}
}
