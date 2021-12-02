package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Record struct {
	Used       int   `json:"used"`
	Available  int   `json:"available"`
	RecordedAt int64 `json:"recorded_at"`
}

type DataCollection struct {
	Records []Record `json:"records"`
}

type GistFile struct {
	Name    string `json:"filename"`
	Type    string `json:"type"`
	Lang    string `json:"language"`
	Content string `json:"content"`
}

type Gist struct {
	Id    string              `json:"id"`
	Files map[string]GistFile `json:"files"`
}

func apiCall(method string, uri string, body io.Reader) []byte {
	username := os.Getenv("GITHUB_USERNAME")
	token := os.Getenv("GITHUB_TOKEN")
	url := "https://api.github.com" + uri
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return data
}

func getGist() *Gist {
	gistId := os.Getenv("GIST_ID")
	res := apiCall("GET", "/gists/"+gistId, nil)
	gist := Gist{}
	json.Unmarshal(res, &gist)
	return &gist
}

func saveGist(g *Gist) {
	data, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}
	apiCall("PATCH", "/gists/"+g.Id, bytes.NewReader(data))
}

func latestData(d *DataCollection) (int, int) {
	fmt.Println("Checking latest")
	max := int64(0)
	var lu, lr int

	for _, item := range d.Records {
		if item.RecordedAt > max {
			max = item.RecordedAt
			lu = item.Used
			lr = item.Available
		}
	}

	return lu, lr
}

func addData(used int, remain int, gist *Gist) {

	fileName := os.Getenv("FILENAME")

	for _, f := range gist.Files {
		dc := DataCollection{}
		if f.Name == fileName {
			err := json.Unmarshal([]byte(f.Content), &dc)
			if err != nil {
				panic(err)
			}

			if len(dc.Records) > 0 {
				lu, lr := latestData(&dc)
				if used == lu && remain == lr {
					fmt.Println("No need to update")
					continue
				}
			}
			fmt.Println("Updating data")

			rec := Record{}
			rec.Available = remain
			rec.Used = used
			rec.RecordedAt = time.Now().UnixNano() / 1000000
			dc.Records = append(dc.Records, rec)

			cnt, err := json.Marshal(dc)
			if err != nil {
				panic(err)
			}
			f.Content = string(cnt)
			gist.Files[fileName] = f
		}
	}

}

func GetData() *DataCollection {
	gist := getGist()
	fileName := os.Getenv("FILENAME")
	for _, f := range gist.Files {
		dc := DataCollection{}
		if f.Name == fileName {
			err := json.Unmarshal([]byte(f.Content), &dc)
			if err != nil {
				panic(err)
			}
			if len(dc.Records) > 0 {
				return &dc
			}
		}
	}
	return nil
}

func StoreData(used int, remain int) {
	gist := getGist()
	addData(used, remain, gist)
	saveGist(gist)
}
