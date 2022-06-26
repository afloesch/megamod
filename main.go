package main

import (
	"fmt"

	"github.com/afloesch/megamod/internal"
)

func main() {
	test, err := internal.FetchManifest("afloesch/megamod/main/.consoleUtilSSE.mm.yml")
	if err != nil {
		fmt.Println(fmt.Errorf("parse error: %s", err))
		return
	}

	fmt.Println("manifest", test)

	v1 := internal.SemVersion(">=v2.0.1")
	v2 := internal.SemVersion("<=v1.6.22-1-dirty")
	v3 := internal.SemVersion(">=v1.5.23")

	fmt.Println(v1.Compare(v2))
	fmt.Println(v2.OpCompare(v3))
}
