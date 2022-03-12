package offchain

import (
	"fmt"
	"github.com/Nik-U/pbc"
	"github.com/W1llyu/ourjson"
)

type PublicParam struct {
	pairing pbc.Pairing
	g       pbc.Element
	u       pbc.Element
	v       pbc.Element
	w       pbc.Element
}

func (PublicParam) PublicParam1(sharedParams string, sharedG []byte) {
	params := pbc.GenerateA(160, 512)
	pairing := params.NewPairing()
	fmt.Println(pairing)
	pairing1, _ := pbc.NewPairingFromString(sharedParams) // loads pairing parameters from a string and instantiates a pairing
	fmt.Println(sharedG)
	x := pairing1.NewG2().SetBytes(sharedG)
	//g := pairing.NewG2().SetBytes(sharedG)
	fmt.Println(x)
}
func NewPublicParam() PublicParam {
	pairing, _ := pbc.NewPairingFromString(Utils{}.readFile("../config/a.properties")) // loads pairing parameters from a string and instantiates a pairing
	jsonStringPub := Utils{}.readFile("../config/pubkey")

	pubkeyObject, _ := ourjson.ParseObject(jsonStringPub)
	gk, _ := pubkeyObject.GetString("g")

	g := pairing.NewG2().SetBytes(Utils{}.base64StringToElementBytes(gk))

	uk, _ := pubkeyObject.GetString("u")
	u := pairing.NewG1().SetBytes(Utils{}.base64StringToElementBytes(uk))

	vk, _ := pubkeyObject.GetString("v")
	v := pairing.NewG2().SetBytes(Utils{}.base64StringToElementBytes(vk))

	wk, _ := pubkeyObject.GetString("w")
	w := pairing.NewG1().SetBytes(Utils{}.base64StringToElementBytes(wk))

	return PublicParam{pairing: *pairing, g: *g, u: *u, v: *v, w: *w}

}

func (PublicParam) PublicParam() PublicParam {
	pairing, _ := pbc.NewPairingFromString(Utils{}.readFile("../config/a.properties")) // loads pairing parameters from a string and instantiates a pairing
	jsonStringPub := Utils{}.readFile("../config/pubkey")

	pubkeyObject, _ := ourjson.ParseObject(jsonStringPub)
	gk, _ := pubkeyObject.GetString("g")
	g := pairing.NewG2().SetBytes(Utils{}.base64StringToElementBytes(gk))
	//fmt.Println(g)
	uk, _ := pubkeyObject.GetString("u")
	u := pairing.NewG1().SetBytes([]byte(uk))
	vk, _ := pubkeyObject.GetString("v")
	v := pairing.NewG2().SetBytes(Utils{}.base64StringToElementBytes(vk))
	//fmt.Println(g)
	wk, _ := pubkeyObject.GetString("w")
	w := pairing.NewG1().SetBytes([]byte(wk))

	return PublicParam{pairing: *pairing, g: *g, u: *u, v: *v, w: *w}

	//g := pairing.NewG2().SetBytes(sharedG)

	/*
		var data map[string]interface{}
		err := json.Unmarshal([]byte(Utils{}.readFile("../config/pubkey")), &data)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Printf("map字典: %v\n", data)
		fmt.Printf("count值:%v\n", data["g"])
		fmt.Printf(reflect.TypeOf([]byte(data["g"])))
		g := pairing.NewG2().SetBytes([]byte(data("g")))
		fmt.Println(g)
	*/
}
func (PublicParam) test() {
	fmt.Println("helloworld")

	var params = pbc.GenerateA(160, 512)

	var pairing = params.NewPairing()

	// Initialize group elements. pbc automatically handles garbage collection.
	var g = pairing.NewG1()
	var h = pairing.NewG2()
	var x = pairing.NewGT()

	// Generate random group elements and pair them
	g.Rand()
	h.Rand()
	fmt.Printf("g = %s\n", g)
	fmt.Printf("h = %s\n", h)
	x.Pair(g, h)
	fmt.Printf("e(g,h) = %s\n", x)

}
