package main

import ( //imports for all the supporting modules and what looks like probably GO specefic functions
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Relative
	"github.com/pterm/pterm"
	"github.com/taubyte/dreamland/cli/common"
	inject "github.com/taubyte/dreamland/cli/inject"
	"github.com/taubyte/dreamland/cli/kill"
	"github.com/taubyte/dreamland/cli/new"
	"github.com/taubyte/dreamland/cli/status"

	// Actual imports

	client "github.com/taubyte/dreamland/service"
	"github.com/taubyte/tau/libdream/services"
	"github.com/urfave/cli/v2"

	// Empty imports for initializing fixtures, and client/service run methods"
	_ "github.com/taubyte/tau/clients/p2p/auth"
	_ "github.com/taubyte/tau/clients/p2p/hoarder"
	_ "github.com/taubyte/tau/clients/p2p/monkey"
	_ "github.com/taubyte/tau/clients/p2p/patrick"
	_ "github.com/taubyte/tau/clients/p2p/seer"
	_ "github.com/taubyte/tau/clients/p2p/tns"
	_ "github.com/taubyte/tau/libdream/common/fixtures"
	_ "github.com/taubyte/tau/protocols/auth"
	_ "github.com/taubyte/tau/protocols/hoarder"
	_ "github.com/taubyte/tau/protocols/monkey"
	_ "github.com/taubyte/tau/protocols/monkey/fixtures/compile"
	_ "github.com/taubyte/tau/protocols/patrick"
	_ "github.com/taubyte/tau/protocols/seer"
	_ "github.com/taubyte/tau/protocols/substrate"
	_ "github.com/taubyte/tau/protocols/tns"
)

func main() { //main function that is essentially 'the program', creates the instance and runs the 'loop' of the program. Probably somewhat wrong in that guess but something similar is happening here. Seems to be either specefic to how GO works or a convention of some sort since main.go is a common file sort of like a index.html or a index.js
	ctx, ctxC := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {//function to catch a cancel or shutdown input
		sig := <-signals
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			pterm.Info.Println("Received signal... Shutting down.")
			ctxC()
		}
	}()

	defer func() {//if the variable for 'DoDaemon' is true (default false) then cancel. Zeno seems to mass close all instances
		if common.DoDaemon {
			ctxC()
			services.Zeno()
		}
	}()

	ops := []client.Option{client.URL(common.DefaultDreamlandURL), client.Timeout(300 * time.Second)} //set some options
	multiverse, err := client.New(ctx, ops...)
	if err != nil {//error catches for crash reporting
		log.Fatalf("Starting new dreamland client failed with: %s", err.Error())
	}

	err = defineCLI(&common.Context{Ctx: ctx, Multiverse: multiverse}).RunContext(ctx, os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func defineCLI(ctx *common.Context) *(cli.App) {//startup the CLI and set some options
	app := &cli.App{
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			new.Command(ctx),
			inject.Command(ctx),
			kill.Command(ctx),
			status.Command(ctx),
		},
		Suggest:              true,
		EnableBashCompletion: true,
	}

	return app
}
