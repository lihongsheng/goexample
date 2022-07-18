package limit

type Limiter interface {
  Allow() bool
  DescString() string
}
