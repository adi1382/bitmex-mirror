package pool

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func TestRate(t *testing.T) {
	l := rate.NewLimiter(rate.Every(time.Minute/120), 120)
	l.WaitN(context.Background(), 120)
	l.SetLimit(0)
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("set limit")
		l.SetLimit(rate.Every(time.Minute / 120))
	}()
	err := l.Wait(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
