package armeria

type EquipmentSlot string

const (
	EquipSlotWallet EquipmentSlot = "wallet"
)

// ValidEquipmentSlots returns the valid slots for equippable items.
func ValidEquipmentSlots() []EquipmentSlot {
	return []EquipmentSlot{
		EquipSlotWallet,
	}
}

// ValidEquipmentSlotsAsString returns the valid slots for equippable items as strings.
func ValidEquipmentSlotsAsString() []string {
	s := make([]string, 0)
	for _, eq := range ValidEquipmentSlots() {
		s = append(s, string(eq))
	}
	return s
}

// EquipSlotMax returns the number of items that can be equipped to a given slot.
func EquipSlotMax(slot EquipmentSlot) int {
	switch slot {
	case EquipSlotWallet:
		return 3
	}

	return 1
}

// EquipSlotFormalName returns the formal name, with proper capitalization, for a given slot.
func EquipSlotFormalName(slot EquipmentSlot) string {
	switch slot {
	case EquipSlotWallet:
		return "Wallet"
	}

	return string(slot)
}
