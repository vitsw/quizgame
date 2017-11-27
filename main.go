package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	pFile := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	shuf := flag.Bool("shuffle", false, "Shuffles quiz problems")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	flag.Parse()

	f, err := os.Open(*pFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	problems, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if *shuf {
		problems = shuffle(problems)
	}

	var (
		total       = len(problems)
		resp        string
		correct     int
		doneChannel = make(chan bool)
	)

	fmt.Print("Press any key to start...")
	_, err = bufio.NewReader(os.Stdin).ReadByte()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for _, problem := range problems {
			fmt.Printf("%s = ", problem[0])
			if _, err := fmt.Scanln(&resp); err != nil {
				log.Fatal(err)
			}
			if resp == strings.ToLower(strings.TrimSpace(problem[1])) {
				correct++
			}
		}
		doneChannel <- true
	}()

	select {
	case <-time.NewTimer(time.Second * time.Duration(*limit)).C:
		fmt.Printf("\nYour score %d from %d", correct, total)
	case <-doneChannel:
		fmt.Printf("\nYour score %d from %d", correct, total)
	}
}

func shuffle(slice [][]string) [][]string {
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	res := make([][]string, len(slice))
	for ind, rI := range rand.Perm(len(slice)) {
		res[ind] = slice[rI]
	}
	return res
}
