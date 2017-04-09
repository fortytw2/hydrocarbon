package testutil

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

// RunContainer runs a given docker container and returns a port on which the
// container can be reached
func RunContainer(ctx context.Context, container string, port string, waitFunc func(addr string) error) (string, error) {
	free := freePort()
	addr := fmt.Sprintf("localhost:%d", free)
	cmd := exec.CommandContext(ctx, "docker", "run", "-p", fmt.Sprintf("%d:%s", free, port), container)
	// run this in the background
	go func() {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			fmt.Printf("could not run container, %s\n", err)
		}
	}()

	for {
		err := waitFunc(addr)
		if err == nil {
			time.Sleep(time.Millisecond * 150)
			break
		}
	}

	return addr, nil
}

func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}
