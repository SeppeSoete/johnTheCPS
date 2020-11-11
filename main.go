package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"strings"
)

var piecesWithoutKing [4]string = [4]string{"R", "N", "B", "Q"}
var stupid bool

func main() {

	stupidFlag := flag.Bool("s", false, "use stupid pedantic rules that are extreme edge cases if someone wants to bully this tool. Looking at you tzlil and eyeoh")
	nthreads := flag.Int("n", 4, "Amount of goroutines to spawn for the md5summing")
	flag.Parse()

	stupid = *stupidFlag
	moves := make(chan string, 10)
	answer := make(chan string, 1)
	go generatePieceMoves(moves)
	if *nthreads <= 0 {
		*nthreads = 1
	}
	for i := 0; i < *nthreads; i++ {
		go md5sum(flag.Args()[0], moves, answer)
	}
	fmt.Printf("The answer was %v\n", <-answer)
}

func generatePieceMoves(c chan string) {
	//Main piece moves
	pieces := [5]string{"R", "N", "B", "Q", "K"}
	files := [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}
	ranks := [8]string{"1", "2", "3", "4", "5", "6", "7", "8"}
	squares := new([64]string)
	i := 0
	for _, file := range files {
		for _, rank := range ranks {
			squares[i] = file + rank
			i++
		}
	}

	//random stuff
	withPostfixes(c, "O-O", false)
	withPostfixes(c, "O-O-O", false)
	withPostfixes(c, "0-0", false)
	withPostfixes(c, "0-0-0", false)
	c <- "1/2-1/2"
	c <- "1-0"
	c <- "0-1"
	c <- "resign"
	c <- "resigns"
	c <- "½–½"

	//all piece moves
	for _, piece := range pieces {
		for _, square := range squares {
			withPostfixes(c, square, false) //regular pawn moves
			for i := 1; i <= 8; i++ {
				withPostfixes(c, piece+square, false)
				withPostfixes(c, piece+"x"+square, false)
				for _, ambigRank := range ranks { //ambiguous ranks
					withPostfixes(c, piece+ambigRank+square, false)
					withPostfixes(c, piece+ambigRank+"x"+square, false)
				}
				for _, ambigFile := range files { //ambiguous files
					withPostfixes(c, piece+ambigFile+square, false)
					withPostfixes(c, piece+ambigFile+"x"+square, false)
				}
				if stupid { //not really stupid, but rare enough to not attempt in first pass
					for _, ambigSquare := range squares { //ambiguous ranks and files
						withPostfixes(c, piece+ambigSquare+square, false)
						withPostfixes(c, piece+ambigSquare+"x"+square, false)
					}
				}
			}
		}
	}

	//HERE BE PAWNS, ALL HOPE SHALL BE LOST BEYOND THIS POINT IN THE CODE

	//capture right
	for file := 0; file < 7; file++ {
		for rank := 1; rank < 7; rank++ {
			withPostfixes(c, files[file]+"x"+files[file+1]+ranks[rank], true)
		}
	}

	//capture left
	for file := 1; file < 8; file++ {
		for rank := 1; rank < 7; rank++ {
			withPostfixes(c, files[file]+"x"+files[file-1]+ranks[rank], true)
		}
	}

	//capture right with promotion
	for file := 0; file < 7; file++ {
		promote(c, files[file]+"x"+files[file+1]+ranks[7], true)
		promote(c, files[file]+"x"+files[file+1]+ranks[0], true)
	}

	//capture left with promotion
	for file := 1; file < 8; file++ {
		promote(c, files[file]+"x"+files[file-1]+"8", true)
		promote(c, files[file]+"x"+files[file-1]+"0", true)
	}

	//pawn pushes with promotion
	for _, file := range files {
		promote(c, file+"8", false)
		promote(c, file+"1", false)
	}

	// if normal mode does not work, try stupid unicode shennanigans
	if !stupid {
		stupid = true
		generatePieceMoves(c)
	}
}

//calculate md5sum of all moves that appear on the move channel, if they match the given checksum put then on the result channel
func md5sum(sum string, move chan string, result chan string) {
	for {
		h := md5.New()
		currMove := <-move
		io.WriteString(h, currMove+"\n")
		ourSum := fmt.Sprintf("%x", h.Sum(nil))
		if strings.Compare(ourSum, sum) == 0 {
			result <- currMove
		}
	}
}

//a pawn promotion
func promote(c chan string, move string, pawnCapture bool) {
	for _, piece := range piecesWithoutKing {
		withPostfixes(c, move+"="+piece, pawnCapture)
		if stupid {
			withPostfixes(c, move+piece, pawnCapture)
		}
	}
}

func withPostfixes(c chan string, move string, pawnCapture bool) {
	c <- move
	c <- move + "#"
	c <- move + "+"
	if stupid {
		c <- move + "†"
		c <- move + "‡"
		if pawnCapture {
			c <- move + "e.p."
			c <- move + " e.p."

			c <- move + "e.p.#"
			c <- move + " e.p.#"
			c <- move + "e.p. #"
			c <- move + " e.p. #"

			c <- move + "e.p.+"
			c <- move + " e.p.+"
			c <- move + "e.p. +"
			c <- move + " e.p. +"

			c <- move + "e.p.†"
			c <- move + " e.p.†"
			c <- move + "e.p. †"
			c <- move + " e.p. †"

			c <- move + "e.p.‡"
			c <- move + " e.p.‡"
			c <- move + "e.p. ‡"
			c <- move + " e.p. ‡"
		}
	}
}
