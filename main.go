package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	docs, err := loadDocuments("enwiki-latest-abstract1.xml")
	fmt.Printf("Finish loading %d documents.\n", len(docs))

	filePtr, err := os.Open("./idx.json")
	if err != nil {
		idx := createIdx(docs)
		filePtr, err := os.Create("./idx.json")
		if err != nil {
			log.Fatal("文件创建失败", err.Error())
		}
		defer filePtr.Close()

		enc := json.NewEncoder(filePtr)
		err = enc.Encode(idx)
		if err != nil {
			log.Fatal("编码错误", err.Error())
		}
	}

	defer filePtr.Close()
	var idx index
	dec := json.NewDecoder(filePtr)
	err = dec.Decode(&idx)
	if err != nil {
		log.Fatal("解码失败", err.Error())
	}

	fmt.Printf("Please provide keyword:\n")

	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		r := idx.search(scan.Text())
		if len(r) == 0 {
			fmt.Println("Sorry, no document contains the keyword.\n")
		} else {
			for _, id := range r {
				fmt.Printf("%d\t%s\n\n", id, docs[id].Text)
			}
		}
		fmt.Printf("Please provide keyword:\n")
	}

}
