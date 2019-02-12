package engine

type Piece int

const (
	none = iota
	I
	J
	L
	S
	Z
	O
)

type Rotation int

const (
	r0 = iota
	r90
	r180
	r270
)

const (
	GameRows= 18
	GameColumns = 10
)

type Game struct {
	Piece Piece
	Rotation Rotation
	Field [GameRows][GameColumns] byte
}

func NewGame() *Game {
	g := Game{
		Piece: none,
		Rotation: r0,
	}

	for i := 0 ; i < GameRows ; i++ {
		for j := 0 ; j < GameColumns ; j++ {
			g.Field[i][j] = 0
		}
	}

	return &g
}
