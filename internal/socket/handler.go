package socket

import (
	"fmt"
)

func AppendOctet(octetString []byte, octet byte) []byte {
	fmt.Println("Received OctetString:", octetString)

	return append(octetString, octet)
}
