package ws

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

type Message struct {
	To      string
	Content []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client
		case client := <-h.Unregister:
			delete(h.Clients, client.ID)
			client.Conn.Close()
		case message := <-h.Broadcast:
			if client, ok := h.Clients[message.To]; ok {
				client.Send <- message.Content
			}
		}
	}
}

func (h *Hub) Notify(to string, msg string) {
	h.Broadcast <- Message{
		To:      to,
		Content: []byte(msg),
	}
}
