package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	serverName = "caffeinate-mcp"
)

type CaffeinateServer struct {
	srv             *server.MCPServer
	activeProcesses map[string]*exec.Cmd
	mu              sync.Mutex
}

func NewCaffeinateServer() *CaffeinateServer {
	// Create server with tool capabilities enabled
	srv := server.NewMCPServer(
		serverName,
		version,
		server.WithToolCapabilities(true),
	)

	cs := &CaffeinateServer{
		srv:             srv,
		activeProcesses: make(map[string]*exec.Cmd),
	}

	// Define and add tools
	cs.setupTools()

	return cs
}

func (s *CaffeinateServer) setupTools() {
	// caffeinate_start tool
	startTool := mcp.NewTool(
		"caffeinate_start",
		mcp.WithDescription("Start caffeinate to prevent system sleep"),
		mcp.WithBoolean("display", mcp.Description("Prevent display from sleeping (-d flag)")),
		mcp.WithBoolean("idle", mcp.Description("Prevent system from idle sleeping (-i flag)")),
		mcp.WithBoolean("disk", mcp.Description("Prevent disk from idle sleeping (-m flag)")),
		mcp.WithBoolean("system", mcp.Description("Prevent system from sleeping when on AC power (-s flag)")),
		mcp.WithBoolean("user", mcp.Description("Declare user is active (-u flag)")),
		mcp.WithNumber("timeout", mcp.Description("Timeout in seconds (-t flag)")),
		mcp.WithNumber("pid", mcp.Description("Wait for process with specified PID to exit (-w flag)")),
	)
	s.srv.AddTool(startTool, s.handleCaffeinateStart)

	// caffeinate_stop tool
	stopTool := mcp.NewTool(
		"caffeinate_stop",
		mcp.WithDescription("Stop a caffeinate process"),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the caffeinate process to stop")),
	)
	s.srv.AddTool(stopTool, s.handleCaffeinateStop)

	// caffeinate_list tool
	listTool := mcp.NewTool(
		"caffeinate_list",
		mcp.WithDescription("List active caffeinate processes"),
	)
	s.srv.AddTool(listTool, s.handleCaffeinateList)
}

func (s *CaffeinateServer) handleCaffeinateStart(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{}

	// Use the new API methods to parse parameters
	if display := request.GetBool("display", false); display {
		args = append(args, "-d")
	}
	if idle := request.GetBool("idle", false); idle {
		args = append(args, "-i")
	}
	if disk := request.GetBool("disk", false); disk {
		args = append(args, "-m")
	}
	if system := request.GetBool("system", false); system {
		args = append(args, "-s")
	}
	if user := request.GetBool("user", false); user {
		args = append(args, "-u")
	}

	if timeout := request.GetInt("timeout", 0); timeout > 0 {
		args = append(args, "-t", strconv.Itoa(timeout))
	}

	if pid := request.GetInt("pid", 0); pid > 0 {
		args = append(args, "-w", strconv.Itoa(pid))
	}

	cmd := exec.Command("caffeinate", args...)
	if err := cmd.Start(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start caffeinate: %v", err)), nil
	}

	id := fmt.Sprintf("%d_%d", cmd.Process.Pid, time.Now().Unix())

	s.mu.Lock()
	s.activeProcesses[id] = cmd
	s.mu.Unlock()

	go func() {
		cmd.Wait()
		s.mu.Lock()
		delete(s.activeProcesses, id)
		s.mu.Unlock()
	}()

	return mcp.NewToolResultText(fmt.Sprintf("Started caffeinate with ID: %s (PID: %d)\nFlags: %s", id, cmd.Process.Pid, strings.Join(args, " "))), nil
}

func (s *CaffeinateServer) handleCaffeinateStop(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required"), nil
	}

	s.mu.Lock()
	cmd, exists := s.activeProcesses[id]
	s.mu.Unlock()

	if !exists {
		return mcp.NewToolResultError(fmt.Sprintf("No caffeinate process found with ID: %s", id)), nil
	}

	if err := cmd.Process.Kill(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop caffeinate: %v", err)), nil
	}

	s.mu.Lock()
	delete(s.activeProcesses, id)
	s.mu.Unlock()

	return mcp.NewToolResultText(fmt.Sprintf("Stopped caffeinate process with ID: %s", id)), nil
}

func (s *CaffeinateServer) handleCaffeinateList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.activeProcesses) == 0 {
		return mcp.NewToolResultText("No active caffeinate processes"), nil
	}

	var result strings.Builder
	result.WriteString("Active caffeinate processes:\n")
	for id, cmd := range s.activeProcesses {
		result.WriteString(fmt.Sprintf("- ID: %s, PID: %d\n", id, cmd.Process.Pid))
	}

	return mcp.NewToolResultText(result.String()), nil
}

func main() {
	// Check for --version flag
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("%s version %s (commit: %s, built at: %s)\n", serverName, version, commit, date)
		os.Exit(0)
	}

	cs := NewCaffeinateServer()

	// Start the stdio server
	if err := server.ServeStdio(cs.srv); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	// Clean up any remaining processes
	cs.mu.Lock()
	for _, cmd := range cs.activeProcesses {
		cmd.Process.Kill()
	}
	cs.mu.Unlock()

	os.Exit(0)
}
