package service

import (
	"fmt"
	"xApp/pkg/service/context"
)

func HandleCompareRES(RES []byte) bool {
	// Compare RES and XRES
	// if same, true. else false

	UEid := 1
	res, ok := context.GetFirstRESByUEid(UEid)
	if !ok {
		fmt.Println("No firstRES found for UEid:", UEid)
	} else {
		fmt.Println("FirstRES for UEid", UEid, ":", res)
	}
	fmt.Println("RES: ", RES)
	DeleteFirstres := context.DeleteFirstRESByUEid(UEid)
	if DeleteFirstres {
		fmt.Println("Delete First RES Success")
	} else {
		fmt.Println("Delete First RES Failed")
	}

	// Compare lengths of slices first
	if len(res) != len(RES) {
		// Slices are not equal if their lengths are different
		// Handle the case where slices have different lengths
		return false
	} else {
		// Compare each byte of the slices
		equal := true
		for i := 0; i < len(res); i++ {
			if res[i] != RES[i] {
				equal = false
				break
			}
		}
		// At this point, 'equal' is true if all bytes match, false otherwise
		if equal {
			// Slices are equal
			return true
		} else {
			// Slices are not equal
			return false
		}
	}
}

func HandleNORAAKACompareRES(RES []byte) bool {
	// Compare RES and XRES
	// if same, true. else false

	UEid := 1
	ue := &context.AmfUe{}
	NORAakaRES := ue.GetAUTN(UEid, 3)

	fmt.Println("RES: ", NORAakaRES)
	fmt.Println("RES from UE: ", RES)

	// Compare lengths of slices first
	if len(NORAakaRES) != len(RES) {
		// Slices are not equal if their lengths are different
		// Handle the case where slices have different lengths
		return false
	} else {
		// Compare each byte of the slices
		equal := true
		for i := 0; i < len(NORAakaRES); i++ {
			if NORAakaRES[i] != RES[i] {
				equal = false
				break
			}
		}
		// At this point, 'equal' is true if all bytes match, false otherwise
		if equal {
			// Slices are equal
			return true
		} else {
			// Slices are not equal
			return false
		}
	}
}

