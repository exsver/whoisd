// Copyright 2017 Openprovider Authors. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package client

import (
	"bytes"
	"log"
	"net"
	"os"
        "os/exec"
        "regexp"

        "github.com/openprovider/whoisd/pkg/config"
	"github.com/openprovider/whoisd/pkg/storage"
	"golang.org/x/net/idna"
)

const (
	queryBufferSize = 256
)

// Record - standard record (struct) for client package
type Record struct {
	Conn  net.Conn
	Query []byte
}

// simplest logger, which initialized during starts of the application
var (
	stdlog = log.New(os.Stdout, "[CLIENT]: ", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "[CLIENT:ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// HandleClient - Sends a client data into the channel
func (client *Record) HandleClient(channel chan<- Record) {
	defer func() {
		if recovery := recover(); recovery != nil {
			errlog.Println("Recovered in HandleClient:", recovery)
			channel <- *client
		}
	}()
	buffer := make([]byte, queryBufferSize)
	numBytes, err := client.Conn.Read(buffer)
	if numBytes == 0 || err != nil {
		return
	}
	client.Query = bytes.ToLower(bytes.Trim(buffer, "\u0000\u000a\u000d"))
	channel <- *client
}

// ProcessClient - Asynchronous a client handling
func ProcessClient(channel <-chan Record, repository *storage.Record, conf *config.Record) {
	message := Record{}
	defer func() {
		if recovery := recover(); recovery != nil {
			errlog.Println("Recovered in ProcessClient:", recovery)
			if message.Conn != nil {
				message.Conn.Close()
			}
		}
	}()
	for {
		message = <-channel
		query, err := idna.ToASCII(string(message.Query))
		if err != nil {
			query = string(message.Query)
		}

                primaryData, primaryStatus := repository.Search(query)
                if conf.SecondaryWhois == "" || primaryStatus == true {
                        stdlog.Println(query, "exists in primary server", primaryStatus)
                        message.Conn.Write([]byte(primaryData))
                        message.Conn.Close()
                        continue;
                }

                secondaryData, err := exec.Command("whois", "-h", conf.SecondaryWhois, query).Output()
                if err != nil {
                        stdlog.Println(query, "error when call secondary server", err)
                        message.Conn.Write([]byte(primaryData))
                        message.Conn.Close()
                        continue;
                }
                matched, err := regexp.Match(conf.SecondaryRegexp, secondaryData)
                if err != nil {
                        stdlog.Println(query, "error on analyzing secondary response", err)
                        message.Conn.Write([]byte(primaryData))
                        message.Conn.Close()
                        continue;
                }
                if matched {
                        stdlog.Println(query, "NOT exists in secondary server", err)
                        message.Conn.Write([]byte(primaryData))
                        message.Conn.Close()
                        continue;
                }
                stdlog.Println(query, "exists in secondary server", err)
		message.Conn.Write(secondaryData)
		message.Conn.Close()
	}
}
