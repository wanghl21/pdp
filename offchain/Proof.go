package main

import "github.com/Nik-U/pbc"

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
