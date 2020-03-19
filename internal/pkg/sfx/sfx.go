package sfx

type ClientSoundEffect string

const (
	InventoryDragStart ClientSoundEffect = "INVENTORY_DRAG_START"
	InventoryDragStop  ClientSoundEffect = "INVENTORY_DRAG_STOP"
	PickupItem         ClientSoundEffect = "PICKUP_ITEM"
	SellBuyItem        ClientSoundEffect = "SELL_BUY_ITEM"
)
