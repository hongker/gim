package bucket

import (
	"github.com/ebar-go/ego/utils/runtime"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	bucket := NewBucket()
	session1 := NewSession("1001", nil)
	bucket.AddSession(session1)

	session2 := NewSession("1002", nil)
	bucket.AddSession(session2)
	bucket.AddChannel("1")
	bucket.Broadcast([]byte("hello"))

	channel1 := bucket.GetChannel("1")
	bucket.SubscribeChannel(channel1, session1, session2)
	channel1.Broadcast([]byte("world"))

	time.Sleep(time.Second)
	bucket.UnsubscribeChannel(channel1, session1)
	channel1.Broadcast([]byte("world, again"))

	runtime.Shutdown(bucket.Stop)
}
