package model

type Navigator struct {
	Categories []Category
}

type Category struct {
	Id      int
	Pid     int
	Node    bool
	Name    string
	Public  bool
	Trading bool
	MinRank int
}
