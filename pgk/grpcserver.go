package main

import (
	"context"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strconv"
	"testtask2/pgk/api"
)

type GRPCServer struct {
}

var (
	maxResults int64 = 50
)

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
func playlistsItems(service *youtube.Service, part []string, maxResults int64, pageToken string, playlistId string) *youtube.PlaylistItemListResponse {
	call := service.PlaylistItems.List(part)
	call = call.MaxResults(maxResults)
	call = call.PlaylistId(playlistId)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	response, err := call.Do()
	handleError(err, "")
	return response
}

func (s *GRPCServer) Check(ctx context.Context, request *api.CheckRequest) (*api.CheckResponse, error) {
	list := ""
	client := &http.Client{Transport: &transport.APIKey{Key: request.GetKey()}}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("when creating Youtube client: %v", err)
	}

	listSnippet := []string{"snippet"}
	pageToken := ""
	i := 0
	for {
		playlistResponse := playlistsItems(service, listSnippet, maxResults, pageToken, request.GetIdPlay())
		for _, playlistItem := range playlistResponse.Items {
			i++
			list += (strconv.Itoa(i) + ": ")
			list += (playlistItem.Snippet.Title + "\n")
		}
		pageToken = playlistResponse.NextPageToken
		if pageToken == "" {
			break
		}
	}
	list += ("Итого: " + strconv.Itoa(i))

	return &api.CheckResponse{List: list}, nil
}

var port = "8081"

func main() {
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...) //сервер
	srvst := &GRPCServer{}
	api.RegisterCheckYoutubeServer(s, srvst)
	lisener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {

		log.Fatal(err)
	}
	if err := s.Serve(lisener); err != nil {
		log.Fatalln(err)
	}
}
