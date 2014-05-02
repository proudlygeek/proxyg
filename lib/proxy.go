package proxyg

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func HandleConn(w http.ResponseWriter, r *http.Request) {
	path := strings.Join([]string{"http://", strings.TrimLeft(r.URL.Path, "/")}, "")
	header, body := proxyCall(path)
	copyHeader(w.Header(), header)
	fmt.Fprint(w, string(body))
}

func Listen() {
	http.HandleFunc("/", HandleConn)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func proxyCall(name string) (header http.Header, body []byte) {
	res, err := http.Get(name)

	if err != nil {
		log.Panic(err)
	}

	header = res.Header
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()

	return header, body
}
