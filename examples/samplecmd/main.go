package main

import (
	"fmt"

	"github.com/rmatsuoka/env"
	"github.com/rmatsuoka/env/examples/samplecmd/internal/pkg1"
)

var str = env.String("MAIN_CONFIG", "default.main", "main config for main")

func main() {
	env.Parse()
	fmt.Println(*str)
	fmt.Println(*pkg1.String)
}
