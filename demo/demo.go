package main

import (
	"cmp"
	"flag"
	"fmt"
	"github.com/pronovic/go-apologies/engine"
	"github.com/pronovic/go-apologies/model"
	"github.com/pronovic/go-apologies/render"
	"github.com/pronovic/go-apologies/source"
	"github.com/rthornton128/goncurses"
	"log"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"
)

const minCols = 155
const minRows = 60

func main() {
	players, delay, exit, mode, cis := parseArgs()

	characters := make([]engine.Character, players)
	for player := 0; player < players; player++ {
		name := fmt.Sprintf("Player %d", player)
		characters[player] = engine.NewCharacter(name, cis)
	}

	runtime, err := engine.NewEngine(mode, characters, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = runtime.StartGame()
	if err != nil {
		log.Fatal(err)
	}

	forceMinimumSize()
	cursesMain(cis, runtime, delay, exit)
}

func parseArgs() (int, int, bool, model.GameMode, source.CharacterInputSource) {
	players := flag.Int("players", 2, "number of players")
	delay := flag.Int("delay", 200, "delay between moves (milliseconds)")
	adult := flag.Bool("adult", false, "run in adult mode")
	input := flag.String("input", "random", "'random' or 'reward' for input source")
	exit := flag.Bool("exit", false, "exit immediately upon completion")

	flag.Parse()

	mode := model.StandardMode
	if *adult {
		mode = model.AdultMode
	}

	cis := source.RandomInputSource()
	if *input == "reward" {
		cis = source.RewardInputSource(nil, nil)
	}

	return *players, *delay, *exit, mode, cis
}

// forceMinimimumSize Force an xterm to resize via a control sequence.
func forceMinimumSize() {
	fmt.Printf("\u001b[8;%d;%dt", minRows, minCols)
	time.Sleep(time.Duration(500)*time.Millisecond)
}

// cursesMain is the ncurses main routine
func cursesMain(cis source.CharacterInputSource, runtime engine.Engine, delay int, exit bool) {
	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer goncurses.End()

	rows, columns := stdscr.MaxYX()
	if columns < minCols || rows < minRows {
		log.Fatalf("Minimum terminal size is %dx%d, but yours is %dx%d", minCols, minRows, columns, rows)
	}

	board, err := goncurses.NewWindow(53, 90, 1, 3)
	if err != nil {
		log.Fatal(err)
	}

	state, err := goncurses.NewWindow(52, 59, 2, 94)
	if err != nil {
		log.Fatal(err)
	}

	history, err := goncurses.NewWindow(5, 150, 54, 3)
	if err != nil {
		log.Fatal(err)
	}

	complete := false

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		_ = <- interrupt
		complete = true
		goncurses.End()
	}()

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	go func() {
		_ = <- resize
		draw(stdscr, board, state, history)
	}()


	for {
		if complete {
			break
		}

		if runtime.Completed() {
			if exit {
				complete = true
			}
		} else {
			game, _ := runtime.PlayNext()
			refresh(cis, runtime, game, delay, stdscr, board, state, history)
		}

		time.Sleep(time.Duration(delay)*time.Millisecond)
	}
}

