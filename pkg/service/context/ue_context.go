package context

type Subscriber struct {
	/* Ue Identity*/
	MSIN     string
	OpcValue string
	kValue   string
	RES      []byte
	Active   bool
}

func NewSubscriber(MSIN string) *Subscriber {
	return &Subscriber{
		MSIN:     MSIN,
		OpcValue: "",
		kValue:   "",
		RES:      nil,
		Active:   false,
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

// GetOpcValueByMSIN returns the firstRES value corresponding to the given UEid.
func GetOpcValueByMSIN(MSIN string) string {
	ue, ok := SubscriberMap[MSIN]
	if !ok {
		return ""
	}
	return ue.OpcValue
}

// GetkValueByMSIN returns the firstRES value corresponding to the given UEid.
func GetkValueByMSIN(MSIN string) string {
	ue, ok := SubscriberMap[MSIN]
	if !ok {
		return ""
	}
	return ue.kValue
}

// GetRESValueByMSIN returns the firstRES value corresponding to the given UEid.
func GetRESValueByMSIN(MSIN string) []byte {
	ue, ok := SubscriberMap[MSIN]
	if !ok {
		return nil
	}
	return ue.RES
}

// SetRESByUEid returns the firstRES value corresponding to the given UEid.
func SetRESByMSIN(MSIN string, res []byte) bool {
	subscriber, ok := SubscriberMap[MSIN]
	if !ok {
		return false
	}
	subscriber.RES = res
	StoreSubscriber(subscriber)
	return true
}

// SetOpcValueByUEid returns the firstRES value corresponding to the given UEid.
func SetOpcValueByMSIN(MSIN string, OPcStr string) bool {
	subscriber, ok := SubscriberMap[MSIN]
	if !ok {
		return false
	}
	subscriber.OpcValue = OPcStr
	StoreSubscriber(subscriber)
	return true
}

// SetkValueByUEid returns the firstRES value corresponding to the given UEid.
func SetkValueByMSIN(MSIN string, Kstr string) bool {
	subscriber, ok := SubscriberMap[MSIN]
	if !ok {
		return false
	}
	subscriber.kValue = Kstr
	StoreSubscriber(subscriber)
	return true
}

func SetSubscriberActiveValue(MSIN string) bool {
	subscriber, ok := SubscriberMap[MSIN]
	if !ok {
		return false
	}
	subscriber.Active = true
	StoreSubscriber(subscriber)
	return subscriber.Active
}

func GetSubscriberActiveValue(MSIN string) bool {
	subscriber, ok := SubscriberMap[MSIN]
	if !ok {
		return false
	}
	return subscriber.Active
}
