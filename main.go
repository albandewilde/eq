package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// SRCDIR is the directory to download images before checking if we already got it
	SRCDIR = "./dl/"
	// DTSDIR is the directory where images are saved
	DSTDIR = "./img/"
)

var (
	// TKN is the discord token
	TKN string
)

func init() {
	TKN = os.Getenv("tkn")
	if TKN == "" {
		log.Fatal("No discord token provided in the environment variable `tkn`")
	}
}

func main() {
	// Start a goroutine that watch files in the `SRCDIR` then put them in the `DSTDIR` if there not already present
	go WatchFiles(SRCDIR, DSTDIR, time.Second*3)

	// Start a discord bot that download files in the `SRCDIR`
	bot, err := NewBot(TKN, SRCDIR)
	if err != nil {
		log.Fatal(err)
	}

	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}
