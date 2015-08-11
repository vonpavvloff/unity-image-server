package main

import (
    "io/ioutil"
    "net/http"
    "strings"
    "encoding/base64"
)

type PNGServer struct {
	SavePath string
}

const (
	PublicKey = "asfsjntvwot2352kofm245m"
)

func (server PNGServer)ServeHTTP(writer http.ResponseWriter,request *http.Request) {
	switch request.Method {
	case "GET":
		writer.WriteHeader(200)
		reader := strings.NewReader("It Works!")
		reader.WriteTo(writer)
	case "POST":
		request.ParseForm()
		vals := request.PostForm
		filename := vals["filename"][0]
		key := vals["key"][0]
		texture := vals["texture"][0]
		if key == PublicKey {
			tex,_ := base64.StdEncoding.DecodeString(texture)
			ioutil.WriteFile(server.SavePath + filename, tex, 0666)
			
			writer.WriteHeader(200)
			reader := strings.NewReader("It Works!")
			reader.WriteTo(writer)
		} else {
			writer.WriteHeader(404)
			reader := strings.NewReader("Not Found!")
			reader.WriteTo(writer)
		}
	default:
		writer.WriteHeader(404)
		reader := strings.NewReader("Not Found!")
		reader.WriteTo(writer)
	}

}


func main() {
	s := PNGServer {"./"}
	http.ListenAndServe(":8080", s)
}
