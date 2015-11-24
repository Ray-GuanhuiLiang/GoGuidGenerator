package server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestServerWithClient(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:5588", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	client := NewGuidClient(conn)
	ctx := context.Background()
	var req Req
	for i := 0; i < 3; i++ {
		resp, err := client.GetGuid(ctx, &req)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("response %s\n", resp)
	}
}
