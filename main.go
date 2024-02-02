package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

const (
	inputFile  = "./.env"
	outputFile = "./runtime-env.js"
)

func main() {

	in, err := os.Open(inputFile)

	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(outputFile)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	_, err = out.WriteString(generate(os.Environ(), in))

	if err != nil {
		log.Fatal(err)
	}
}

func generate(env []string, cfg io.Reader) string {
	i, err := io.ReadAll(cfg)

	if err != nil {
		log.Fatal(err)
	}

	c := parseCfg(i)
	e := parseEnv(env)

	for k := range e {
		if !slices.Contains(c, k) {
			delete(e, k)
		}
	}

	j, err := json.Marshal(e)

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("window.__RUNTIME_CONFIG__ = %s;", j)
}

func parseCfg(cfg []byte) []string {
	a := []string{}

	for _, v := range strings.Split(strings.ReplaceAll(string(cfg), "\r\n", "\n"), "\n") {
		p := strings.SplitN(v, "=", 2)
		a = append(a, p[0])
	}
	return a
}

func parseEnv(env []string) map[string]string {
	m := make(map[string]string, len(env))

	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		m[p[0]] = p[1]
	}

	return m
}
