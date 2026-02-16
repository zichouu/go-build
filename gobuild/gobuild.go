package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/zichouu/go-pkg/exe"
	"golang.org/x/sync/errgroup"
)

func main() {
	_, err := exec.LookPath("go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buildDirList := []string{"."}
	if len(os.Args) >= 2 {
		buildDirList = os.Args[1:]
	}

	var g errgroup.Group
	for _, v := range buildDirList {
		g.Go(func() error {
			err := Build(v)
			return err
		})
	}
	g.Wait()
}

func Build(dir string) error {
	var g errgroup.Group
	for _, v := range Ver {
		g.Go(func() error {
			envGOOS := fmt.Sprintf("GOOS=%v", v.GOOS)
			envGOARCH := fmt.Sprintf("GOARCH=%v", v.GOARCH)
			o := fmt.Sprintf("build/%v-%v/", v.GOOS, v.GOARCH)
			arg := []string{"go", "build", "-trimpath", "-ldflags", "-s -w", "-o", o, "./..."}
			err := exe.Run(dir, []string{envGOOS, envGOARCH, "CGO_ENABLED=0"}, arg...)
			return err
		})
	}
	err := g.Wait()
	return err
}
