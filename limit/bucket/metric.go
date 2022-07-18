package bucket

type Metric struct {
  AllowReqCount int64
  StartTime     int64
  ReqCount      int64
  Next          *Metric
}

func (m *Metric) Reset() {
  m.AllowReqCount = 0
  m.ReqCount = 0
  m.StartTime = 0
}
