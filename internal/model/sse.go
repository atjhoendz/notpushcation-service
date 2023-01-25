package model

import log "github.com/sirupsen/logrus"

type (
	SSEBrokerUsecase interface {
		//Start()
		Serve()
	}

	SSEHandler func(msgChan chan *SSEMessage)

	SSEBroker struct {
		clients       map[chan *SSEMessage]bool
		newClients    chan chan *SSEMessage
		closedClients chan chan *SSEMessage
		message       chan *SSEMessage
	}

	SSEMessage struct {
		Event EventType
		Data  string
	}

	EventType string
)

const (
	Update       EventType = "UPDATE"
	LiveBlogPost EventType = "LIVE_BLOG_POST"
)

func NewSSEBroker() *SSEBroker {
	return &SSEBroker{
		clients:       make(map[chan *SSEMessage]bool),
		newClients:    make(chan chan *SSEMessage),
		closedClients: make(chan chan *SSEMessage),
		message:       make(chan *SSEMessage),
	}
}

func (b *SSEBroker) AddNewClient(c chan *SSEMessage) {
	b.newClients <- c

	log.Info("add new client")
}

//func (b *SSEBroker) AddNewClient(c chan *SSEMessage) {
//	b.clients[c] = true
//	log.Info("added new client")
//}

func (b *SSEBroker) CloseClient(c chan *SSEMessage) {
	delete(b.clients, c)
	log.Info("removed client")
}

func (b *SSEBroker) BroadcastMessage(msg *SSEMessage) {
	b.message <- msg
	log.Info("broadcast nich")
}

func (b *SSEBroker) Start() {
	go func() {
		for {
			select {
			case c := <-b.newClients:
				b.clients[c] = true
				log.Info("added new client")
			case c := <-b.closedClients:
				b.CloseClient(c)
			case msg := <-b.message:
				for c, _ := range b.clients {
					c <- msg
				}
				log.Infof("broadcast message to %d clients", len(b.clients))
			}
		}
	}()
}
