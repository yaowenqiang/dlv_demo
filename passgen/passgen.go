package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"
)

const (
	alphanum = "apcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	symbols = "~!@#$%^&*(){}[]<>?/"
	alphanumsym = alphanum + symbols
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type logMessage struct {
	level string
	message string
}

const (
	logLevelInfo = "INFO"
	logLevelWarning = "WARNING"
	logLevelError = "ERROR"
)

var wg = sync.WaitGroup{}

func main() {
	logCh := make(chan logMessage, 2)
	wg.Add(1)
	go logger(logCh)

	logCh <- logMessage{logLevelInfo, "Starting app"}

	length := flag.Int("len", 10, "Length of the password to generate")
	useSymbols := flag.Bool("sym", true, "Include symbols in generated password")

	flag.Parse()

	password := genPassword(*length, *useSymbols)
	fmt.Println("The generated password is:")
	for _, val := range password {
		fmt.Print(string(val))
	}


	fmt.Println()

	logCh <- logMessage{logLevelInfo, "Closing app"}

	close(logCh)

	wg.Wait()
}


func logger(logCh <-chan logMessage) {
	const fileName = "log.txt"
	var logFile *os.File
	logFile, err := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}

	l := log.New(logFile, "", log.LstdFlags)

	for msg := range logCh {
		l.Printf("[%v] %v", msg.level, msg.message)
	}

	logFile.Close()
	wg.Done()

}

func genPassword(length int, useSymbols bool) []byte {
	result := make([]byte, length)
	order := make([]float64, len(result))

	var source string
	if useSymbols {
		source = alphanumsym
	} else {
		source = alphanum
	}


	for i := range result {
		result[i] = source[int(rnd.Float64() * float64(len(source)))]
		order[i] = rnd.Float64()
	}

	sort.Slice(result, func(l, r int) bool  {
		return order[l] < order[r]

	})

	return result
}
