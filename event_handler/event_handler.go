package eventhandler

import rl "github.com/gen2brain/raylib-go/raylib"
import "tetris/entity"

const (
	LEFT  = 1
	RIGHT = 2
	DOWN  = 3
)

type UpdateEvent struct {
	MovingDirection int
	RotateDirection int
	GameState       int
}

func HandleEvent() UpdateEvent {
	updateEvent := UpdateEvent{
		RotateDirection: 0,
	}
	if rl.IsKeyPressed(rl.KeyA) {
		updateEvent.MovingDirection = LEFT
	} else if rl.IsKeyPressed(rl.KeyD) {
		updateEvent.MovingDirection = RIGHT
	} else if rl.IsKeyDown(rl.KeyS) {
		updateEvent.MovingDirection = DOWN
	}

	if rl.IsKeyPressed(rl.KeyR) {
		updateEvent.RotateDirection = entity.CLOCKWISE
	} else if rl.IsKeyPressed(rl.KeyL) {
		updateEvent.RotateDirection = entity.ANTI_CLOCKWISE
	}

	return updateEvent
}