func draw(
		stdscr *goncurses.Window,
		board *goncurses.Window,
		state *goncurses.Window,
		history *goncurses.Window) {
	err := stdscr.Clear()
	if err != nil {
		log.Fatal(err)
	}

	err = stdscr.Box(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	err = board.Box(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	err = state.Box(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	err = history.Box(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	stdscr.Refresh()
	board.Refresh()
	state.Refresh()
	history.Refresh()
}

func refresh(
		cis source.CharacterInputSource,
		runtime engine.Engine,
		game model.Game,
		delay int,
		stdscr *goncurses.Window,
		board *goncurses.Window,
		state *goncurses.Window,
		history *goncurses.Window) {
	refreshScreen(stdscr)
	refreshBoard(game, board)
	refreshState(cis, runtime, game, delay, state)
	refreshHistory(game, history)
}

func refreshScreen(stdscr *goncurses.Window) {
	err := stdscr.Box(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	stdscr.MovePrint(1, 95, "APOLOGIES DEMO")
	stdscr.MovePrint(1, 138, "CTRL-C TO EXIT")
	stdscr.Move(minRows-2, minCols-2)  // bottom-right corner

	stdscr.Refresh()
}

func refreshBoard(game model.Game, board *goncurses.Window) {
	err := board.Clear()
	if err != nil {
		log.Fatal(err)
	}

	err = board.Box(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	rendered, err := render.Board(game)
	if err != nil {
		log.Fatal(err)
	}

	// Unfortunately, the ncurses implementation for GoLang does not work with some of the
	// unicode characters that I've been using in the board definition in the Python
	// implementation.  In some cases, I can translate to a roughly equivalent result, and
	// in other cases I'm just clearing out the junk.  It's not as legible as I had
	// intended, but still serves the purpose of letting someone watch the code play a
	// game.

	row := 0
	char := 0
	runes := []rune(rendered)
	for _, r := range runes {
		if r == '\n' {
			row += 1
			char = 0
		} else {
			char += 1
			if r == '┌' {
				board.MoveAddChar(row, char, goncurses.ACS_ULCORNER)
			} else if r == '┐' {
				board.MoveAddChar(row, char, goncurses.ACS_URCORNER)
			} else if r == '└' {
				board.MoveAddChar(row, char, goncurses.ACS_LLCORNER)
			} else if r == '┘' {
				board.MoveAddChar(row, char, goncurses.ACS_LRCORNER)
			} else if r == '─' {
				board.MoveAddChar(row, char, goncurses.ACS_HLINE)
			} else if r == '│' {
				board.MoveAddChar(row, char, goncurses.ACS_VLINE)
			} else if r == '◼' {
				board.MovePrint(row, char, string(' '))
			} else if r == '●' {
				board.MovePrint(row, char, string(' '))
			} else if r == '▶' {
				board.MovePrint(row, char, string(' '))
			} else if r == '▼' {
				board.MovePrint(row, char, string(' '))
			} else if r == '◀' {
				board.MovePrint(row, char, string(' '))
			} else if r == '▲' {
				board.MovePrint(row, char, string(' '))
			} else {
				board.MovePrint(row, char, string(r))
			}
		}
	}

	board.Refresh()
}

func refreshState(
		cis source.CharacterInputSource,
		runtime engine.Engine,
		game model.Game,
		delay int,
		state *goncurses.Window) {
	err := state.Clear()
	if err != nil {
		log.Fatal(err)
	}

	err = state.Box(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	state.MovePrint(1, 2, "CONFIGURATION")
	state.MovePrintf(3, 3, "Players..: %d", runtime.Players())
	state.MovePrintf(4, 3, "Mode.....: %s", runtime.Mode().Value())
	state.MovePrintf(5, 3, "Source...: %s", cis.Name())
	state.MovePrintf(6, 3, "Delay....: %d ms", delay)
	state.MovePrintf(7, 3, "State....: %s", runtime.State())

	players := make([]model.Player, 0)
	for _, player := range game.Players() {
		players = append(players, player)
	}

	slices.SortStableFunc(players, func(i, j model.Player) int {
		return cmp.Compare(i.Color().Value(), j.Color().Value())
	})

	row := 10
	for _, player := range players {
		state.MovePrintf(row + 0, 2, "%s PLAYER", strings.ToUpper(player.Color().Value()))
		state.MovePrintf(row + 2, 3, "Hand.....: %s", renderHand(player))
		state.MovePrintf(row + 3, 3, "Pawns....:")
		state.MovePrint(row + 4, 6, player.Pawns()[0])
		state.MovePrint(row + 5, 6, player.Pawns()[1])
		state.MovePrint(row + 6, 6, player.Pawns()[2])
		state.MovePrint(row + 7, 6, player.Pawns()[3])
		row += 10
	}

	state.Refresh()
}

func refreshHistory(game model.Game, history *goncurses.Window) {
	err := history.Clear()
	if err != nil {
		log.Fatal(err)
	}

	err = history.Box(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	length := len(game.History())
	var entries []model.History
	if length < 4 {
		entries = game.History()
	} else {
		entries = game.History()[length-3:length]
	}

	row := 1
	for _, entry := range entries {
		history.MovePrint(row, 2, entry)
		row += 1
	}

	history.Refresh()
}

func renderHand(player model.Player) string {
	if len(player.Hand()) == 0 {
		return "n/a"
	} else {
		hand := ""
		for _, card := range player.Hand() {
			hand = card.Type().Value() + " " + hand
		}
		return hand
	}
}