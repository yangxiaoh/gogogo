package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func getServer(port int) http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ans"))
	})

	return http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", port),
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	done := make(chan error, 2)
	stop := make(chan struct{})

	server1 := getServer(8080)
	server2 := getServer(8081)

	g.Go(func() error {
		g.Go(func() error {
			<- stop
			fmt.Println("stop serve 1")
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
			return server1.Shutdown(ctx)
		})
		done <- server1.ListenAndServe()
		return nil
	})

	g.Go(func() error {
		g.Go(func() error {
			<- stop
			fmt.Println("stop serve 2")
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
			return server2.Shutdown(ctx)
		})
		done <- server2.ListenAndServe()
		return nil
	})

	quit := make(chan os.Signal, 0)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	g.Go(func() error {
		select {
		case <- ctx.Done():
			fmt.Printf("err group exit [%+v]\n", ctx.Err())
		case sig := <- quit:
			fmt.Printf("get os signal: [%+v]\n", sig)
		case <- done:
			fmt.Printf("server err\n")
		}

		close(stop)
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("get err [%+v]", err)
	}

	fmt.Println("finish")
}