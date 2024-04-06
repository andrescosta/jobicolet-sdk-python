package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/andrescosta/goico/pkg/runtimes/wasm"
)

func main() {
	ctx := context.Background()
	runtime, err := wasm.NewRuntimeWithCompilationCache("./cache")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing Wazero: %v\n", err)
		os.Exit(1)
	}
	defer runtime.Close(ctx)
	dir := os.Args[1]
	wasmf, err := os.ReadFile(dir + "/python.wasm")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading wasm binary: %v\n", err)
		os.Exit(1)
	}
	mounts := []string{
		dir + "/lib/python3.13:/usr/local/lib/python3.13:ro",
		dir + "/sdk:/usr/local/lib/jobico:ro",
		dir + "/s:/s",
	}
	args := []string{
		dir + "/python.wasm",
		dir + "s/main.py",
		"aaa",
	}
	buffIn := &bytes.Buffer{}
	buffOut := &bytes.Buffer{}
	buffErr := &bytes.Buffer{}
	e := []wasm.EnvVar{{Key: "PYTHONPATH", Value: "/usr/local/lib/jobico"}}
	modi, err := wasm.NewIntModule(ctx, runtime, wasmf, Log, mounts, args, e, buffIn, buffOut, buffErr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error instantiating the module: %v\n", err)
		fmt.Printf("Dump Error: %s\n", buffErr.String())
		fmt.Printf("Dump Std out: %s\n", buffOut.String())
		os.Exit(1)
	}
	defer modi.Close(ctx)
	err = modi.Run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error instantiating the module: %v\n", err)
		fmt.Printf("Dump Error: %s\n", buffErr.String())
		fmt.Printf("Dump Std out: %s\n", buffOut.String())
		os.Exit(1)
	}
	var res int32
	binary.Read(buffOut, binary.LittleEndian, &res)
	fmt.Printf("Result:\n%d\n", res)
	msgRes, err := io.ReadAll(buffOut)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Text:\n%s\n", msgRes)
	var lvl uint8
	var size uint32
	binary.Read(buffErr, binary.LittleEndian, &lvl)
	binary.Read(buffErr, binary.LittleEndian, &size)
	fmt.Printf("Level:%d\n", lvl)
	fmt.Printf("Size:%d\n", size)
	msgLog := make([]byte, size)
	_, err = buffErr.Read(msgLog)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Text:%s\n", msgLog)
}

func Log(context.Context, uint32, string) error {
	return nil
}
