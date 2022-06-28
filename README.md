Tetris Game
================

Console game
----------------
To play the game use:

```go run tetris_console_game.go```

Or build it via:

```go build -o tetris_console_game tetris_console_game.go```

Screenshot:

![Console Game Screenshot](https://raw.githubusercontent.com/superfrink/tetris/master/doc/tetris-screenshot.png)

[Video of game play](https://youtu.be/E1sI_jp-vLU "Video of game play") (A human is playing on the left, the right is showing a game with PRNG input.)

Bucket game
----------------

There is a ```-b``` flag to the ```tetris_console_game``` to play a simple bucket tetris-like game.  The bucket game only has one piece and that piece is a single block.  The game only has 10 rows and 3 columns.

![Bucket Game Screenshot](https://raw.githubusercontent.com/superfrink/tetris/master/doc/bucket-game-screenshot.png)

Build status
----------------

![Check-In Build](https://github.com/superfrink/tetris/actions/workflows/push-check.yml/badge.svg)
