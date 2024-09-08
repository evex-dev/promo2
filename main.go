package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/evex-dev/promo2/src/connect"
	"github.com/manifoldco/promptui"
)

var wg sync.WaitGroup

func main() {
	validatePath := func(input string) error {
		if _, err := os.Stat(input); err != nil {
			return err
		}

		return nil
	}

	proxiesPathPrompt := promptui.Prompt{
		Label:    "Proxies Path",
		Validate: validatePath,
	}

	proxiesPath, err := proxiesPathPrompt.Run()

	if err != nil {
		fmt.Println(err)
		return
	}

	proxiesBody, err := os.Open(proxiesPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(proxiesBody)
	var proxylist []string
	for scanner.Scan() {
		text := scanner.Text()
		proxylist = append(proxylist, text)
	}

	resultPathPrompt := promptui.Prompt{
		Label:    "Result Path",
		Validate: validatePath,
	}

	resultPath, err := resultPathPrompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		oneGenerate := func() {
			defer wg.Done()

			fmt.Println("\x1b[33m[=] Requesting...\x1b[0m")

			promos, err := connect.GetPromoUrls(proxylist)

			if err != nil {
				fmt.Println("\x1b[31m[-] Error: ", err, "\x1b[0m")
				return
			}

			fmt.Println("\x1b[32m[+] Found: ", len(promos), "urls\x1b[0m")

			resultContent := strings.Join(promos, "\n")
			resultFile, err := os.Create(resultPath)

			if err != nil {
				fmt.Println("\x1b[31m[-] Error: ", err, "\x1b[0m")
				return
			}

			defer resultFile.Close()

			if len(promos) != 0 {
				resultList, err := getResultFileContent(resultPath)

				if err != nil {
					fmt.Println("\x1b[31m[-] Error: ", err, "\x1b[0m")
					return
				}

				fmt.Println(len(resultList))
				resultFile.WriteString(resultList + "\n" + resultContent)
			}

			for _, promo := range promos {
				fmt.Println("\x1b[32m[+] Found: ", promo, "\x1b[0m")
			}
		}

		thread := 1

		for i := 0; i < thread; i++ {
			wg.Add(1)
			go oneGenerate()
		}

		wg.Wait()
	}
}

func getResultFileContent(resultPath string) (string, error) {
	resultBody, err := os.Open(resultPath)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer resultBody.Close()

	resultScanner := bufio.NewScanner(resultBody)
	var resultList []string
	for resultScanner.Scan() {
		text := resultScanner.Text()
		resultList = append(resultList, text)
	}

	return strings.Join(resultList, "\n"), nil
}
