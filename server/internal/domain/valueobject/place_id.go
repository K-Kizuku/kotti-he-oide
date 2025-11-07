package valueobject

import (
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// PlaceID は、ゲーム内の場所を表すValue Object
type PlaceID struct {
	value string
}

// 固定5箇所の場所ID
const (
	PlaceSpiralStairs  = "spiral_stairs"   // 螺旋階段を見上げる高い天井
	PlaceFireplace     = "fireplace"       // メインホールの暖炉のレンガ
	PlaceBackDoorHinge = "back_door_hinge" // 裏玄関の扉の蝶番
	PlaceEntranceDoor  = "entrance_door"   // 入口エントランスの扉
	PlacePianoRoom     = "piano_room"      // 階上応接室のピアノ
)

// validPlaces は、有効な場所IDのマップ
var validPlaces = map[string]bool{
	PlaceSpiralStairs:  true,
	PlaceFireplace:     true,
	PlaceBackDoorHinge: true,
	PlaceEntranceDoor:  true,
	PlacePianoRoom:     true,
}

// NewPlaceID は、文字列からPlaceIDを作成する
func NewPlaceID(id string) (PlaceID, error) {
	if !validPlaces[id] {
		return PlaceID{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"invalid place ID",
			nil,
		)
	}
	return PlaceID{value: id}, nil
}

// String は、PlaceIDを文字列として返す
func (p PlaceID) String() string {
	return p.value
}

// Equals は、2つのPlaceIDが等しいかどうかを判定する
func (p PlaceID) Equals(other PlaceID) bool {
	return p.value == other.value
}

// GetAllPlaces は、すべての有効な場所IDを返す
func GetAllPlaces() []string {
	return []string{
		PlaceSpiralStairs,
		PlaceFireplace,
		PlaceBackDoorHinge,
		PlaceEntranceDoor,
		PlacePianoRoom,
	}
}
