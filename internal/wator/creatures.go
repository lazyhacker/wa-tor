package wator

var sharkID int

type Creature struct {
	age      int
	lastMove uint
}

func (c *Creature) SetAge(a int) {
	c.age = a
}

func (c *Creature) Age() int {
	return c.age
}

func (c *Creature) LastMove() uint {
	return c.lastMove
}

func (c *Creature) SetLastMove(t uint) {

	c.lastMove = t
}

// ----------- Sharks -------------------
type Shark struct {
	health int
	Creature
	id int
}

func NewShark() *Shark {
	id := sharkID
	sharkID++
	return &Shark{
		sharkHealth,
		Creature{},
		id,
	}
}

func (s *Shark) Spawn() bool {

	if s.age%sharkSpawnRate == 0 && s.age > 0 {
		return true
	}
	return false
}

// ----------- Fish -------------------
type Fish struct {
	Creature
}

func NewFish() *Fish {
	return &Fish{
		Creature{},
	}
}

func (f *Fish) Spawn() bool {

	if f.age%fishSpawnRate == 0 && f.age > 0 {
		return true
	}
	return false
}
