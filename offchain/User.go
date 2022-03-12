package offchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Nik-U/pbc"
	"github.com/W1llyu/ourjson"
	"os"
)

type User struct {
	username string
	x        pbc.Element
	keys     PublicParam
}

func newUser(username string) User {
	//params := pbc.GenerateA(160, 512) //rBit是Zr阶数，qBit是G1阶数
	//writeFile("../config/a.properties", params.String())

	pairing, _ := pbc.NewPairingFromString(readFile("../config/a.properties")) // loads pairing parameters from a string and instantiates a pairing
	x := pairing.NewZr().Rand()
	file, err := os.Stat("../config/pubkey")
	if err != nil {
		fmt.Println(err)
	}
	var g, u, v, w *pbc.Element

	if file != nil {

		jsonStringPub := readFile("../config/pubkey")
		//fmt.Println(jsonStringPub)
		pubkeyObject, _ := ourjson.ParseObject(jsonStringPub)
		gk, _ := pubkeyObject.GetString("g")
		g = pairing.NewG2().SetBytes(base64StringToElementBytes(gk))
		uk, _ := pubkeyObject.GetString("u")
		u = pairing.NewG1().SetBytes(base64StringToElementBytes(uk))
		vk, _ := pubkeyObject.GetString("v")
		v = pairing.NewG2().SetBytes(base64StringToElementBytes(vk))
		//fmt.Println(g)
		wk, _ := pubkeyObject.GetString("w")
		w = pairing.NewG1().SetBytes(base64StringToElementBytes(wk))

		jsonStringPri := readFile("../config/prikey")
		prikeyObject, _ := ourjson.ParseObject(jsonStringPri)
		//fmt.Println(prikeyObject)
		xk, _ := prikeyObject.GetString("x")
		//fmt.Println(xk)
		x = pairing.NewZr().SetBytes(base64StringToElementBytes(xk))

	} else {

		g = pairing.NewG1().Rand()
		x = pairing.NewZr().Rand()
		u = pairing.NewG2().Rand()

		v = pairing.NewG1().NewFieldElement()
		v = v.PowZn(g, x)
		w = pairing.NewG2().NewFieldElement()
		w = w.PowZn(u, x)
		pubkey := Pubkey{
			G: elementToBase64(*g),
			U: elementToBase64(*u),
			V: elementToBase64(*v),
			W: elementToBase64(*w),
		}
		prikey := Prikey{
			X: elementToBase64(*x),
		}
		pubkeyObject, err := json.Marshal(pubkey)
		prikeyObject, _ := json.Marshal(prikey)
		if err != nil {
			fmt.Println(err)
		}
		writeFile("../config/pubkey", string(pubkeyObject))
		writeFile("../config/prikey", string(prikeyObject))

	}
	//jsonStringPub := readFile("../config/pubkey")
	return User{
		username: username,
		x:        *x,
		keys:     PublicParam{pairing: *pairing, g: *g, u: *u, v: *v, w: *w},
	}

}

func (usr *User) tagGen(filepath string) {
	fmt.Println("User generate tags ...")
	var nSectors [][]byte = proProcessFile(filepath)
	//fmt.Println(len(nSectors))
	//fmt.Println(string(nSectors[0]))
	var tags []pbc.Element = usr.metaGen(nSectors)

	//fmt.Println(len(tags))
	usr.storeTag(tags)
	fmt.Println("ok.")
}

func (usr *User) metaGen(file [][]byte) []pbc.Element {
	var count int = len(file)
	blockTags := make([]pbc.Element, count)
	var i = 0
	for i < count {
		hash256 := sha256.New()
		hash256.Write(file[i])
		hashBytes := hash256.Sum(nil)
		hs := usr.keys.pairing.NewG1().NewFieldElement().SetFromHash(hashBytes)
		mi := usr.keys.pairing.NewZr().NewFieldElement().SetFromHash(file[i])
		tmpU := &usr.keys.u
		tmp := usr.keys.pairing.NewG1().NewFieldElement()
		//var tmp *pbc.Element = tmpU
		tmp = tmp.PowZn(tmpU, mi)
		tmp = tmp.Mul(tmp, hs)
		tmp = tmp.PowZn(tmp, &usr.x)
		blockTags[i] = *tmp
		i++
		//blockTags[i] = keys.u.duplicate().powZn(mi.duplicate()).mul(hash).powZn(x);

		//signature := hex.EncodeToString(hash256)

		//var hs pbc.Element
		//hs =
		//signature := hex.EncodeToString(hash256)
		//var hash pbc.Element
		//pairing := User{}.key.pairing
		//hash = pairing.NewG2().SetFromHash()
	}
	return blockTags

}

func (u *User) storeTag(tags []pbc.Element) {
	//var tagsObject1 ourjson.JsonObject
	var i = 0
	tmp := make(map[string]string, len(tags))
	for i < len(tags) {
		var keyString string = string(i)
		//fmt.Println(elementToBase64(tags[i]))
		//fmt.Println(&tags[i])
		tmp[keyString] = elementToBase64(tags[i])
		i++
	}
	//fmt.Println(tmp)
	tagsObject1, err := json.Marshal(tmp)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(tagsObject1)
	writeFile("../data/"+u.username+"tag", string(tagsObject1))
	writeFile("../config/tags", string(tagsObject1))
	//writeFile("./config/tags", tagsObject1.toString())
	//tagsObject1.Put(keyString, elementToBase64(tags[i]))
}
