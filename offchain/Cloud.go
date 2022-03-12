package main

type Cloudp struct {
	keys PublicParam
}

func NewCloudp() Cloudp {
	return Cloudp{NewPublicParam()}

}
