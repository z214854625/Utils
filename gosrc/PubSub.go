package main

import (
	"fmt"
	"sync"
)

// Subscriber 包含订阅者的通道以及一个取消订阅的函数
type Subscriber struct {
	Ch   chan interface{}
	Done chan struct{}
}

// PubSub struct to hold the subscribers and the sync mechanisms
type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string]map[*Subscriber]struct{} //Subscriber 作为键，struct{}作为值不占内存，可以通过Subscriber的key删除该项
}

// NewPubSub creates a new PubSub instance
func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string]map[*Subscriber]struct{}),
	}
}

// Subscribe subscribes to a given topic and returns a Subscriber struct
func (ps *PubSub) Subscribe(topic string) *Subscriber {
	sub := &Subscriber{
		Ch:   make(chan interface{}),
		Done: make(chan struct{}),
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.subscribers[topic]; !exists {
		ps.subscribers[topic] = make(map[*Subscriber]struct{})
	}
	ps.subscribers[topic][sub] = struct{}{}
	return sub
}

// Unsubscribe unsubscribes a given subscriber from a topic
func (ps *PubSub) Unsubscribe(topic string, sub *Subscriber) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if subs, exists := ps.subscribers[topic]; exists {
		if _, exists := subs[sub]; exists {
			delete(subs, sub)
			close(sub.Ch)
			close(sub.Done)
		}
		// 如果某个订阅者列表为空了，则删除主题
		if len(subs) == 0 {
			delete(ps.subscribers, topic)
		}
	}
}

// Publish publishes a message to a given topic
func (ps *PubSub) Publish(topic string, msg interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if subs, exists := ps.subscribers[topic]; exists {
		for s := range subs {
			s.Ch <- msg
		}
	} else {
		fmt.Printf("subscribers empty! topic= %v, msg= %v\n", topic, msg)
	}
}

// func main() {
// 	ps := NewPubSub()

// 	// Subscriber 1
// 	sub1 := ps.Subscribe("topic1")
// 	go func() {
// 		for {
// 			select {
// 			case msg, ok := <-sub1.Ch:
// 				if !ok {
// 					fmt.Println("Subscriber 1 channel closed")
// 					return
// 				}
// 				fmt.Printf("Subscriber 1 received: %v\n", msg)
// 			case <-sub1.Done:
// 				fmt.Println("Subscriber 1 done")
// 				return
// 			}
// 		}
// 	}()

// 	// Subscriber 2
// 	sub2 := ps.Subscribe("topic1")
// 	go func() {
// 		for {
// 			select {
// 			case msg, ok := <-sub2.Ch:
// 				if !ok {
// 					fmt.Println("Subscriber 2 channel closed")
// 					return
// 				}
// 				fmt.Printf("Subscriber 2 received: %v\n", msg)
// 			case <-sub2.Done:
// 				fmt.Println("Subscriber 2 done")
// 				return
// 			}
// 		}
// 	}()

// 	// Publish messages
// 	ps.Publish("topic1", "Hello, World!")
// 	ps.Publish("topic1", "Another message")

// 	// Unsubscribe Subscriber 1
// 	ps.Unsubscribe("topic1", sub1)

// 	// Publish another message
// 	ps.Publish("topic1", "Message after unsubscribe")

// 	// Unsubscribe Subscriber 2
// 	ps.Unsubscribe("topic1", sub2)

// 	// Publish another message, this time no subscribers should receive it
// 	ps.Publish("topic1", "Final message")

// 	// Sleep to allow goroutines to finish
// 	select {} // Keep the main function alive
// }
