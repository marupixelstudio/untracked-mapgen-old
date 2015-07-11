package elevation

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func add(file *os.File, shapefile string, profile int, i int) {
	s := `
<Layer name="srtm100-` + strconv.Itoa(i) + `" status="on" srs="+proj=latlong +datum=WGS84">
  <StyleName>contours` + strconv.Itoa(profile) + `</StyleName>
  <StyleName>contours-text` + strconv.Itoa(profile) + `</StyleName>
  <Datasource>
    <Parameter name="type">shape</Parameter>
    <Parameter name="file">` + (processedDir + shapefile) + `</Parameter>
  </Datasource>
</Layer>
`
	file.Write([]byte(s))
}

func Styleize() {
	i := 0
	files, err := ioutil.ReadDir(processedDir)
	if err != nil {
		log.Fatal(err)
	}
	contours, err := os.Create(styleDir + "layers-contours.include")
	defer contours.Close()

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".shp") {
			i += 1
			if strings.HasSuffix(f.Name(), "c100.shp") {
				add(contours, f.Name(), 100, i)
			} else if strings.HasSuffix(f.Name(), "c50.shp") {
				add(contours, f.Name(), 50, i)
			} else if strings.HasSuffix(f.Name(), "c10.shp") {
				add(contours, f.Name(), 10, i)
			}
		}
	}
}
