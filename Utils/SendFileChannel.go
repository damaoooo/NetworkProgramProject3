package Utils

import (
	"errors"
	"sync"
)

type FileChanUUId struct {
	Channel chan string
	Uuid    string
}

type SendFileManager struct {
	ChannelList []FileChanUUId
	lock        sync.RWMutex
}

func InitManager() *SendFileManager {
	s := SendFileManager{}
	s.New()
	return &s
}

func (c *SendFileManager) New() *SendFileManager {
	c.ChannelList = []FileChanUUId{}
	return c
}

func (c *SendFileManager) IsEmpty() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.ChannelList) == 0
}

func (c *SendFileManager) Add(channel chan string, uuid string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	fileChanUuid := FileChanUUId{
		Channel: channel,
		Uuid:    uuid,
	}
	c.ChannelList = append(c.ChannelList, fileChanUuid)
}

func (c *SendFileManager) IsExist(uuid string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, fileChan := range c.ChannelList {
		if fileChan.Uuid == uuid {
			return true
		}
	}
	return false
}

func (c *SendFileManager) FindChanByUuid(uuid string) (chan string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, fileChan := range c.ChannelList {
		if fileChan.Uuid == uuid {
			return fileChan.Channel, nil
		}
	}
	return nil, errors.New("no_such_uuid")
}

func (c *SendFileManager) RemoveChanByUuid(uuid string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	idx := -1
	for i, fileChan := range c.ChannelList {
		if fileChan.Uuid == uuid {
			idx = i
		}
	}
	if idx >= 0 {
		c.ChannelList = append(c.ChannelList[:idx], c.ChannelList[:idx+1]...)
	}
}
