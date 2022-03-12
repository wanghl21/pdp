package offchain

import (
	"fmt"
	"github.com/Nik-U/pbc"
	"github.com/W1llyu/ourjson"
)

type Proof struct {
	Miu     pbc.Element `json:"miu"`
	R       pbc.Element `json:"R"`
	HashMul pbc.Element `json:"hashMul"`
}

type ProofObject struct {
	Miu     string `json:"miu"`
	R       string `json:"R"`
	HashMul string `json:"hashMul"`
}

func getProofObeject(verifier Verifier, proofobject []byte) ProofObject {
	fmt.Print("Blockchain getProofObeject......\t\t\t")
	//proofString := Utils{}.readFile("../config/proof")
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
