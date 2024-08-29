package spawner

import (
	"math/rand"
	"tetris/entity"
)

type BlockSpawner struct {
	MaxWidth   int
	Randomizer rand.Rand
}

func (bs BlockSpawner) Spawn() (entity.BlockEntity, error) {
	randomBlock := bs.Randomizer.Intn(entity.Z)
	randomXCoordinate := bs.Randomizer.Intn(bs.MaxWidth)
	randomColor := bs.Randomizer.Intn(entity.GREEN)

	for _, location := range entity.BLOCK_OCCUPYING_LOCATION[randomBlock] {
		if randomXCoordinate+int(location[0]) > bs.MaxWidth {
			excessLocation := randomXCoordinate + int(location[0]) - bs.MaxWidth
			randomXCoordinate -= excessLocation
		}

		if randomXCoordinate+int(location[0]) < 0 {
			excessLocation := randomXCoordinate + int(location[0])
			randomXCoordinate -= excessLocation
		}
	}

	newEntity, err := entity.New(randomBlock, randomColor, [2]int{randomXCoordinate, 0})

	if err != nil {
		return entity.BlockEntity{}, err
	}

	return newEntity, nil
}
