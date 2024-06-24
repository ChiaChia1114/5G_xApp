package service

import (
	"encoding/hex"
	"fmt"
	"xApp/pkg/service/context"
)

var SubscriberMSIN string = ""

func HandleCompareRES(RES []byte) bool {
	// Compare RES and XRES
	// if same, true. else false

	res := context.GetRESValueByMSIN(SubscriberMSIN)
	if res == nil {
		fmt.Println("No RES found for MSIN: ", SubscriberMSIN)
	}

	xAppToken := context.GlobalToken
	NewRES, err := context.XorBytes(RES, xAppToken)
	if err != nil {
		fmt.Println("Error:", err)
	}
	RES = NewRES
	fmt.Println("RES: ", RES)
	fmt.Println("XRES: ", res)

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

	NORAakaRES := context.GetRESValueByMSIN(SubscriberMSIN)
	if NORAakaRES == nil {
		fmt.Println("Error for getting RES from NORA-AKA procedure. MSIN: ", SubscriberMSIN)
	}

	xAppToken := context.GlobalToken
	NewRES, err := context.XorBytes(RES, xAppToken)
	if err != nil {
		fmt.Println("Error:", err)
	}
	RES = NewRES
	fmt.Println("RES: ", RES)
	fmt.Println("XRES: ", NORAakaRES)

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
	receivedBytes := octet
	DetectByte := receivedBytes[2:3]
	length := len(receivedBytes)

	// 0x56 means Authentication Request
	// 0x57 means Authentication Response
	if DetectByte[0] == 86 {
		fmt.Println("NORA-Box received Authentication Request")
	} else if DetectByte[0] == 87 {
		fmt.Println("NORA-Box received Authentication Response")
	} else if DetectByte[0] == 65 {
		fmt.Println("NORA-Box received Registration Request")
	} else {
		fmt.Println("NORA-Box received unknown message")
	}

	switch DetectByte[0] {
	case byte(0x56):
		if length < 10 {
			fmt.Println("Error: Insufficient bytes in receivedBytes.")
		}

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

		fmt.Println("T-RAND: ", RANDnewBytes)
		fmt.Println("T-AUTN: ", AutnnewBytes)
		fmt.Println("T-XRES: ", XREStartnexBytes)

		RANDElementID := []byte{0x21}
		AUTNElementID := []byte{0x20}

		// Generate the NAS packet and send back to CU
		OriginalNASMessage := []byte{}
		OriginalNASMessage = append(OriginalNASMessage, Header...)
		OriginalNASMessage = append(OriginalNASMessage, RANDElementID...)
		OriginalNASMessage = append(OriginalNASMessage, RANDnewBytes...)
		OriginalNASMessage = append(OriginalNASMessage, AUTNElementID...)
		OriginalNASMessage = append(OriginalNASMessage, AutnnewBytes...)

		// Store Subscriber's information

		context.SetOpcValueByMSIN(SubscriberMSIN, opcStr)
		context.SetkValueByMSIN(SubscriberMSIN, kStr)
		context.SetRESByMSIN(SubscriberMSIN, XREStartnexBytes)

		return OriginalNASMessage
	case byte(0x57):
		if length < 10 {
			fmt.Println("Error: Insufficient bytes in receivedBytes.")
		}
		var OriginalNASMessage []byte
		Header := receivedBytes[:4]
		RES := receivedBytes[5:21]

		RESLength := []byte{0x01}
		OriginalNASMessage = append(OriginalNASMessage, Header...)
		OriginalNASMessage = append(OriginalNASMessage, RESLength...)

		// Check the status of the UE
		// Check if it is triggger NORA-AKA or not.
		checkresult := context.GetSubscriberActiveValue(SubscriberMSIN)
		if !checkresult {
			ResultOfCompare := HandleCompareRES(RES)

			if ResultOfCompare == true {
				CompareResultTrue := []byte{0x01}
				OriginalNASMessage = append(OriginalNASMessage, CompareResultTrue...)
				context.SetSubscriberActiveValue(SubscriberMSIN)

				return OriginalNASMessage
			} else {
				CompareResultFalse := []byte{0x00}
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)

				return OriginalNASMessage
			}
		} else {
			//Handle NORA-AKA Authentication Response
			NORAResultOfCompare := HandleNORAAKACompareRES(RES)

			if NORAResultOfCompare == true {
				CompareResultTrue := []byte{0x01}
				OriginalNASMessage = append(OriginalNASMessage, CompareResultTrue...)

				return OriginalNASMessage
			} else {
				CompareResultFalse := []byte{0x00}
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)

				return OriginalNASMessage
			}
		}

	case byte(0x41):
		// 0x41: Registration Request
		// In order to trigger second Authentication which means NORA-AKA

		// Trigger NORA-AKA
		var OriginalNASMessage []byte

		OctetstringMSIN := receivedBytes[14:19]
		MSIN := hex.EncodeToString(OctetstringMSIN)
		SubscriberMSIN = MSIN
		fmt.Println("UE MSIN: ", MSIN)

		if !context.GetSubscriberActive(MSIN) {
			// Create an initial UE information
			newSubscriber := context.NewSubscriber(MSIN)
			context.StoreSubscriber(newSubscriber)

			if !context.GetSubscriberActive(MSIN) {
				// Still fault with create subscriber
				CreateSubscriberfalied := []byte{0x00}
				OriginalNASMessage = append(OriginalNASMessage, CreateSubscriberfalied...)

				return OriginalNASMessage
			} else {
				// Still Success with create subscriber
				CreateSubscriberSuccess := []byte{0x01}
				OriginalNASMessage = append(OriginalNASMessage, CreateSubscriberSuccess...)
				return OriginalNASMessage
			}
		} else {
			// Trigger NORA-AKA
			RANDElementID := []byte{0x21}
			AUTNElementID := []byte{0x20}

			if length < 10 {
				fmt.Println("Error: Insufficient bytes in receivedBytes.")
			}

			opcValue := context.GetOpcValueByMSIN(MSIN)
			if opcValue == "" {
				fmt.Println("Error for empty opcValue")
			}
			kValue := context.GetkValueByMSIN(MSIN)
			if kValue == "" {
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

			fmt.Println("T-RAND: ", RANDnewBytes)
			fmt.Println("T-AUTN: ", AutnnewBytes)
			fmt.Println("T-XRES: ", XREStartnexBytes)

			// Start to compose the Nora authentication packet.
			NORAheader := []byte{0x7e, 0x00, 0x56, 0x00, 0x02, 0x00, 0x00}
			OriginalNASMessage = append(OriginalNASMessage, NORAheader...)
			OriginalNASMessage = append(OriginalNASMessage, RANDElementID...)
			OriginalNASMessage = append(OriginalNASMessage, RANDnewBytes...)
			OriginalNASMessage = append(OriginalNASMessage, AUTNElementID...)
			OriginalNASMessage = append(OriginalNASMessage, AutnnewBytes...)

			context.SetRESByMSIN(MSIN, XREStartnexBytes)

			return OriginalNASMessage
		}

	default:

		return nil
	}
}
