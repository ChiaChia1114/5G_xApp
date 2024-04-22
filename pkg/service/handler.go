package service

import (
	"encoding/hex"
	"fmt"
	"time"
	"xApp/pkg/service/context"
)

func HandleCompareRES(RES []byte) bool {
	// Compare RES and XRES
	// if same, true. else false

	UEid := 1
	res, ok := context.GetRESValueByUEid(UEid)
	if !ok {
		fmt.Println("No RES found for UEid:", UEid)
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

	UEid := 1
	NORAakaRES, GetRESResult := context.GetRESValueByUEid(UEid)
	if !GetRESResult {
		fmt.Println("Error for getting RES from NORA-AKA procedure.")
	}

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

func HandleMessageSelection(octet []byte) []byte {
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

		// Separate opc and k
		opcValue := receivedBytes[7:39]
		kValue := receivedBytes[39:71]

		// Generate the res with the function of xApp_auth
		opcStr := string(opcValue)
		kStr := string(kValue)

		av, result := XAppAKAGenerateAUTH(opcStr, kStr)
		if !result {
			fmt.Println("Error for generate the Authentication vector.")
		}

		RANDhexString := av.Rand
		RANDnewBytes, err := hex.DecodeString(RANDhexString)
		if err != nil {
			fmt.Println("Error decoding hex string:", err)
		}

		AutnhexString := av.Autn
		AutnnewBytes, err := hex.DecodeString(AutnhexString)
		if err != nil {
			fmt.Println("Error decoding hex string:", err)
		}

		XREStarthexString := av.XresStar
		XREStartnexBytes, err := hex.DecodeString(XREStarthexString)
		if err != nil {
			fmt.Println("Error decoding hex string:", err)
		}

		RANDElementID := []byte{0x21}
		AUTNElementID := []byte{0x20}

		// Generate the NAS packet and send back to CU
		OriginalNASMessage := []byte{}
		OriginalNASMessage = append(OriginalNASMessage, Header...)
		OriginalNASMessage = append(OriginalNASMessage, RANDElementID...)
		OriginalNASMessage = append(OriginalNASMessage, RANDnewBytes...)
		OriginalNASMessage = append(OriginalNASMessage, AUTNElementID...)
		OriginalNASMessage = append(OriginalNASMessage, AutnnewBytes...)

		// Create an initial UE information
		newUe := context.NewAmfUe(1, opcStr, kStr, XREStartnexBytes)
		context.StoreAmfUe(newUe)

		return OriginalNASMessage
	case byte(0x57):
		if length < 10 {
			fmt.Println("Error: Insufficient bytes in receivedBytes.")
		}
		var OriginalNASMessage []byte
		Header := receivedBytes[:4]
		RES := receivedBytes[5:21]

		UEid := 1

		RESLength := []byte{0x01}
		OriginalNASMessage = append(OriginalNASMessage, Header...)
		OriginalNASMessage = append(OriginalNASMessage, RESLength...)

		// Check the status of the UE
		// Check if it is triggger NORA-AKA or not.
		checkresult := context.GetStatusByUEid(UEid)
		if !checkresult {
			ResultOfCompare := HandleCompareRES(RES)
			context.SetStatusByUEid(UEid)

			if ResultOfCompare == true {
				CompareResultTrue := []byte{0x01}

				// Convert string to []byte
				OriginalNASMessage = append(OriginalNASMessage, CompareResultTrue...)

				endTime = time.Now()
				serviceTime := endTime.Sub(startTime)
				fmt.Println("First Authentication Transmission time: %v", serviceTime)

				return OriginalNASMessage
			} else {
				CompareResultFalse := []byte{0x00}

				// Convert string to []byte
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)

				endTime = time.Now()
				serviceTime := endTime.Sub(startTime)
				fmt.Println("First Authentication Transmission time: %v", serviceTime)

				return OriginalNASMessage
			}
		} else {
			//Handle NORA-AKA Authentication Response
			NORAResultOfCompare := HandleNORAAKACompareRES(RES)

			if NORAResultOfCompare == true {
				CompareResultTrue := []byte{0x01}

				// Convert string to []byte
				OriginalNASMessage = append(OriginalNASMessage, CompareResultTrue...)

				endTime = time.Now()
				serviceTime := endTime.Sub(startTime)
				fmt.Println("NORA-AKA Transmission time: %v", serviceTime)

				return OriginalNASMessage
			} else {
				CompareResultFalse := []byte{0x00}

				// Convert string to []byte
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)

				endTime = time.Now()
				serviceTime := endTime.Sub(startTime)
				fmt.Println("NORA-AKA Transmission time: %v", serviceTime)

				return OriginalNASMessage
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

		UEid := 1
		opcValue, GetOpcResult := context.GetOpcValueByUEid(UEid)
		if !GetOpcResult {
			fmt.Println("Error for empty opcValue")
		}
		kValue, GetkResult := context.GetkValueByUEid(UEid)
		if !GetkResult {
			fmt.Println("Error for empty kValue")
		}

		av, result := XAppAKAGenerateAUTH(opcValue, kValue)
		if !result {
			fmt.Println("Error for generate the Authentication vector.")
		}

		RANDhexString := av.Rand
		RANDnewBytes, err := hex.DecodeString(RANDhexString)
		if err != nil {
			fmt.Println("Error decoding hex string:", err)
		}

		AutnhexString := av.Autn
		AutnnewBytes, err := hex.DecodeString(AutnhexString)
		if err != nil {
			fmt.Println("Error decoding hex string:", err)
		}

		XREStarthexString := av.XresStar
		XREStartnexBytes, err := hex.DecodeString(XREStarthexString)
		if err != nil {
			fmt.Println("Error decoding hex string:", err)
		}

		// Start to compose the Nora authentication packet.
		NORAheader := []byte{0x7e, 0x00, 0x56, 0x00, 0x02, 0x00, 0x00}
		OriginalNASMessage = append(OriginalNASMessage, NORAheader...)
		OriginalNASMessage = append(OriginalNASMessage, RANDElementID...)
		OriginalNASMessage = append(OriginalNASMessage, RANDnewBytes...)
		OriginalNASMessage = append(OriginalNASMessage, AUTNElementID...)
		OriginalNASMessage = append(OriginalNASMessage, AutnnewBytes...)

		context.SetRESByUEid(UEid, XREStartnexBytes)

		return OriginalNASMessage
	default:

		return nil
	}
}
