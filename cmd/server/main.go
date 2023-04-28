package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rossmerr/events-server/config"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/server/v3/embed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	configFile = flag.String("config", "config.yaml", "The yaml config file to load")
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)
}

func main() {
	flag.Parse()
	flag.PrintDefaults()

	conf := config.NewConfig(*configFile)
	if level, err := log.ParseLevel(conf.LogLevel); err == nil {
		log.SetLevel(level)
	}

	certFilePath, err := checkFile(conf.TLS.CertFile)
	if err != nil {
		log.Fatal(err)
	}

	keyFilePath, err := checkFile(conf.TLS.KeyFile)
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption
	if conf.TLS.UseTlS {
		cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
		if err != nil {
			log.Fatalf("failed to load key pair: %v", err)
		}

		opts = append(opts, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	}
	grpcServer := grpc.NewServer(opts...)

	ctx := context.Background()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	server := &http.Server{Addr: fmt.Sprintf(":%d", conf.HTTP), Handler: nil}

	// http
	go func() {
		fmt.Printf("listening for HTTP on: %d\n", conf.HTTP)
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPC))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// grpc
	go func() {
		fmt.Printf("listening for GRPC on: %d\n", conf.GRPC)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// etcd
	var e *embed.Etcd

	go func() {
		cfg := embed.NewConfig()
		cfg.Dir = conf.Dir // "default.etcd"

		e, err = embed.StartEtcd(cfg)
		if err != nil {
			log.Fatal(err)
		}

		select {
		case <-e.Server.ReadyNotify():
			log.Println("Server is ready!")
		case <-time.After(60 * time.Second):
			e.Server.Stop() // trigger a shutdown
			log.Println("Server took too long to start!")
			done <- true
		}
	}()

	<-done
	fmt.Println("exiting")
	grpcServer.GracefulStop()
	e.Close()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		// handle err
		log.Fatalf("failed to shutdown http server: %v", err)
	}

}

func checkFile(file string) (string, error) {
	if !filepath.IsAbs(file) {
		dir, err := os.Getwd()
		if err != nil {
			return file, fmt.Errorf("failed to current working directory: %w", err)
		}

		file = path.Join(dir, file)
	}
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return file, fmt.Errorf("failed to read file: %s, %w", file, err)
	}
	return file, err
}
