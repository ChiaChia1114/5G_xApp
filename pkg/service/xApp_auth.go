package service

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/milenage"
	"github.com/free5gc/util/ueauth"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

type xAppAuthenticationParameterRAND struct {
	Iei   int
	Octet []uint8
}

type AuthenticationVector struct {
	AvType   string `json:"avType" yaml:"avType" bson:"avType" mapstructure:"AvType"`
	Rand     string `json:"rand" yaml:"rand" bson:"rand" mapstructure:"Rand"`
	Xres     string `json:"xres" yaml:"xres" bson:"xres" mapstructure:"Xres"`
	Autn     string `json:"autn" yaml:"autn" bson:"autn" mapstructure:"Autn"`
	XresStar string `json:"xresStar" yaml:"xresStar" bson:"xresStar" mapstructure:"XresStar"`
}

const (
	keyStrLen int = 32
	opcStrLen int = 32
)

func (r xAppAuthenticationParameterRAND) SetxAppRANDValue(rANDValue []uint8) {
	copy(r.Octet[0:16], rANDValue[:])
}

func strictHex(s string, n int) string {
	l := len(s)
	if l < n {
		return fmt.Sprintf(strings.Repeat("0", n-l) + s)
	} else {
		return s[l-n : l]
	}
}

func XAppAKAGenerateAUTH(OpcValue string, kValue string) (response *AuthenticationVector, result bool) {
	var authInfoRequest models.AuthenticationInfoRequest
	authInfoRequest.ServingNetworkName = "5G:mnc093.mcc208.3gppnetwork.org"
	authInfoRequest.ResynchronizationInfo = nil
	authInfoRequest.SupportedFeatures = ""
	fmt.Println("In GenerateAuthDataProcedure")

	rand.Seed(time.Now().UnixNano())
	//var err error

	/*
		K, RAND, CK, IK: 128 bits (16 bytes) (hex len = 32)
		SQN, AK: 48 bits (6 bytes) (hex len = 12) TS33.102 - 6.3.2
		AMF: 16 bits (2 bytes) (hex len = 4) TS33.102 - Annex H
	*/

	hasOPC := false
	var kStr, opcStr string
	var k, opc []byte
	var err error

	//------------------------ Terry Modify Start --------------------------//
	//	Goals: Generate a OUT-X nas packet.                                 //
	//  Method:                                                             //
	//     1. Connect to mongo DB                                           //
	//     2. Get the basic op & K* in the mongo DB                         //
	//     3. Generate the parameter with Authentication                    //
	//----------------------------------------------------------------------//

	kStr = kValue
	if len(kStr) == keyStrLen {
		k, err = hex.DecodeString(kStr)
		if err != nil {
			fmt.Println("err:", err)
		}
	}

	opcStr = OpcValue
	if len(opcStr) == opcStrLen {
		opc, err = hex.DecodeString(opcStr)
		if err != nil {
			fmt.Println("err:", err)
		} else {
			hasOPC = true
		}
	}
	if !hasOPC {
		response = nil
		return nil, false
	}

	RAND := make([]byte, 16)
	_, err = cryptoRand.Read(RAND)
	if err != nil {
		fmt.Println("err:", err)
		return nil, false
	}

	sqnStr := "16f3b3f70ff2"
	sqn, err := hex.DecodeString(sqnStr)
	if err != nil {
		fmt.Println("err:", err)
		return nil, false
	}
	//fmt.Println("sqnStr: ", sqnStr)
	//fmt.Println("K=[%x], sqn=[%x], OPC=[%x]", k, sqn, opc)

	hexString := "8000"
	//AMF, err := hex.DecodeString(authSubs.AuthenticationManagementField)
	AMF, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println("err:", err)
		return nil, false
	}
	//fmt.Println("RAND=[%x], AMF=[%x]", RAND, AMF)

	// increment sqn
	bigSQN := big.NewInt(0)
	sqn, err = hex.DecodeString(sqnStr)
	if err != nil {
		fmt.Println("err:", err)
		return nil, false
	}

	bigSQN.SetString(sqnStr, 16)

	bigInc := big.NewInt(1)
	bigSQN = bigInc.Add(bigSQN, bigInc)

	SQNheStr := fmt.Sprintf("%x", bigSQN)
	SQNheStr = strictHex(SQNheStr, 12)

	// Run milenage
	macA, macS := make([]byte, 8), make([]byte, 8)
	CK, IK := make([]byte, 16), make([]byte, 16)
	RES := make([]byte, 8)
	AK, AKstar := make([]byte, 6), make([]byte, 6)

	// Generate macA, macS
	err = milenage.F1(opc, k, RAND, sqn, AMF, macA, macS)
	if err != nil {
		fmt.Println("milenage F1 err:", err)
	}

	// Generate RES, CK, IK, AK, AKstar
	// RES == XRES (expected RES) for server
	err = milenage.F2345(opc, k, RAND, RES, CK, IK, AK, AKstar)
	if err != nil {
		fmt.Println("milenage F2345 err:", err)
	}
	//fmt.Println("milenage RES=[%s]", hex.EncodeToString(RES))

	// Generate AUTN
	//fmt.Println("SQN=[%x], AK=[%x]", sqn, AK)
	//fmt.Println("AMF=[%x], macA=[%x]", AMF, macA)
	SQNxorAK := make([]byte, 6)
	for i := 0; i < len(sqn); i++ {
		SQNxorAK[i] = sqn[i] ^ AK[i]
	}
	//fmt.Println("SQN xor AK=[%x]", SQNxorAK)
	AUTN := append(append(SQNxorAK, AMF...), macA...)
	//fmt.Println("AUTN=[%x]", AUTN)

	var av AuthenticationVector

	// derive XRES*
	key := append(CK, IK...)
	FC := ueauth.FC_FOR_RES_STAR_XRES_STAR_DERIVATION
	P0 := []byte(authInfoRequest.ServingNetworkName)
	P1 := RAND
	P2 := RES

	kdfValForXresStar, err := ueauth.GetKDFValue(
		key, FC, P0, ueauth.KDFLen(P0), P1, ueauth.KDFLen(P1), P2, ueauth.KDFLen(P2))
	if err != nil {
		fmt.Println("Get kdfValForXresStar err: %+v", err)
	}
	xresStar := kdfValForXresStar[len(kdfValForXresStar)/2:]
	//fmt.Println("xresStar=[%x]", xresStar)

	// derive Kausf
	FC = ueauth.FC_FOR_KAUSF_DERIVATION
	P0 = []byte(authInfoRequest.ServingNetworkName)
	P1 = SQNxorAK

	// Fill in rand, xresStar, autn, kausf
	av.Rand = hex.EncodeToString(RAND)
	av.XresStar = hex.EncodeToString(xresStar)
	av.Autn = hex.EncodeToString(AUTN)
	av.AvType = "5G_HE_AKA"

	response = &av
	return response, true
}
