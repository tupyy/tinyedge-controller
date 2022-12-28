package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	gitRepo "github.com/tupyy/tinyedge-controller/internal/repo/git"
)

func main() {
	tmpDir, _ := os.MkdirTemp("/home/cosmin/tmp", "git-")
	repo := entity.Repository{
		Url:       "/home/cosmin/tmp/manifestwork",
		LocalPath: "/home/cosmin/tmp/git/manifest",
	}

	g := gitRepo.New(tmpDir)
	newRepo, err := g.Open(context.TODO(), repo)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("repo: %+v", newRepo)

	for {
		err := g.Pull(context.TODO(), newRepo)
		log.Printf("error: %v", err)
		<-time.After(5 * time.Second)
	}

}
