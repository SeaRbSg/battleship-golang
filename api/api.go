package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dan-v/seattlerb-battleship/common"
)

type API struct {
	httpClient *http.Client
	baseURL    string
	GameID     int
}

func NewAPI(baseURL string) (*API, error) {
	client := &API{
		httpClient: &http.Client{Timeout: (time.Second * 20)},
		baseURL:    baseURL,
	}
	err := client.setup()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *API) setup() error {
	resp, err := c.httpClient.Get(c.baseURL + "new_game")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	game := &common.NewGameResponse{}
	if err := json.NewDecoder(resp.Body).Decode(game); err != nil {
		return err
	}
	c.GameID = game.ID
	return nil
}

func (c *API) Attack(guess string, opponentResponse *common.Response) (*common.Turn, error) {
	if opponentResponse == nil {
		opponentResponse = &common.Response{
			Hit:  false,
			Lost: false,
		}
	}
	data := common.Turn{
		GameID:   c.GameID,
		Response: *opponentResponse,
		Guess:    common.Guess{guess},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payload)

	req, err := http.NewRequest("POST", c.baseURL+"turn", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	turn := &common.Turn{}
	if err := json.NewDecoder(strings.NewReader(string(responseData))).Decode(turn); err != nil {
		return nil, err
	}
	return turn, nil
}
