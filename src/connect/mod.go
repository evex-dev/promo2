package connect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type RequestParms struct {
	UserUUID   string `json:"userUuid"`
	CampaignID string `json:"campaignId"`
}

type Games struct {
	Players []Player `json:"players"`
}

type Player struct {
	UUID string `json:"uuid"`
}

func GetPromoUrls(proxylist []string) ([]string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	setupReq, err := http.NewRequest("GET", "https://www.chess.com/service/gamelist/top?limit=50&from="+generateRandomIntString(1040), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to connect setupReqURL")
	}

	resp, err := client.Do(setupReq)
	if err != nil {
		return nil, fmt.Errorf("failed to connect setupReqURL")
	}

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response")
	}
	_ = resp.Body.Close()

	var games []Games

	err = json.Unmarshal(payload, &games)
	if err != nil {
		return nil, fmt.Errorf("games unmarshal error: %w", err)
	}

	promoUrls := []string{}

	var wg sync.WaitGroup

	for _, game := range games {
		for _, user := range game.Players {
			wg.Add(1)
			go func(user Player) {
				defer wg.Done()

				fmt.Println("\x1b[34m[/] Fetching UUID: " + user.UUID + "\x1b[0m")

				body, err := json.Marshal(RequestParms{
					UserUUID:   user.UUID,
					CampaignID: "4daf403e-66eb-11ef-96ab-ad0a069940ce",
				})
				if err != nil {
					fmt.Println(err)
					return
				}

				proxyUrl, err := url.Parse(proxylist[rand.Intn(len(proxylist))])
				if err != nil {
					fmt.Println("\x1b[31m[-] Error: ", err, "\x1b[0m")
					return
				}

				client := &http.Client{
					Timeout: time.Second * 10,
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyUrl),
					},
				}

				req, err := http.NewRequest("POST", "https://www.chess.com/rpc/chesscom.partnership_offer_codes.v1.PartnershipOfferCodesService/RetrieveOfferCode", bytes.NewBuffer(body))
				if err != nil {
					return
				}

				req.Header.Set("accept", "application/json, text/plain, */*")
				req.Header.Set("accept-language", "ja,en-US;q=0.9,en;q=0.8")
				req.Header.Set("cache-control", "no-cache")
				req.Header.Set("content-type", "application/json")
				req.Header.Set("pragma", "no-cache")
				req.Header.Set("priority", "u=1, i")
				req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Google Chrome\";v=\"128\"")
				req.Header.Set("sec-ch-ua-mobile", "?0")
				req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
				req.Header.Set("sec-fetch-dest", "empty")
				req.Header.Set("sec-fetch-mode", "cors")
				req.Header.Set("sec-fetch-site", "same-origin")
				req.Header.Set("referrer", "https://www.chess.com/play/computer/discord-wumpus?utm_source=chesscom&utm_medium=homepagebanner&utm_campaign=discord2024")

				promoResp, err := client.Do(req)

				if err != nil {
					return
				}

				if promoResp.StatusCode != 200 {
					fmt.Println("\x1b[31m[-] Fetching Error: " + strconv.Itoa(resp.StatusCode) + "\x1b[0m" + "Status: " + resp.Status)

					return
				}

				fmt.Println("\x1b[32m[+] Found: " + user.UUID + "\x1b[0m")
				defer promoResp.Body.Close()
				responseBody, err := io.ReadAll(promoResp.Body)
				if err != nil {
					fmt.Println(err)
					return
				}

				var data map[string]interface{}
				err = json.Unmarshal(responseBody, &data)
				if err != nil {
					fmt.Println(err)
					return
				}
				if data["codeValue"] != nil {
					content_data := "https://discord.com/billing/promotions/" + data["codeValue"].(string)

					promoUrls = append(promoUrls, content_data)
					return
				}
				
			}(user)
		}
	}

	wg.Wait()

	return promoUrls, nil
}

func generateRandomIntString(max int) string {
	return strconv.Itoa(rand.Intn(max))
}
