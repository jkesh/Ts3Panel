package core

import (
	"context"
	"log"
	"sync"
	"time"

	"Ts3Panel/config"

	"github.com/jkesh/ts3-go/ts3"
)

type DebugLogger struct{}

func (l *DebugLogger) Print(v ...interface{})                 { log.Print(v...) }
func (l *DebugLogger) Printf(format string, v ...interface{}) { log.Printf(format, v...) }
func (l *DebugLogger) Debug(v ...interface{})                 { log.Print(v...) }
func (l *DebugLogger) Debugf(format string, v ...interface{}) { log.Printf(format, v...) }

var (
	Client  *ts3.Client
	once    sync.Once
	sseChan = make(chan map[string]string, 100)
)

func InitTS3() {
	once.Do(func() {
		conf := config.GlobalConfig.TS3

		var err error
		if conf.Protocol == "ssh" {
			log.Println("[Core] Connecting via SSH...")
			Client, err = ts3.NewSSHClient(conf.Host, conf.Port, conf.User, conf.Password)
		} else {
			log.Println("[Core] Connecting via TCP (Raw)...")
			cfg := ts3.Config{
				Host:            conf.Host,
				Port:            conf.Port,
				Timeout:         10 * time.Second,
				KeepAlivePeriod: 1 * time.Minute,
			}
			Client, err = ts3.NewClient(cfg)
		}

		if err != nil {
			log.Fatalf("[Core] Connect failed: %v", err)
		}
		log.Println("[Core] Connected. Waiting 1s before login...")
		time.Sleep(1 * time.Second)
		Client.SetLogger(&DebugLogger{})

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if conf.Protocol == "tcp" {
			if err := Client.Login(ctx, conf.User, conf.Password); err != nil {
				log.Fatalf("[Core] Login failed: %v", err)
			}
		}

		if err := Client.Use(ctx, 1); err != nil {
			log.Fatalf("[Core] Select server failed: %v", err)
		}

		go registerEvents()
		log.Println("[Core] TS3 Service Ready.")
	})
}

func registerEvents() {
	ctx := context.Background()
	_ = Client.OnTextMessage(ctx, func(msg string) {
		BroadcastToSSE("message", msg)
	})
	_ = Client.OnClientEnter(ctx, func(msg string) {
		BroadcastToSSE("enter", msg)
	})
}

func BroadcastToSSE(evtType, msg string) {
	select {
	case sseChan <- map[string]string{"type": evtType, "data": msg}:
	default:
	}
}

func GetSSEChannel() <-chan map[string]string {
	return sseChan
}
