package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"encoding/gob"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"sigs.k8s.io/external-dns/endpoint"
)

func startServerToServeTargets(endpoints []*endpoint.Endpoint, cancel func()) net.Listener {
	ln, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Println("Error listening on ", ln.Addr().String(), ": ", err.Error())
		os.Exit(1)
	}

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			ln.Close()
			return
		}
		enc := gob.NewEncoder(conn)
		enc.Encode(endpoints)
		ln.Close()
		cancel()
	}()
	fmt.Printf("Server listening on %s\n", ln.Addr().String())
	return ln
}

func handleSigterm(cancel func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	signal.Notify(signals, syscall.SIGINT)
	<-signals
	fmt.Println("Received SIGTERM. Terminating...")
	cancel()
}

func serveMetrics(address string) {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.Handle("/metrics", promhttp.Handler())

	fmt.Println(http.ListenAndServe(address, nil).Error())
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	go handleSigterm(cancel)
	go serveMetrics("0.0.0.0:9099")

	startServerToServeTargets(
		[]*endpoint.Endpoint{
			{
				DNSName:    "abc.example.org",
				Targets:    endpoint.Targets{"1.2.3.4"},
				RecordType: endpoint.RecordTypeA,
				RecordTTL:  180,
			},
		},
		cancel,
	)

	select {
	case <-ctx.Done():
		fmt.Println("exiting")
		return
	}
}
