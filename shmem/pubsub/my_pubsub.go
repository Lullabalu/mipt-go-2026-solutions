//go:build !solution

package pubsub

import (
	"context"
	"errors"
	"sync"
)

var _ Subscription = (*MySubscription)(nil)

type MySubscription struct {
	mu   sync.Mutex
	msgs chan any

	alive bool
	send  MsgHandler
}

func (s *MySubscription) Send(msg any) {
	s.msgs <- msg
}

func (s *MySubscription) Unsubscribe() {
	s.mu.Lock()
	s.alive = false
	s.mu.Unlock()
}

var _ PubSub = (*MyPubSub)(nil)

type MyPubSub struct {
	mu  sync.Mutex
	ctx context.Context

	done chan struct{}

	alive bool
	mp    map[string][]*MySubscription
}

func Run(sub *MySubscription, done chan struct{}, send chan any, doneSenler chan struct{}) {
	list := []any{}
LOOP:
	for {
		select {
		case <-done:
			break LOOP
		case msg := <-sub.msgs:
			list = append(list, msg)
		default:
			if len(list) > 0 {
				select {
				case send <- list[0]:
					list = list[1:]
				default:
				}
			}
		}
	}

	for _, msg := range list {
		send <- msg
	}
	close(doneSenler)

}

func Sendler(done chan struct{}, send MsgHandler, get chan any) {
	for {
		select {
		case msg := <-get:
			send(msg)
		case <-done:
			return
		}
	}

}

func NewPubSub() PubSub {
	pubsub := MyPubSub{ctx: context.TODO(), alive: true, mp: make(map[string][]*MySubscription), done: make(chan struct{})}
	return &pubsub
}

func (p *MyPubSub) Subscribe(subj string, cb MsgHandler) (Subscription, error) {
	p.mu.Lock()
	if !p.alive {
		defer p.mu.Unlock()
		return nil, errors.New("Паб закрыт")
	}
	sub := MySubscription{alive: true, send: cb, msgs: make(chan any)}
	p.mp[subj] = append(p.mp[subj], &sub)
	p.mu.Unlock()

	send := make(chan any)
	donesend := make(chan struct{})
	go Run(&sub, p.done, send, donesend)
	go Sendler(donesend, sub.send, send)
	return &sub, nil
}

func (p *MyPubSub) Publish(subj string, msg any) error {
	p.mu.Lock()
	if !p.alive {
		defer p.mu.Unlock()
		return errors.New("Паб закрыт")
	}

	for _, v := range p.mp[subj] {
		if v.alive {
			v.msgs <- msg
		}
	}

	p.mu.Unlock()
	return nil
}

func (p *MyPubSub) Close(ctx context.Context) error {
	p.mu.Lock()
	p.alive = false
	p.ctx = ctx
	close(p.done)
	p.mu.Unlock()

	return nil
}
