package badge

type Badge struct {
	Id   uint
	code string
}

func (b Badge) String() string {
	return b.code
}
