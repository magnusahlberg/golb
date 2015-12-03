package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"strconv"
	"time"
	"github.com/knieriem/markdown"
)

type Article struct {
	name string
	date time.Time
	path string
	html string
}

type ByDate []Article

func (a ByDate) Len() int { return len(a ) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].date.After(a[j].date) }

var (
	articles []Article
)

func fetchArticles() {
	walkErr := filepath.Walk("posts/", visit)
	if walkErr != nil {
		fmt.Printf("Walk error: %s", walkErr)
	}
	sort.Sort(ByDate(articles))
}

func visit(path string, info os.FileInfo, err error) error {
	splitPath := strings.Split(path, "/")
	if len(splitPath) > 4 {
		name := splitPath[len(splitPath)-1]
		name = name[:len(name)-len(".md")]
		year, err := strconv.Atoi(splitPath[1])
		month, err := strconv.Atoi(splitPath[2])
		day, err := strconv.Atoi(splitPath[3])
		if err != nil {
			return nil //TODO: Error handling
		}
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		var buf bytes.Buffer
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Open error: %s", err)
			return nil //TODO: Error handling
		}
		defer file.Close()
		p := markdown.NewParser(nil)
		p.Markdown(file, markdown.ToHTML(&buf))

		article := Article{name, date, path, buf.String()}
		articles = append(articles, article)
		fmt.Printf("\nPath: %s\nDate: %s\nName: %s\nHTML:%s\n\n", path, date.String(), name, buf)
	}
	return nil
}

func renderPage(page bytes.Buffer, w http.ResponseWriter) {
	header, _ := ioutil.ReadFile("templates/header.html")
	footer, _ := ioutil.ReadFile("templates/footer.html")
	fmt.Fprintf(w, "%s %s %s", header, page.String(), footer)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var page bytes.Buffer
	if len(articles) == 0 {
		fetchArticles()
	}
	for _, article := range articles {
		fmt.Printf("Article: %s", article.name)
		page.Write([]byte(article.html))
	}
	renderPage(page, w)
}

func clearCache() {
	articles = nil
}

func main() {
	ticker := time.NewTicker(time.Minute * 30)
	go func() {
		for _ = range ticker.C {
			clearCache()
		}
	}()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8181", nil)
}