func HandleMessageSelection(octet []byte) ([]byte, []byte) {
	receivedBytes := octet
	DetectByte := receivedBytes[2:3]
	length := len(receivedBytes)

	// 0x56 means Authentication Request
	// 0x57 means Authentication Response
	fmt.Println("DetectByte: ", DetectByte[0])
	switch DetectByte[0] {
	case byte(0x56):
		if length < 10 {
			fmt.Println("Error: Insufficient bytes in receivedBytes.")
		}
		Header := receivedBytes[:7]
		firstRand := receivedBytes[7:23]
		firstAutn := receivedBytes[23:39]
		firstRES := receivedBytes[39:55]

		fmt.Println("Header: ", Header)
		fmt.Println("firstRand: ", firstRand)
		fmt.Println("firstAutn: ", firstAutn)
		fmt.Println("firstRES: ", firstRES)

		//count := 55
		//for i := 0; i <= 8; i++ {
		//	fmt.Println("AUTN:", i, receivedBytes[count:count+16])
		//	fmt.Println("RAND:", i, receivedBytes[count+16:count+32])
		//	fmt.Println("RES:", i, receivedBytes[count+32:count+48])
		//	count = count + 48
		//}

		RANDElementID := []byte{0x21}
		AUTNElementID := []byte{0x20}

		OriginalNASMessage := []byte{}
		OriginalNASMessage = append(OriginalNASMessage, Header...)
		OriginalNASMessage = append(OriginalNASMessage, RANDElementID...)
		OriginalNASMessage = append(OriginalNASMessage, firstRand...)
		OriginalNASMessage = append(OriginalNASMessage, AUTNElementID...)
		OriginalNASMessage = append(OriginalNASMessage, firstAutn...)
		OtherNASMessage := receivedBytes[55:487]

		// Create an initial UE information
		newUe := context.NewAmfUe(1, firstRES)
		context.StoreAmfUe(newUe)

		return OriginalNASMessage, OtherNASMessage
	case byte(0x57):
		if length < 10 {
			fmt.Println("Error: Insufficient bytes in receivedBytes.")
		}
		var OriginalNASMessage []byte
		Header := receivedBytes[:4]
		RES := receivedBytes[5:21]

		// Check if it is triggger NORA-AKA or not.
		UEid := 1
		checkStatus := context.CheckUserStatus(UEid)
		if checkStatus {
			ResultOfCompare := HandleCompareRES(RES)

			RESLength := []byte{0x01}
			OriginalNASMessage = append(OriginalNASMessage, Header...)
			OriginalNASMessage = append(OriginalNASMessage, RESLength...)

			if ResultOfCompare == true {
				CompareResultTrue := []byte{0x01}

				// Convert string to []byte
				OriginalNASMessage = append(OriginalNASMessage, CompareResultTrue...)
				//fmt.Println("ResultOfCompare: ", OriginalNASMessage)

				return OriginalNASMessage, nil
			} else {
				CompareResultFalse := []byte{0x00}

				// Convert string to []byte
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)
				//fmt.Println("ResultOfCompare: ", OriginalNASMessage)

				return OriginalNASMessage, nil
			}
		} else {
			// Handle NORA-AKA Authentication Response
			ResultOfCompare := HandleNORAAKACompareRES(RES)

			RESLength := []byte{0x01}
			OriginalNASMessage = append(OriginalNASMessage, Header...)
			OriginalNASMessage = append(OriginalNASMessage, RESLength...)

			if ResultOfCompare == true {
				CompareResultTrue := []byte{0x01}

				// Convert string to []byte
				OriginalNASMessage = append(OriginalNASMessage, CompareResultTrue...)
				//fmt.Println("ResultOfCompare: ", OriginalNASMessage)

				return OriginalNASMessage, nil
			} else {
				CompareResultFalse := []byte{0x00}

				// Convert string to []byte
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)
				//fmt.Println("ResultOfCompare: ", OriginalNASMessage)

				return OriginalNASMessage, nil
			}
		}

	case byte(0x41):
		// 0x41: Registration Request
		// In order to trigger second Authentication which means NORA-AKA

		// Trigger NORA-AKA
		var OriginalNASMessage []byte
		RANDElementID := []byte{0x21}
		AUTNElementID := []byte{0x20}

		if length < 10 {
			fmt.Println("Error: Insufficient bytes in receivedBytes.")
		}
		ue := &context.AmfUe{}
		UEid := 1

		NORAakaRAND := ue.GetAUTN(UEid, 1)
		NORAakaAUTN := ue.GetAUTN(UEid, 2)

		fmt.Println("NORA-AKA RAND: ", NORAakaRAND, " ,NORA-AKA AUTN: ", NORAakaAUTN)

		// Start to compose the Nora authentication packet.
		NORAheader := []byte{0x7e, 0x00, 0x56, 0x00, 0x02, 0x00, 0x00}
		OriginalNASMessage = append(OriginalNASMessage, NORAheader...)
		OriginalNASMessage = append(OriginalNASMessage, RANDElementID...)
		OriginalNASMessage = append(OriginalNASMessage, NORAakaRAND...)
		OriginalNASMessage = append(OriginalNASMessage, AUTNElementID...)
		OriginalNASMessage = append(OriginalNASMessage, NORAakaAUTN...)

		return OriginalNASMessage, nil
	default:

		return nil, nil
	}
}

func HandleOtherMessage(OtherMessage []byte) {
	Message := OtherMessage
	fmt.Println("Other Message: ", Message)

	//Separate and store the AUTN, RAND and RES
	ue := &context.AmfUe{}

	UEid := 1
	count := 0
	fmt.Println("AuthForRAND 1: ", Message[count:count+16])
	fmt.Println("AuthForAUTN 1: ", Message[count+16:count+32])
	fmt.Println("AuthForRES 1: ", Message[count+32:count+48])

	for counter := 0; counter <= 8; counter++ {
		ue.SetAuthParam(UEid, Message[count:count+16], 1, counter)
		ue.SetAuthParam(UEid, Message[count+16:count+32], 2, counter)
		ue.SetAuthParam(UEid, Message[count+32:count+48], 3, counter)
		count += 48
	}

	//fmt.Println("AuthForRAND: ", ue.GetAUTN(UEid, 1))
	//fmt.Println("AuthForAUTN: ", ue.GetAUTN(UEid, 2))
	//fmt.Println("AuthForRES: ", ue.GetAUTN(UEid, 3))
	//
	//fmt.Println("AuthForRAND: ", ue.GetAUTN(UEid, 1))
	//fmt.Println("AuthForAUTN: ", ue.GetAUTN(UEid, 2))
	//fmt.Println("AuthForRES: ", ue.GetAUTN(UEid, 3))

}
