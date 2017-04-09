package testutil

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Container struct {
	Name string
	Addr string

	cmd *exec.Cmd
}

func (c *Container) Shutdown() {
	c.cmd.Process.Signal(syscall.SIGTERM)
}

// RunContainer runs a given docker container and returns a port on which the
// container can be reached
func RunContainer(container string, port string, waitFunc func(addr string) error) (*Container, error) {

	free := freePort()
	addr := fmt.Sprintf("localhost:%d", free)
	cmd := exec.Command("docker", "run", "-p", fmt.Sprintf("%d:%s", free, port), container)
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

	return &Container{
		Name: container,
		Addr: addr,
		cmd:  cmd,
	}, nil
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
