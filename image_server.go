package main

import (
    "io/ioutil"
    "net/http"
    "strings"
    "encoding/base64"
    "html/template"
    "path"
)

type PNGServer struct {
}

const (
	PublicKey = "asfsjntvwot2352kofm245m"
	ImageSavePath = "/home/pavvloff/image_server"
)

const pageHTML = `
<html>
  <head>
  <title>Created Images</title>
  </head>
  <body>
  <ul>
    {{range $file := .files}}
      {{if not $file.IsDir}}
        <li>
          <img src="/image/{{$file.Name}}">
        </li>
      {{end}}
    {{end}}
  </ul>
  </body>
 </html>
`

func SaveImage(writer http.ResponseWriter,request *http.Request) {
	switch request.Method {
	case "POST":
		request.ParseForm()

		vals := request.PostForm
		filename := vals["filename"][0]
		key := vals["key"][0]
		texture := vals["texture"][0]
		if key == PublicKey {
			tex,_ := base64.StdEncoding.DecodeString(texture)
			ioutil.WriteFile(path.Join(ImageSavePath,filename), tex, 0666)
			writer.WriteHeader(200)
			reader := strings.NewReader("It Works!")
			reader.WriteTo(writer)
		} else {
			writer.WriteHeader(404)
			reader := strings.NewReader("Not Found - key doesn't match!")
			reader.WriteTo(writer)
		}
	default:
		writer.WriteHeader(404)
		reader := strings.NewReader("Not Found, but keep trying!")
		reader.WriteTo(writer)
	}
}

func ListImages(writer http.ResponseWriter,request *http.Request) {
	switch request.Method {
	case "GET":
		html,err := template.New("results_page").Parse(pageHTML)
		if err != nil {
			writer.WriteHeader(500)
			reader := strings.NewReader("Template broken!")
			reader.WriteTo(writer)
			return
		}
		writer.WriteHeader(200)
		files,_ := ioutil.ReadDir(ImageSavePath)
		html.Execute(writer,map[string]interface{}{"files":files})
	default:
		writer.WriteHeader(404)
		reader := strings.NewReader("Not Found!")
		reader.WriteTo(writer)
	}
}


func main() {
	http.Handle("/image/",http.StripPrefix("/image/",http.FileServer(http.Dir(ImageSavePath))))
	http.HandleFunc("/upload",SaveImage)
	http.HandleFunc("/list/",ListImages)
	http.ListenAndServe(":8080",nil)
}
