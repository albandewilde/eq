package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// TKN is the discord token
	TKN string
	// SRCDIR is the directory to download images before checking if we already got it
	SRCDIR string
	// DSTDIR is the directory where images are saved
	DSTDIR string
)

func init() {
	// Check if we have a discord token
	TKN = os.Getenv("TKN")
	if TKN == "" {
		log.Fatal("No discord token provided in the environment variable `TKN`")
	}

	// Check if directories are given as environment variable
	SRCDIR = os.Getenv("SRC_DIR")
	if SRCDIR == "" {
		log.Fatal("No source directory provided in the environment variable `SRC_DIR`")
	}
	DSTDIR = os.Getenv("DST_DIR")
	if DSTDIR == "" {
		log.Fatal("No destination directory provided in the environment variable `DST_DIR`")
	}
}

func main() {
	// Check if directories exist
	_, err := os.Stat(SRCDIR)
	if os.IsNotExist(err) {
		log.Fatalf("Directory `%s` doesn't exist.", SRCDIR)
	}
	_, err = os.Stat(DSTDIR)
	if os.IsNotExist(err) {
		log.Fatalf("Directory `%s` doesn't exist.", DSTDIR)
	}

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
