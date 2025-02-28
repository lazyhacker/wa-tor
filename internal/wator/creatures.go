package wator

var sharkID int

type creature struct {
	chronon int
	turn    uint
}

func (c *creature) setAge(a int) {
	c.chronon = a
}

func (c *creature) age() int {
	return c.chronon
}

func (c *creature) lastMove() uint {
	return c.turn
}

func (c *creature) setLastMove(t uint) {

	c.turn = t
}

// ----------- Sharks -------------------
type shark struct {
	health int
	creature
	id int
}

func NewShark() *shark {
	id := sharkID
	sharkID++
	return &shark{
		sharkHealth,
		creature{},
		id,
	}
}

func (s *shark) spawn() bool {

	if s.chronon%sharkSpawnRate == 0 && s.chronon > 0 {
		return true
	}
	return false
}

// ----------- fish -------------------
type fish struct {
	creature
}

func NewFish() *fish {
	return &fish{
		creature{},
	}
}

func (f *fish) spawn() bool {

	if f.chronon%fishSpawnRate == 0 && f.chronon > 0 {
		return true
	}
	return false
}
