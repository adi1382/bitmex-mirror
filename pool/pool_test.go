package pool

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/bitmex-mirror/auth"
)

func TestSingleClient(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	p := NewPool(true, logger)

	go func() {
		for {
			select {
			case msg := <-p.Removals:
				fmt.Println("###########################################################")
				fmt.Println("###########################################################")
				fmt.Println("###########################################################")
				fmt.Println("Removal:", msg)
				fmt.Println("###########################################################")
				fmt.Println("###########################################################")
				fmt.Println("###########################################################")
			}
		}
	}()

	fmt.Println("pool created")
	config := auth.NewConfig("YXTfYOwQn0e7gq3bziPyeNDa", "6tbsbjU1GjCGipthTHuLb2PrVswCuGBrA10JQ5IrhlGy0Ytz")

	p.Put(config, "", config.Key, logger)
	acc, err := p.Get(config.Key)
	if err != nil {
		t.Error(err)
	}

	acc.WaitForPartial()
}

func TestNewPool(t *testing.T) {
	//defer profile.Start(profile.ProfilePath("."), profile.TraceProfile).Stop()
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	p := NewPool(true, logger)

	go func() {
		for {
			select {
			case msg := <-p.Removals:
				fmt.Println("Removal:", msg)
			}
		}
	}()

	fmt.Println("pool created")
	config := auth.NewConfig("YXTfYOwQn0e7gq3bziPyeNDa", "6tbsbjU1GjCGipthTHuLb2PrVswCuGBrA10JQ5IrhlGy0Ytz")
	//config2 := auth.NewConfig("YXTfYOwQn0e7gq3bziPyeND", "6tbsbjU1GjCGipthTHuLb2PrVswCuGBrA10JQ5IrhlGy0Ytz")

	num := 30

	//for i := 0; i < num; i++ {
	//	if i%3 == 0 {
	//		p.Put(config2, "", config.Key+strconv.Itoa(i), logger)
	//	} else {
	//		p.Put(config, "", config.Key+strconv.Itoa(i), logger)
	//	}
	//}

	fmt.Println("num go routine 1: ", runtime.NumGoroutine())

	for i := 0; i < num; i++ {
		p.Put(config, "", config.Key+strconv.Itoa(i), logger)
	}

	fmt.Println("num go routine 2: ", runtime.NumGoroutine())

	accounts := make([]*Account, 0, num)
	for i := 0; i < num; i++ {
		a, err := p.Get(config.Key + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		accounts = append(accounts, a)
	}

	fmt.Println(len(accounts))

	for i := 0; i < num; i++ {
		accounts[i].WaitForPartial()
	}

	fmt.Println("num go routine 3: ", runtime.NumGoroutine())

	select {}
}
