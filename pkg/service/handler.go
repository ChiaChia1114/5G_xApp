package service

import (
	"encoding/hex"
	"fmt"
	"time"
	"xApp/pkg/service/context"
	filer "xApp/pkg/service/exportfile"
	Authtimer "xApp/pkg/service/timer"
)

var SubscriberMSIN string = ""

func HandleCompareRES(RES []byte) bool {
	// Compare RES and XRES
	// if same, true. else false

	subscriber := &context.Subscriber{}
	res := subscriber.GetSubscriberFirstRES(SubscriberMSIN)
	if res == nil {
		fmt.Println("No firstRES found for UE")
	}
	DeleteFirstres := subscriber.SubscriberDeleteFirstRES(SubscriberMSIN)
	if !DeleteFirstres {
		fmt.Println("Delete First RES Failed")
	}

	// Encryption with the xApp token
	xAppToken := context.GlobalToken
	// XOR the two byte slices
	NewRES, err := context.XorBytes(RES, xAppToken)
	if err != nil {
		fmt.Println("Error:", err)
	}
	RES = NewRES

	// fmt.Println("XRES: ", res)
	// fmt.Println("RES: ", RES)

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

	subscriber := &context.Subscriber{}
	NORAakaRES := subscriber.GetAuthenticationVectors(SubscriberMSIN, 3)

	// Encryption with the xApp token
	xAppToken := context.GlobalToken
	// XOR the two byte slices
	NewNORAakaRES, err := context.XorBytes(NORAakaRES, xAppToken)
	NewRES, err := context.XorBytes(RES, xAppToken)
	if err != nil {
		fmt.Println("Error:", err)
	}
	RES = NewRES
	NORAakaRES = NewNORAakaRES

	// fmt.Println("XRES: ", NORAakaRES)
	// fmt.Println("RES: ", RES)

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
	var OriginalNASMessage []byte
	subscriber := &context.Subscriber{}

	receivedBytes := octet
	DetectByte := receivedBytes[2:3]
	length := len(receivedBytes)
	if length < 10 {
		fmt.Println("Error: Insufficient bytes in receivedBytes.")
	}

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
		// Handle Authentication Request

		startTime := time.Now()
		ResultofSetTimer := Authtimer.SetStartTime(1, startTime)
		if !ResultofSetTimer {
			fmt.Println("Set Start time failed.")
		}

		Header := receivedBytes[:7]
		firstRand := receivedBytes[7:23]
		firstAutn := receivedBytes[23:39]
		firstRES := receivedBytes[39:55]

		// Encryption with the xApp token
		xAppToken := context.GlobalToken
		// XOR the two byte slices
		NewXRES, err := context.XorBytes(firstRES, xAppToken)
		if err != nil {
			fmt.Println("Error:", err)
		}
		firstRES = NewXRES

		RANDElementID := []byte{0x21}
		AUTNElementID := []byte{0x20}

		OriginalNASMessage = append(OriginalNASMessage, Header...)
		OriginalNASMessage = append(OriginalNASMessage, RANDElementID...)
		OriginalNASMessage = append(OriginalNASMessage, firstRand...)
		OriginalNASMessage = append(OriginalNASMessage, AUTNElementID...)
		OriginalNASMessage = append(OriginalNASMessage, firstAutn...)
		OtherNASMessage := receivedBytes[55:487]

		// Store frist RES
		StoreFirstRESResult := subscriber.SetSubscriberFirstRES(SubscriberMSIN, firstRES)
		if StoreFirstRESResult == false {
			fmt.Println("Store First RES Failed")
		}

		fmt.Println("OtherNASMessage:", OtherNASMessage)

		return OriginalNASMessage, OtherNASMessage
	case byte(0x57):
		// Handle Authentication Response

		Header := receivedBytes[:4]
		RES := receivedBytes[5:21]

		// Check if it is triggger NORA-AKA or not.

		checkStatus := context.GetSubscriberCount(SubscriberMSIN)
		if checkStatus == 0 {
			ResultOfCompare := HandleCompareRES(RES)

			RESLength := []byte{0x01}
			OriginalNASMessage = append(OriginalNASMessage, Header...)
			OriginalNASMessage = append(OriginalNASMessage, RESLength...)

			if ResultOfCompare == true {
				CompareResultTrue := []byte{0x01}
				OriginalNASMessage = append(OriginalNASMessage, CompareResultTrue...)

				endTime := time.Now()
				err := filer.ReadTimeFromFile(1, endTime)
				if err != nil {
					fmt.Println(err)
				}

				AuthCount := context.GetSubscriberCount(SubscriberMSIN)
				if AuthCount == 9 {
					// Delete the UE MAP for the creation
					context.DeleteSubscriber(SubscriberMSIN)
					CheckUE := context.CheckSubscriberStatus(SubscriberMSIN)
					if !CheckUE {
						fmt.Println("Delete UE status context success.")
					}
				} else {
					// UE count plus 1
					context.SubscriberCountPlus(SubscriberMSIN)
				}

				if AuthCount == 0 {
					fmt.Println("Create a new UE map")
				}

				return OriginalNASMessage, nil
			} else {
				CompareResultFalse := []byte{0x00}
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)

				endTime := time.Now()
				err := filer.ReadTimeFromFile(1, endTime)
				if err != nil {
					fmt.Println(err)
				}

				AuthCount := context.GetSubscriberCount(SubscriberMSIN)
				// fmt.Println("Count: ", AuthCount)
				if AuthCount == 9 {
					// Delete the UE MAP for the creation
					context.DeleteSubscriber(SubscriberMSIN)
					CheckUE := context.CheckSubscriberStatus(SubscriberMSIN)
					if !CheckUE {
						fmt.Println("Delete UE status context success.")
					}
				} else {
					// UE count plus 1
					context.SubscriberCountPlus(SubscriberMSIN)
				}

				if AuthCount == 0 {
					fmt.Println("Create a new UE map")
				}

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

				endTime := time.Now()
				err := filer.ReadTimeFromFile(2, endTime)
				if err != nil {
					fmt.Println(err)
				}

				AuthCount := context.GetSubscriberCount(SubscriberMSIN)
				fmt.Println("Count: ", AuthCount)
				if AuthCount == 9 {
					// Delete the UE MAP for the creation
					context.DeleteSubscriber(SubscriberMSIN)
					CheckUE := context.CheckSubscriberStatus(SubscriberMSIN)
					if !CheckUE {
						fmt.Println("Delete UE status context success.")
					}
				} else {
					// UE count plus 1
					context.SubscriberCountPlus(SubscriberMSIN)
				}

				if AuthCount == 0 {
					fmt.Println("Create a new UE map")
				}

				return OriginalNASMessage, nil
			} else {
				CompareResultFalse := []byte{0x00}
				OriginalNASMessage = append(OriginalNASMessage, CompareResultFalse...)

				endTime := time.Now()
				err := filer.ReadTimeFromFile(2, endTime)
				if err != nil {
					fmt.Println(err)
				}

				AuthCount := context.GetSubscriberCount(SubscriberMSIN)
				fmt.Println("Count: ", AuthCount)
				if AuthCount == 9 {
					// Delete the UE MAP for the creation
					context.DeleteSubscriber(SubscriberMSIN)
					CheckUE := context.CheckSubscriberStatus(SubscriberMSIN)
					if !CheckUE {
						fmt.Println("Delete UE status context success.")
					}
				} else {
					// UE count plus 1
					context.SubscriberCountPlus(SubscriberMSIN)
				}

				if AuthCount == 0 {
					fmt.Println("Create a new UE map")
				}

				return OriginalNASMessage, nil
			}
		}

	case byte(0x41):
		// 0x41: Registration Request
		// Check UE status to trigger NORA-AKA or not.

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

				return OriginalNASMessage, nil
			} else {
				// Still Success with create subscriber
				CreateSubscriberSuccess := []byte{0x01}
				OriginalNASMessage = append(OriginalNASMessage, CreateSubscriberSuccess...)
				return OriginalNASMessage, nil
			}
		} else {
			// Trigger NORA-AKA
			startTime := time.Now()
			ResultofSetTimer := Authtimer.SetStartTime(1, startTime)
			if !ResultofSetTimer {
				fmt.Println("Set Start time failed.")
			}

			RANDElementID := []byte{0x21}
			AUTNElementID := []byte{0x20}
			NORAakaRAND := subscriber.GetAuthenticationVectors(SubscriberMSIN, 1)
			NORAakaAUTN := subscriber.GetAuthenticationVectors(SubscriberMSIN, 2)

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
		}
	default:

		return nil, nil
	}
}

func HandleOtherMessage(OtherMessage []byte) {
	//Separate and store the AUTN, RAND and RES
	Message := OtherMessage
	subscriber := &context.Subscriber{}
	count := 0

	for counter := 0; counter <= 9; counter++ {
		subscriber.SetAuthenticationVectors(SubscriberMSIN, Message[count:count+16], 1, counter)
		subscriber.SetAuthenticationVectors(SubscriberMSIN, Message[count+16:count+32], 2, counter)
		subscriber.SetAuthenticationVectors(SubscriberMSIN, Message[count+32:count+48], 3, counter)
		count += 48
	}

}
