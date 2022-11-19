package ranks

//go:generate go run golang.org/x/tools/cmd/stringer -type=Rank

type Rank int

const (
	None Rank = iota
	Normal
	CommunityManager
	Guide
	Hobba
	SuperHobba
	Moderator
	Administrator
)

func Ranks() []Rank {
	return []Rank{
		None,
		Normal,
		CommunityManager,
		Guide,
		Hobba,
		SuperHobba,
		Moderator,
		Administrator,
	}
}
