package main

import "fmt"

type RealSeq struct {
	seq int
}

func (seqData *RealSeq) GetSeq() int {
	return seqData.seq
}

func (seqData *RealSeq) SetSeq(seq int) {
	seqData.seq = seq
}

func main() {
	seqSize := 100
	seqMap := NewSeqMap(seqSize)

	seqData := new(RealSeq)
	seqMap.NewSeq(seqData)

	gotData := seqMap.GetData(seqData.GetSeq())
	fmt.Printf("%v\n", gotData)

	delData := seqMap.DelSeq(seqData.GetSeq())
	fmt.Printf("%v\n", delData)

	gotData = seqMap.GetData(seqData.GetSeq())
	fmt.Printf("%v\n", gotData)

	for i := 0; i < 100; i++ {
		seqData = new(RealSeq)
		seqMap.NewSeq(seqData)
		fmt.Printf("%d, ", seqData.seq)
	}

	gotData = seqMap.GetData(30)
	fmt.Printf("\ngot one test: %v\n", gotData)

	for i := 20; i < 90; i++ {
		delData := seqMap.DelSeq(i)
		fmt.Printf("del i=%d, %v\n", i, delData)
	}

	fmt.Printf("test got\n")

	for i := 0; i < 100; i++ {
		gotData = seqMap.GetData(i)
		fmt.Printf("got i=%d, %v\n", i, gotData)
	}

	fmt.Printf("\n")
}
