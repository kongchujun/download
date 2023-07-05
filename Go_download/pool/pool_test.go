package pool

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ io.Closer = &SampleConn{}

type SampleConn struct {
}

func (c *SampleConn) Close() error {
	fmt.Println("conn has been close!")
	return nil
}

func TestNewGenericPool(t *testing.T) {
	minOpen := 1
	maxOpen := 2
	maxLifetime := time.Second
	factory := func() (io.Closer, error) {
		return &SampleConn{}, nil
	}
	pool, err := NewGenericPool(minOpen, maxOpen, maxLifetime, factory)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	conn1, err := pool.Acquire()
	assert.NoError(t, err)
	if err != nil {
		return
	}
	fmt.Println(conn1)

	conn2, err := pool.Acquire()
	assert.NoError(t, err)
	if err != nil {
		return
	}
	fmt.Println(conn2)

	conn3, err := pool.Acquire()
	assert.NoError(t, err)
	if err != nil {
		return
	}
	fmt.Println(conn3)
}
