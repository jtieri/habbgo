package item

import (
	"strconv"
	"strings"
)

// Definition represents an Item's actual metadata.
type Definition struct {
	ID          int
	Sprite      string
	TopHeight   float64
	Length      int
	Width       int
	Color       string
	Behaviors   Behaviors
	Name        string
	Description string
	DrinkIDs    []int
	Tradeable   bool
	Recyclable  bool
	Interaction Interactioner
}

// DefaultItemDefinition returns the default Definition for items which contains a 1x1 size (width x length) and a
// default height of DefaultTopHeight.
func DefaultItemDefinition() *Definition {
	return &Definition{
		TopHeight: DefaultTopHeight,
		Length:    1,
		Width:     1,
	}
}

// NewItemDefinition will return a new Definition using the specified information.
func NewItemDefinition(
	ID, length, width int,
	sprite, name, desc, behaviorData, interaction, color, drinkIDData string,
	topHeight float64,
	tradable, recyclable bool,
) Definition {

	// Don't let the top height be 0 so that the collision map can be built correctly.
	// We want the height to be taller than default room tile height.
	if topHeight == 0 {
		topHeight = DefaultTopHeight
	}

	// Split the string representation of the drink ID list and build the slice of ints.
	drinkIDStrings := strings.Split(drinkIDData, ",")
	drinkIDs := make([]int, len(drinkIDStrings))
	if len(drinkIDStrings) > 0 {
		for _, id := range drinkIDStrings {
			drinkID, _ := strconv.Atoi(id)
			drinkIDs = append(drinkIDs, drinkID)
		}
	}

	itemDef := Definition{
		ID:          ID,
		Sprite:      sprite,
		TopHeight:   topHeight,
		Length:      length,
		Width:       width,
		Color:       color,
		Behaviors:   parseBehaviorData(behaviorData),
		Name:        name,
		Description: desc,
		DrinkIDs:    drinkIDs,
		Tradeable:   tradable,
		Recyclable:  recyclable,
		Interaction: interactionByName(interaction),
	}

	// Ensure the item can be walked into if it is a gate by setting the top height to 0.
	if !itemDef.Behaviors.Contains(CanSitOnTop) && !itemDef.Behaviors.Contains(CanStackOnTop) && !itemDef.Behaviors.Contains(CanLayOnTop) {
		itemDef.TopHeight = 0
	}

	return itemDef
}

// ContainsBehavior returns true if the Definition's Behaviors contains the specified Behavior.
func (d *Definition) ContainsBehavior(behavior Behavior) bool {
	return d.Behaviors.Contains(behavior)
}

// AddBehavior will add a Behavior to Definition's Behaviors if it is not present already.
func (d *Definition) AddBehavior(behavior Behavior) {
	if !d.Behaviors.Contains(behavior) {
		d.Behaviors = append(d.Behaviors, behavior)
	}
}

// Icon creates the catalogue icon.
func (d *Definition) Icon(spriteID int) string {
	icon := d.Sprite
	if spriteID > 0 {
		icon += " " + strconv.Itoa(spriteID)
	}
	return icon
}
