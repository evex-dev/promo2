package connect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
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

func GetPromoUrls(proxyUrl string) ([]string, error) {
	prxoyURL, err := url.Parse(proxyUrl)

	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %s", err)
	}

	transport := http.DefaultTransport

	if proxyUrl != "" {
		transport = &http.Transport{
			Proxy: http.ProxyURL(prxoyURL),
		}
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	headers := http.Header{}
	headers.Add("accept", "application/json, text/plain, */*")
	headers.Add("accept-language", "ja,en-US;q=0.9,en;q=0.8")
	headers.Add("cache-control", "no-cache")
	headers.Add("content-type", "application/json")
	headers.Add("pragma", "no-cache")
	headers.Add("priority", "u=1, i")
	headers.Add("sec-ch-ua", "\"Chromium\";v=\""+generateRandomIntString(200)+"\", \"Not;A=Brand\";v=\"24\", \"Google Chrome\";v=\"128\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "\"Windows\"")
	headers.Add("sec-fetch-dest", "empty")
	headers.Add("sec-fetch-mode", "cors")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("referrer", "https://www.chess.com/play/computer/discord-wumpus?utm_source=chesscom&utm_medium=homepagebanner&utm_campaign=discord2024")

	setupReq, err := http.NewRequest("GET", "https://www.chess.com/service/gamelist/top?limit=50&from=" + generateRandomIntString(1040), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect setupReqURL")
	}

	setupReq.Header = headers

	resp, err := client.Do(setupReq)
	if err != nil {
		return nil, fmt.Errorf("failed to connect setupReqURL")
	}

	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response")
	}

	var games []Games

	err = json.Unmarshal(payload, &games)
	if err != nil {
		return nil, fmt.Errorf("games unmarshal error")
	}

	promoUrls := []string{}

	for _, game := range games {
		for _, user := range game.Players {
			body, err := json.Marshal(RequestParms{
				UserUUID:   user.UUID,
				CampaignID: "4daf403e-66eb-11ef-96ab-ad0a069940ce",
			})
			if err != nil {
				return nil, fmt.Errorf("unknown Error")
			}

			req, err := http.NewRequest("POST", "https://www.chess.com/rpc/chesscom.partnership_offer_codes.v1.PartnershipOfferCodesService/RetrieveOfferCode", bytes.NewBuffer(body))
			if err != nil {
				return nil, fmt.Errorf("unknown Error")
			}

			req.Header = headers

			resp, err = client.Do(req)
			if err != nil {
				return nil, fmt.Errorf("connection Error")
			}

			if resp.Status == "200" {
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("")
				}
				var data map[string]interface{}
				err = json.Unmarshal(body, &data)
				if err != nil {
					return nil, fmt.Errorf("unknown Error")
				}
				if data["codeValue"] != nil {
					content_data := "https://discord.com/billing/promotions/" + data["codeValue"].(string)

					promoUrls = append(promoUrls, content_data)
				}
			}
		}
	}

	return promoUrls, nil
}

func generateRandomIntString(max int) string {
	return string(rand.Intn(max))
}
