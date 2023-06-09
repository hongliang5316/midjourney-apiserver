package service

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

var KeyChan *keyChan

func init() {
	KeyChan = &keyChan{
		Mu:    &sync.Mutex{},
		Inner: map[string](chan MessageInfo){},
	}
}

type keyChan struct {
	Mu    *sync.Mutex
	Inner map[string](chan MessageInfo)
}

func (k *keyChan) Init(key string) {
	k.Mu.Lock()
	defer k.Mu.Unlock()
	k.Inner[key] = make(chan MessageInfo, 1)
}

func (k *keyChan) Del(key string) {
	k.Mu.Lock()
	defer k.Mu.Unlock()
	delete(k.Inner, key)
}

func (k *keyChan) Get(key string) chan MessageInfo {
	k.Mu.Lock()
	defer k.Mu.Unlock()
	c, ok := k.Inner[key]
	if !ok {
		return nil
	}

	return c
}

func (k *keyChan) Len() int {
	k.Mu.Lock()
	defer k.Mu.Unlock()
	return len(k.Inner)
}

type MessageInfo struct {
	ID        string
	StartTime int64
	Error     *discordgo.MessageEmbed
}
