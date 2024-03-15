module superfrink.net/tetris

go 1.22.1

require (
	github.com/nsf/termbox-go v1.1.1
	superfrink.net/tetris/engine v0.0.0-00010101000000-000000000000
	superfrink.net/tetris/streamer v0.0.0-00010101000000-000000000000
)

require (
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/nats-io/nats.go v1.33.1 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

replace superfrink.net/tetris/engine => ./engine

replace superfrink.net/tetris/streamer => ./streamer
