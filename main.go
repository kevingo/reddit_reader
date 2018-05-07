package main

import (
	"fmt"
	"flag"
	"errors"
	"net/http"
	"encoding/json"
	"strconv"
)

type Response struct {
	Data struct {
		Children []struct {
			Data Item
		}
	}
}

type Item struct {
	Title    string
	URL      string
	Comments int `json:"num_comments"`
}

var (
	Topic   = flag.String("t", "programming", "")
	Number    = flag.Int("n", 2, "Default pages to fetched")
)

func main() {
	flag.Parse()
	items, err := Get(*Topic, *Number)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(items)
}

func Get(topic string, number int) ([]Item, error) {
	n := strconv.Itoa(number)
	url := fmt.Sprintf("http://reddit.com/r/%s.json?limit=%s", topic, n)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return nil, err
	}

	items := make([]Item, len(r.Data.Children))
	for i, child := range r.Data.Children {
		items[i] = child.Data
	}
	return items, nil
}