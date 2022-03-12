package offchain

import (
	"encoding/json"
	"fmt"
	"github.com/Nik-U/pbc"
)

type Challenge struct {
	num    int
	random pbc.Element
}

func getChallenges(verifier Verifier, challengesObject []byte) []Challenge {
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
