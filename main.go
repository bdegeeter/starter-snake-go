package main

import (
	"encoding/json"
	"fmt"
	"log"
  "time"
	"math/rand"
	"net/http"
	"os"
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
		Author:     "bd",        // TODO: Your Battlesnake username
		Color:      "#888800", // TODO: Personalize
		Head:       "default", // TODO: Personalize
		Tail:       "default", // TODO: Personalize
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleStart is called at the start of each game your Battlesnake is playing.
// The GameRequest object contains information about the game that's about to start.
// TODO: Use this function to decide how your Battlesnake is going to look on the board.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("START\n")
}

// HandleMove is called for each turn of each game.
// Valid responses are "up", "down", "left", or "right".
// TODO: Use the information in the GameRequest object to determine your next move.
func HandleMove(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}
  /*
  request.chooseMove()
  
	// Choose a random direction to move in
	possibleMoves := []string{"up", "down", "left", "right"}
	move := possibleMoves[rand.Intn(len(possibleMoves))]
  */
  move := request.chooseMove()
  fmt.Printf("I chose: %s\n", move)
	response := MoveResponse{
		Move: move,
	}

	fmt.Printf("MOVE: %s\n", response.Move)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}
func (gr GameRequest) offBoard(coord Coord) bool {
  if coord.X < 0 { return true }
  if coord.Y < 0 { return true }
  if coord.X >= gr.Board.Width { return true }
  if coord.Y >= gr.Board.Height { return true }
  return false
}

func coordEquals(a, b Coord) bool {
  if a.X == b.X && a.Y == b.Y {
    return true
  }
  return false
}

func (gr GameRequest) chooseMove() string {
  possibleMoves := []string{"up", "down", "left", "right"}
  r := rand.New(rand.NewSource(time.Now().Unix()))

  for _, i := range r.Perm(len(possibleMoves)) {
    mv := possibleMoves[i]
    coord := gr.coordAsMove(mv)
    fmt.Printf("%v\n", coord)
    if !gr.offBoard(coord) && !coordEquals(coord, gr.You.Body[1]) {
      return mv
    }
  }
  return ""
}


func (gr GameRequest) coordAsMove(move string) Coord {
  fmt.Printf("head X:%d Y:%d\n", gr.You.Head.X, gr.You.Head.Y)
  switch move {
    case "up":
      fmt.Printf("up move X:%d Y:%d\n", gr.You.Head.X,gr.You.Head.Y+1)
      return Coord{gr.You.Head.X, gr.You.Head.Y+1}
    case "down":
      fmt.Printf("down move X:%d Y:%d\n", gr.You.Head.X, gr.You.Head.Y-1)
      return Coord{gr.You.Head.X, gr.You.Head.Y-1}
    case "left":
      fmt.Printf("left move X:%d Y:%d\n", gr.You.Head.X-1, gr.You.Head.Y)
      return Coord{gr.You.Head.X-1, gr.You.Head.Y}
    case "right":
      fmt.Printf("right move X:%d Y:%d\n", gr.You.Head.X+1, gr.You.Head.Y)
      return Coord{gr.You.Head.X+1, gr.You.Head.Y}
  }
  return Coord{}
}

// HandleEnd is called when a game your Battlesnake was playing has ended.
// It's purely for informational purposes, no response required.
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("gg\nEND\n")
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)

	fmt.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
