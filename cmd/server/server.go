package main

import (
	"os"

	"github.com/seal-io/seal/pkg/server"
	"github.com/seal-io/seal/utils/clis"
	"github.com/seal-io/seal/utils/log"
)

func main() {
	cmd := server.Command()
	app := clis.AsApp(cmd)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
