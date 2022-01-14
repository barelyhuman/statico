package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// handle flags
	// - watch
	// - config
	// - init
	// - serve
	appRef := &Statico{}
	initApp := flag.Bool("init", false, "Initialize the most minimal version of statico")
	enableWatch := flag.Bool("watch", false, "Start statico in watch mode")
	enableWatchAlias := flag.Bool("w", false, "alias `-watch`")
	enableServe := flag.Bool("serve", false, "Enable file server")
	enableServeAlias := flag.Bool("s", false, "alias `-serve`")
	configFileFlag := flag.String("config", "", "Config file to use")
	configFileFlagAlias := flag.String("c", "", "alias `-config`")

	flag.Parse()

	if *initApp {
		runAppInitScript()
		return
	}

	appRef.config = readConfigFlags(*configFileFlag, *configFileFlagAlias)

	appRef.Build()

	if *enableWatch || *enableWatchAlias || *enableServe || *enableServeAlias {
		waitForKill := make(chan int)

		if *enableWatch || *enableWatchAlias {
			go appRef.WatchFiles()
		}

		if *enableServe || *enableServeAlias {
			go appRef.ServeFiles()
		}

		<-waitForKill
	}

	// collect all files to read and their data
	// go through each file one by one to start mapping them
	// add these files to a cache so it's easier to update singular files

	// things needed

	// site meta, can be taken from config
	// post meta, will be dynamic and left for the template creator to decide
	// post fixed meta will be `title`,`published` and `date` since we will be
	// tranferring these with needed functions

	// AtomFeed Generation

	fmt.Printf(
		Dim(logTime()) + Success("Compiled! Static files generated to: "+Bullet(appRef.config.OutPath)),
	)
}

func bail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
