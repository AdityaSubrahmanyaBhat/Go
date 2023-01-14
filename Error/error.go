package error

import (
	"fmt"
)

func HandleError(err error) {
	if err != nil {
		fmt.Printf("Error : %v",err)
	}
}