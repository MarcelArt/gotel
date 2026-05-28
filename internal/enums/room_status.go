package enums

const (
	RoomVacant     = "VACANT"
	RoomOccupied   = "OCCUPIED"
	RoomDirty      = "DIRTY"
	RoomCleaning   = "CLEANING"
	RoomOutOfOrder = "OUT_OF_ORDER"
)

var RoomStatusList = []string{
	RoomVacant,
	RoomOccupied,
	RoomDirty,
	RoomCleaning,
	RoomOutOfOrder,
}
