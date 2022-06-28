package main

import (
	"context"
	"fmt"

	"github.com/afloesch/megamod/internal"
)

func main() {
	path := "./tmp"
	ctx := context.Background()
	test, err := internal.FetchManifest(ctx, "afloesch/megamod/main/testlist.mm.yml")
	if err != nil {
		fmt.Println(fmt.Errorf("parse error: %s", err))
		return
	}

	fmt.Println("mod:", test.Name)

	if len(test.Dependency) > 0 {
		deps, err := test.FetchDepManifests(ctx)
		if err != nil {
			fmt.Println(fmt.Errorf("parse error: %s", err))
			return
		}

		for _, d := range deps {
			fmt.Println("mod dep:", d.Name)
		}

		for _, d := range deps {
			fmt.Println("downloading mod:", d.Name)
			err = d.FetchRelease(ctx, path)
			if err != nil {
				fmt.Println(fmt.Errorf("fetch deps error: %s", err))
				return
			}
		}
	}

	if test.URL != "" {
		fmt.Println("downloading mod:", test.Name)
		err = test.FetchRelease(ctx, path)
		if err != nil {
			fmt.Println(fmt.Errorf("fetch error: %s", err))
			return
		}
	}

	fmt.Println("mod files downloaded to ./tmp")
}
