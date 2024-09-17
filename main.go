package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Game struct {
	ID      string `json:"id"`
	Timeout int32  `json:"timeout"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Battlesnake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int32   `json:"health"`
	Body   []Coord `json:"body"`
	Head   Coord   `json:"head"`
	Length int32   `json:"length"`
	Shout  string  `json:"shout"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`
}

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}

type GameRequest struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

// HandleIndex is called when your Battlesnake is created and refreshed
// by play.battlesnake.com. BattlesnakeInfoResponse contains information about
// your Battlesnake, including what it should look like on the game board.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "bd",
		Color:      "#888800",
		Head:       "sand-worm",
		Tail:       "round-bum",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal().Err(err)
	}
}

// HandleStart is called at the start of each game your Battlesnake is playing.
// The GameRequest object contains information about the game that's about to start.
// TODO: Use this function to decide how your Battlesnake is going to look on the board.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal().Err(err)
	}

	// Nothing to respond with here
	log.Info().Str("name", request.You.Name).Str("id", request.You.ID).Msg("game start")
}

// HandleMove is called for each turn of each game.
// Valid responses are "up", "down", "left", or "right".
// TODO: Use the information in the GameRequest object to determine your next move.
func HandleMove(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal().Err(err)
	}
	move := request.chooseMove()
	log.Debug().Str("name", request.You.Name).Str("id", request.You.ID).Str("move", move).Msg("move chosen")
	response := MoveResponse{
		Move: move,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal().Err(err)
	}
}
func (gr GameRequest) offBoard(coord Coord) bool {
	log.Debug().Str("name", gr.You.Name).Str("id", gr.You.ID).Msg("off board - start")
	defer log.Debug().Str("name", gr.You.Name).Str("id", gr.You.ID).Msg("off board - end")
	if coord.X < 0 {
		return true
	}
	if coord.Y < 0 {
		return true
	}
	if coord.X >= gr.Board.Width {
		return true
	}
	if coord.Y >= gr.Board.Height {
		return true
	}
	return false
}

func (gr GameRequest) collideSnake(futureCoord Coord) bool {
	log.Debug().Str("name", gr.You.Name).Str("id", gr.You.ID).Msg("snake collision - start")
	defer log.Debug().Str("name", gr.You.Name).Str("id", gr.You.ID).Msg("snake collision - end")
	for _, snake := range gr.Board.Snakes {
		log.Debug().Str("name", gr.You.Name).Str("id", gr.You.ID).
			Str("other_name", snake.Name).Str("other_id", snake.ID).Msg("snake collision - check")
		for _, sb := range snake.Body {
			if coordEquals(futureCoord, sb) {
				return true
			}
		}
	}
	return false
}

func coordEquals(a, b Coord) bool {
	if a.X == b.X && a.Y == b.Y {
		return true
	}
	return false
}

func (gr GameRequest) chooseMove() string {
	defaultMove := "left"
	possibleMoves := []string{"up", "down", "left", "right"}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	log.Debug().Str("name", gr.You.Name).Str("id", gr.You.ID).Msg("choose move - start")
	for _, i := range r.Perm(len(possibleMoves)) {
		mv := possibleMoves[i]
		futureCoord := gr.coordAsMove(mv)
		log.Debug().Str("move", mv).Interface("coord", futureCoord).Msg("check future move")
		if !gr.offBoard(futureCoord) && !gr.collideSnake(futureCoord) {
			return mv
		}
	}
	log.Warn().Str("move", defaultMove).Msg("no options. using default move")
	log.Debug().Str("name", gr.You.Name).Str("id", gr.You.ID).Msg("choose move - end")
	return defaultMove

}

//coordAsMove input move (up,down,left, right) get back
//coordinate of your snake head for the move
func (gr GameRequest) coordAsMove(move string) Coord {
	switch move {
	case "up":
		return Coord{gr.You.Head.X, gr.You.Head.Y + 1}
	case "down":
		return Coord{gr.You.Head.X, gr.You.Head.Y - 1}
	case "left":
		return Coord{gr.You.Head.X - 1, gr.You.Head.Y}
	case "right":
		return Coord{gr.You.Head.X + 1, gr.You.Head.Y}
	}
	return Coord{}
}

// HandleEnd is called when a game your Battlesnake was playing has ended.
// It's purely for informational purposes, no response required.
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal().Err(err)
	}

	// Nothing to respond with here
	log.Info().Str("name", request.You.Name).Str("id", request.You.ID).Msg("gg - game end")
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)
	log.Info().Str("url", fmt.Sprintf("http://0.0.0.0:%s", port)).Msg("starting battlesnake server")
	log.Fatal().Err(http.ListenAndServe(":"+port, nil))
}
