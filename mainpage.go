package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testtask2/pgk/api"
)

var homePage string
var port = "8081"

const homePageFileName = "homePage.html"

func init() {
	homePageBytes, err := ioutil.ReadFile(homePageFileName)
	if err != nil {
		log.Fatalln(err)
	}
	homePage = string((homePageBytes))
}
func main() {
	var playlistId, key string
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, homePage)
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		playlistId = r.FormValue("playlistID")
		key = r.FormValue("key")

		conn, err := grpc.Dial("grpc-server:"+port, grpc.WithInsecure()) //0.0.0.0 если нужно запустить 2 файла через go run, если через docker-copmose то так grpc-server:
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()
		client := api.NewCheckYoutubeClient(conn)

		response, err := client.Check(context.Background(), &api.CheckRequest{Key: key, IdPlay: playlistId})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Fprint(w, response.GetList())
	})
	log.Fatalln(http.ListenAndServe(":1234", nil))
}
