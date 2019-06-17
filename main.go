package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/exporter"
)

// 干员信息
type Agent struct {
	Name        string  `json:"name"`
	EnglishName string  `json:"english_name"`
	Slogan      string  `json:"slogan"`
	Class       string  `json:"class"`
	Star        string  `json:"star"`
	Group       string  `json:"group"`
	Profile     profile `json:"profile"`
	Tag         string  `json:"tag"`
	Character   string  `json:"character"`
	Record      string  `json:"record"`
	Image       image   `json:"image"`
	Link        string  `json:"link"`
}

// 干员档案信息
type profile struct {
	Position string `json:"position"`
	Mastery  string `json:"mastery"`
	XP       string `json:"xp"`
	From     string `json:"from"`
	Birthday string `json:"birthday"`
	Race     string `json:"race"`
	Height   string `json:"height"`
	Status   string `json:"status"`
}

type image struct {
	Image1 string `json:"image_1"`
	Image2 string `json:"image_2"`
}

// 文件名
var File = "arknight_agents.json"

func main() {
	file, err := os.Open(File)
	defer file.Close()
	if os.IsNotExist(err) {
		fetchAgents()
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var agent Agent
		err = json.Unmarshal(scanner.Bytes(), &agent)
		DownloadImage(agent.Name, "image_1", agent.Image.Image1)
		DownloadImage(agent.Name, "image_2", agent.Image.Image2)
		log.Println("Downloaded", agent.Name)
	}
}

func fetchAgents() {
	geziyor.NewGeziyor(geziyor.Options{
		StartURLs: []string{"http://wiki.joyme.com/arknights/%E5%9B%BE%E9%89%B4%E4%B8%80%E8%A7%88"},
		Exporters: []geziyor.Exporter{exporter.JSONExporter{
			FileName: File,
		}},
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36",
		ParseFunc: func(r *geziyor.Response) {
			r.DocHTML.Find("#Contentbox2 ").Children().Each(func(_ int, s *goquery.Selection) {
				s.Find("table").Each(func(_ int, s *goquery.Selection) {
					link := s.Find("tr").First().Find("a")
					title, _ := link.Attr("title")
					href, _ := link.Attr("href")
					agent := Agent{
						Name: title,
						Link: "http://wiki.joyme.com" + href,
					}
					r.Geziyor.Get(agent.Link, func(r *geziyor.Response) {
						getAgentDetail(&agent, r)
					})
					r.Exports <- agent
				})
			})
		},
	}).Start()
}

func getAgentDetail(agent *Agent, r *geziyor.Response) {
	// Slogan 和 英文名
	slogan := r.DocHTML.Find("#mw-content-text > div:nth-child(8) > table > tbody > tr > td:nth-child(2) > div:nth-child(1) > big > big > b")
	englishName := r.DocHTML.Find("#mw-content-text > div.mwiki_hide > div > div:nth-child(3) > div > div:nth-child(1) > table > tbody > tr:nth-child(1) > td")
	agent.Slogan = slogan.Text()
	agent.EnglishName = strings.TrimSpace(englishName.Text())

	// 基本信息
	basicTable := r.DocHTML.Find("#mw-content-text > div.mwiki_hide > div > div:nth-child(3) > div > div:nth-child(2) > table > tbody")
	class := basicTable.Find("tr:nth-child(1) > td:nth-child(2)")
	group := basicTable.Find("tr:nth-child(2) > td:nth-child(2)")
	star := basicTable.Find("tr:nth-child(1) > td:nth-child(4)")
	character := basicTable.Find("tr:nth-child(5) > td")
	tag := basicTable.Find("tr:nth-child(6) > td")

	agent.Class = strings.TrimSpace(class.Text())
	agent.Group = strings.TrimSpace(group.Text())
	agent.Star = strings.TrimSpace(star.Text())
	agent.Character = strings.TrimSpace(character.Text())
	agent.Tag = strings.TrimSpace(tag.Text())

	// 履历信息
	record := r.DocHTML.Find("#mw-content-text > div.mwiki_hide > div > div:nth-child(4) > div > div:nth-child(2) > table > tbody > tr:nth-child(2) > td")
	agent.Record = strings.TrimSpace(record.Text())

	// 档案信息
	profileTable := r.DocHTML.Find("#mw-content-text > div.tj-big > div.tj-bg-right > div.tj-bgs.wiki_hide > table.wikitable > tbody")
	position := profileTable.Find("tr:nth-child(1) > td")
	mastery := profileTable.Find("tr:nth-child(2) > td")
	xp := profileTable.Find("tr:nth-child(3) > td")
	from := profileTable.Find("tr:nth-child(4) > td")
	birthday := profileTable.Find("tr:nth-child(5) > td")
	race := profileTable.Find("tr:nth-child(6) > td")
	height := profileTable.Find("tr:nth-child(7) > td")
	status := profileTable.Find("tr:nth-child(8) > td")

	agent.Profile.Position = strings.TrimSpace(position.Text())
	agent.Profile.Mastery = strings.TrimSpace(mastery.Text())
	agent.Profile.XP = strings.TrimSpace(xp.Text())
	agent.Profile.From = strings.TrimSpace(from.Text())
	agent.Profile.Birthday = strings.TrimSpace(birthday.Text())
	agent.Profile.Race = strings.TrimSpace(race.Text())
	agent.Profile.Height = strings.TrimSpace(height.Text())
	agent.Profile.Status = strings.TrimSpace(status.Text())

	// 图片信息
	img1 := r.DocHTML.Find("#con_1 > div > div > a > img")
	if img, exists := img1.Attr("src"); exists {
		agent.Image.Image1 = strings.ReplaceAll(img, "/dr/1120__", "")
	}
	img2 := r.DocHTML.Find("#con_2 > div > div > a > img")
	if img, exists := img2.Attr("src"); exists {
		agent.Image.Image2 = strings.ReplaceAll(img, "/dr/1120__", "")
	}
}
