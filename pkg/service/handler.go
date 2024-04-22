package service

import (
	"fmt"
	"time"
	"xApp/pkg/service/context"
)

func HandleCompareRES(RES []byte) bool {
	// Compare RES and XRES
	// if same, true. else false

	UEid := 1
	res, ok := context.GetFirstRESByUEid(UEid)
	if !ok {
		fmt.Println("No firstRES found for UEid:", UEid)
	}
	fmt.Println("RES: ", RES)
	DeleteFirstres := context.DeleteFirstRESByUEid(UEid)
	if !DeleteFirstres {
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
	var startTime time.Time
	var endTime time.Time
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
		startTime = time.Now()
		Header := receivedBytes[:7]
		firstRand := receivedBytes[7:23]
		firstAutn := receivedBytes[23:39]
		firstRES := receivedBytes[39:55]

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
				OriginalNASMessage = append(OriginalNASMessage, CompareResultTrue...)

				endTime = time.Now()
				serviceTime := endTime.Sub(startTime)
				fmt.Println("First Authentication transmission time: %v", serviceTime)

				return OriginalNASMessage, nil
			} else {
				CompareResultFalse := []byte{0x00}
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)

				endTime = time.Now()
				serviceTime := endTime.Sub(startTime)
				fmt.Println("First Authentication transmission time: %v", serviceTime)

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
				OriginalNASMessage = append(OriginalNASMessage, CompareResultTrue...)

				endTime = time.Now()
				serviceTime := endTime.Sub(startTime)
				fmt.Println("NORA-AKA transmission time: %v", serviceTime)

				return OriginalNASMessage, nil
			} else {
				CompareResultFalse := []byte{0x00}
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)

				endTime = time.Now()
				serviceTime := endTime.Sub(startTime)
				fmt.Println("NORA-AKA transmission time: %v", serviceTime)

				return OriginalNASMessage, nil
			}
		}

	case byte(0x41):
		// 0x41: Registration Request
		// In order to trigger second Authentication which means NORA-AKA
		// Trigger NORA-AKA

		var OriginalNASMessage []byte
		startTime = time.Now()
		RANDElementID := []byte{0x21}
		AUTNElementID := []byte{0x20}
		ue := &context.AmfUe{}
		UEid := 1
		NORAakaRAND := ue.GetAUTN(UEid, 1)
		NORAakaAUTN := ue.GetAUTN(UEid, 2)

		if length < 10 {
			fmt.Println("Error: Insufficient bytes in receivedBytes.")
		}

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
	//Separate and store the AUTN, RAND and RES
	Message := OtherMessage
	ue := &context.AmfUe{}
	UEid := 1
	count := 0

	for counter := 0; counter <= 8; counter++ {
		ue.SetAuthParam(UEid, Message[count:count+16], 1, counter)
		ue.SetAuthParam(UEid, Message[count+16:count+32], 2, counter)
		ue.SetAuthParam(UEid, Message[count+32:count+48], 3, counter)
		count += 48
	}

}
