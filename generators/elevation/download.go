package elevation

import (
	"archive/zip"
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const baseUrl string = `http://dds.cr.usgs.gov/srtm/version2_1/SRTM1/`

var pool *Pool

func download(url string, filename string) string {
	// Skip if downloaded
	hgt := sourceDir + strings.Replace(filename, ".zip", "", 1)
	if _, err := os.Stat(hgt); err == nil {
		log.Println("file exists: %s", hgt)
		return strings.Replace(filename, ".zip", "", 1)
	}

	// Download
	log.Println("Downloading", filename)
	resp, err := http.Get(url + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Unpack response
	body, err := ioutil.ReadAll(resp.Body)
	r, err := zip.NewReader(bytes.NewReader(body), resp.ContentLength)
	if err != nil {
		log.Fatal(err)
	}

	// Read and save entries
	for _, zf := range r.File {
		filename := sourceDir + zf.Name
		dst, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer dst.Close()
		src, err := zf.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()

		io.Copy(dst, src)
		return zf.Name
	}
	return ""
}

func scanFolder(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		if strings.HasSuffix(href, ".zip") {
			c := pool.Borrow()
			go func() {
				defer pool.Return(c)
				filename := download(url, href)
				if filename != "" {
					process(filename)
				}
			}()
		}

		if !strings.HasPrefix(href, "/") && strings.HasSuffix(href, "/") {
			log.Printf("%+v", href)
			scanFolder(url + href)
		}
	})
}

func Download(parallelism int) {
	pool = NewPool(parallelism)

	scanFolder(baseUrl)
}
