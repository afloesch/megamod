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
	ver := swizzle.SemVer("v0.1.1")
	mod, err := repo.FetchManifest(ctx, ver)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("mod:", mod)

	if len(mod.Files) > 0 {
		fmt.Println("downloading mod:", mod.Name)
		err = mod.DownloadReleaseFiles(ctx, path)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("mod files downloaded to ./tmp")

	err = mod.AddDependency(ctx, swizzle.Repo("afloesch/sse-skse"), swizzle.SemVer("v2.0.20"))
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
