package main

import (
	"os"

	"github.com/seal-io/walrus/pkg/server"
	"github.com/seal-io/walrus/utils/clis"
	"github.com/seal-io/walrus/utils/log"
)

func main() {
	cmd := server.Command()

	app := clis.AsApp(cmd)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
