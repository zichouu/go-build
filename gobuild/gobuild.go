package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zichouu/go-pkg/exe"
	"github.com/zichouu/go-pkg/file"
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
			buildDir := "build"
			link := "_"
			osArch := fmt.Sprintf("%v%v%v", v.GOOS, link, v.GOARCH)
			o := fmt.Sprintf("%v/%v/", buildDir, osArch)
			aenv := []string{envGOOS, envGOARCH, envCGO}
			arg := []string{"go", "build", "-trimpath", "-ldflags", `"-s -w"`, "-o", o, "./..."}
			_, err := exe.Run(dir, aenv, arg...)
			if err != nil {
				return err
			}
			files, err := os.ReadDir(filepath.Join(dir, o))
			if err != nil {
				return err
			}
			for _, v := range files {
				name := v.Name()
				ext := ""
				if strings.HasSuffix(v.Name(), ".exe") {
					name = strings.Split(v.Name(), ".exe")[0]
					ext = ".exe"
				}
				src := filepath.Join(dir, buildDir, osArch, v.Name())
				dst := filepath.Join(dir, buildDir, fmt.Sprintf("%v%v%v%v", name, link, osArch, ext))
				err = file.Copy(src, dst)
				if err != nil {
					return err
				}
			}
			return err
		})
	}
	err := g.Wait()
	return err
}
