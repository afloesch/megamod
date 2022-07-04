package main

import (
	"context"
	"fmt"
	"math"

	"github.com/afloesch/megamod/swizzle"
)

func dlModfiles(ctx context.Context, mod *swizzle.Manifest, path string) error {
	done, prog, err := mod.DownloadReleaseFile(ctx, mod.Files[0], path)

	var p float64
	var e error
	for {
		select {
		case <-done:
			fmt.Println()
			return nil
		case p = <-prog:
			fmt.Printf(fmt.Sprintf("\r%v percent of %v", math.Floor(p), mod.Files[0].Size()))
		case e = <-err:
			return e
		}
	}
}

func main() {
	path := "./tmp/archive"
	ctx := context.Background()
	repo := swizzle.Repo("afloesch/megamod")
	rel, err := repo.Release(ctx, "v0.1.1")
	if err != nil {
		fmt.Println(err)
		return
	}

	mod, err := repo.Manifest(ctx, rel)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("mod:", mod)

	if len(mod.Files) > 0 {
		fmt.Println("downloading mod:", mod.Name)
		err := dlModfiles(ctx, mod, path)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("mod files downloaded to ./tmp")

	err = mod.AddDependency(ctx, "afloesch/sse-skse", ">=v2.0.20")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Mod deps", mod.Dependency)
	err = mod.WriteFile("./tmp/test.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	newmod, err := swizzle.New().ReadFile("./swiz.zle")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("mod:", newmod)
}
