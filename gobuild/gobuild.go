package main

import (
	"fmt"
	"log/slog"
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
	err = g.Wait()
	if err != nil {
		slog.Error("Build err", "err", err)
	}
}

func Build(dir string) error {
	var g errgroup.Group
	for _, v := range Ver {
		g.Go(func() error {
			envCGO := "CGO_ENABLED=0"
			envGOOS := fmt.Sprintf("GOOS=%v", v.GOOS)
			envGOARCH := fmt.Sprintf("GOARCH=%v", v.GOARCH)
			o := fmt.Sprintf("build/%v-%v/", v.GOOS, v.GOARCH)
			aenv := []string{envGOOS, envGOARCH, envCGO}
			arg := []string{"go", "build", "-trimpath", "-ldflags", `"-s -w"`, "-o", o, "./..."}
			err := exe.Run(dir, aenv, arg...)
			return err
		})
	}
	err := g.Wait()
	return err
}
