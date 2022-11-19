package position

import (
	"fmt"
	"math"
)

// Position represents an in game position which contains an (x,y,z) coordinate set along with,
// the body and head rotation directions for the entity occupying the position.
type Position struct {
	X            int
	Y            int
	Z            float64
	BodyRotation int
	HeadRotation int
}

// Equals returns true if the calling Position has the same (x,y) coordinates as the specified Position.
func (p *Position) Equals(pos Position) bool {
	return p.X == pos.X && p.Y == pos.Y
}

// Add returns a new Position that is the sum of the (x,y,z) coordinates from the calling Position's p and pos.
func (p *Position) Add(pos Position) Position {
	return Position{
		X: p.X + pos.X,
		Y: p.Y + pos.Y,
		Z: p.Z + pos.Z,
	}
}

// Squared returns an integer representing the distance squared between the Position's p and pos.
func (p *Position) Squared(pos Position) int {
	var deltaX, deltaY float64
	deltaX = float64(p.X - pos.X)
	deltaY = float64(p.Y - pos.Y)
	return int(math.Sqrt(math.Pow(deltaX, 2) + math.Pow(deltaY, 2)))
}

// Touches returns true if the Position p is touching the Position pos.
func (p *Position) Touches(pos Position) bool {
	return p.Squared(pos) <= 1
}

// PositionInFront returns the Position that is in front of the calling Position.
func (p *Position) PositionInFront() Position {
	pos := Position{
		X:            p.X,
		Y:            p.Y,
		Z:            p.Z,
		BodyRotation: p.BodyRotation,
		HeadRotation: p.HeadRotation,
	}

	switch {
	case p.BodyRotation == 0:
		pos.Y -= 1
	case p.BodyRotation == 1:
		pos.X += 1
		pos.Y -= 1
	case p.BodyRotation == 2:
		pos.X += 1
	case p.BodyRotation == 3:
		pos.X += 1
		pos.Y += 1
	case p.BodyRotation == 4:
		pos.Y += 1
	case p.BodyRotation == 5:
		pos.X -= 1
		pos.Y += 1
	case p.BodyRotation == 6:
		pos.X -= 1
	case p.BodyRotation == 7:
		pos.X -= 1
		pos.Y -= 1
	}

	return pos
}

// PositionInBack returns the Position that is behind the calling Position.
func (p *Position) PositionInBack() Position {
	pos := Position{
		X:            p.X,
		Y:            p.Y,
		Z:            p.Z,
		BodyRotation: p.BodyRotation,
		HeadRotation: p.HeadRotation,
	}

	switch {
	case p.BodyRotation == 0:
		pos.Y += 1
	case p.BodyRotation == 1:
		pos.X -= 1
		pos.Y += 1
	case p.BodyRotation == 2:
		pos.X -= 1
	case p.BodyRotation == 3:
		pos.X -= 1
		pos.Y -= 1
	case p.BodyRotation == 4:
		pos.Y -= 1
	case p.BodyRotation == 5:
		pos.X += 1
		pos.Y -= 1
	case p.BodyRotation == 6:
		pos.X += 1
	case p.BodyRotation == 7:
		pos.X += 1
		pos.Y += 1
	}

	return pos
}

// PositionToLeft returns the Position to the left of the calling Position.
func (p *Position) PositionToLeft() Position {
	pos := Position{
		X:            p.X,
		Y:            p.Y,
		Z:            p.Z,
		BodyRotation: p.BodyRotation,
		HeadRotation: p.HeadRotation,
	}

	switch {
	case p.BodyRotation == 0:
		pos.X -= 1
	case p.BodyRotation == 1:
		pos.X -= 1
		pos.Y -= 1
	case p.BodyRotation == 2:
		pos.Y -= 1
	case p.BodyRotation == 3:
		pos.X += 1
		pos.Y -= 1
	case p.BodyRotation == 4:
		pos.X += 1
	case p.BodyRotation == 5:
		pos.X += 1
		pos.Y += 1
	case p.BodyRotation == 6:
		pos.Y += 1
	case p.BodyRotation == 7:
		pos.X -= 1
		pos.Y += 1
	}

	return pos
}

// PositionToRight returns the Position to the right of the calling Position.
func (p *Position) PositionToRight() Position {
	pos := Position{
		X:            p.X,
		Y:            p.Y,
		Z:            p.Z,
		BodyRotation: p.BodyRotation,
		HeadRotation: p.HeadRotation,
	}

	switch {
	case p.BodyRotation == 0:
		pos.X += 1
	case p.BodyRotation == 1:
		pos.X += 1
		pos.Y += 1
	case p.BodyRotation == 2:
		pos.Y += 1
	case p.BodyRotation == 3:
		pos.X -= 1
		pos.Y += 1
	case p.BodyRotation == 4:
		pos.X -= 1
	case p.BodyRotation == 5:
		pos.X -= 1
		pos.Y -= 1
	case p.BodyRotation == 6:
		pos.Y -= 1
	case p.BodyRotation == 7:
		pos.X += 1
		pos.Y -= 1
	}

	return pos
}

// String returns a string representing the (x,y) coordinates of the calling Position.
// The returned string will be in the format [x, y].
func (p *Position) String() string {
	return fmt.Sprintf("[%d, %d]", p.X, p.Y)
}

// CalculateHeadRotation returns the direction that a player needs to rotate to face the specified Position relative to
// the players current Position.
func (p *Position) CalculateHeadRotation(lookToPos Position) int {
	distance := p.BodyRotation - p.PositionDirection(lookToPos)

	if p.BodyRotation%2 == 0 {
		if distance < 0 {
			return p.BodyRotation + 1
		} else if distance > 0 {
			return p.BodyRotation - 1
		}
	}

	return p.BodyRotation
}

// PositionDirection returns the direction of a Position relative to the calling Position.
func (p *Position) PositionDirection(lookToPos Position) int {
	switch {
	case p.X < lookToPos.X && p.Y > lookToPos.Y:
		return 1
	case p.X < lookToPos.X:
		return 2
	case p.X < lookToPos.X && p.Y < lookToPos.Y:
		return 3
	case p.Y < lookToPos.Y:
		return 4
	case p.X > lookToPos.X && p.Y < lookToPos.Y:
		return 5
	case p.X > lookToPos.X:
		return 6
	case p.X > lookToPos.X && p.Y > lookToPos.Y:
		return 7
	}

	return 0 // We should never be falling through the switch statement to end up here.
}

// WalkDirection returns the direction that a player needs to rotate to walk to a specified Position.
func WalkDirection(currentPos, goToPos Position) int {
	switch {
	case currentPos.X == goToPos.X:
		if currentPos.Y < goToPos.Y {
			return 4
		}

		return 0
	case currentPos.X > goToPos.X:
		if currentPos.Y == goToPos.Y {
			return 6
		}

		if currentPos.Y < goToPos.Y {
			return 5
		}

		return 7
	default:
		if currentPos.Y == goToPos.Y {
			return 2
		}

		if currentPos.Y < goToPos.Y {
			return 3
		}

		return 1
	}
}
