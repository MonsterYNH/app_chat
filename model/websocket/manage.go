package socket

import "sync"

type Manage struct {
	clients map[string]*Client
	lock *sync.Mutex
}

func CreateManage() *Manage {
	manage := &Manage{clients: make(map[string]*Client), lock: &sync.Mutex{}}
	return manage
}

func(manage *Manage) AddClient(userId string, client *Client) {
	manage.lock.Lock()
	manage.clients[userId] = client
	manage.lock.Unlock()
}

func(manage *Manage) GetClient(userId string) (*Client, bool) {
	manage.lock.Lock()
	client, exist := manage.clients[userId]
	manage.lock.Unlock()
	return client, exist
}