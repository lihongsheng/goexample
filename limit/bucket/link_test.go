package bucket

import (
  "sync"
  "testing"
  "time"
)

func TestBucket_Allow(t *testing.T) {
  l := NewLinkBucket(10, 1000, 100)
  wg := sync.WaitGroup{}
  wg.Add(10)
  i := 0
  for i < 10 {
    go func() {
      defer wg.Done()
      tt := time.After(900 * time.Millisecond)
      for {
        select {
        case <-tt:
          return
        default:
          l.Allow()

        }
      }
    }()
    i++
  }
  wg.Wait()
  t.Log(l.DescString())
  if l.AllowCount > 100 {
    t.Fatalf("err")
  }
  time.Sleep(1 * time.Second)
  if !l.Allow() {
    t.Fatalf("err")
  }
  t.Log(l.DescString())
  time.Sleep(1100 * time.Millisecond)
  if !l.Allow() {
    t.Fatalf("err")
  }
  t.Log(l.DescString())
}
