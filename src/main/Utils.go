package main

import (
	"encoding/base64"
	"fmt"
	"github.com/Nik-U/pbc"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type Utils struct {
}

// element to base64 Stringstr sting / e element
func (Utils) elementToBase64(e pbc.Element) string {
	return base64.StdEncoding.EncodeToString(e.Bytes())
}

//base64 to elementByte[]
func (Utils) base64StringToElementBytes(elementBase64String string) []byte {
	elementByte, _ := base64.StdEncoding.DecodeString(elementBase64String)
	//fmt.Println(elementByte)
	return elementByte
}

//可以直接使用os.writefile
func (Utils) writeFile(outputfile string, content string) {
	var err error = os.WriteFile(outputfile, []byte(content), 0666)
	if err != nil {
		fmt.Println(err)
		fmt.Println("writeFile")
		return
	}
}

func (Utils) readFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		fmt.Println("readFile")
	}
	var content string = string(data)
	return content

}
func (Utils) readAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	fmt.Println(f)
	if err != nil {
		return nil, err
	}
	//var data []byte
	//data := json.Unmarshal(ioutil.ReadAll(f))
	return ioutil.ReadAll(f)
}

func (Utils) GenRandom(len int, total int) []int {

	rand.Seed(time.Now().UnixNano())
	hashSet := NewHashSet("int")
	//fmt.Println(len)
	//fmt.Println(hashSet.Size())

	for hashSet.Size() < len {
		var tmp int = int(rand.Int63n(int64(total)))
		//fmt.Println(tmp)
		err := hashSet.Add(tmp)
		if err != nil {
			fmt.Println(err)
		}

	}

	var ints []int
	for _, v1 := range hashSet.data {
		switch v := v1.(type) {
		case int:
			//fmt.Println("整型", v)
			var s int
			s = v
			//fmt.Println(s)
			ints = append(ints, s)

		}
	}
	return ints

}
