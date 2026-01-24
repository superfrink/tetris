package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/klauspost/compress/gzip"
	"superfrink.net/tetris/streamer"
)

type MoveRecord struct {
	Field [][]int
	Move  uint8
	Piece uint8
	PCol  uint8
	PRot  uint8
	PRow  uint8
}

func MessageToMoveRecord(message streamer.Message) MoveRecord {
	move := MoveRecord{
		Field: message.Game.Field,
		Move:  message.Move,
		Piece: uint8(message.Game.Piece),
		PCol:  uint8(message.Game.PiecePosCol),
		PRot:  uint8(message.Game.PieceRotation),
		PRow:  uint8(message.Game.PiecePosRow),
	}
	return move
}

type GameRecorder struct {
	mu    sync.Mutex
	moves []MoveRecord
}

func (r *GameRecorder) QueueMove(rec MoveRecord) {
	r.mu.Lock()
	r.moves = append(r.moves, rec)
	r.mu.Unlock()
}

func (r *GameRecorder) RecordToFile() {
	r.mu.Lock()

	count := len(r.moves)
	slog.Info("messages to write to file", "count", count)

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

		for _, message := range r.moves {
			slog.Debug("message", "move", message.Move)
			err := encoder.Encode(message)
			if err != nil {
				slog.Error("encoding message", "err", fmt.Sprintf("%s", err))
			}
		}

		r.moves = r.moves[:0]
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

	// third decode the record
	decoder := json.NewDecoder(fileGzip)

	for {
		var rec MoveRecord
		err := decoder.Decode(&rec)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decoding move record from file %q %w", filename, err)
		}
		slog.Info("decoded", "rec", rec)
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
		moves: make([]MoveRecord, 0),
	}

	// GOAL: connect to NATS

	stream := streamer.Streamer{}
	stream.Connect(*natsUrl, *natsCredFile, "fixme")
	streamMesgChan := stream.RecvChan()
	slog.Info("Connected to NATS", "url", *natsUrl)

	// GOAL: write the game moves to file periodically

	go func(rec *GameRecorder) {
		ticker := time.NewTicker(600 * time.Second)

		for {
			<-ticker.C
			slog.Info("tick")

			recorder.RecordToFile()
		}
	}(&recorder)

	// GOAL: write the game moves to file on signals USR1, HUP, INT, TERM

	go func(rec *GameRecorder) {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
		for {
			sig := <-sigChan
			slog.Info("signal arrived", "sig", sig.String())
			rec.RecordToFile()
			if sig == syscall.SIGTERM || sig == syscall.SIGINT {
				os.Exit(0)
			}
		}
	}(&recorder)

	// GOAL: receive and record messages

	for {
		message := <-streamMesgChan
		slog.Debug("message", "message", fmt.Sprintf("%+v", message))

		if message.Type == streamer.Move {
			record := MessageToMoveRecord(message)
			recorder.QueueMove(record)
		}
	}
}
