package jackychan_test

import (
	"fmt"
	jackychan "github.com/rolancia.jackychan"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestChan(t *testing.T) {
	t.Run("channel", func(t *testing.T) {
		ch := jackychan.Chan[int]("")
		go func() {
			ch.Push(1)
			ch.Push(2)
		}()
		res := sum(ch.Pop(), ch.Pop())
		assert.Equal(t, 3, res)
	})

	t.Run("async channel", func(t *testing.T) {
		ch := jackychan.ChanAsync[int]("", 5)
		ch.Push(1)
		ch.Push(2)
		res := sum(ch.Pop(), ch.Pop())
		assert.Equal(t, 3, res)
	})

	t.Run("select channel", func(t *testing.T) {
		n := 2
		wg := sync.WaitGroup{}
		wg.Add(n)
		ch1 := jackychan.Chan[int]("")
		ch2 := jackychan.Chan[string]("")
		go func() {
			cnt := 0
			for i := 0; i < n; i++ {
				select {
				case v := <-ch1.Selector():
					assert.Equal(t, 1, v)
					wg.Done()
					cnt++
				case v := <-ch2.Selector():
					assert.Equal(t, "1", v)
					wg.Done()
					cnt++
				}
			}
		}()
		ch1.Push(1)
		ch2.Push("1")
		wg.Wait()
	})

	t.Run("select async channel", func(t *testing.T) {
		n := 3
		ch1 := jackychan.ChanAsync[int]("", n)
		ch2 := jackychan.ChanAsync[string]("", n)
		for i := 0; i < n; i++ {
			ch1.Push(i)
			ch2.Push(fmt.Sprintf("%d", i))
		}
		cnt := 0
		for cnt < 2*n {
			select {
			case v := <-ch1.Selector():
				_ = v
				cnt++
			case v := <-ch2.Selector():
				_ = v
				cnt++
			}
		}
	})
}

func sum(a, b int) int {
	return a + b
}
