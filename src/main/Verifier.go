package main

import (
	"fmt"
)

type Verifier struct {
	keys PublicParam
}

func NewVerifier() Verifier {
	return Verifier{NewPublicParam()}

}

func (v *Verifier) ChallengeGen(c int, allBlocks int) []Challenge {
	fmt.Println("ChallengeGen")
	var challenges []Challenge
	var idx = Utils{}.GenRandom(c, allBlocks)
	var i = 0
	for i < c {
		e := v.keys.pairing.NewZr().Rand()
		var c = Challenge{idx[i], *e}
		challenges = append(challenges, c)
		i++
	}

	return challenges
}

func (vr *Verifier) VeriProof(challenges []Challenge, proof Proof, id string) bool {

	miu := proof.Miu
	r := proof.R
	hashMul := proof.HashMul
	tags := FileOperation{}.readTags(vr.keys, id, challenges)
	sigma := vr.keys.pairing.NewG1().Set1()

	var i = 0
	for i < len(challenges) {
		tags[i] = tags[i].PowZn(tags[i], &challenges[i].random)
		sigma = sigma.Mul(sigma, tags[i])
		i++
	}

	sigma = sigma.Mul(sigma, &r)
	tmpu := vr.keys.u
	
	tmppowMiu := tmpu.PowZn(&tmpu, &miu)
	powMiu := tmppowMiu.Mul(tmppowMiu, &hashMul)
	tmp1 := vr.keys.pairing.NewGT().Pair(powMiu, &vr.keys.v)
	tmp2 := vr.keys.pairing.NewGT().Pair(sigma, &vr.keys.g)
	fmt.Println(tmp1)
	fmt.Println(tmp2)
	return tmp1.Equals(tmp2)
}
