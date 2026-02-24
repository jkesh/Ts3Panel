package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"Ts3Panel/config"

	"github.com/jkesh/ts3-go/v2/ts3"
)

type DebugLogger struct{}

func (l *DebugLogger) Print(v ...interface{})                 { log.Print(v...) }
func (l *DebugLogger) Printf(format string, v ...interface{}) { log.Printf(format, v...) }
func (l *DebugLogger) Debug(v ...interface{})                 { log.Print(v...) }
func (l *DebugLogger) Debugf(format string, v ...interface{}) { log.Printf(format, v...) }

type SSEMessage struct {
	Type string
	Data string
}

var (
	client           *ts3.Client
	clientMu         sync.Mutex
	once             sync.Once
	initErr          error
	ErrClientDown    = errors.New("ts3 client is not initialized")
	sseMu            sync.RWMutex
	sseSubscribers   = make(map[int]chan SSEMessage)
	nextSubscriberID int
)

func InitTS3() error {
	once.Do(func() {
		conf := config.GlobalConfig.TS3
		runtimeCfg := ts3.Config{
			Host:            conf.Host,
			Port:            conf.Port,
			Timeout:         10 * time.Second,
			KeepAlivePeriod: 1 * time.Minute,
		}

		switch conf.Protocol {
		case "ssh":
			if conf.User == "" || conf.Password == "" {
				initErr = errors.New("ts3 ssh requires user/password")
				return
			}
			log.Println("[Core] Connecting via SSH...")
			client, initErr = ts3.NewSSHClientWithConfig(conf.Host, conf.Port, conf.User, conf.Password, runtimeCfg)
			if initErr != nil && shouldFallbackToTCP(initErr) {
				fallbackPort := conf.Port
				if fallbackPort == 0 || fallbackPort == 10022 {
					fallbackPort = 10011
				}

				log.Printf("[Core] SSH handshake failed (%v), fallback to TCP Raw on %s:%d...", initErr, conf.Host, fallbackPort)
				runtimeCfg.Port = fallbackPort
				client, initErr = ts3.NewClient(runtimeCfg)
				if initErr == nil {
					log.Println("[Core] TCP fallback connected successfully.")
				}
			}
		case "tcp":
			if conf.User == "" || conf.Password == "" {
				initErr = errors.New("ts3 tcp requires user/password")
				return
			}
			log.Println("[Core] Connecting via TCP (Raw)...")
			client, initErr = ts3.NewClient(runtimeCfg)
		case "webquery":
			if conf.APIKey == "" {
				initErr = errors.New("ts3 webquery requires api_key")
				return
			}
			log.Println("[Core] Connecting via HTTP WebQuery...")
			wqCfg := ts3.WebQueryConfig{
				Host:            conf.Host,
				Port:            conf.Port,
				HTTPS:           conf.HTTPS,
				APIKey:          conf.APIKey,
				BasePath:        conf.BasePath,
				Timeout:         10 * time.Second,
				KeepAlivePeriod: 1 * time.Minute,
				VirtualServerID: conf.ServerID,
			}
			client, initErr = ts3.NewWebQueryClient(wqCfg)
		default:
			initErr = fmt.Errorf("unsupported ts3 protocol: %s", conf.Protocol)
		}

		if initErr != nil {
			return
		}

		log.Println("[Core] Connected.")
		client.SetLogger(&DebugLogger{})

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if conf.Protocol == "tcp" {
			if err := client.Login(ctx, conf.User, conf.Password); err != nil {
				initErr = fmt.Errorf("login failed: %w", err)
				return
			}
		}

		if err := client.Use(ctx, conf.ServerID); err != nil {
			initErr = fmt.Errorf("select server %d failed: %w", conf.ServerID, err)
			return
		}

		if conf.Protocol != "webquery" {
			go registerEvents()
		} else {
			log.Println("[Core] WebQuery mode: event subscribe is disabled.")
		}
		log.Println("[Core] TS3 Service Ready.")
	})
	return initErr
}

func shouldFallbackToTCP(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "handshake failed: eof")
}

func registerEvents() {
	ctx := context.Background()
	if err := client.OnTextMessage(ctx, func(msg string) {
		BroadcastToSSE("message", msg)
	}); err != nil {
		log.Printf("[Core] register text event failed: %v", err)
	}
	if err := client.OnClientEnter(ctx, func(msg string) {
		BroadcastToSSE("enter", msg)
	}); err != nil {
		log.Printf("[Core] register enter event failed: %v", err)
	}
}

func BroadcastToSSE(evtType, msg string) {
	event := SSEMessage{Type: evtType, Data: msg}

	sseMu.RLock()
	defer sseMu.RUnlock()

	for _, sub := range sseSubscribers {
		select {
		case sub <- event:
		default:
		}
	}
}

func SubscribeSSE() (int, <-chan SSEMessage) {
	sseMu.Lock()
	defer sseMu.Unlock()

	nextSubscriberID++
	id := nextSubscriberID
	ch := make(chan SSEMessage, 64)
	sseSubscribers[id] = ch
	return id, ch
}

func UnsubscribeSSE(id int) {
	sseMu.Lock()
	defer sseMu.Unlock()

	ch, ok := sseSubscribers[id]
	if !ok {
		return
	}
	delete(sseSubscribers, id)
	close(ch)
}

// WithTS3 ensures calls run under a single critical section and with initialized client.
func WithTS3(fn func(c *ts3.Client) error) error {
	clientMu.Lock()
	defer clientMu.Unlock()

	if client == nil {
		return ErrClientDown
	}
	return fn(client)
}

func WithTS3Value[T any](fn func(c *ts3.Client) (T, error)) (T, error) {
	clientMu.Lock()
	defer clientMu.Unlock()

	var zero T
	if client == nil {
		return zero, ErrClientDown
	}
	return fn(client)
}
