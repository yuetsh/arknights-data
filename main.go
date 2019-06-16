package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/exporter"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

// 干员信息
type Agent struct {
	//基本信息
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
	Slogan      string `json:"slogan"`
	Class       string `json:"class"`
	Star        string `json:"star"`
	Group       string `json:"group"`
	Tag         string `json:"tag"`
	Character   string `json:"character"`
	Record      string `json:"record"`
	Avatar      string `json:"avatar"`
	Image1      string `json:"image_1"`
	Image2      string `json:"image_2"`
	Link        string `json:"link"`
}

// json
type Agents struct {
	Agents []Agent `json:"agents"`
}

// 全部干员
var AllAgents []Agent
var File = "arknight_agents.json"

func main() {
	fetchAgents()
	readAgents()
	downloadAgentsImage()
}

func fetchAgents() {
	_ = os.Remove(File)
	geziyor.NewGeziyor(geziyor.Options{
		StartURLs: []string{"http://wiki.joyme.com/arknights/%E5%9B%BE%E9%89%B4%E4%B8%80%E8%A7%88"},
		Exporters: []geziyor.Exporter{exporter.JSONExporter{
			FileName: File,
		}},
		UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B137 Safari/601.1 AlipayClient MicroMessenger",
		ParseFunc: func(r *geziyor.Response) {
			r.DocHTML.Find("#Contentbox2 ").Children().Each(func(_ int, s *goquery.Selection) {
				s.Find("table").Each(func(_ int, s *goquery.Selection) {
					link := s.Find("tr").First().Find("a")
					title, _ := link.Attr("title")
					href, _ := link.Attr("href")
					img, _ := link.Find("img").Attr("src")
					img = strings.Replace(img, "/dr/150__", "", -1)
					agent := Agent{Name: title, Link: "http://wiki.joyme.com" + href, Avatar: img}
					r.Geziyor.Get(agent.Link, func(r *geziyor.Response) {
						// 基本信息
						slogan := r.DocHTML.Find("#mw-content-text > div:nth-child(8) > table > tbody > tr > td:nth-child(2) > div:nth-child(1) > big > big > b")
						englishName := r.DocHTML.Find("#mw-content-text > div.mwiki_hide > div > div:nth-child(3) > div > div:nth-child(1) > table > tbody > tr:nth-child(1) > td")
						table := r.DocHTML.Find("#mw-content-text > div.mwiki_hide > div > div:nth-child(3) > div > div:nth-child(2) > table > tbody")
						class := table.Find("tr:nth-child(1) > td:nth-child(2)")
						group := table.Find("tr:nth-child(2) > td:nth-child(2)")
						star := table.Find("tr:nth-child(1) > td:nth-child(4)")
						character := table.Find("tr:nth-child(5) > td")
						tag := table.Find("tr:nth-child(6) > td")

						agent.Slogan = slogan.Text()
						agent.EnglishName = strings.TrimSpace(englishName.Text())
						agent.Class = strings.TrimSpace(class.Text())
						agent.Group = strings.TrimSpace(group.Text())
						agent.Star = strings.TrimSpace(star.Text())
						agent.Character = strings.TrimSpace(character.Text())
						agent.Tag = strings.TrimSpace(tag.Text())

						// 履历信息
						record := r.DocHTML.Find("#mw-content-text > div.mwiki_hide > div > div:nth-child(4) > div > div:nth-child(2) > table > tbody > tr:nth-child(2) > td")
						agent.Record = strings.TrimSpace(record.Text())

						// 图片信息
						img1 := r.DocHTML.Find("#con_1 > div > div > a > img")
						if img, exists := img1.Attr("src"); exists {
							agent.Image1 = strings.Replace(img, "/dr/1120__", "", -1)
						}
						img2 := r.DocHTML.Find("#con_2 > div > div > a > img")
						if img, exists := img2.Attr("src"); exists {
							agent.Image2 = strings.Replace(img, "/dr/1120__", "", -1)
						}
					})
					AllAgents = append(AllAgents, agent)
				})
			})
			r.Exports <- Agents{AllAgents}
		},
	}).Start()
}

func readAgents() {
	bytes, err := ioutil.ReadFile(File)
	if err != nil {
		panic(err)
	}
	var data Agents
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		panic(err)
	}
	AllAgents = data.Agents
}

func downloadAgentsImage() {
	var wg sync.WaitGroup
	for _, agent := range AllAgents {
		wg.Add(1)
		go func(agent Agent) {
			DownloadImage(agent.Name, "avatar", agent.Avatar)
			DownloadImage(agent.Name, "image_1", agent.Image1)
			DownloadImage(agent.Name, "image_2", agent.Image2)
			wg.Done()
			log.Println("Download", agent.Name)
		}(agent)
	}
	wg.Wait()
}
