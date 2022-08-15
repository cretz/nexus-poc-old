package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("missing command name")
	}
	switch os.Args[1] {
	case "gen-protos":
		return genProtos()
	default:
		return fmt.Errorf("invalid commmand")
	}
}

func genProtos() error {
	_, currFile, _, _ := runtime.Caller(0)
	nexusGoDir := filepath.Join(currFile, "../../")
	apiDir := filepath.Join(currFile, "../../../../api")

	frontendProtos, err := filepath.Glob(filepath.Join(apiDir, "nexus/frontend/v1/*.proto"))
	if err != nil {
		return err
	}
	backendProtos, err := filepath.Glob(filepath.Join(apiDir, "nexus/backend/v1/*.proto"))
	if err != nil {
		return err
	}

	protocArgs := []string{
		"-I", apiDir,
		"--go_out", nexusGoDir,
		"--go_opt", "module=github.com/cretz/nexus-poc/sdk/go",
		"--go-grpc_out", nexusGoDir,
		"--go-grpc_opt", "module=github.com/cretz/nexus-poc/sdk/go",
		"--grpc-gateway_out", nexusGoDir,
		"--grpc-gateway_opt", "module=github.com/cretz/nexus-poc/sdk/go",
	}
	for _, proto := range frontendProtos {
		relProto := strings.TrimPrefix(filepath.ToSlash(strings.TrimPrefix(proto, apiDir)), "/")
		arg := "M" + relProto + "=github.com/cretz/nexus-poc/sdk/go/nexus/frontend/frontendpb"
		protocArgs = append(protocArgs, "--go_opt", arg, "--go-grpc_opt", arg, "--grpc-gateway_opt", arg)
	}
	for _, proto := range backendProtos {
		relProto := strings.TrimPrefix(filepath.ToSlash(strings.TrimPrefix(proto, apiDir)), "/")
		arg := "M" + relProto + "=github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
		protocArgs = append(protocArgs, "--go_opt", arg, "--go-grpc_opt", arg, "--grpc-gateway_opt", arg)
	}
	protocArgs = append(protocArgs, frontendProtos...)
	protocArgs = append(protocArgs, backendProtos...)

	log.Printf("Running protoc with args %v", protocArgs)
	cmd := exec.Command("protoc", protocArgs...)
	cmd.Stdin, cmd.Stderr, cmd.Stdout = os.Stdin, os.Stderr, os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("protoc failed: %w", err)
	}
	return nil
}
