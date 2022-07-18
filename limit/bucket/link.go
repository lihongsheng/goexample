package bucket

import (
"fmt"
"sync"
"time"
)


type Bucket struct {
  WindowLengthInMs int64
  Capacity         int
  start            *Metric
  end              *Metric
  AllowMax         int64
  nodeSize         int
  lock             sync.Mutex
  AllowCount       uint64
  ReqCount         uint64
  InitTime         int64
}

func NewLinkBucket(capacity int, intervalInMs int64, allowMax int64) *Bucket {
  return &Bucket{
    Capacity:         capacity,
    WindowLengthInMs: intervalInMs / int64(capacity),
    AllowMax:         allowMax,
    lock:             sync.Mutex{},
  }
}

func (l *Bucket) Allow() bool {
  l.lock.Lock()
  defer l.lock.Unlock()
  nowTime := time.Now().UnixNano() / 1e6
  if l.InitTime <= 0 {
    l.InitTime = nowTime
  }
  allow := false
  currentNode := l.getCurrentWindow(nowTime)
  allowCount := int64(0)
  node := l.start
  for node != nil {
    allowCount += node.AllowReqCount
    node = node.Next
  }
  currentNode.ReqCount++
  l.ReqCount++
  if allowCount < l.AllowMax {
    allow = true
    l.AllowCount++
    currentNode.AllowReqCount++
  }

  return allow
}

func (l *Bucket) DescString() string {
  desc := ""
  node := l.start
  for node != nil {
    desc += fmt.Sprintf("windowStartTime[%d], windowReqCount[%d], windowAllowReqCount[%d] \n", node.StartTime, node.ReqCount, node.AllowReqCount)
    node = node.Next
  }
  ExecTime := (time.Now().UnixNano() / 1e6) - l.InitTime
  desc += fmt.Sprintf("ReqCount[%d], AllowCount[%d] InitTime[%d] ExecTime[%d]", l.ReqCount, l.AllowCount, l.InitTime, ExecTime)
  return desc
}

func (l *Bucket) getCurrentWindow(nowTime int64) *Metric {
  oldWindow := l.end
  if l.nodeSize == 0 {
    l.reset(nowTime)
  } else {
    if nowTime == oldWindow.StartTime || nowTime < (oldWindow.StartTime+l.WindowLengthInMs) {

      return oldWindow
    } else if nowTime > oldWindow.StartTime { //超过一个时间周期重置即可
      if nowTime > (oldWindow.StartTime + l.WindowLengthInMs*int64(l.Capacity)) {
        l.reset(nowTime)
      } else {
        // 时间周期内，递增end,直到当前请求时间在最后一个时间槽范围内
        for (oldWindow.StartTime + l.WindowLengthInMs) <= nowTime {
          window := &Metric{
            AllowReqCount: 0,
            StartTime:     oldWindow.StartTime + l.WindowLengthInMs,
            ReqCount:      0,
            Next:          nil,
          }
          oldWindow.Next = window
          oldWindow = window
          l.nodeSize++
          if l.nodeSize > l.Capacity {
            l.start = l.start.Next
            l.nodeSize--
          }
        }
        l.end = oldWindow
        return oldWindow
      }
    }
  }

  return l.end
}

func (l *Bucket) reset(nowTime int64) {
  oldWindow := &Metric{
    AllowReqCount: 0,
    StartTime:     nowTime,
    ReqCount:      0,
    Next:          nil,
  }
  l.start = oldWindow
  l.end = oldWindow
  l.nodeSize = 1
}

func main() {


}
