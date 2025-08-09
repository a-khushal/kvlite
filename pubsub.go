package main

func (kv *KVStore) Subscribe(channel string, client *Client) {
	kv.subsMu.Lock()
	defer kv.subsMu.Unlock()

	if kv.subscribers[channel] == nil {
		kv.subscribers[channel] = make(map[*Client]bool)
	}
	kv.subscribers[channel][client] = true
}

func (kv *KVStore) Unsubscribe(channel string, client *Client) {
	kv.subsMu.Lock()
	defer kv.subsMu.Unlock()

	if clients, ok := kv.subscribers[channel]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(kv.subscribers, channel)
		}
	}
}

func (kv *KVStore) Publish(channel, message string) {
	kv.subsMu.RLock()
	defer kv.subsMu.RUnlock()

	if clients, ok := kv.subscribers[channel]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				// drop message if client send buffer full (to avoid blocking)
			}
		}
	}
}
