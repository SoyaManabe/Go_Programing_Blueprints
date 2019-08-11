package trace
import (
  "testing"
  "bytes"
)

func TestNew(t *testing.T) {
  var buf bytes.Buffer
  tracer := New(&buf)
  if tracer == nil {
    t.Error("Return from New is nil")
  } else {
    tracer.Trace("Hello, trace package")
    if buf.String() != "Hello, trace package\n" {
      t.Errorf("'%s', wrong sentence was output", buf.String())
    }
  }
}
