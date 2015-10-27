package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"net/http"
	"github.com/knieriem/markdown"	
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	header, _ := ioutil.ReadFile("templates/header.html")
	footer, _ := ioutil.ReadFile("templates/footer.html")
	var buf, page bytes.Buffer
	file, err := os.Open("posts/index.md")
	if err != nil {
		return
	}
	defer file.Close()
	p := markdown.NewParser(nil)
	p.Markdown(file, markdown.ToHTML(&buf))
	page.Write(header)
	page.Write(buf.Bytes())
	page.Write(footer)
	fmt.Fprintf(w, "%s", page.String())
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":80", nil)
}

