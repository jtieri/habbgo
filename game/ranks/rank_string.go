// Code generated by "stringer -type=Rank"; DO NOT EDIT.

package ranks

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None-0]
	_ = x[Normal-1]
	_ = x[CommunityManager-2]
	_ = x[Guide-3]
	_ = x[Hobba-4]
	_ = x[SuperHobba-5]
	_ = x[Moderator-6]
	_ = x[Administrator-7]
}

const _Rank_name = "NoneNormalCommunityManagerGuideHobbaSuperHobbaModeratorAdministrator"

var _Rank_index = [...]uint8{0, 4, 10, 26, 31, 36, 46, 55, 68}

func (i Rank) String() string {
	if i < 0 || i >= Rank(len(_Rank_index)-1) {
		return "Rank(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Rank_name[_Rank_index[i]:_Rank_index[i+1]]
}
