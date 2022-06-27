package main

import (
	"context"
	"fmt"

	"github.com/afloesch/megamod/internal"
)

func main() {
	ctx := context.Background()
	test, err := internal.FetchManifest(ctx, "afloesch/megamod/main/.consoleUtilSSE.mm.yml")
	if err != nil {
		fmt.Println(fmt.Errorf("parse error: %s", err))
		return
	}

	fmt.Println("mod file:", test.URL)

	if len(test.Dependency) > 0 {
		err = test.FlattenDeps(ctx)
		if err != nil {
			fmt.Println(fmt.Errorf("parse error: %s", err))
			return
		}

		for _, v := range test.AllDependencies() {
			fmt.Println("dependency:", v.URL)
		}

		err = test.FetchDeps(ctx, "./tmp")
		if err != nil {
			fmt.Println(fmt.Errorf("fetch deps error: %s", err))
			return
		}
	}

	if test.URL != "" {
		err = test.FetchRelease(ctx, "./tmp")
		if err != nil {
			fmt.Println(fmt.Errorf("fetch error: %s", err))
			return
		}
	}

	fmt.Println("mod files downloaded to ./tmp")
}
