package offchain

import (
	"encoding/json"
	"fmt"

	"github.com/W1llyu/ourjson"
)

type Verifier struct {
	keys PublicParam
}

func NewVerifier() Verifier {
	return Verifier{NewPublicParam()}

}

func (vr *Verifier) getChallenges(verifier Verifier, challengesObject []byte) []Challenge {
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
		valueBytes := Utils{}.base64StringToElementBytes(v)
		valueElement := verifier.keys.pairing.NewZr().SetBytes(valueBytes)
		challenges = append(challenges, Challenge{num: keyInt, random: *valueElement})
	}
	return challenges

}

func (vr *Verifier) getProofObeject(verifier Verifier, proofobject []byte) ProofObject {
	fmt.Print("Blockchain getProofObeject......\t\t\t")
	//proofString := Utils{}.readFile("../config/proof")
	proofObject, _ := ourjson.ParseObject(string(proofobject))

	miuString, _ := proofObject.GetString("miu")
	miu := verifier.keys.pairing.NewZr().SetBytes(Utils{}.base64StringToElementBytes(miuString))
	wString, _ := proofObject.GetString("R")
	r := verifier.keys.pairing.NewG2().SetBytes(Utils{}.base64StringToElementBytes(wString))

	hashMulString, _ := proofObject.GetString("hashMul")
	hashMul := verifier.keys.pairing.NewG1().SetBytes(Utils{}.base64StringToElementBytes(hashMulString))
	proof := Proof{
		Miu:     *miu,
		R:       *r,
		HashMul: *hashMul,
	}

	miuob := Utils{}.elementToBase64(proof.Miu)
	rob := Utils{}.elementToBase64(proof.R)
	hashMulob := Utils{}.elementToBase64(proof.HashMul)
	var proofOb = ProofObject{
		Miu:     miuob,
		R:       rob,
		HashMul: hashMulob,
	}
	return proofOb

}

func (vr *Verifier) VeriProof(challenges []Challenge, proofObject ProofObject, id string) bool {

	miuString := proofObject.Miu
	miu := vr.keys.pairing.NewZr().SetBytes(Utils{}.base64StringToElementBytes(miuString))
	wString := proofObject.R
	r := vr.keys.pairing.NewG2().SetBytes(Utils{}.base64StringToElementBytes(wString))

	hashMulString := proofObject.HashMul
	hashMul := vr.keys.pairing.NewG1().SetBytes(Utils{}.base64StringToElementBytes(hashMulString))

	//miu := proof.Miu
	//	r := proof.R
	//hashMul := proof.HashMul
	tags := FileOperation{}.readTags(vr.keys, id, challenges) // 相当于读取tags
	sigma := vr.keys.pairing.NewG1().Set1()

	var i = 0
	for i < len(challenges) {
		tags[i] = tags[i].PowZn(tags[i], &challenges[i].random)
		sigma = sigma.Mul(sigma, tags[i])
		i++
	}

	sigma = sigma.Mul(sigma, r)
	tmpu := vr.keys.u

	tmppowMiu := tmpu.PowZn(&tmpu, miu)
	powMiu := tmppowMiu.Mul(tmppowMiu, hashMul)
	tmp1 := vr.keys.pairing.NewGT().Pair(powMiu, &vr.keys.v)
	tmp2 := vr.keys.pairing.NewGT().Pair(sigma, &vr.keys.g)
	fmt.Println(tmp1)
	fmt.Println(tmp2)
	return tmp1.Equals(tmp2)
}
