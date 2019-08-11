package trace
import (
  "io"
  "fmt"
)
// Trace is an Interface showing objects which can record logs in codes.
type Tracer interface {
  Trace(...interface{})
}

func New(w io.Writer) Tracer {
  return &tracer{out: w}
}

type tracer struct {
  out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
  t.out.Write([]byte(fmt.Sprint(a...)))
  t.out.Write([]byte("\n"))
}
