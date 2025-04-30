package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"log"
	"net/http"
)

const (
	baseUrl       = "http://localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/%d"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	ID       int64    `json:"id"`
	Info     NoteInfo `json:"info"`
	CreateAt string   `json:"create_at"`
	UpdateAt string   `json:"update_at"`
}

func createNoteClient() (Note, error) {
	note := NoteInfo{
		Title:    gofakeit.BeerName(),
		Context:  gofakeit.IPv4Address(),
		Author:   gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}

	data, err := json.Marshal(note)
	if err != nil {
		return Note{}, err
	}

	resp, err := http.Post(baseUrl+createPostfix, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Note{}, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf(err.Error())
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		return Note{}, err
	}

	var createdNote Note
	if err = json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
		return Note{}, err
	}

	return createdNote, nil

}

func getNoteClient(id int64) (Note, error) {
	resp, err := http.Get(fmt.Sprintf(baseUrl+getPostfix, id))
	if err != nil {
		log.Fatal("Failed to get note", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf(err.Error())
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return Note{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Note{}, fmt.Errorf("failed to get note: %d", resp.StatusCode)
	}

	var note Note

	if err = json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return Note{}, err
	}
	return note, nil

}

func main() {

	note, err := createNoteClient()
	if err != nil {
		log.Fatal("failed to create note:", err)
	}

	log.Println(color.RedString("Note created:"))
	log.Println(color.GreenString("%+v", note))

	note, err = getNoteClient(note.ID)
	if err != nil {
		log.Fatal("failed to get note:", err)
	}

	log.Printf(color.GreenString("Note created:\n"), color.GreenString("%+v", note))
}
