package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/bogdanovandreycode/agentboard/internal/httpapi"
	"github.com/bogdanovandreycode/agentboard/internal/mcpserver"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/bogdanovandreycode/agentboard/internal/webui"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var version = "0.1.0"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "agentboard:", err)
		os.Exit(1)
	}
}
func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}
	switch args[0] {
	case "init":
		return initCommand(args[1:])
	case "serve":
		return serveCommand(args[1:], false)
	case "open":
		return serveCommand(args[1:], true)
	case "mcp":
		return mcpCommand(args[1:])
	case "version", "--version", "-v":
		fmt.Println("AgentBoard", version)
		return nil
	default:
		return usage()
	}
}
func usage() error {
	fmt.Fprintln(os.Stderr, "Usage: agentboard <init|serve|open|mcp|version>")
	return errors.New("unknown or missing command")
}

func initCommand(args []string) error {
	f := flag.NewFlagSet("init", flag.ContinueOnError)
	db := f.String("db", "", "database path")
	if err := f.Parse(args); err != nil {
		return err
	}
	projectPath := "."
	if f.NArg() > 0 {
		projectPath = f.Arg(0)
	}
	st, err := persistence.Open(*db)
	if err != nil {
		return err
	}
	defer st.Close()
	p, err := st.InitProject(context.Background(), projectPath)
	if err != nil {
		return err
	}
	fmt.Printf("Initialized AgentBoard project %s\nProject ID: %s\nPath: %s\n", p.Name, p.ID, p.Path)
	return nil
}

func serveCommand(args []string, open bool) error {
	f := flag.NewFlagSet("serve", flag.ContinueOnError)
	addr := f.String("addr", "127.0.0.1:7337", "listen address")
	db := f.String("db", "", "database path")
	if err := f.Parse(args); err != nil {
		return err
	}
	if f.NArg() > 1 || (!open && f.NArg() != 0) {
		return errors.New("usage: agentboard open [--addr address] [--db path] [project-path]")
	}
	url := "http://" + *addr
	if open {
		path := "."
		if f.NArg() == 1 {
			path = f.Arg(0)
		}
		var err error
		url, err = launcherURL(url, path)
		if err != nil {
			return err
		}
		// Reuse an existing AgentBoard server without starting another process.
		if agentBoardRunning("http://" + *addr) {
			return openURL(url)
		}
	}
	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	st, err := persistence.Open(*db)
	if err != nil {
		return err
	}
	defer st.Close()
	svc := service.New(st)
	apiHandler, uiHandler := httpapi.New(svc), webui.Handler()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			apiHandler.ServeHTTP(w, r)
			return
		}
		uiHandler.ServeHTTP(w, r)
	})
	server := &http.Server{Addr: *addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
	if open {
		go func() {
			time.Sleep(250 * time.Millisecond)
			if err := openURL(url); err != nil {
				log.Printf("open browser: %v", err)
			}
		}()
	}
	log.Printf("AgentBoard %s listening on %s", version, url)
	err = server.Serve(listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func mcpCommand(args []string) error {
	f := flag.NewFlagSet("mcp", flag.ContinueOnError)
	project := f.String("project", "", "managed project path")
	worker := f.String("worker", "", "worker id or slug")
	db := f.String("db", "", "database path")
	if err := f.Parse(args); err != nil {
		return err
	}
	if *project == "" || *worker == "" {
		return errors.New("mcp requires --project and --worker")
	}
	abs, err := filepath.Abs(*project)
	if err != nil {
		return err
	}
	st, err := persistence.Open(*db)
	if err != nil {
		return err
	}
	defer st.Close()
	svc := service.New(st)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	agent, err := svc.BeginAgentSession(ctx, abs, *worker, "stdio MCP")
	if err != nil {
		return err
	}
	defer svc.EndAgentSession(context.Background(), agent)
	server := mcpserver.New(svc, agent, version)
	return server.Run(ctx, &mcp.StdioTransport{})
}

func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
