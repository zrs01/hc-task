package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	exec "github.com/alexellis/go-execute/pkg/v1"
	"github.com/urfave/cli/v2"
)

var version = "development"
var wg sync.WaitGroup
var lock sync.Mutex

func main() {
	cliapp := cli.NewApp()
	cliapp.Name = "task"
	cliapp.Usage = "Execute concurrent tasks"
	cliapp.Version = version

	var ifile, ofile, cmd string
	var max int
	cliapp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "input",
			Aliases:     []string{"i"},
			Usage:       "input file with content separated by enter",
			Required:    true,
			Destination: &ifile,
		},
		&cli.StringFlag{
			Name:        "command",
			Aliases:     []string{"c"},
			Usage:       "command to be executed",
			Required:    true,
			Destination: &cmd,
		},
		&cli.IntFlag{
			Name:        "tasks",
			Usage:       "max tasks execute concurrently",
			Required:    false,
			Value:       5,
			Destination: &max,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "output file",
			Required:    false,
			Value:       "output.txt",
			Destination: &ofile,
		},
	}
	cliapp.Action = func(ctx *cli.Context) error {
		return execute(max, ifile, ofile, cmd)
	}
	if err := cliapp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func execute(max int, ifile, ofile, cmd string) error {
	fbytes, err := os.ReadFile(ifile)
	if err != nil {
		return err
	}
	if len(fbytes) == 0 {
		return fmt.Errorf("empty content in %s", ifile)
	}

	// delete the output file if it exists
	os.Remove(ofile)

	// make the max concurrent tasks
	log.Printf("max concurrent: %d", max)
	guard := make(chan int, max)
	items := strings.Split(strings.ReplaceAll(string(fbytes), "\r\n", "\n"), "\n")
	for _, item := range items {
		if strings.TrimSpace(item) != "" {
			guard <- 1
			wg.Add(1)
			go executeTask(guard, item, cmd, ofile)
		}
	}
	wg.Wait()
	return nil
}

func executeTask(guard chan int, item, cmd, ofile string) {
	defer func() {
		<-guard
		wg.Done()
	}()

	cmdSplit := strings.Split(cmd, " ")
	cmdSplit = append(cmdSplit, item)
	execCmd := exec.ExecTask{
		Command:     cmdSplit[0],
		Args:        cmdSplit[1:],
		StreamStdio: true,
		// PrintCommand: true,
	}

	start := time.Now()
	// run the command
	res, err := execCmd.Execute()
	if err != nil {
		log.Printf("Error: %s", err)
	}
	// randomSleep(item)
	execResult := "Success"
	if res.ExitCode > 0 {
		execResult = "Failure"
	}

	end := time.Now()
	elapsed := time.Since(start)

	lock.Lock()
	defer lock.Unlock()
	f, err := os.OpenFile(ofile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%-20s: [%s] %s - %s (%s)\n", item, execResult, start.Format("2006-01-02 15:04:05.000"), end.Format("2006-01-02 15:04:05.000"), elapsed)); err != nil {
		log.Printf("Error: %s", err)
	}
}

// func randomSleep(item string) {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	n := r.Intn(5) // n will be between 0 and 5
// 	log.Printf("[%s] Sleeping %d seconds...\n", item, n)
// 	time.Sleep(time.Duration(n) * time.Second)
// 	log.Printf("[%s] Done\n", item)
// }
