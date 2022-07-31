package counter

import (
	"fmt"
	"sync"
	"time"
)

type WindowSliderCounter struct {
	bucketCount int64 //bucket的数量
	bucketNS    int64 // 每个bucket的纳秒数
	windowNS    int64 // 窗口的纳秒数
	maxCount    int64 // 限流的上线
	buckets     map[int64]*Bucket
	mutext      sync.Mutex
}

type Bucket struct {
	Total  int64
	Permit int64
	Reject int64
}

func NewWindowSliderCounter(bucketCount, budgetMS, maxCount int64) *WindowSliderCounter {
	if bucketCount < 1 || budgetMS < 1 {
		return nil
	}
	res := &WindowSliderCounter{
		bucketCount: bucketCount,
		bucketNS:    budgetMS * 1e6,
		windowNS:    bucketCount * budgetMS * 1e6,
		maxCount:    maxCount,
		buckets:     make(map[int64]*Bucket, bucketCount),
	}
	go func() {
		for {
			select {
			case <-time.After(time.Duration(budgetMS) * time.Millisecond):
				end := res.getRemoveEnd()
				for k, _ := range res.buckets {
					if k >= end {
						continue
					}
					delete(res.buckets, k)
				}
				fmt.Printf("after remove, len:%v \n", len(res.buckets)) //仅用于测试，不准确
			}
		}
	}()
	return res
}

func (c *WindowSliderCounter) getRemoveEnd() int64 {
	return (time.Now().UnixNano() - c.windowNS) / c.bucketNS
}

func (c *WindowSliderCounter) Check() bool {
	key := time.Now().UnixNano() / c.bucketNS
	c.mutext.Lock()
	defer c.mutext.Unlock()
	if _, ok := c.buckets[key]; !ok {
		c.buckets[key] = &Bucket{}
	}
	var count int64
	k := key
	for i := (int)(c.bucketCount) - 1; i >= 0; i-- {
		if v, ok := c.buckets[k]; ok {
			count += v.Permit
		}
		k--
	}
	fmt.Println("count:", count)
	if count < c.maxCount {
		c.buckets[key].permit()
		return true
	}
	c.buckets[key].reject()
	return false
}

func (b *Bucket) permit() {
	b.Total++
	b.Permit++
}

func (b *Bucket) reject() {
	b.Total++
	b.Reject++
}
