package context

type Subscriber struct {
	/* Ue Identity*/
	MSIN       string
	UEid       int
	RAN_ID     string
	AmfAssocID string
	firstRES   []byte
	RANDMap    map[int][]byte
	AUTNMap    map[int][]byte
	RESMap     map[int][]byte
	count      int
}

func NewSubscriber(MSIN string) *Subscriber {
	return &Subscriber{
		MSIN:       MSIN,
		UEid:       1,
		RAN_ID:     "",
		AmfAssocID: "",
		firstRES:   nil,
		RANDMap:    make(map[int][]byte), // Initialize the slice of byte slices
		AUTNMap:    make(map[int][]byte),
		RESMap:     make(map[int][]byte),
		count:      0,
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
