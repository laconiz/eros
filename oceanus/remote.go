package oceanus

import (
	"github.com/laconiz/eros/network"
	"sync"
)

type Courses struct {
	burls map[string]*Burl

	//
	courses map[string]*Course

	mutex *sync.Mutex
}

func (c *Courses) addNode(node *Node, session network.Session) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	burl, ok := c.burls[node.ID]
	if ok {
		burl.node = node
		logger.Infof("node updated: %+v", node)
	} else {
		burl = NewBurl(node)
		logger.Infof("node join: %+v", node)
	}

	burl.session = session
	for _, course := range c.courses {
		course.router.expired()
	}
}
