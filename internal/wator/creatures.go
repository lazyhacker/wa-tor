package wator

type Creature struct {
	age int
}

func (c *Creature) SetAge(a int) {
	c.age = a
}

func (c *Creature) Age() int {
	return c.age
}

// ----------- Sharks -------------------
type Shark struct {
	health int
	Creature
}

func NewShark() *Shark {
	return &Shark{
		sharkHealth,
		Creature{},
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
