package pool

import (
	"context"
	"fmt"
	"github.com/bitmex-mirror/auth"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestNewPool(t *testing.T) {
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
	for i := 0; i < 30; i++ {
		p.Put(context.Background(), config, "", logger, config.Key+strconv.Itoa(i))
	}
	accounts := make([]*account, 0, 30)
	for i := 0; i < 30; i++ {
		a, err := p.Get(config.Key + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		accounts = append(accounts, a)
	}

	fmt.Println(len(accounts))

	for i := 0; i < 30; i++ {
		accounts[i].WaitForPartial()
	}
}
