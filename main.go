//
// Copyright 2016 Gregory Trubetskoy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Tgres is a tool for receiving and reporting on simple time series
// written in Go which uses PostgreSQL for storage.
package main

import (
	"flag"
	"github.com/tgres/tgres/daemon"
)

func parseFlags() (textCfgPath, gracefulProtos, join string) {

	// Parse the flags, if any
	flag.StringVar(&textCfgPath, "c", "./etc/tgres.conf", "path to config file")
	flag.StringVar(&join, "join", "", "List of add:port,addr:port,... of nodes to join")
	flag.StringVar(&gracefulProtos, "graceful", "", "list of fds (internal use only)")
	flag.Parse()

	return
}

func main() {
	if cfg := daemon.Init(parseFlags()); cfg != nil {
		daemon.Finish(cfg)
	}
}
