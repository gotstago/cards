package table_test

import (
	//"encoding/json"
	"testing"

	"github.com/gotstago/cards/hand"
	"github.com/gotstago/cards/table"
)

func TestTarabishSeating(t *testing.T) {
	t.Parallel()

	opts := table.Config{
		Game: table.Tarabish,
		// Stakes: table.Stakes{
		// 	SmallBet: 1,
		// 	BigBet:   2,
		// 	Ante:     0,
		// },
		NumOfSeats: 4,
	}

	p1 := Player("1", []PlayerAction{})
	p1Dup := Player("1", []PlayerAction{})
	p2 := Player("2", []PlayerAction{})

	tbl := table.New(opts, hand.NewDealer())

	// sit player 1
	if err := tbl.Sit(p1, table.Parameters{Seat:0}); err != nil {
		t.Fatal(err)
	}

	// can't sit dup player 1

	if err := tbl.Sit(p1Dup, table.Parameters{Seat:1}); err != table.ErrAlreadySeated {
		t.Fatal("should already be seated")
	}

	// can't sit player 2 in invalid seat
	if err := tbl.Sit(p2, table.Parameters{Seat:4}); err != table.ErrInvalidSeat {
		t.Fatal("can't sit in invalid seat")
	}

	// can't sit player 2 in occupied seat
	if err := tbl.Sit(p2, table.Parameters{Seat:0}); err != table.ErrSeatOccupied {
		t.Fatal("can't sit in occupied seat")
	}
}

func TestTarRaises(t *testing.T) {
	t.Parallel()

	opts := table.Config{
		Game: table.Tarabish,
		// Stakes: table.Stakes{
		// 	SmallBet: 1,
		// 	BigBet:   2,
		// 	Ante:     0,
		// },
		NumOfSeats: 4,
	}

	p1 := Player("1", []PlayerAction{})
	p2 := Player("2", []PlayerAction{})
	p3 := Player("3", []PlayerAction{})
	p4 := Player("4", []PlayerAction{})

	tbl := table.New(opts, hand.NewDealer())

	if err := tbl.Sit(p1, table.Parameters{Seat:0}); err != nil {
		t.Fatal(err)
	}
	if err := tbl.Sit(p2, table.Parameters{Seat:1}); err != nil {
		t.Fatal(err)
	}
	if err := tbl.Sit(p3, table.Parameters{Seat:2}); err != nil {
		t.Fatal(err)
	}
	if err := tbl.Sit(p4, table.Parameters{Seat:3}); err != nil {
		t.Fatal(err)
	}
    
    //bidding
    p1.Bid("pass")
    p2.Bid("pass")
    // p2.Fold()
    // p3.Fold()
    // p4.Fold()
	// preflop
	// p1.Call()
	// p2.Call()
	// p3.Call()
	// p4.Check()

	// // flop
	// p3.Check()
	// p4.Check()
	// p1.Bet(48)
	// p2.Call()
	// p3.Raise(2)
	// p4.Raise(8)

	for i := 0; i < 2; i++ {
		if _, _, err := tbl.Next(); err != nil {
			t.Fatal(err)
		}else{
            t.Logf("Next succeeded : action is %v number of seats is %d and %d",
            tbl.Action(),
            tbl.NumOfSeats(), 
            tbl.Round(),
            )
        }
	}

	if tbl.Action() != 1 {
		t.Fatal("action should be on player 2",tbl.Action())
	}

	players := tbl.Players()
	if players[1].CanRaise() {
		t.Fatal("player 2 shouldn't be able to raise")
	}

	// p2.Call()
	// _, _, err := tbl.Next()
	// _, _, err = tbl.Next()
	// if err != nil {
	// 	t.Fatal(err)
	// }
}

func TestTableSetup(t *testing.T) {
	t.Parallel()

	opts := table.Config{
		Game: table.Tarabish,
		// Stakes: table.Stakes{
		// 	SmallBet: 1,
		// 	BigBet:   2,
		// 	Ante:     0,
		// },
		NumOfSeats: 4,
	}

	p1 := Player("1", []PlayerAction{})
	p2 := Player("2", []PlayerAction{})
	p3 := Player("3", []PlayerAction{})
	p4 := Player("4", []PlayerAction{})

	tbl := table.New(opts, hand.NewDealer())

	if err := tbl.Sit(p1, table.Parameters{Seat:0}); err != nil {
		t.Fatal(err)
	}
	if err := tbl.Sit(p2, table.Parameters{Seat:1}); err != nil {
		t.Fatal(err)
	}
	if err := tbl.Sit(p3, table.Parameters{Seat:2}); err != nil {
		t.Fatal(err)
	}
	if err := tbl.Sit(p4, table.Parameters{Seat:3}); err != nil {
		t.Fatal(err)
	}
    
    //bidding
    p1.Bid("pass")
    p2.Bid("pass")

	for i := 0; i < 2; i++ {
		if _, _, err := tbl.Next(); err != nil {
			t.Fatal(err)
		}else{
            t.Logf("Next succeeded : action is %v number of seats is %d and %d",
            tbl.Action(),
            tbl.NumOfSeats(), 
            tbl.Round(),
            )
        }
	}

	if tbl.Action() != 1 {
		t.Fatal("action should be on player 2",tbl.Action())
	}

	players := tbl.Players()
	if players[1].CanRaise() {
		t.Fatal("player 2 shouldn't be able to raise")
	}

	// p2.Call()
	// _, _, err := tbl.Next()
	// _, _, err = tbl.Next()
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
