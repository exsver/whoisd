// Copyright 2017 Openprovider Authors. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/openprovider/whoisd/pkg/config"
	"github.com/openprovider/whoisd/pkg/service"
	"github.com/openprovider/whoisd/pkg/version"
)

var (
	stdlog, errlog *log.Logger
)

// Init "Usage" helper
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	stdlog = log.New(os.Stdout, "[SETUP]: ", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "[SETUP:ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
	flag.Usage = func() {
		fmt.Println(config.Usage())
	}
}

func main() {
	daemonName, daemonDescription := "whoisd", "Whois Daemon"
	daemon, err := service.New(daemonName, daemonDescription)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	flag.Parse()
	if daemon.Config.ShowVersion {
		buildTime, err := time.Parse(time.RFC3339, version.DATE)
		if err != nil {
			buildTime = time.Now()
		}
		fmt.Println(daemonName, version.RELEASE, buildTime.Format(time.RFC1123), version.REPO, version.COMMIT)
		os.Exit(0)
	}
	status, err := daemon.Run()
	if err != nil {
		errlog.Printf("%s - %s", status, err)
		os.Exit(1)
	}
	fmt.Println(status)
}
