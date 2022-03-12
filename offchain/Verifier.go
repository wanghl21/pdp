package offchain

import (
	"fmt"
)

type Verifier struct {
	keys PublicParam
}

func NewVerifier() Verifier {
	return Verifier{NewPublicParam()}

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
