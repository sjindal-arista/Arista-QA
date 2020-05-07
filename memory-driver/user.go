package driver

type User struct {
	handle     string
	fullName   string
	passHash   string
	reputation int // keep it simple for now. ideally it deserves own object with various attributes and stats
	asked      []*Question
	answered   []*Answer
}
