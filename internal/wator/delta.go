package wator

import "fmt"

// Delta describes the changes of a creature between two Chronon.
type Delta struct {
	Object int // type of creature: FISH, SHARK
	From   int // position in previous Chronon
	To     int // position in current Chronon
	Action int // Action = NO_ACTION, MOVE, DEATH, BIRTH
}

// Dump prints out the content of a delta.
func (d *Delta) Dump() {

	var obj, action string

	switch d.Object {

	case NONE:
		obj = "NONE"
	case FISH:
		obj = "FISH"
	case SHARK:
		obj = "SHARK"
	default:
		obj = "UNKNOWN"

	}

	switch d.Action {
	case NO_ACTION:
		action = "NO_ACTION"

	case MOVE:
		action = "MOVE"
	case MOVE_NONE:
		action = "MOVE_NONE"
	case MOVE_NORTH:
		action = "MOVE_NORTH"
	case MOVE_SOUTH:
		action = "MOVE_SOUTH"
	case MOVE_EAST:
		action = "MOVE_EAST"
	case MOVE_WEST:
		action = "MOVE_WEST"
	case DEATH:
		action = "DEATH"
	case BIRTH:
		action = "BIRTH"
	case ATE:
		action = "ATE"
	}

	fmt.Printf("Animal = %v from %d to %d Action=%v\n", obj, d.From, d.To, action)

}
