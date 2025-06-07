package queue

import "fmt"

type testItem struct {
	n int
}

func TestItem(n int) *testItem {
	return &testItem{n: n}
}

func (t *testItem) process() error {
	fmt.Println(t.n)
	return nil
}
