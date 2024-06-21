package context

import (
	"encoding/hex"
	"fmt"
)

type Subscriber struct {
	/* Ue Identity*/
	MSIN     string
	UEid     int
	firstRES []byte
	RANDMap  map[int][]byte
	AUTNMap  map[int][]byte
	RESMap   map[int][]byte
	count    int
}

func NewSubscriber(MSIN string) *Subscriber {
	return &Subscriber{
		MSIN:     MSIN,
		UEid:     1,
		firstRES: nil,
		RANDMap:  make(map[int][]byte), // Initialize the slice of byte slices
		AUTNMap:  make(map[int][]byte),
		RESMap:   make(map[int][]byte),
		count:    0,
	}
}

// AmfUeMap maps UEid to corresponding AmfUe instances.
var SubscriberMap map[string]*Subscriber

func init() {
	SubscriberMap = make(map[string]*Subscriber)
}

// StoreSubscriber stores the Subscriber instance in the SubscriberMap.
func StoreSubscriber(ue *Subscriber) {
	SubscriberMap[ue.MSIN] = ue
}

func GetSubscriberActive(MSIN string) bool {
	_, ok := SubscriberMap[MSIN]
	if !ok {
		return false
	}
	return true
}

// AddByteSlice adds a new byte slice to the AmfUe instance.
func (subscriber *Subscriber) SetAuthenticationVectors(MSIN string, res []byte, AuthType int, counter int) {
	if _, ok := SubscriberMap[MSIN]; !ok {
		// Create a new AmfUe object if the UEid doesn't exist
		SubscriberMap[MSIN] = NewSubscriber(MSIN)
	}
	switch AuthType {
	case 1:
		SubscriberMap[MSIN].RANDMap[counter] = res
	case 2:
		SubscriberMap[MSIN].AUTNMap[counter] = res
	case 3:
		SubscriberMap[MSIN].RESMap[counter] = res
	default:
		fmt.Println("Set Auth Parameter failed.")
	}
}

func (subscriber *Subscriber) GetSubscriberMSIN(MSIN string) string {
	if _, ok := SubscriberMap[MSIN]; !ok {
		fmt.Println("Can't find the SubscriberMap, RanID: ", MSIN)
	}
	subscirberMSIN := SubscriberMap[MSIN].MSIN

	return subscirberMSIN
}

func (subscriber *Subscriber) GetSubscriberCount(MSIN string) int {
	if _, ok := SubscriberMap[MSIN]; !ok {
		fmt.Println("Can't find the SubscriberMap, MSIN: ", MSIN)
	}
	subscirberCount := SubscriberMap[MSIN].count

	return subscirberCount
}

func (subscriber *Subscriber) GetAuthenticationVectors(MSIN string, AuthType int) []byte {
	// Check if the UEid exists in the map
	if _, ok := SubscriberMap[MSIN]; ok {
		var AuthResponse []byte
		switch AuthType {
		case 1:
			for i := 0; i <= 8; i++ {
				if SubscriberMap[MSIN].RANDMap[i] != nil {
					AuthResponse = SubscriberMap[MSIN].RANDMap[i]
					subscriber.SetAuthenticationVectors(MSIN, nil, 1, i)
					//fmt.Println("Get Auth RAND Parameter Success.")
					return AuthResponse
				}
			}
			return nil
		case 2:
			for i := 0; i <= 8; i++ {
				if SubscriberMap[MSIN].AUTNMap[i] != nil {
					AuthResponse = SubscriberMap[MSIN].AUTNMap[i]
					subscriber.SetAuthenticationVectors(MSIN, nil, 2, i)
					//fmt.Println("Get Auth AUTN Parameter Success.")
					return AuthResponse
				}
			}
			return nil
		case 3:
			for i := 0; i <= 8; i++ {
				if SubscriberMap[MSIN].RESMap[i] != nil {
					AuthResponse = SubscriberMap[MSIN].RESMap[i]
					subscriber.SetAuthenticationVectors(MSIN, nil, 3, i)
					//fmt.Println("Get Auth RES Parameter Success.")
					return AuthResponse
				}
			}
			return nil
		default:
			fmt.Println("Get Auth Parameter failed.")
		}

	}
	return nil // Return nil if UEid not found
}

func (subscriber *Subscriber) SetSubscriberFirstRES(MSIN string, FirstRES []byte) bool {
	if _, ok := SubscriberMap[MSIN]; !ok {
		fmt.Println("Can't find the SubscriberMap, MSIN: ", MSIN)
	}
	SubscriberMap[MSIN].firstRES = FirstRES

	OctetstringCheckFirstRES := subscriber.GetSubscriberFirstRES(MSIN)
	CheckFirstRES := hex.EncodeToString(OctetstringCheckFirstRES)
	TrueFirstRES := hex.EncodeToString(FirstRES)

	if CheckFirstRES != TrueFirstRES {
		return false
	}

	return true
}

func (subscriber *Subscriber) GetSubscriberFirstRES(MSIN string) []byte {
	if _, ok := SubscriberMap[MSIN]; !ok {
		fmt.Println("Can't find the SubscriberMap, MSIN: ", MSIN)
	}
	FristRES := SubscriberMap[MSIN].firstRES

	return FristRES
}

func (subscriber *Subscriber) SubscriberDeleteFirstRES(MSIN string) bool {
	subscribers, ok := SubscriberMap[MSIN]
	if ok {
		subscribers.firstRES = nil
		StoreSubscriber(subscribers)
		SubscriberCheckFirstRESresult := subscriber.GetSubscriberFirstRES(MSIN)
		if SubscriberCheckFirstRESresult == nil {
			return true
		}
	}
	return false
}

//Check Subscriber  Status return the booling value.
func CheckSubscriberStatus(MSIN string) bool {
	_, ok := SubscriberMap[MSIN]
	if !ok {
		return false
	}
	return true
}

func GetSubscriberCount(MSIN string) int {
	subscriber, ok := SubscriberMap[MSIN]
	if ok {
		Count := subscriber.count
		return Count
	}
	return 0
}

func SubscriberGetCount(MSIN string) int {
	subscriber, ok := SubscriberMap[MSIN]
	if ok {
		Count := subscriber.count
		return Count
	}
	return 0
}

func SubscriberCountPlus(MSIN string) bool {
	subscriber, ok := SubscriberMap[MSIN]
	if ok {
		Count := subscriber.count
		Count++
		subscriber.count = Count
		StoreSubscriber(subscriber)
		CheckResult := SubscriberGetCount(MSIN)
		if CheckResult != Count {
			return false
		}
	}
	return true
}

func DeleteSubscriber(MSIN string) {
	delete(SubscriberMap, MSIN)
}
