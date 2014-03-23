package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	runcount int = 1 // # of times the exe was run
    command      = flag.String("cmd", "", "The command to run. ex: go run index.go")
)

func main() {
	flag.Parse()
	clearConsole()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File to run: %v \n", *command)

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				//log.Println(ev)
				if ev.IsModify() {
					clearConsole()

					installCmd := exec.Command("go", "install")
					output, err := installCmd.CombinedOutput()
					if err != nil {
						fmt.Printf("\x1b[38;5;%dm%s\x1b[0m\n", 1, err)
						if len(output) > 0 {
							fmt.Println(string(output))
						}
					}

                    commandAndArgs := strings.Split(*command, " ");

					startTime := time.Now()
                    // pass an slice of strings to a varidict param
                    cmd := exec.Command(commandAndArgs[0], commandAndArgs[1:]...)
					output, err = cmd.CombinedOutput()
					endTime := time.Now()

					runMessage := fmt.Sprintf("Run #%d took %v", runcount, endTime.Sub(startTime))
					if err != nil {
						fmt.Printf("\x1b[38;5;%dm%s\x1b[0m\n", 1, runMessage)
						fmt.Print(err)
						if len(output) > 0 {
							fmt.Println(string(output))
						}
					} else {
						fmt.Printf("\x1b[38;5;%dm%s\x1b[0m\n", 2, runMessage)
						fmt.Println(string(output))
					}

					runcount++
				}
			case err := <-watcher.Error:
				log.Println("Error with watcher:", err)
			}
		}

	}()

	path := ".gorunner.tmp"
	err = watcher.Watch(path)
	if err != nil {
		log.Println("Failed to watch", path, ":", err)
	}
	fmt.Println("Watching:", path)

	// loop through all subdirectories and watch them
	//filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
	//	if err != nil {
	//		log.Println("Error walking path:", err)
	//	}

	//	//if !info.IsDir() {
	//	//	if strings.HasSuffix(path, ".go") {
	//	//		err = watcher.Watch(path)
	//	//		if err != nil {
	//	//			log.Println("Failed to watch", path, ":", err)
	//	//		}
	//	//		fmt.Println("Watching:", path)
	//	//	}
	//	//}

	//	if info.IsDir() {
	//		// ignore .git
	//		if strings.Contains(path, ".git") {
	//			return nil
	//		}
	//		err = watcher.Watch(path)
	//		if err != nil {
	//			log.Println("Failed to watch", path, ":", err)
	//		}
	//		fmt.Println("Watching:", path)
	//	}
	//	return nil
	//})

	<-done

	watcher.Close()
}

func clearConsole() {
	os.Stdout.Write([]byte("\033[2J"))
	os.Stdout.Write([]byte("\033[H'"))
	// cmd := exec.Command("clear")
	// output, err := cmd.CombinedOutput()
	// handleCmdError(output, err, "Can't clear screen")
	// log.Printf("%s\n", output)
}

func handleCmdError(output []byte, err error, errMsg string) {
	if err != nil {
		log.Printf("%s.", errMsg)
		if len(output) > 0 {
			log.Println(string(output))
		}
		logRed(err)
		log.Fatal()
	}
}

func logGreen(msg string) {
	logWithColor(msg, 2)
}

func logBlue(msg string) {
	logWithColor(msg, 4)
}

func logRed(msg interface{}) {
	logWithColor(msg, 1)
}

func logWithColor(msg interface{}, colorCode int) {
	log.Printf("\x1b[38;5;%dm%s\x1b[0m\n", colorCode, msg)
}
