package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

// NewBot return a discord session that download images in the `dlDir`
func NewBot(tkn, dlDir string) (*dgo.Session, error) {
	bot, err := dgo.New("Bot " + tkn)
	if err != nil {

		return nil, err
	}

	dlCallback := configureDlCallback(dlDir)
	// Register callback functions
	bot.AddHandler(dlCallback)

	return bot, nil
}

func configureDlCallback(dlDir string) func(*dgo.Session, *dgo.MessageCreate) {
	return func(s *dgo.Session, m *dgo.MessageCreate) {
		for _, attachemnts := range m.Attachments {
			err := dlFile(attachemnts.URL, dlDir, uuid.NewString())
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func dlFile(URL, dir, fileName string) error {
	// Fetch the file
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the body request
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Write the file to the disk
	err = ioutil.WriteFile(filepath.Join(dir, fileName), content, 0644)
	if err != nil {
		return err
	}
	return nil
}
