package common
import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type slotJob struct {
	id      int64
	do      func()
	ttl     int //当前的ttl
	initTTL int //初始化的ttl
	ch      chan bool
	loop    bool
}

type Slot struct {
	index int
	hooks *list.List
}

type TimeWheel struct {
	autoId         int64
	tick           *time.Ticker
	wheel          []*Slot
	hashWheel      map[int64]*list.Element
	ticksPerwheel  int
	tickPeriod     time.Duration
	currentTick    int
	slotJobWorkers chan bool
	lock           *sync.RWMutex
}

//超时时间及每个timewheel所需要的tick数
func NewTimeWheel(tickPeriod time.Duration, ticksPerwheel int, slotJobWorkers int) *TimeWheel {
	tw := &TimeWheel{
		lock:           &sync.RWMutex{},
		tickPeriod:     tickPeriod,
		hashWheel:      make(map[int64]*list.Element, 10000),
		tick:           time.NewTicker(tickPeriod),
		slotJobWorkers: make(chan bool, slotJobWorkers),
		wheel: func() []*Slot {
			//ticksPerWheel make ticksPerWheel+1 slide
			w := make([]*Slot, 0, ticksPerwheel+1)
			for i := 0; i < ticksPerwheel+1; i++ {
				w = append(w, func() *Slot {
					return &Slot{
						index: i,
						hooks: list.New()}
				}())
			}
			return w
		}(),
		ticksPerwheel: ticksPerwheel + 1,
		currentTick:   0}

	go func() {
		for i := 0;; i++ {
			i = i % tw.ticksPerwheel
			<-tw.tick.C
			tw.lock.Lock()
			tw.currentTick = i
			tw.lock.Unlock()
			//notify expired
			tw.notifyExpired(i)
		}
	}()

	return tw
}

func (self *TimeWheel) Monitor() string {
	ticks := 0
	for _, v := range self.wheel {
		ticks += v.hooks.Len()
	}
	return fmt.Sprintf("[%s] TimeWheel  total-tick:[%d]  workers:[%d/%d]", time.Now().Format("2006-01-02 15:04:05"),
		ticks, len(self.slotJobWorkers), cap(self.slotJobWorkers))
}

//notifyExpired func
func (self *TimeWheel) notifyExpired(idx int) {
	var remove *list.List
	self.lock.RLock()
	slots := self.wheel[idx]
	for e := slots.hooks.Back(); nil != e; e = e.Prev() {
		sj := e.Value.(*slotJob)
		sj.ttl--
		//ttl expired
		if sj.ttl <= 0 {

			//fmt.Printf("%d run..\n", sj.id)

			if sj.loop {//需要继续loop的job，则重置ttl
				sj.ttl=sj.initTTL
			}else {//非loop类型的job则加入删除列表
				if nil == remove {
					remove = list.New()
				}
				//记录删除
				remove.PushFront(e)
			}

			//fmt.Println("use worker")
			self.slotJobWorkers <- true //控制worker数量
			//fmt.Println("user worker --")

			//async
			go func() {
				defer func() {
					if err := recover(); nil != err {
						//ignored
						fmt.Printf("job exec error: %s\n", err.(error))
					}
					//fmt.Println("start release worker")
					<-self.slotJobWorkers
					//fmt.Println("release worker ok")

				}()

				sj.do()

				if !sj.loop {
					sj.ch <- true //非loop类型任务，返回结果通知
					close(sj.ch)
				}else {

				}

				//				select {
				//				case sj.ch <- true:
				//					fmt.Printf("任务%d成功返回", sj.id)
				//				case <-time.After(2 * time.Second):
				//					fmt.Printf("任务%d执行超时\n", sj.id)
				//					<-sj.ch
				//				}

			}()
		}
		//fmt.Printf("job detail, id: %d ttl: %d\n", sj.id, sj.ttl)
	}

	self.lock.RUnlock()

	if nil != remove {
		//remove
		for e := remove.Back(); nil != e; e = e.Prev() {
			re := e.Value.(*list.Element)
			self.lock.Lock()
			slots.hooks.Remove(e.Value.(*list.Element))
			delete(self.hashWheel, re.Value.(*slotJob).id)
			self.lock.Unlock()
		}
	}

}

//add timeout func
func (self *TimeWheel) After(timeout time.Duration, do func()) (int64, chan bool) {

	idx := self.preTickIndex()

	self.lock.Lock()
	slots := self.wheel[idx]
	ttl := int(int64(timeout) / (int64(self.tickPeriod) * int64(self.ticksPerwheel)))
	// log.Debug("After|TTL:%d|%d\n", ttl, timeout)
	id := self.timerId(idx)
	job := &slotJob{
		id :id,
		do : do,
		ttl :ttl,
		initTTL:ttl,
		ch : make(chan bool, 1),
		loop :false,
	}//设置loop为false，即本任务只执行一次

	e := slots.hooks.PushFront(job)
	self.hashWheel[id] = e
	self.lock.Unlock()
	return id, job.ch
}

func (self *TimeWheel) Loop(timeout time.Duration, do func()) (int64, chan bool) {

	idx := self.preTickIndex()

	self.lock.Lock()
	slots := self.wheel[idx]
	ttl := int(int64(timeout) / (int64(self.tickPeriod) * int64(self.ticksPerwheel-1)))
	fmt.Printf("timeout %d tickPeriod %d  ticksPerwheel %d  ttl %d \n", timeout, self.tickPeriod, self.ticksPerwheel, ttl)

	id := self.timerId(idx)
	job := &slotJob{
		id :id,
		do : do,
		ttl :ttl,
		initTTL:ttl,
		ch : make(chan bool, 1),
		loop :true,
	}//设置loop为true，即本任务需在不断地执行
	e := slots.hooks.PushFront(job)
	self.hashWheel[id] = e
	self.lock.Unlock()
	return id, job.ch
}

func (self *TimeWheel) Remove(timerId int64) {
	self.lock.Lock()
	e, ok := self.hashWheel[timerId]
	if ok {
		sid := self.decodeSlot(timerId)
		sl := self.wheel[sid]
		sl.hooks.Remove(e)
		delete(self.hashWheel, timerId)
	}
	self.lock.Unlock()
}

func (self *TimeWheel) decodeSlot(timerId int64) int {
	return int(timerId >> 32)
}

func (self *TimeWheel) timerId(idx int) int64 {
	id := int64(int64(idx<<32) |((atomic.AddInt64(&self.autoId, 1))>>32))
    //id := int64(int64(idx<<32) |int64 (int64(atomic.AddInt64(&self.autoId, 1)<<32 ))>>32)
	//fmt.Println(id)
	return id
}

func (self *TimeWheel) preTickIndex() int {
	self.lock.RLock()
	idx := self.currentTick
	if idx > 0 {
		idx -= 1
	} else {
		idx = self.ticksPerwheel - 1
	}
	self.lock.RUnlock()
	return idx
}