package streamer

import (
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"superfrink.net/tetris/engine"
)

type MessageType string

const (
	StateUpdate = "StateUpdate"
	Move        = "Move"
)

// Message FIXME
type Message struct {
	Type MessageType
	Game engine.Game
	Move byte
}

// Streamer FIXME
type Streamer struct {
	nc       *nats.Conn
	ec       *nats.EncodedConn
	chanName string
	recvCh   chan Message
	sendCh   chan Message
}

// Connect FIXME
func (s *Streamer) Connect(url string, credFile string, streamName string) error {
	var err error

	s.chanName = streamName

	if s.chanName == "" {
		err = errors.New("invalid channel name")
		return fmt.Errorf("%q %w", s.chanName, err)
	}
	s.nc, err = nats.Connect(url, nats.UserCredentials(credFile))
	if err != nil {
		return fmt.Errorf("Connecting to NATS: %s, %w", url, err)
	}

	s.ec, err = nats.NewEncodedConn(s.nc, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("Creating JSON encoded NATS channel: %s, %w", url, err)
	}

	s.recvCh = make(chan Message)
	s.ec.BindRecvChan(s.chanName, s.recvCh)

	s.sendCh = make(chan Message)
	s.ec.BindSendChan(s.chanName, s.sendCh)

	return nil
}

// RecvChan FIXME
func (s *Streamer) RecvChan() chan Message {
	return s.recvCh
}

// SendGameState FIXME
func (s *Streamer) SendGameState(game engine.Game) {
	m := Message{
		Type: StateUpdate,
		Game: game,
	}

	s.sendCh <- m
}

// SendMove FIXME
func (s *Streamer) SendMove(move byte, game engine.Game) {

	m := Message{
		Type: Move,
		Game: game,
		Move: move,
	}

	s.sendCh <- m
}
