// 线程安全的消息队列

package queue

import (
	"container/list"
	"errors"
	"sync"
)

var (
	ErrClosed   = errors.New("add message to a closed queue")
	ErrOverflow = errors.New("add message to a full queue")
)

type Queue struct {
	lst    list.List  // 消息列表
	cond   *sync.Cond // 同步条件变量
	cap    int        // 最大容量 限制内存 <=0时无限制
	closed bool       // 是否关闭
	mutex  sync.Mutex //
}

// 向队列中追加一条消息
func (q *Queue) Add(msg interface{}) error {

	q.mutex.Lock()

	// 向已关闭的队列追加消息
	if q.closed {
		q.mutex.Unlock()
		return ErrClosed
	}

	// 队列消息堆积数量已到最大值
	if q.cap > 0 && q.lst.Len() >= q.cap {
		q.mutex.Unlock()
		return ErrOverflow
	}

	// 追加消息并同步状态
	q.lst.PushBack(msg)
	q.mutex.Unlock()
	q.cond.Signal()

	return nil
}

// 从队列中取出一组消息及当前队列状态
func (q *Queue) Pick() (ret []interface{}, exit bool) {

	q.mutex.Lock()

	// 队列为空且未关闭时等待状态通知
	if q.lst.Len() == 0 && !q.closed {
		q.cond.Wait()
	}

	exit = q.closed

	// 获取队列中的所有消息
	for q.lst.Len() > 0 {
		ret = append(ret, q.lst.Remove(q.lst.Front()))
	}

	q.mutex.Unlock()

	return
}

// 关闭消息队列
func (q *Queue) Close() {

	q.mutex.Lock()

	// 队列已经关闭
	if q.closed {
		q.mutex.Unlock()
		return
	}

	// 更新队列并同步状态
	q.closed = true
	q.mutex.Unlock()
	q.cond.Signal()
}

// 创建一个消息队列
func NewQueue(cap int) *Queue {

	q := &Queue{cap: cap}
	q.cond = sync.NewCond(&q.mutex)

	return q
}
