package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func die(message string) {
	_, _ = os.Stderr.WriteString(message)
	os.Exit(1)
}

type dep struct {
	source string
	target string
	live   bool
}

func main() {
	filter := func(key string) bool {
		return true
	}

	if len(os.Args) == 2 && len(os.Args[1]) > 0 {
		exp, err := regexp.Compile(os.Args[1])
		if err != nil {
			die("Invalid filter express. Must be an RE2 regex.")
		}
		filter = func(key string) bool {
			return exp.MatchString(key)
		}
	}

	// depMap := map[string]string{}
	deps := []dep{}

	cmd := exec.Command("go", "mod", "graph")
	r, err := cmd.StdoutPipe()
	if err != nil {
		die(err.Error())
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)

	toCheck := []string{}
	checked := map[string]bool{}

	go func() {
		if err := cmd.Start(); err != nil {
			die(err.Error())
		}
	}()

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		live := filter(parts[1])
		deps = append(deps, dep{parts[0], parts[1], live})
		// need to check any parents of items that match the filter
		// unless that item is already live
		if live && !filter(parts[0]) {
			toCheck = append(toCheck, parts[0])
			checked[parts[1]] = true
			checked[parts[0]] = true
		}
	}

	for len(toCheck) > 0 {
		key := toCheck[len(toCheck)-1]
		toCheck = toCheck[:len(toCheck)-1]
		checked[key] = true
		for i, d := range deps {
			if d.target == key && !d.live {
				deps[i].live = true
				if _, ok := checked[d.source]; !ok {
					toCheck = append(toCheck, d.source)
				}
			}
		}
	}

	fmt.Println("digraph {")
	for _, d := range deps {
		if !d.live {
			continue
		}
		fmt.Print("    \"")
		fmt.Print(d.source)
		fmt.Print("\" -> \"")
		fmt.Print(d.target)
		fmt.Print("\"\n")
	}
	fmt.Println("}")
}
