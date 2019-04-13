package main

import (
	"fmt"
	"testing"
	"time"
)

func TestPermissionPerSecond(t *testing.T) {
	calls := 100
	start := time.Now()
	done := make(chan bool, calls)
	for i := 0; i < calls; i++ {
		go func() {
			ctx, err := createContext()
			if err != nil {
				t.Fatalf("failed to create context")
			}

			conn, err := connectToResponder(ctx)
			if err != nil {
				t.Fatalf("-> failed to connect to responder")
			}

			if err := preTransact(conn, ctx); err != nil {
				t.Fatalf("-> pre transaction failed: %s", err)
			}

			if err := transact(conn, ctx); err != nil {
				t.Fatalf("-> transaction failed: %s", err)
			}

			if err := conn.Close(); err != nil {
				t.Fatalf("-> failed to close connection to the responder")
			}
			done <- true
		}()
	}

	for i := 0; i < calls; i++ {
		<-done
	}

	end := time.Since(start)
	fmt.Println(calls, "request per", end.String())
}
