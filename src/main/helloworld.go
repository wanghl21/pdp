package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
)

func main() {

	filePath := "../data/whl.txt"               //audit file
	filePath2 := "../data/whl1.txt"             //audit file
	n := FileOperation{}.blockNumbers(filePath) // the number of total blocks
	c := n                                      // the number of blocks which you want to audit

	u := newUser("user0")

	var i = 45
	for i < 46 {
		chal1 := Challenge{i, *u.keys.pairing.NewZr().Rand()}
		chal2 := Challenge{i, *u.keys.pairing.NewZr().Rand()}
		var chals1 []Challenge
		chals1 = append(chals1, chal1)
		var chals2 []Challenge
		chals2 = append(chals2, chal2)
		block1 := FileOperation{}.readDate(chals1, filePath)
		block2 := FileOperation{}.readDate(chals2, filePath2)

		blockStr1 := string(block1[0])
		blockStr2 := string(block2[0])
		//fmt.Println(blockStr1)
		//fmt.Println(blockStr2)
		fmt.Print("strings.EqualFold(blockStr2, blockStr1):  ")
		fmt.Println(strings.EqualFold(blockStr2, blockStr1))

		hash256 := sha256.New()
		hash256.Write(block1[0])
		hashBytes1 := hash256.Sum(nil)
		/*
			s1 := hex.EncodeToString(hashBytes1)
			fmt.Println(hash256)
			fmt.Println(hashBytes1)
			fmt.Println(s1)

		*/

		hash2562 := sha256.New()
		hash2562.Write(block2[0])
		hashBytes2 := hash2562.Sum(nil)

		fmt.Print("bytes.Equal(hashBytes1, hashBytes2)： ")
		fmt.Println(bytes.Equal(hashBytes1, hashBytes2))
		/*
			s2 := hex.EncodeToString(hashBytes2)
			fmt.Println(hash2562)
			fmt.Println(hashBytes2)
			fmt.Println(s2)

		*/

		hash1 := u.keys.pairing.NewG1().NewFieldElement().SetBytes(hashBytes1)
		hash2 := u.keys.pairing.NewG1().NewFieldElement().SetBytes(hashBytes2)
		//fmt.Println(hash1.String())
		//fmt.Println(hash2.String())
		fmt.Println(hash1.Equals(hash2))
		/*
			mblock1 := u.keys.pairing.NewZr().NewFieldElement().SetBytes(block1[0])
			mblock2 := u.keys.pairing.NewZr().NewFieldElement().SetBytes(block2[0])

			var j = 0
			for j < 1024 {
				if block1[0][j] != block2[0][j] {
					fmt.Println(j)
				}
				j++
			}
			if mblock1 != mblock2 {
				fmt.Println(i)
			}

		*/
		i++
	}
	u.tagGen(filePath)
	Blockchain{}.challengeGen(c, n)
	cloud := NewCloud()
	cloud.proofGen(filePath2)
	Blockchain{}.verify(u.username)

	/*
		var f FileOperation
		f.proProcessFile("../data/whl.txt")
		var c []Challenge
		ch := Challenge{num: 1}
		c = append(c, ch)
		f.readDate(c, "../data/whl.txt")
		//var user User
		//user = user.User("user0")
		//fmt.Println(user.key)
		//fmt.Println(user.key.pairing.NewZr().Rand())

		//var bls messageData
		//bls.main1()
		//var p Blockchain

		//var param PublicParam
		//var u Utils
		/*
			sharedParams := Utils{}.readFile("../config/b.properties")
			jsonStringPub := Utils{}.readFile("../config/pubkey")

			pubkeyObject, _ := ourjson.ParseObject(jsonStringPub)
			g1, _ := pubkeyObject.GetString("g")
			fmt.Printf(g1)
			var u Utils
			sharedG := u.base64StringToElementBytes(g1)
			fmt.Println(sharedG)
			fmt.Println(sharedParams)
			params := pbc.GenerateA(160, 512)
			fmt.Println(params)
			Utils{}.writeFile("../config/b.properties", params.String())
			fmt.Println(reflect.TypeOf(params))
			sharedParams1 := Utils{}.readFile("../config/b.properties")
			pairing := pbc.NewPairing(params)
			fmt.Println(pairing)
			pairing, _ = pbc.NewPairingFromString(params.String()) // loads pairing parameters from a string and instantiates a pairing
			fmt.Println(pairing)
			pairing2, _ := pbc.NewPairingFromString(sharedParams1) // loads pairing parameters from a string and instantiates a pairing
			fmt.Println(pairing2)
			fmt.Println(sharedG)
			x := pairing.NewG2().SetBytes(sharedG)
			//g := pairing.NewG2().SetBytes(sharedG)
			fmt.Println(x)

	*/
	//param.PublicParam()
	/*
		fmt.Println(s)
		param.test()
		var c int = 1
		fmt.Print(c)
		p.ringing()
		p.chalFile("s", User{})
		p.challengeGen(1, 3)
		p.verify("s")
		//u。tils := Utils{}
		//var data []byte

	*/

	//data = json.Unmarshal(utils.ReadAll("../config/proof"), &data)

}
