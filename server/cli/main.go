package main

import (
	"context"
	gwProto "github.com/eelf/gitweb/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	c, err := grpc.Dial("localhost:2004", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	cli := gwProto.NewGitwebClient(c)
	resp, err := cli.RepoList(context.Background(), &gwProto.RepoListRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
