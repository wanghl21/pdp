package main

import (
	"encoding/json"
	"fmt"
	"github.com/W1llyu/ourjson"
)

type Blockchain struct {
}

func (bc Blockchain) chalFile(filepath string, u User) {
	fmt.Println("Blockchain computes challenges......\t\t")
	fmt.Println("[ok].")
}

func (bc Blockchain) challengeGen(c int, n int) {
	fmt.Println("Blockchain computes challenges......\t\t")
	var v = NewVerifier()
	var challenges []Challenge = v.ChallengeGen(c, n)
	storeChallenges(challenges)
	//fmt.Print(challenges)
	fmt.Println("[ok].")
}

func storeChallenges1(challenges []Challenge) {
	var cmap = make([]map[string]string, len(challenges))
	var i = 0
	for i < len(challenges) {
		ctmp := challenges[i]
		keyString := string(ctmp.num)
		valueString := Utils{}.elementToBase64(ctmp.random)
		cmap[i][keyString] = valueString
		i++
	}
	challengesObject, _ := json.Marshal(cmap)
	Utils{}.writeFile("../config/challenges", string(challengesObject))
}
func storeChallenges(challenges []Challenge) {
	var cmap = make(map[int]string, len(challenges))
	var i = 0
	for i < len(challenges) {
		ctmp := challenges[i]
		keyString := ctmp.num
		valueString := Utils{}.elementToBase64(ctmp.random)
		cmap[keyString] = valueString
		i++
	}

	challengesObject, _ := json.Marshal(cmap)
	Utils{}.writeFile("../config/challenges", string(challengesObject))
}

func getChallenges(verifier Verifier) []Challenge {
	chalstring := Utils{}.readFile("../config/challenges")
	var challenges []Challenge
	var cmap map[int]string
	err := json.Unmarshal([]byte(chalstring), &cmap)
	if err != nil {
		fmt.Println(err)
	}

	for k, v := range cmap {
		keyInt := k
		valueBytes := Utils{}.base64StringToElementBytes(v)
		valueElement := verifier.keys.pairing.NewZr().SetBytes(valueBytes)
		challenges = append(challenges, Challenge{num: keyInt, random: *valueElement})
	}

	return challenges

}

func getProof(verifier Verifier) Proof {

	proofString := Utils{}.readFile("../config/proof")
	proofObject, _ := ourjson.ParseObject(proofString)

	miuString, _ := proofObject.GetString("miu")
	miu := verifier.keys.pairing.NewZr().SetBytes(Utils{}.base64StringToElementBytes(miuString))
	wString, _ := proofObject.GetString("R")
	r := verifier.keys.pairing.NewG2().SetBytes(Utils{}.base64StringToElementBytes(wString))

	hashMulString, _ := proofObject.GetString("hashMul")
	hashMul := verifier.keys.pairing.NewG1().SetBytes(Utils{}.base64StringToElementBytes(hashMulString))
	return Proof{
		Miu:     *miu,
		R:       *r,
		HashMul: *hashMul,
	}
}

func (Blockchain) verify(id string) {
	fmt.Print("Blockchain verifies proof......\t\t\t")
	verifier := NewVerifier()
	var proof Proof = getProof(verifier)
	var challenges []Challenge = getChallenges(verifier)
	isTrue := verifier.VeriProof(challenges, proof, id)
	fmt.Print("[ok].")
	fmt.Println(isTrue)
}
