package item

import "strings"

//go:generate go run golang.org/x/tools/cmd/stringer -type=Behavior

// Behavior specifies a particular trait that an in game item may possess (e.g. can you sit on it, can you lay on it, etc).
type Behavior int

// The set of Behavior's that are supported.
const (
	Invalid Behavior = iota
	Solid
	CanStackOnTop
	CanSitOnTop
	CanStandOnTop
	CanLayOnTop
	CustomDataNumericOnOff
	RequiresTouchingForInteraction
	CustomDataTrueFalse
	PublicSpaceObject
	ExtraParameter
	Dice
	CustomDataOnOff
	CustomDataNumericState
	Teleporter
	RequiresRightsForInteraction
	Gate
	PrizeTrophy
	Roller
	Redeemable
	SoundMachine
	SoundMachineSampleSet
	Jukebox
	WallItem
	PostIt
	Decoration
	WheelOfFortune
	RoomDimmer
	Present
	Photo
	PlaceRollerOnTop
	Invisible
	Effect
	RedirectRotation0
	RedirectRotation2
	RedirectRotation4
	SongDisk
	PetWaterBowl
	PetFood
	PetCatFood
	PetDogFood
	PetCrocFood
	PrivateFurniture
)

type Behaviors []Behavior

// Contains checks the Behaviors and returns true if the specified Behavior is present.
func (b Behaviors) Contains(behavior Behavior) bool {
	for _, be := range b {
		if be == behavior {
			return true
		}
	}
	return false
}

// parseBehaviorData will parse a ',' delimited string representation of Behavior's and return a slice
// of the appropriate Behavior types.
func parseBehaviorData(data string) Behaviors {
	var behaviors []Behavior

	if data == "" {
		return nil
	}

	for _, b := range strings.Split(data, ",") {
		bType := behaviorFromString(b)
		behaviors = append(behaviors, bType)
	}

	return behaviors
}

// behaviorFromString will return the appropriate Behavior that corresponds to the specified string.
func behaviorFromString(behavior string) Behavior {
	switch strings.ToLower(behavior) {
	case "solid":
		return Solid
	case "can_stack_on_top":
		return CanStackOnTop
	case "can_sit_on_top":
		return CanSitOnTop
	case "can_stand_on_top":
		return CanStandOnTop
	case "can_lay_on_top":
		return CanLayOnTop
	case "custom_data_numeric_on_off":
		return CustomDataNumericOnOff
	case "requires_touching_for_interaction":
		return RequiresTouchingForInteraction
	case "custom_data_true_false":
		return CustomDataTrueFalse
	case "public_space_object":
		return PublicSpaceObject
	case "extra_parameter":
		return ExtraParameter
	case "dice":
		return Dice
	case "custom_data_on_off":
		return CustomDataOnOff
	case "custom_data_numeric_state":
		return CustomDataNumericState
	case "teleporter":
		return Teleporter
	case "requires_rights_for_interaction":
		return RequiresRightsForInteraction
	case "gate":
		return Gate
	case "prize_trophy":
		return PrizeTrophy
	case "roller":
		return Roller
	case "redeemable":
		return Redeemable
	case "sound_machine":
		return SoundMachine
	case "sound_machine_sample_set":
		return SoundMachineSampleSet
	case "jukebox":
		return Jukebox
	case "wall_item":
		return WallItem
	case "post_it":
		return PostIt
	case "decoration":
		return Decoration
	case "wheel_of_fortune":
		return WheelOfFortune
	case "room_dimmer":
		return RoomDimmer
	case "present":
		return Present
	case "photo":
		return Photo
	case "place_roller_on_top":
		return PlaceRollerOnTop
	case "invisible":
		return Invisible
	case "effect":
		return Effect
	case "redirect_rotation_0":
		return RedirectRotation0
	case "redirect_rotation_2":
		return RedirectRotation2
	case "redirect_rotation_4":
		return RedirectRotation4
	case "song_disk":
		return SongDisk
	case "pet_water_bowl":
		return PetWaterBowl
	case "pet_food":
		return PetFood
	case "pet_cat_food":
		return PetCatFood
	case "pet_dog_food":
		return PetDogFood
	case "pet_croc_food":
		return PetCrocFood
	case "private_furniture":
		return PrivateFurniture
	default:
		return Invalid
	}
}
