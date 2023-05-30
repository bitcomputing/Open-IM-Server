package fcm

import (
	push "Open_IM/internal/push/providers"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Push(t *testing.T) {
	offlinePusher := NewFcm()
	resp, err := offlinePusher.Push(context.Background(), []string{"test_uid"}, "test", "test", "12321", push.PushOpts{})
	assert.Nil(t, err)
	fmt.Println(resp)
}
