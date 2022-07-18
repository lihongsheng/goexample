package token

import (
  "fmt"
  "goexample/limit"
  "sync"
  "time"
)

// Limiter
// 比如一秒内最大允许 500, 一秒内打满，那么一毫秒就是 0.5个
type Limiter struct {
  // 最大容量
  Capacity int // 容量
  // 速率
  Rate float64
  // 最后的令牌数
  LastToken int
  // 最后请求时间
  LastTime time.Time
  // ttl 生存周期
  Ttl time.Time
  // maxTimeOut
  timeOut time.Duration
  // lock
  lock sync.Mutex
  // TotalCount
  TotalCount int64
  // currentCount
  CurrentCount int64
  // allowCount
  AllowCount int64
}

func NewTokenLimit(capacity int, rate float64) limit.Limiter {
  ttl := time.Duration((float64(capacity)/rate)*2) * time.Millisecond
  l := &Limiter{
    Capacity:  capacity,
    Rate:      rate,
    LastToken: 0,
    LastTime:  time.Time{},
  }
  l.timeOut = ttl
  l.Ttl = l.LastTime.Add(ttl)

  return l
}

func (l *Limiter) Allow() bool {
  l.lock.Lock()
  defer l.lock.Unlock()
  nowTime := time.Now()
  l.TotalCount++
  l.CurrentCount++
  if l.Ttl.Before(nowTime) { // 这段时间内没有请求,重置记录
    l.LastToken = 0
    l.LastTime = nowTime
    l.Ttl = nowTime.Add(l.timeOut)
    l.CurrentCount = 0
    l.AllowCount = 0
  }
  // 这段时间内产生的令牌桶数
  delta := nowTime.UnixMilli() - l.LastTime.UnixMilli()
  deltaToken := int(float64(delta) * l.Rate)
  if deltaToken <= 0 {
    deltaToken = 0
  }
  allow := false
  l.LastToken += deltaToken
  // 不要超过总量
  if l.LastToken > l.Capacity {
    l.LastToken = l.Capacity
  }
  // 记录这次请求的时间
  l.LastTime = nowTime
  l.Ttl = nowTime.Add(l.timeOut)
  if l.LastToken >= 1 {
    l.AllowCount++
    l.LastToken--
    allow = true
  }

  return allow
}

func (l *Limiter) DescString() string {
  l.lock.Lock()
  defer l.lock.Unlock()
  return fmt.Sprintf("TotalCount[%d] CurrentCount[%d] AllowCount[%d] LastTime[%d] LastToken[%d]", l.TotalCount, l.CurrentCount, l.AllowCount, l.LastTime.UnixMilli(), l.LastToken)
}
