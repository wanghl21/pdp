package offchain

import (
	"fmt"
	"io"
	"os"

	"github.com/Nik-U/pbc"
	"github.com/W1llyu/ourjson"
)

type FileOperation struct {
}

const blockSize int = 1024

func (FileOperation) blockNumbers(filePath string) int {
	file, err := os.Stat(filePath)
	if err != nil {
		fmt.Println(err)
	}
	var fileLength int = int(file.Size())
	var number int = fileLength / (blockSize)
	var remain int = 0
	if fileLength%blockSize > 0 {
		remain = 1

	}
	return number + remain
}

func (FileOperation) proProcessFile(filePath string) [][]byte {
	var fileBlocks int = FileOperation{}.blockNumbers(filePath)
	nSectors := make([][]byte, fileBlocks)
	//var nSectors [][]byte
	inputFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("open file faild,err:%s\n", err)
	}
	var i int = 0
	for i < fileBlocks {
		nSectors[i] = make([]byte, blockSize)
		n, err := inputFile.Read(nSectors[i][:])
		//fmt.Println(nSectors)
		if err == io.EOF {
			break
		}
		if err == io.EOF {
			break
		}
		if n < blockSize {
			//nSectors[i][n:blockSize] = 0
		}
		i++
	}
	defer inputFile.Close()
	return nSectors

}

func (FileOperation) readDate(challenges []Challenge, filePath string) [][]byte {
	c := len(challenges)
	data := make([][]byte, c)
	inputFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("open file faild,err:%s\n", err)
	}
	var i int = 0
	for i < c {
		data[i] = make([]byte, blockSize)
		inputFile.Seek(int64(challenges[i].num*blockSize), 0)
		//tmp := make([]byte, challenges[i].num*blockSize)
		//inputFile.Read(tmp[:])

		n, err1 := inputFile.Read(data[i][:])

		if err1 == io.EOF {
			break
			fmt.Println(err1)
		}

		if n < blockSize {
			//nSectors[i][n:blockSize] = 0
		}
		i++
	}

	return data
}
func (FileOperation) readTags(keys PublicParam, id string, challenges []Challenge) []*pbc.Element {
	var jsonString string = Utils{}.readFile("../data/" + id + "tag")
	tagsObject, _ := ourjson.ParseObject(jsonString)
	count := len(challenges)
	tags := make([]*pbc.Element, count)
	var i = 0
	for i < count {
		jsonStr, _ := tagsObject.GetString(string(challenges[i].num))
		tags[i] = keys.pairing.NewG2().SetBytes(Utils{}.base64StringToElementBytes(jsonStr))
		i++
	}
	return tags

}
