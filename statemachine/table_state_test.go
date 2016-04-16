package statemachine

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/kr/pretty"
)

// action represents an intent .
type action struct {
	typ actionType // The type of this action.
	val string     // The value of this action.
}

// itemType identifies the type of lex items.
type actionType int

const (
	actionError    actionType = iota // error occurred; value is text of error
	actionBid                        // player bid
	actionAnnounce                   // player announcement - eg. Bella
	actionPlayCard                   // player submitting a card to play
	actionDeal                       // player request to deal
	actionAccuse                     // accuse another player of a misplay
	actionEOG                        //end of game
)

// StateMachine provides a simple state machine for testing the Executor.
type StateMachine struct {
	err       bool
	callTrace []string
	bids      []string
	actions   <-chan action
}

// Start implements StateFn.
func (s *StateMachine) Start() (StateFn, error) {
	s.trace()
	return s.Bid, nil
}

// Bid implements StateFn.
func (s *StateMachine) Bid() (StateFn, error) {
	//fmt.Println("bid action : ", <-s.actions)
	s.trace()
	currentAction := <-s.actions
    if currentAction.typ != actionBid {
        fmt.Println("equal? ",currentAction.typ,actionBid)
        //return s.Bid, nil//retry Bid
    }
    if len(s.bids) > 3 {
        return s.Middle, nil
    }
	if len(s.bids) == 0 {
		//fmt.Println("bid action : ", <-s.actions)
		fmt.Println("bid action : ", currentAction)
		s.bids = append(s.bids, currentAction.val)
		return s.Bid, nil
	}
	currentBid := s.bids[len(s.bids)-1]
	if currentBid == "pass" && len(s.bids) < 4 {
		if len(s.bids) == 3 {
			s.bids = append(s.bids, "hearts")
			return s.Middle, nil
		}
		s.bids = append(s.bids, "pass")
		return s.Bid, nil
	}
	return s.Middle, nil
}

// Middle implements StateFn.
func (s *StateMachine) Middle() (StateFn, error) {
	s.trace()
	if s.err {
		return s.Error, nil
	}
	return s.End, nil
}

// End implements StateFn.
func (s *StateMachine) End() (StateFn, error) {
	s.trace()
	return nil, nil
}

// Error implements StateFn.
func (s *StateMachine) Error() (StateFn, error) {
	s.trace()
	return nil, fmt.Errorf("error")
}

func (s *StateMachine) reset() {
	s.callTrace = nil
}

// trace adds the caller's name to s.callTrace.
func (s *StateMachine) trace() {
	pc, _, _, _ := runtime.Caller(1)
	s.callTrace = append(s.callTrace, fScrub(runtime.FuncForPC(pc).Name()))
}

type logging struct {
	msgs []string
}

func (l *logging) Log(s string, i ...interface{}) {
	l.msgs = append(l.msgs, fmt.Sprintf(s, i...))
}

func gen(actions []action) <-chan action {
	out := make(chan action)
	go func() { //goroutine allows code to block but does not block main thread
		for _, n := range actions {
			out <- n
			fmt.Println("writing action ", n)
		}
		close(out)
	}()
	return out
}

func TestExecutor(t *testing.T) {
	tests := []struct {
		desc      string
		err       bool
		shouldLog bool
		log       []string
		actions   []action
	}{
		{
			desc: "With error in state machine execution",
			err:  true,
		},
		{
			desc:      "Success",
			shouldLog: true,
			log: []string{
				"StateMachine[tester]: StateFn(Start) starting",
				"StateMachine[tester]: StateFn(Start) finished",
				"StateMachine[tester]: StateFn(Bid) starting",
				"StateMachine[tester]: StateFn(Bid) finished",
				"StateMachine[tester]: StateFn(Bid) starting",
				"StateMachine[tester]: StateFn(Bid) finished",
				"StateMachine[tester]: StateFn(Bid) starting",
				"StateMachine[tester]: StateFn(Bid) finished",
				"StateMachine[tester]: StateFn(Bid) starting",
				"StateMachine[tester]: StateFn(Bid) finished",
				"StateMachine[tester]: StateFn(Middle) starting",
				"StateMachine[tester]: StateFn(Middle) finished",
				"StateMachine[tester]: StateFn(End) starting",
				"StateMachine[tester]: StateFn(End) finished",
				"StateMachine[tester]: Execute() completed with no issues",
				"StateMachine[tester]: The following is the StateFn's called with this execution:",
				"StateMachine[tester]: \tStart",
				"StateMachine[tester]: \tBid",
				"StateMachine[tester]: \tBid",
				"StateMachine[tester]: \tBid",
				"StateMachine[tester]: \tBid",
				"StateMachine[tester]: \tMiddle",
				"StateMachine[tester]: \tEnd",
			},
			actions: []action{
				action{typ: actionBid, val: "pass"},
				action{typ: actionBid, val: "pass"},
				action{typ: actionBid, val: "pass"},
				action{typ: actionBid, val: "hearts"},
			},
		},
	}
	// mapD := map[string]int{"apple": 5, "lettuce": 7}
	// mapB, _ := json.Marshal(mapD)
	// fmt.Println(string(mapB))
	sm := &StateMachine{}
	for _, test := range tests {
		sm.err = test.err
		sm.bids = []string{}
		actionChan := gen(test.actions)
		// for elem := range actionChan {
		// 	t.Log(elem)
		// }
		sm.actions = actionChan
		l := &logging{}
		exec := New("tester", sm.Start, Reset(sm.reset), LogFacility(l.Log))
		if test.shouldLog {
			exec.Log(true)
		} else {
			exec.Log(false)
		}
		err := exec.Execute()
		t.Logf("bids %v", sm.bids)
		switch {
		case err == nil && test.err:
			t.Errorf("Test %q: got err == nil, want err != nil", test.desc)
			continue
		case err != nil && !test.err:
			t.Errorf("Test %q: got err != %q, want err == nil", test.desc, err)
			continue
		}

		if diff := pretty.Diff(sm.callTrace, exec.Nodes()); len(diff) != 0 {
			t.Errorf("Test %q: node trace was no accurate got/want diff:\n%s", test.desc, strings.Join(diff, "\n"))
		}

		if diff := pretty.Diff(l.msgs, test.log); len(diff) != 0 {
			t.Errorf("Test %q: log was not as expected:\n%s", test.desc, strings.Join(diff, "\n"))
		}
		t.Log("logging.....", l.msgs)
	}
}

func TestMock(t *testing.T) {
	var _ Executor = &MockExecutor{}
}
