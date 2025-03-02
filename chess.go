package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ChessGameNormal struct {
	game_id    string // 6 letter code
	players    []string
	colors     []string
	times      []int
	timesadd   []int
	isFinished bool
	board      board
	PGN        string
	mycolor    int
}

type board struct {
	desk          [8][8]int
	wCastleK      bool
	wCastleQ      bool
	bCastleK      bool
	bCastleQ      bool
	turnOfPlayer  bool
	enPassant     [2]int
	halfMoveCount int
	fullMoveCount int
}

var pieceMap = map[rune]int{
	'P': 1, 'N': 2, 'B': 3, 'R': 4, 'Q': 5, 'K': 6, // White pieces
	'p': -1, 'n': -2, 'b': -3, 'r': -4, 'q': -5, 'k': -6, // Black pieces
}

func FENtoBoard(fen string) board {
	var b board
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			b.desk[i][j] = 0
		}
	}
	parts := strings.Split(fen, " ")

	rows := strings.Split(parts[0], "/")
	for r, row := range rows {
		file := 0
		for _, ch := range row {
			if ch >= '1' && ch <= '8' {
				file += int(ch - '0')
			} else {
				b.desk[r][file] = pieceMap[ch]
				file++
			}
		}
	}
	b.turnOfPlayer = (parts[1] == "w")
	b.wCastleK = strings.Contains(parts[2], "K")
	b.wCastleQ = strings.Contains(parts[2], "Q")
	b.bCastleK = strings.Contains(parts[2], "k")
	b.bCastleQ = strings.Contains(parts[2], "q")
	if parts[3] == "-" {
		b.enPassant = [2]int{-1, -1}
	} else {
		b.enPassant = [2]int{int(parts[3][1] - '1'), int(parts[3][0] - 'a')}
	}
	b.halfMoveCount, _ = strconv.Atoi(parts[4])
	b.fullMoveCount, _ = strconv.Atoi(parts[5])
	return b
}

func BoardToFEN(b board) string {
	var fenBuilder strings.Builder
	for r := 0; r < 8; r++ {
		emptyCount := 0
		for f := 0; f < 8; f++ {
			piece := b.desk[r][f]
			if piece == 0 {
				emptyCount++
			} else {
				if emptyCount > 0 {
					fenBuilder.WriteString(strconv.Itoa(emptyCount))
					emptyCount = 0
				}
				for ch, val := range pieceMap {
					if val == piece {
						fenBuilder.WriteRune(ch)
						break
					}
				}
			}
		}
		if emptyCount > 0 {
			fenBuilder.WriteString(strconv.Itoa(emptyCount))
		}
		if r < 7 {
			fenBuilder.WriteString("/")
		}
	}
	if b.turnOfPlayer {
		fenBuilder.WriteString(" w ")
	} else {
		fenBuilder.WriteString(" b ")
	}
	var castleRights string
	if b.wCastleK {
		castleRights += "K"
	}
	if b.wCastleQ {
		castleRights += "Q"
	}
	if b.bCastleK {
		castleRights += "k"
	}
	if b.bCastleQ {
		castleRights += "q"
	}
	if castleRights == "" {
		castleRights = "-"
	}
	fenBuilder.WriteString(castleRights + " ")
	if b.enPassant[0] == -1 {
		fenBuilder.WriteString("- ")
	} else {
		fenBuilder.WriteString(fmt.Sprintf("%c%d ", 'a'+b.enPassant[1], b.enPassant[0]+1))
	}
	fenBuilder.WriteString(strconv.Itoa(b.halfMoveCount) + " ")
	fenBuilder.WriteString(strconv.Itoa(b.fullMoveCount))
	return fenBuilder.String()
}

func isValidMove(b board, from [2]int, to [2]int, mycolor int) bool {
}

func getColorOfCell(b board, coord [2]int) {
}
