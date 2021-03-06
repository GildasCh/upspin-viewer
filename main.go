package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gildasch/upspin-viewer/anonymous"
	_ "upspin.io/transports"
	"upspin.io/upspin"
)

func main() {
	client, err := anonymous.NewClient()
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := strings.TrimPrefix(r.URL.Path, "/")

		if url == "" {
			fmt.Fprintf(w, "<p>Looking for inspiration? Try out <a href='/augie@upspin.io'>augie@upspin.io's account</a></p>")
			return
		}

		de, err := client.Lookup(upspin.PathName(url), true)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if de.IsDir() {
			des, err := client.Glob(upspin.AllFilesGlob(upspin.PathName(url)))
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			for _, dess := range des {
				fmt.Fprintf(w, "<a href='/%s'>%s</a><br />", dess.SignedName, dess.SignedName)
			}
			return
		}

		f, err := client.Open(upspin.PathName(url))
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		_, err = io.Copy(w, f)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Listening on %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
