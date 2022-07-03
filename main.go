package main

import (
	"context"
	"fmt"

	"github.com/afloesch/megamod/swizzle"
)

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
		err = mod.DownloadReleaseFile(ctx, mod.Files[0], path)
		if err != nil {
			fmt.Println(err)
			return
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
