package sfx

type ClientSoundEffect string

const (
	InventoryDragStart ClientSoundEffect = "INVENTORY_DRAG_START"
	InventoryDragStop  ClientSoundEffect = "INVENTORY_DRAG_STOP"
	PickupItem         ClientSoundEffect = "PICKUP_ITEM"
	SellBuyItem        ClientSoundEffect = "SELL_BUY_ITEM"
	CatMeow            ClientSoundEffect = "CAT_MEOW"
	Teleport           ClientSoundEffect = "TELEPORT"
)

// List returns a slice containing all of the valid sound effects.
func List() []string {
	return []string{
		string(InventoryDragStart),
		string(InventoryDragStop),
		string(PickupItem),
		string(SellBuyItem),
		string(CatMeow),
		string(Teleport),
	}
}
