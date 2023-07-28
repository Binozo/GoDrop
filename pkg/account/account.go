package account

import "github.com/Binozo/GoDrop/internal/utils"

// AppleAccount doesn't really do anything.
// Just use the default one
// AppleAccount is being used for BLE advertising
type AppleAccount struct {
	ID    string
	Email string
	Phone string
}

var Default = AppleAccount{
	ID:    "support@apple.com",
	Email: "support@apple.com",
	Phone: "11111111111",
}

func (account *AppleAccount) BuildManufacturerData() []byte {
	id := utils.GetFirstTwoBytesFromSha256(account.ID)
	email := utils.GetFirstTwoBytesFromSha256(account.Email)
	phone := utils.GetFirstTwoBytesFromSha256(account.Phone)

	return []byte{0x5, 0x12, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1,
		id[0], id[1], phone[0], phone[1], email[0], email[1], email[0], email[1], 0x0}
}
