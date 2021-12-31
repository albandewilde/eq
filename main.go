package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
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

	// HOST is the host of the static server
	HOST string
	// PORT is the port where the static server listen
	PORT int64
	// BASEURL of the static server
	BASEURL string
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

	// Check if the port and host are given
	HOST = os.Getenv("HOST")
	if HOST == "" {
		log.Fatal("No Host provided in the environment variable `HOST`")
	}
	port := os.Getenv("PORT")
	var err error
	PORT, err = strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	BASEURL = os.Getenv("BASEURL")
	if BASEURL == "" {
		BASEURL = "/"
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
	go WatchFiles(SRCDIR, DSTDIR, time.Hour*3)

	// Start a discord bot that download files in the `SRCDIR`
	bot, err := NewBot(TKN, SRCDIR)
	if err != nil {
		log.Fatal(err)
	}
	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Bot is now running.  Press CTRL-C to exit.")

	// Create the server to serve files in the `DSTDIR` directory
	srv := CreateFileServer(HOST, PORT, BASEURL, DSTDIR)

	// Gracefully close the discord bot and http server
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	go func() {
		<-sc
		bot.Close()
		srv.Shutdown(context.Background())
	}()

	// Start static file server
	log.Printf("Server running on %s:%d%s\n", HOST, PORT, BASEURL)
	srv.ListenAndServe()
}
