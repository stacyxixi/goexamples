package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	//"sync"
	"time"
)

var shared int

type HomePageSize struct {
	URL  string
	Size int
}

func f(n int) {
	shared++
	//fmt.Println(n, ":", shared)
	for i := 0; i < 10; i++ {
		//fmt.Println(n, ":", i, shared)
		amt := time.Duration(rand.Intn(250))
		time.Sleep(time.Millisecond * amt)
	}
}

func pinger(c chan<- string) { //bidirectional pinger(c chan string)
	for i := 0; ; i++ {
		//time.Sleep(time.Second * 5)
		c <- "ping"
		//time.Sleep(time.Second * 5)
	}
}
func ponger(c chan<- string) {
	for i := 0; ; i++ {
		c <- "pong"
		//time.Sleep(time.Second * 5)
	}
}

func push(c chan<- string, s string) {
	for i := 0; ; i++ {
		c <- s
	}
}
func printer(c <-chan string) {
	for {
		msg := <-c
		fmt.Println(msg)
		time.Sleep(time.Second * 1)
	}
}

func selector() {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		for {
			c1 <- "from 1"
			time.Sleep(time.Second * 2)
		}
	}()
	go func() {
		for {
			c2 <- "from 2"
			time.Sleep(time.Second * 3)
		}
	}()
	go func() {
		for {
			select {
			case msg1 := <-c1:
				fmt.Println(msg1)
			case msg2 := <-c2:
				fmt.Println(msg2)
			case <-time.After(time.Second):
				fmt.Println("timeout")
				//default:
				//fmt.Println("nothing ready")
			}
		}
	}()
	var input string
	fmt.Scanln(&input)
}

func pingpong() {

	for i := 0; i < 10; i++ {
		go f(i)
	}
	//go f(0)
	//fmt.Println(shared)

	c := make(chan string)
	for i := 0; i < 10; i++ {
		go push(c, strconv.Itoa(i)) //order is determistic???
	}

	//go pinger(c)
	//go ponger(c)
	go printer(c)

	var input string
	fmt.Scanln(&input)
	fmt.Println(shared)
}

func main() {
	//selector()
	//pingpong()
	//bufferedChan()
	buggy()
}

func buggy() {
	var x int
	//threads := runtime.GOMAXPROCS(0)
	//fmt.Println("number of threads=", threads)
	for i := 0; i < 20; i++ {
		go func() {
			for {
				//fmt.Println(x)
				x++ //tight loop will clog the code
			}
		}()
	}
	time.Sleep(time.Second)
	fmt.Println("x =", x)
}

func bufferedChan() { // a buffered channel is asynchorous
	urls := []string{
		"http://www.apple.com",
		"http://www.amazon.com",
		"http://www.google.com",
		"http://www.microsoft.com",
	}
	results := make(chan HomePageSize)
	for _, url := range urls {
		go func(url string) { //notice that we pass url instead of directly referencing url because of the rules of closure
			res, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			bs, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			results <- HomePageSize{
				URL:  url,
				Size: len(bs),
			}
		}(url)
	}
	var biggest HomePageSize

	for range urls {

		result := <-results
		if result.Size > biggest.Size {
			biggest = result
		}
	}
	fmt.Println("The biggest home page:", biggest.URL, " size is ", biggest.Size)

}
