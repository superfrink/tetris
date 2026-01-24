module superfrink.net/tetris

go 1.22.1

require (
	github.com/klauspost/compress v1.17.2
	github.com/nsf/termbox-go v1.1.1
	github.com/pa-m/sklearn v0.0.0-20200711083454-beb861ee48b1
	superfrink.net/tetris/engine v0.0.0-00010101000000-000000000000
	superfrink.net/tetris/streamer v0.0.0-00010101000000-000000000000
)

require (
	github.com/chewxy/math32 v1.0.4 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/nats-io/nats.go v1.33.1 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pa-m/optimize v0.0.0-20190612075243-15ee852a6d9a // indirect
	github.com/pa-m/randomkit v0.0.0-20191001073902-db4fd80633df // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/exp v0.0.0-20191129062945-2f5052295587 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/tools v0.0.0-20200225230052-807dcd883420 // indirect
	gonum.org/v1/gonum v0.6.1 // indirect
)

replace superfrink.net/tetris/engine => ./engine

replace superfrink.net/tetris/streamer => ./streamer
