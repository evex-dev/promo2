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
		resultFileContent, err := getResultFileContent(resultPath);
		if err != nil {
			fmt.Println("\x1b[31m[-] Error: ", err, "\x1b[0m")
			return
		}

		promosList := strings.Split(resultFileContent, "\n");

		oneGenerate := func() {
			defer wg.Done()

			fmt.Println("\x1b[33m[=] Requesting...\x1b[0m")

			promos, err := connect.GetPromoUrls(proxylist)

			if err != nil {
				fmt.Println("\x1b[31m[-] Error: ", err, "\x1b[0m")
				return
			}

			promosLength := len(promos)

			fmt.Println("\x1b[32m[+] Found: ", promosLength, "urls\x1b[0m")

			for _, promo := range promos {
				fmt.Println("\x1b[32m[+] Found: ", promo, "\x1b[0m")
			}

			promosList = append(promosList, promos...)

			if promosLength != 0 {
				resultFile, err := os.Create(resultPath)

				if err != nil {
					fmt.Println("\x1b[31m[-] Error: ", err, "\x1b[0m")
					return
				}

				resultFile.WriteString(strings.Join(promosList, "\n"))

				resultFile.Close()
			}
		}

		thread := 30

		for i := 0; i < thread; i++ {
			wg.Add(1)
			go oneGenerate()
		}

		wg.Wait()
	}
}

func getResultFileContent(resultPath string) (string, error) {
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		return "", err
	}

	content, err := os.ReadFile(resultPath)
    if err != nil {
        return "", err
    }

    return string(content), nil
}
