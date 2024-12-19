package main

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "dry run to ensure everything is working",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			go func() {
				log.Printf("about to start main()...")
				main()
			}()
			time.Sleep(time.Second * 10)
		})
	}
}
