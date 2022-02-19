package mirror

import (
	"crypto/rand"
	"fmt"
	"github.com/bitmex-mirror/pool"
	"log"
	"math"
	"strconv"
	"time"
)

type RatioType int

const seed string = "00f2"

const (
	RatioWalletBalance RatioType = iota
	RatioAvailableBalance
	RatioMarginBalance
)

type symbol struct {
	Symbol    string
	LotSize   float64
	TickSize  float64
	LastPrice float64
}

func NewMirror(host, sub *pool.Account, subId string) {
	//m := mirror{
	//	host:  host,
	//	sub:   sub,
	//	subId: subId,
	//}
}

type mirror struct {
	host             *pool.Account
	sub              *pool.Account
	subId            string
	hostCh           chan []byte
	balanceRatio     map[string]float64 // key: currency, value: ratio
	balanceType      RatioType
	balanceThreshPct float64
	symbolInfo       map[string]symbol
	manualCalibrate  chan struct{}
	closing          chan string
	Done             <-chan struct{}
}

func (m *mirror) manager(done chan struct{}, logger *log.Logger) {
	defer close(done)

	calibrateTicker := time.NewTicker(time.Minute * 5)
	for {
		select {
		case msg := <-m.hostCh:
			m.handleHostMsg(msg)
		case <-m.host.Done:
			logger.Println("mirror Host account disconnected")
			return
		case <-m.sub.Done:
			return
		case by := <-m.closing:
			m.sub.Close()
			logger.Printf("mirror closing for %s by %s\n", m.subId, by)
			return
		case <-calibrateTicker.C:
			for {
				// flush host channel
				for len(m.hostCh) > 0 {
					<-m.hostCh
				}
				//m.calibrate()
				if len(m.hostCh) > 0 {
					continue // messages were received during calibration, need to recalibrate
				}
				break // calibration success
			}
		}
	}
}

// requestCalibration can be used by anyone to request a recalibration
func (m *mirror) requestCalibration() {
	select {
	case m.manualCalibrate <- struct{}{}:
	default:
	}
}

func (m *mirror) BalanceRatio(currency string) float64 {
	hostMargins := m.host.MarginsCopy()
	subMargins := m.sub.MarginsCopy()
	for i := range hostMargins {
		for j := range subMargins {
			if hostMargins[i].Currency == subMargins[j].Currency {
				var balanceRatio float64
				switch m.balanceType {
				case RatioWalletBalance:
					balanceRatio =
						subMargins[j].WalletBalance.ElseZero() / hostMargins[i].WalletBalance.ElseZero()
				case RatioAvailableBalance:
					balanceRatio =
						subMargins[j].AvailableMargin.ElseZero() / hostMargins[i].AvailableMargin.ElseZero()
				case RatioMarginBalance:
					balanceRatio =
						subMargins[j].MarginBalance.ElseZero() / hostMargins[i].MarginBalance.ElseZero()
				}
				oldBalance := m.balanceRatio[hostMargins[i].Currency]
				if oldBalance == 0 {
					m.balanceRatio[hostMargins[i].Currency] = balanceRatio
				} else {
					if m.balanceThreshPct < (balanceRatio-oldBalance)/oldBalance {
						m.balanceRatio[hostMargins[i].Currency] = balanceRatio
					}
				}
			}
		}
	}

	return m.balanceRatio[currency]
}

func (m *mirror) roundOrderSize(symbol string, size float64) float64 {
	return round(size, m.symbolInfo[symbol].LotSize)
}

func round(x, unit float64) float64 {
	rounded := math.Floor(x/unit) * unit
	formatted, _ := strconv.ParseFloat(fmt.Sprintf("%.10f", rounded), 64)

	return formatted
}

func clOrdId(ordId string) string {
	b := make([]byte, 6)
	rand.Read(b)
	uuid := fmt.Sprintf("%s-%s-%x", ordId[:18], seed, b)
	return uuid
}

func randomGuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid[0:8]
}

func matchIDs(ordId, clOrdId string) bool {
	return ordId != "" && ordId[:18] == clOrdId[:18] && clOrdId[19:23] == seed
}

func (m *mirror) close(by string) {
	select {
	case m.closing <- by:
		<-m.Done
	case <-m.Done:
	}

	//a.closeDone.Do(func() {
	//	close(a.Done)
	//	a.pool.internalRemovals <- Removal{
	//		id:     a.id,
	//		Reason: errors.New(reason),
	//	}
	//	a.resetPartials()
	//})
	//logger.Printf("closing Account: %s | reason: %s", a.id, reason)
}
