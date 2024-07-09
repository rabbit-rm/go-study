package main

import (
	"fmt"
	"sync"
	"testing"
)

type Name struct {
	FirstName string
	LastName  string
	Age       uint8
}

func TestName(t *testing.T) {
}

func TestInter(t *testing.T) {
	var a = []int{1, 2, 3}
	for k, v := range a {
		if k == 0 {
			a[0], a[1] = 100, 200
			fmt.Println(a)
		}
		a[k] = v + 100
	}
	fmt.Println(a)
}

func TestPrint(t *testing.T) {
	dogChan = make(chan struct{}, 1)
	catChan = make(chan struct{}, 1)
	fishChan = make(chan struct{}, 1)
	fishChan <- struct{}{}
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		Cat(fishChan, catChan)
	}()
	go func() {
		defer wg.Done()
		Dog(catChan, dogChan)
	}()
	go func() {
		defer wg.Done()
		Fish(dogChan, fishChan)
	}()
	wg.Wait()
}

var dogChan chan struct{}
var catChan chan struct{}
var fishChan chan struct{}

func Dog(rev, send chan struct{}) {
	/*for range 100 {
		<-rev
		fmt.Println("dog")
		send <- struct{}{}
	}*/
}
func Cat(rev, send chan struct{}) {
	/*for range 100 {
		<-rev
		fmt.Println("cat")
		send <- struct{}{}
	}*/
}

func Fish(rev, send chan struct{}) {
	/*for range 100 {
		<-rev
		fmt.Println("fish")
		send <- struct{}{}
	}*/
}
