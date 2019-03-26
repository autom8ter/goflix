// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"errors"
	"github.com/autom8ter/goflix/client"
	"github.com/autom8ter/goflix/player"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var cfg = &client.ClientConfig{}
var play string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start goflix server",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg.TorrentPath == "" {
			log.Fatalln(rootCmd.UsageString(), errors.New("please provide a valid torrent --torrent | -t"))
		}
		client, err := client.New(cfg)
		if err != nil {
			log.Fatalln(err.Error())
		}
		// Http handler.
		go func() {
			http.HandleFunc("/", client.GetFile)
			log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil))
		}()
		if play != "" {
			go func() {
				for !client.ReadyForPlayback() {
					time.Sleep(time.Second)
				}
				player.OpenPlayer(play, cfg.Port)
			}()
		}
		// Handle exit signals.
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel,
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		go func(interruptChannel chan os.Signal) {
			for range interruptChannel {
				log.Println("Exiting...")
				client.Close()
				os.Exit(0)
			}
		}(interruptChannel)

		// Cli render loop.
		for {
			client.Render()
			time.Sleep(time.Second)
		}
	},
}

func init() {
	serveCmd.PersistentFlags().StringVarP(&cfg.TorrentPath, "torrent", "t", "", "torrent to download")
	serveCmd.PersistentFlags().IntVar(&cfg.Port, "port", 8080, "port to serve torrent")
	serveCmd.PersistentFlags().IntVar(&cfg.TorrentPort, "torrent-port", 50007, "port to torrent from")
	serveCmd.PersistentFlags().IntVar(&cfg.MaxConnections, "max-conn", 200, "max connections")
	serveCmd.PersistentFlags().BoolVar(&cfg.Seed, "seed", false, "to seed or not to seed")
	serveCmd.PersistentFlags().BoolVar(&cfg.TCP, "tcp", true, "to tcp or not to tcp")
	serveCmd.PersistentFlags().StringVarP(&play, "player", "p", "", "Open the stream with a video player ("+player.JoinPlayerNames()+")")
	rootCmd.AddCommand(serveCmd)
}
