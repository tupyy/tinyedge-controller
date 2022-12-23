package main

import (
	"context"
	"log"
	"os"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	gitRepo "github.com/tupyy/tinyedge-controller/internal/repo/git"
)

func main() {
	tmpDir, _ := os.MkdirTemp("/home/cosmin/tmp", "git-")
	repo := entity.Repository{
		Url: "/home/cosmin/tmp/manifestwork",
	}

	g := gitRepo.New(tmpDir)
	newRepo, err := g.Open(context.TODO(), repo)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("repo: %+v", newRepo)

	manifestWorks, err := g.GetManifests(context.TODO(), newRepo)
	for _, m := range manifestWorks {
		log.Printf("---------")
		log.Printf("manifest works: %+v", m)
	}
}
