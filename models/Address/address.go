package address

import (
	"encoding/json"
)

type Address struct{
	City string
	Country string
	State string
	Pincode json.Number
}