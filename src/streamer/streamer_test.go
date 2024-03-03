package streamer

import (
	"log"
	"testing"

	"superfrink.net/tetris/engine"
)

var testURL = "connect.ngs.global"
var testCredFile = "../NGS-Default-CLI.creds"

func TestConnect(t *testing.T) {
	s := Streamer{
		chanName: "fixme" + t.Name(),
	}

	err := s.Connect(testURL, testCredFile)
	if err != nil {
		t.Errorf("Connect failed. %q", err)
	}
}

func TestConnectBadURL(t *testing.T) {
	s := Streamer{
		chanName: "fixme" + t.Name(),
	}

	err := s.Connect("badurl", testCredFile)
	if err == nil {
		t.Errorf("Expected connect to fail.")
	}
}

func TestConnectBadCredFile(t *testing.T) {
	s := Streamer{
		chanName: "fixme" + t.Name(),
	}

	err := s.Connect(testURL, "badcredfilename")
	if err == nil {
		t.Errorf("Expected connect to fail.")
	}
}

func TestInvalidChannel(t *testing.T) {
	s := Streamer{}

	err := s.Connect(testURL, testCredFile)
	if err == nil {
		t.Errorf("Expected connect to fail.")
	}
}

func TestSendGameState(t *testing.T) {
	s := Streamer{
		chanName: "fixme" + t.Name(),
	}

	err := s.Connect(testURL, testCredFile)
	if err != nil {
		t.Errorf("Connect failed. %q", err)
		return
	}

	_, gameInputCh, gameUpdateCh := engine.NewBucketGame()

	gameInputCh <- engine.PlayInputRotate
	game := <-gameUpdateCh

	go func() {
		s.SendGameState(*game)
	}()

	recvCh := s.RecvChan()

	recvMsg := <-recvCh
	log.Printf("%+v", recvMsg)

	if recvMsg.Type != StateUpdate {
		t.Errorf("Message type.  got: %s, want: %s", recvMsg.Type, StateUpdate)
	}

	if recvMsg.Game.Seed != game.Seed {
		t.Errorf("Sent game seed was not received.  got: %d, want: %d", recvMsg.Game.Seed, game.Seed)
	}
}

func TestSendMove(t *testing.T) {
	s := Streamer{
		chanName: "fixme" + t.Name(),
	}

	err := s.Connect(testURL, testCredFile)
	if err != nil {
		t.Errorf("Connect failed. %q", err)
		return
	}

	_, gameInputCh, gameUpdateCh := engine.NewBucketGame()

	gameInputCh <- engine.PlayInputRotate
	game := <-gameUpdateCh

	sentMove := engine.PlayInputMoveRight
	go func() {
		s.SendMove(sentMove, *game)
	}()

	recvCh := s.RecvChan()

	recvMsg := <-recvCh
	log.Printf("%+v", recvMsg)

	if recvMsg.Type != Move {
		t.Errorf("Message type.  got: %s, want: %s", recvMsg.Type, Move)
	}

	if recvMsg.Move != sentMove {
		t.Errorf("Sent move was not received.  got: %d, want: %d", recvMsg.Move, sentMove)
	}
}
