package main

import (
	"context"
	"fmt"

	"github.com/afloesch/megamod/internal"
)

func main() {
	path := "./tmp/archive"
	ctx := context.Background()
	repo := internal.Repo("afloesch/megamod")
	mod, err := repo.FetchManifest(ctx, internal.SemVer("v0.1.1"))
	if err != nil {
		fmt.Println(fmt.Errorf("parse error: %s", err))
		return
	}

	fmt.Println("mod:", mod)

	if len(mod.Files) > 0 {
		fmt.Println("downloading mod:", mod.Name)
		err = mod.DownloadReleaseFiles(ctx, path)
		if err != nil {
			fmt.Println(fmt.Errorf("fetch error: %s", err))
			return
		}
	}

	fmt.Println("mod files downloaded to ./tmp")
}
