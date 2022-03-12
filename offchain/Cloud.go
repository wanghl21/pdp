package offchain

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

func (cd *Cloud) challengeGen(seed int, c int, allBlocks int) []Challenge {
	fmt.Println("Cloud computes challenges......\t\t")
	var challenges []Challenge
	var idx = GenRandom(seed, c, allBlocks) //应该是区块链上的信息
	//fmt.Println(idx)
	var i = 0
	for i < c {
		e := cd.keys.pairing.NewZr().Rand()
		var c = Challenge{idx[i], *e}
		challenges = append(challenges, c)
		i++
	}

	return challenges
}

func (cd *Cloud) storeChallenges(challenges []Challenge) []byte {

	fmt.Println("Cloud storeChallenges......\t\t")
	var cmap = make(map[int]string, len(challenges))
	var i = 0
	for i < len(challenges) {
		ctmp := challenges[i]
		keyString := ctmp.num
		valueString := elementToBase64(ctmp.random)
		cmap[keyString] = valueString
		i++
	}

	challengesObject, _ := json.Marshal(cmap)
	writeFile("../config/challenges", string(challengesObject))
	return challengesObject

}

func (cd *Cloud) genProof(filePath string, challenges []Challenge) Proof {
	fmt.Println("Cloud computes proof......\t\t\t")
	//var proof map[string]pbc.Element
	var cdata [][]byte = readDate(challenges, filePath)
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

func (cd *Cloud) storeProof(proof Proof) []byte {
	fmt.Println("Cloud storeProof......\t\t\t")

	miu := elementToBase64(proof.Miu)
	r := elementToBase64(proof.R)
	hashMul := elementToBase64(proof.HashMul)
	var proofOb = ProofObject{
		Miu:     miu,
		R:       r,
		HashMul: hashMul,
	}

	proofObject, _ := json.Marshal(proofOb)
	writeFile("../config/proof", string(proofObject))
	return proofObject

}
