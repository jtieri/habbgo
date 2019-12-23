package model

type Navigator struct {
	categories []*Category
}

type Category struct {
	id      int
	pid     int
	name    string
	minRank int
	public  bool
	trading bool
	node    bool
}
