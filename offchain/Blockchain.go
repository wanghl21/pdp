package offchain

import (
	"encoding/json"
	"fmt"

	"github.com/W1llyu/ourjson"
)

type Blockchain struct {
}

func NewBlockchain() Blockchain {
	return Blockchain{}

}
func (bc *Blockchain) verify(id string, proofObject ProofObject, challenges []Challenge) bool {
	fmt.Print("Blockchain verifies proof......\t\t\t")
	verifier := NewVerifier()
	//var proof1 Proof = getProof(verifier)

	//	var challenges []Challenge = getChallenges(verifier)
	isTrue := verifier.VeriProof(challenges, proofObject, id)
	fmt.Print("[ok].")
	fmt.Println(isTrue)
	return isTrue
}

func (bc *Blockchain) getChallenges(verifier Verifier, challengesObject []byte) []Challenge {
	//chalstring := Utils{}.readFile("../config/challenges")
	fmt.Print("Blockchain getChallenges......\t\t\t")
	var challenges []Challenge
	var cmap map[int]string
	err := json.Unmarshal(challengesObject, &cmap)
	if err != nil {
		fmt.Println(err)
	}

	for k, v := range cmap {
		keyInt := k
		valueBytes := base64StringToElementBytes(v)
		valueElement := verifier.keys.pairing.NewZr().SetBytes(valueBytes)
		challenges = append(challenges, Challenge{num: keyInt, random: *valueElement})
	}
	return challenges
}

func (bc *Blockchain) getProof(verifier Verifier, proofobject []byte) Proof {

	//proofString := Utils{}.readFile("../config/proof")
	proofObject, _ := ourjson.ParseObject(string(proofobject))

	miuString, _ := proofObject.GetString("miu")
	miu := verifier.keys.pairing.NewZr().SetBytes(base64StringToElementBytes(miuString))
	wString, _ := proofObject.GetString("R")
	r := verifier.keys.pairing.NewG2().SetBytes(base64StringToElementBytes(wString))

	hashMulString, _ := proofObject.GetString("hashMul")
	hashMul := verifier.keys.pairing.NewG1().SetBytes(base64StringToElementBytes(hashMulString))
	return Proof{
		Miu:     *miu,
		R:       *r,
		HashMul: *hashMul,
	}

}

func (bc *Blockchain) getProofObeject(verifier Verifier, proofobject []byte) ProofObject {
	fmt.Print("Blockchain getProofObeject......\t\t\t")
	//proofString := readFile("../config/proof")
	proofObject, _ := ourjson.ParseObject(string(proofobject))

	miuString, _ := proofObject.GetString("miu")
	miu := verifier.keys.pairing.NewZr().SetBytes(base64StringToElementBytes(miuString))
	wString, _ := proofObject.GetString("R")
	r := verifier.keys.pairing.NewG2().SetBytes(base64StringToElementBytes(wString))

	hashMulString, _ := proofObject.GetString("hashMul")
	hashMul := verifier.keys.pairing.NewG1().SetBytes(base64StringToElementBytes(hashMulString))
	proof := Proof{
		Miu:     *miu,
		R:       *r,
		HashMul: *hashMul,
	}

	miuob := elementToBase64(proof.Miu)
	rob := elementToBase64(proof.R)
	hashMulob := elementToBase64(proof.HashMul)
	var proofOb = ProofObject{
		Miu:     miuob,
		R:       rob,
		HashMul: hashMulob,
	}
	return proofOb

}
