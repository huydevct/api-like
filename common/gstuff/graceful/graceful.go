package graceful

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Graceful ..
type Graceful struct {
	sig          os.Signal
	mux          sync.Mutex
	processing   int
	gracefulStop chan os.Signal
}

// New ..
func New(timer time.Duration) (g *Graceful) {
	g = &Graceful{
		mux:          sync.Mutex{},
		gracefulStop: make(chan os.Signal),
	}
	signal.Notify(g.gracefulStop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-g.gracefulStop
		g.mux.Lock()
		g.sig = sig
		g.mux.Unlock()
		log.Printf("Graceful shutdown %v\n", g.sig)
		time.Sleep(timer)
		os.Exit(0)
	}()

	return
}

// ExitProgram : Chủ động kết thúc chuơng trình
func (g *Graceful) ExitProgram() {
	log.Println("Program is existed.")
	g.gracefulStop <- syscall.SIGTERM
}

// CheckAlive ..
func (g *Graceful) CheckAlive() bool {
	g.mux.Lock()
	sig := g.sig
	g.mux.Unlock()
	return sig == nil
}

// CheckAndExit ..
func (g *Graceful) CheckAndExit() {
	for {
		time.Sleep(time.Millisecond * 100)
		if !g.CheckAlive() && g.processing == 0 {
			os.Exit(0)
		}
	}
}

// IncreaseProcessing ..
func (g *Graceful) IncreaseProcessing() {
	g.mux.Lock()
	g.processing++
	g.mux.Unlock()
}

// DescreaseProcessing ..
func (g *Graceful) DescreaseProcessing() {
	g.mux.Lock()
	g.processing--
	g.mux.Unlock()
}
