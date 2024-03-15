package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/klauspost/compress/gzip"
	"superfrink.net/tetris/streamer"
)

type GameRecorder struct {
	mu       sync.Mutex
	messages []streamer.Message
}

func (r *GameRecorder) QueueMove(message streamer.Message) {
	r.mu.Lock()
	r.messages = append(r.messages, message)
	r.mu.Unlock()
}

func (r *GameRecorder) RecordMoveFile() {
	r.mu.Lock()

	count := len(r.messages)
	slog.Info("messages", "count", count)

	if count > 1 {

		// first open the output file
		filename := fmt.Sprintf("move_records_%d.dat", time.Time.Unix(time.Now()))
		slog.Info("creeating file", "filename", filename)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			slog.Error("creating file", "filename", filename, "err", fmt.Sprintf("%s", err))
			goto endOfFunction
		}
		defer file.Close()

		// second use a gzip writer
		fileGzip := gzip.NewWriter(file)
		defer fileGzip.Close()

		// third encode the messages
		// DOC: gob doesn't like the Game.PRNG // encoder := gob.NewEncoder(fileGzip)
		encoder := json.NewEncoder(fileGzip)

		for _, message := range r.messages {
			slog.Info("message", "move", message.Move)
			err := encoder.Encode(message)
			if err != nil {
				slog.Error("encoding message", "err", fmt.Sprintf("%s", err))
			}
		}
	}

endOfFunction:
	r.mu.Unlock()
}

func dumpFile(filename string) error {

	// first open the file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening file %q %w", filename, err)
	}
	defer file.Close()

	// second gunzip the file
	fileGzip, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("gunziping file %q %w", filename, err)
	}
	defer fileGzip.Close()

	// third decode the messages
	decoder := json.NewDecoder(fileGzip)

	for {
		var message streamer.Message
		err := decoder.Decode(&message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decoding message from file %q %w", filename, err)
		}
		slog.Info("decoded", "message", message)
	}

	return nil
}

func main() {

	var natsUrl = flag.String("u", "", "NATS URL")
	var natsCredFile = flag.String("c", "", "NATS credential file")
	var recordGame = flag.Bool("r", false, "Record game moves from stream")
	var dumpRecordingFile = flag.String("d", "", "Dump recorded game")
	flag.Parse()

	if (!*recordGame) && (*dumpRecordingFile == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *recordGame && (*natsUrl == "" || *natsCredFile == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *dumpRecordingFile != "" && (*natsUrl != "" || *natsCredFile != "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *dumpRecordingFile != "" {
		err := dumpFile(*dumpRecordingFile)
		if err != nil {
			slog.Error("dumping file", "filename", *dumpRecordingFile, "err", err)
		}
		os.Exit(0)
	}

	// GOAL: setup the game recorder

	var recorder = GameRecorder{
		messages: make([]streamer.Message, 0),
	}

	// GOAL: connect to NATS and record the game moves

	stream := streamer.Streamer{}
	stream.Connect(*natsUrl, *natsCredFile, "fixme")
	streamMesgChan := stream.RecvChan()
	slog.Info("Connected to NATS", "url", *natsUrl)

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case message := <-streamMesgChan:
			slog.Debug("message", "message", fmt.Sprintf("%+v", message))
			if message.Type == streamer.Move {
				recorder.QueueMove(message)
			}
		case <-ticker.C:
			slog.Info("tick")
			recorder.RecordMoveFile()
		}
	}
}
