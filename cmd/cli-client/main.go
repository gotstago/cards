package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gotstago/cards/hand"
	"github.com/gotstago/cards/table"
)

const (
	fold  = "fold"
	check = "check"
	call  = "call"
	bet   = "bet"
	raise = "raise"
)

var (
	tbl *table.PokerTable
)

type player struct {
	id string
}

func (p *player) ID() string {
	return p.id
}

func (p *player) FromID(id string) (table.Player, error) {
	return &player{p.ID()}, nil
}

func (p *player) Action() (table.Action, int) {
	current := tbl.CurrentPlayer()

	// get action from input
	validActions := []string{}
	for _, a := range tbl.ValidActions() {
		validActions = append(validActions, strings.ToLower(string(a)))
	}

	// show info
	currentInfoFormat := "\nChips %d, Outstanding %d, MinRaise %d, MaxRaise %d"
	fmt.Printf(currentInfoFormat, current.Chips(), tbl.Outstanding(), tbl.MinRaise(), tbl.MaxRaise())

	// get action from input
	var input string
	actionFormat := "\nPlayer %s Action (%s):\n"
	fmt.Printf(actionFormat, p.ID(), strings.Join(validActions, ","))
	if _, err := fmt.Scan(&input); err != nil {
		fmt.Println("Error", err)
		return p.Action()
	}

	// parse action
	action, err := actionFromInput(input)
	if err != nil {
		fmt.Println("Error", err)
		return p.Action()
	}
	if !(action == table.Bet || action == table.Raise) {
		return action, 0
	}

	// get amount from input
	amountFormat := "\nEnter Bet / Raise Amount:\n"
	fmt.Printf(amountFormat)
	if _, err := fmt.Scan(&input); err != nil {
		fmt.Println("Error", err)
		return p.Action()
	}

	// parse amount
	chips, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		fmt.Println("Error", err)
		return p.Action()
	}
	return action, int(chips)
}

func main() {
	p1 := playerFromInput("Player 1")
	p2 := playerFromInput("Player 2")
	p3 := playerFromInput("Player 3")
	p4 := playerFromInput("Player 4")

	opts := table.Config{
		Game:       table.Tarabish,
		Limit:      table.NoLimit,
		Stakes:     table.Stakes{SmallBet: 1, BigBet: 2, Ante: 0},
		NumOfSeats: 4,
	}
	tbl = table.New(opts, hand.NewDealer())
	if err := tbl.Sit(p1, table.Parameters{0,100}); err != nil {
		panic(err)
	}
	if err := tbl.Sit(p2, table.Parameters{1,100}); err != nil {
		panic(err)
	}
	if err := tbl.Sit(p3, table.Parameters{2,100}); err != nil {
		panic(err)
	}
	if err := tbl.Sit(p4, table.Parameters{3,100}); err != nil {
		panic(err)
	}

	runTable(tbl)
	fmt.Println("DONE")
}

func runTable(tbl *table.PokerTable) {
	for {
		results, done, err := tbl.Next()
		if done {
			return
		}
		printTable(tbl)
		if err != nil {
			fmt.Println("Error", err)
		}
		if results != nil {
			printResults(tbl, results)
		}
	}
}

func printTable(tbl *table.PokerTable) {
	players := tbl.Players()
	fmt.Println("")
	fmt.Println("-----Table-----")
	fmt.Println(tbl)
	fmt.Println(players[0])
	fmt.Println(players[1])
	fmt.Println("-----Table-----")
	fmt.Println("")
}

func printResults(tbl *table.PokerTable, results map[int][]*table.Result) {
	players := tbl.Players()
	for seat, resultList := range results {
		for _, result := range resultList {
			fmt.Println(players[seat].Player().ID()+":", result)
		}
	}
}

func playerFromInput(desc string) table.Player {
	var input string
	fmt.Printf("\nPick %s name:\n", desc)
	//wait for input from user at command line.
    if _, err := fmt.Scan(&input); err != nil {
		fmt.Println("Error", err)
		return playerFromInput(desc)
	}
	return &player{id: input}
}

func actionFromInput(input string) (table.Action, error) {
	switch input {
	case fold:
		return table.Fold, nil
	case check:
		return table.Check, nil
	case call:
		return table.Call, nil
	case bet:
		return table.Bet, nil
	case raise:
		return table.Raise, nil
	}
	return table.Fold, errors.New(input + " is not an action.")
}
