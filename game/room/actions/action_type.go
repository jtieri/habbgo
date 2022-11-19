package actions

type ActionType int

const (
	Move ActionType = iota
	Sit
	Lay
	Dance
	Swim
	Wave
	Gesture
	Talk
	PlayerSleep
	Trade
	Sign
	Dead
	Jump
	PetSleep
	Eat
	CarryItem
	CarryDrink
	CarryFood
	UseItem
	UseDrink
	UseFood
	FlatControl
)

/*
   MOVE("mv"),
   SIT("sit"),
   LAY("lay"),
   FLAT_CONTROL("flatctrl"),
   DANCE("dance"),
   SWIM("swim"),
   CARRY_ITEM("cri"),
   CARRY_DRINK("carryd"),
   CARRY_FOOD("carryf"),
   USE_ITEM("usei"),
   USE_FOOD("eat"),
   USE_DRINK("drink"),
   WAVE("wave"),
   GESTURE("gest"),
   TALK("talk"),
   AVATAR_SLEEP("Sleep"),
   TRADE("trd"),
   SIGN("sign"),
   DEAD("ded"),
   JUMP("jmp"),
   PET_SLEEP("slp"),
   EAT("eat");
*/
