package offchain

func main() {
	filePath := "../data/whl.txt"   //audit file
	filePath2 := "../data/whl1.txt" //audit file
	n := blockNumbers(filePath)     // the number of total blocks 把文件分块
	c := n / 2                      // the number of blocks which you want to audit
	u := newUser("user0")           //client发送文件分块，
	u.tagGen(filePath)              //计算tag存到区块链上，
	// c应该是有一个F（b||i）得到的子集
	cloud := NewCloud() // 相当于cloud login
	var seed = 2        //作为挑选数据块的随机函数的种子应该是区块链相关信息

	//云端根据随机的种子选择要审计的数据块 然后生成challenge和proof 发给区块链进行验证
	var challenges []Challenge = cloud.challengeGen(seed, c, n) //根据seed随机生成challenges
	challengesObject := cloud.storeChallenges(challenges)       //把challenges[]转为byte便于传给chaincode

	var proof Proof = cloud.genProof(filePath2, challenges)
	proofObject := cloud.storeProof(proof)

	verifier := NewVerifier()
	bc := NewBlockchain()
	bc_chal := bc.getChallenges(verifier, challengesObject)
	bc_proofObject := bc.getProofObeject(verifier, proofObject)
	bc.verify(u.username, bc_proofObject, bc_chal) //这里的u.uername 就像相当于short recod

}
