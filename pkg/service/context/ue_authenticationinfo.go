package context

type AmfUe struct {
	/* Ue Identity*/
	UEid     int
	OpcValue string
	kValue   string
	RES      []byte
	Active   bool
}

func NewAmfUe(UEid int, opcValue string, kValue string, resResponse []byte) *AmfUe {
	return &AmfUe{
		UEid:     UEid,
		OpcValue: opcValue,
		kValue:   kValue,
		RES:      resResponse,
		Active:   false,
		//firstRES: firstRES,
		//RANDMap:  make(map[int][]byte), // Initialize the slice of byte slices
		//AUTNMap:  make(map[int][]byte),
		//RESMap:   make(map[int][]byte),
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

////Delete FirstRES returns the booling value corresponding to the given UEid.
//func DeleteFirstRESByUEid(UEid int) bool {
//	ue, ok := AmfUeMap[UEid]
//	if ok {
//		ue.firstRES = nil
//		StoreAmfUe(ue)
//		CheckResult := CheckUserStatus(UEid)
//		if CheckResult {
//			// Delete First failed
//			return false
//		} else {
//			// Delete First success
//			return true
//		}
//	}
//	return false
//}

//// AddByteSlice adds a new byte slice to the AmfUe instance.
//func (ue *AmfUe) SetAuthParam(UEid int, res []byte, AuthType int, counter int) {
//	if _, ok := AmfUeMap[UEid]; !ok {
//		// Create a new AmfUe object if the UEid doesn't exist
//		AmfUeMap[UEid] = NewAmfUe(UEid, nil)
//	}
//	switch AuthType {
//	case 1:
//		AmfUeMap[UEid].RANDMap[counter] = res
//		//fmt.Println("Set ", counter, " Auth RAND Parameter Success.")
//		//fmt.Println("Print ", counter, " RAND: ", res)
//	case 2:
//		AmfUeMap[UEid].AUTNMap[counter] = res
//		//fmt.Println("Set ", counter, " Auth AUTN Parameter Success.")
//		//fmt.Println("Print ", counter, " AUTN: ", res)
//	case 3:
//		AmfUeMap[UEid].RESMap[counter] = res
//		//fmt.Println("Set ", counter, " Auth RES Parameter Success.")
//		//fmt.Println("Print ", counter, " RES: ", res)
//	default:
//		fmt.Println("Set Auth Parameter failed.")
//	}
//}

//func (ue *AmfUe) GetAUTN(UEid int, AuthType int) []byte {
//	// Check if the UEid exists in the map
//	if _, ok := AmfUeMap[UEid]; ok {
//		var AuthResponse []byte
//		switch AuthType {
//		case 1:
//			for i := 0; i <= 8; i++ {
//				if AmfUeMap[UEid].RANDMap[i] != nil {
//					AuthResponse = AmfUeMap[UEid].RANDMap[i]
//					ue.SetAuthParam(UEid, nil, 1, i)
//					//fmt.Println("Get Auth RAND Parameter Success.")
//					return AuthResponse
//				}
//			}
//			return nil
//		case 2:
//			for i := 0; i <= 8; i++ {
//				if AmfUeMap[UEid].AUTNMap[i] != nil {
//					AuthResponse = AmfUeMap[UEid].AUTNMap[i]
//					ue.SetAuthParam(UEid, nil, 2, i)
//					//fmt.Println("Get Auth AUTN Parameter Success.")
//					return AuthResponse
//				}
//			}
//			return nil
//		case 3:
//			for i := 0; i <= 8; i++ {
//				if AmfUeMap[UEid].RESMap[i] != nil {
//					AuthResponse = AmfUeMap[UEid].RESMap[i]
//					ue.SetAuthParam(UEid, nil, 3, i)
//					//fmt.Println("Get Auth RES Parameter Success.")
//					return AuthResponse
//				}
//			}
//			return nil
//		default:
//			fmt.Println("Get Auth Parameter failed.")
//		}
//
//	}
//	return nil // Return nil if UEid not found
//}
