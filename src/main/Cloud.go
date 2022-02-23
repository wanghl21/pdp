package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Nik-U/pbc"
)

type Cloud struct {
	keys PublicParam
}

func NewCloud() Cloud {
	return Cloud{NewPublicParam()}

}
func (cd *Cloud) proofGen(filePath string) {
	fmt.Println("Cloud computes proof......\t\t\t")
	var challenges []Challenge = cd.getChallenges()
	proof := cd.genProof(challenges, filePath)
	miu := Utils{}.elementToBase64(proof.Miu)
	r := Utils{}.elementToBase64(proof.R)
	hashMul := Utils{}.elementToBase64(proof.HashMul)
	var proofObject = ProofObject{
		Miu:     miu,
		R:       r,
		HashMul: hashMul,
	}
	storeProof(proofObject)
	fmt.Println("[ok].")
}

func storeProof(proof ProofObject) {
	proofObject, _ := json.Marshal(proof)
	Utils{}.writeFile("../config/proof", string(proofObject))

}

func (cd *Cloud) genProof(challenges []Challenge, filePath string) Proof {

	//var proof map[string]pbc.Element
	var cdata [][]byte = FileOperation{}.readDate(challenges, filePath)
	var pairing pbc.Pairing = cd.keys.pairing
	c := len(challenges)
	miu := pairing.NewZr().Set0()
	r := pairing.NewZr().Rand()
	R := pairing.NewG1().NewFieldElement()
	R = R.PowZn(&cd.keys.w, r)

	hashMul := pairing.NewG2().Set1()

	var i = 0
	for i < c {

		if len(cdata[i]) == 0 {
			fmt.Println("file length mismatch!")
			break
		}
		hash256 := sha256.New()
		hash256.Write(cdata[i])
		//fmt.Println(cdata[i])
		hashBytes := hash256.Sum(nil)

		hash := pairing.NewG1().NewFieldElement().SetFromHash(hashBytes)
		hashMul = hashMul.Mul(hashMul, hash.PowZn(hash, &challenges[i].random))
		mi := pairing.NewZr().SetFromHash(cdata[i])
		miu = miu.Add(miu, mi.MulZn(mi, &challenges[i].random))
		i++
	}
	miu = miu.Add(miu, r)

	return Proof{
		*miu,
		*R,
		*hashMul,
	}

}

func (cd *Cloud) getChallenges() []Challenge {

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
		valueElement := cd.keys.pairing.NewZr().SetBytes(valueBytes)
		challenges = append(challenges, Challenge{num: keyInt, random: *valueElement})
	}
	return challenges

}
