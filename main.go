package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/evex-dev/promo2/src/connect"
)

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
		proxylist = append(proxylist, scanner.Text())
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
		fmt.Println("\x1b[33m[=] Requesting...\x1b[0m")
		
		proxyUrl := ""

		if len(proxylist) > 0 {
			proxyUrl = proxylist[rand.Intn(len(proxylist))]
		}

		promos, err := connect.GetPromoUrls(proxyUrl)

		if err != nil {
			fmt.Println("\x1b[31m[-] Error: ", err, "\x1b[0m")
			continue
		}

		resultContent := strings.Join(promos, "\n")
		resultFile, err := os.Open(resultPath)

		if err != nil {
			fmt.Println("\x1b[31m[-] Error: ", err, "\x1b[0m")
			continue
		}

		scanner := bufio.NewScanner(resultFile)
		var resultlist []string
		for scanner.Scan() {
			resultlist = append(resultlist, scanner.Text())
		}

		resultFile.WriteString(fmt.Sprintf("%s\n%s", strings.Join(resultlist, "\n"), resultContent))

		for _, promo := range promos {
			fmt.Println("\x1b[32m[+] Found: ", promo, "\x1b[0m")
		}
	}
}
