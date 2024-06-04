package context

import (
	"fmt"
)

type AmfUe struct {
	/* Ue Identity*/
	UEid     int
	OpcValue string
	kValue   string
	RES      []byte
	Active   bool
	count    int
	RANDmap  map[int]string
	AUTNmap  map[int]string
	XRESmap  map[int]string
}

func NewAmfUe(UEid int, opcValue string, kValue string, resResponse []byte) *AmfUe {
	return &AmfUe{
		UEid:     UEid,
		OpcValue: opcValue,
		kValue:   kValue,
		RES:      resResponse,
		Active:   false,
		count:    0,
		RANDmap:  make(map[int]string),
		AUTNmap:  make(map[int]string),
		XRESmap:  make(map[int]string),
	}
}

// AmfUeMap maps UEid to corresponding AmfUe instances.
var AmfUeMap map[int]*AmfUe

func init() {
	AmfUeMap = make(map[int]*AmfUe)
}

// StoreAmfUe stores the AmfUe instance in the AmfUeMap.
func StoreAmfUe(ue *AmfUe) {
	AmfUeMap[ue.UEid] = ue
}

// GetOpcValueByUEid returns the firstRES value corresponding to the given UEid.
func GetOpcValueByUEid(UEid int) (string, bool) {
	ue, ok := AmfUeMap[UEid]
	if !ok {
		return "", false
	}
	return ue.OpcValue, true
}

// GetkValueByUEid returns the firstRES value corresponding to the given UEid.
func GetkValueByUEid(UEid int) (string, bool) {
	ue, ok := AmfUeMap[UEid]
	if !ok {
		return "", false
	}
	return ue.kValue, true
}

// GetRESValueByUEid returns the firstRES value corresponding to the given UEid.
func GetRESValueByUEid(UEid int) ([]byte, bool) {
	ue, ok := AmfUeMap[UEid]
	if !ok {
		return nil, false
	}
	return ue.RES, true
}

// GetStatusByUEid returns the firstRES value corresponding to the given UEid.
func GetStatusByUEid(UEid int) bool {
	ue, ok := AmfUeMap[UEid]
	if !ok {
		return false
	}
	return ue.Active
}

// SetStatusByUEid returns the firstRES value corresponding to the given UEid.
func SetStatusByUEid(UEid int) bool {
	ue, ok := AmfUeMap[UEid]
	if !ok {
		return false
	}
	ue.Active = true
	StoreAmfUe(ue)
	CheckResult := GetStatusByUEid(UEid)
	return CheckResult
}

// SetRESByUEid returns the firstRES value corresponding to the given UEid.
func SetRESByUEid(UEid int, res []byte) bool {
	ue, ok := AmfUeMap[UEid]
	if !ok {
		return false
	}
	ue.RES = res
	StoreAmfUe(ue)
	return true
}

// SetCountByUEid
func SetCountByUEid(UEid int, count int) bool {
	ue, ok := AmfUeMap[UEid]
	if !ok {
		return false
	}
	ue.count = count
	StoreAmfUe(ue)
	return true
}

// GetCountByUEid
func GetCountByUEid(UEid int) (int, bool) {
	ue, ok := AmfUeMap[UEid]
	if !ok {
		return 0, false
	}
	return ue.count, true
}

// AddByteSlice adds a new byte slice to the AmfUe instance.
func (ue *AmfUe) SetAuthParam(UEid int, value string, AuthType int, counter int) {
	if _, ok := AmfUeMap[UEid]; !ok {
		// Create a new AmfUe object if the UEid doesn't exist
		fmt.Println("Error for set the Auth parameters.")
	}
	switch AuthType {
	case 1:
		AmfUeMap[UEid].RANDmap[counter] = value
		//fmt.Println("Set ", counter, " Auth RAND Parameter Success.")
		//fmt.Println("Print ", counter, " RAND: ", res)
	case 2:
		AmfUeMap[UEid].AUTNmap[counter] = value
		//fmt.Println("Set ", counter, " Auth AUTN Parameter Success.")
		//fmt.Println("Print ", counter, " AUTN: ", res)
	case 3:
		AmfUeMap[UEid].XRESmap[counter] = value
		//fmt.Println("Set ", counter, " Auth RES Parameter Success.")
		//fmt.Println("Print ", counter, " RES: ", res)
	default:
		fmt.Println("Set Auth Parameter failed.")
	}
}

func (ue *AmfUe) GetRAND(UEid int, count int) string {
	if _, ok := AmfUeMap[UEid]; ok {
		var AuthResponse string
		if AmfUeMap[UEid].RANDmap[count] != "" {
			AuthResponse = AmfUeMap[UEid].RANDmap[count]
			return AuthResponse
		}
	}
	return ""
}

func (ue *AmfUe) GetAUTN(UEid int, count int) string {
	if _, ok := AmfUeMap[UEid]; ok {
		var AuthResponse string
		if AmfUeMap[UEid].AUTNmap[count] != "" {
			AuthResponse = AmfUeMap[UEid].AUTNmap[count]
			return AuthResponse
		}
	}
	return ""
}

func (ue *AmfUe) GetXRES(UEid int, count int) string {
	if _, ok := AmfUeMap[UEid]; ok {
		var AuthResponse string
		if AmfUeMap[UEid].XRESmap[count] != "" {
			AuthResponse = AmfUeMap[UEid].XRESmap[count]
			return AuthResponse
		}
	}
	return ""
}
