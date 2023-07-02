package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/anaskhan96/soup"
)

type Sponsors struct {
	Current int `json:"current"`
	Past    int `json:"past"`
}

type Response struct {
	Sponsors Sponsors `json:"sponsors"`
}

func getSponsorCount(u string) string {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return `{"sponsors": {"count":"Error: Invalid URL."}}`
	}

	queryParams := parsedUrl.Query()
	usr := queryParams.Get("u")

	if usr == "" {
		return `{"sponsors": {"count":"Error: No user specified."}}`
	}

	url := fmt.Sprintf("https://github.com/sponsors/%s", usr)
	htmlResponse, err := soup.Get(url)

	if err != nil {
		return `{"sponsors": {"count":"Error: Unable to fetch the page."}}`
	}

	doc := soup.HTMLParse(htmlResponse)

	section := doc.Find("div", "id", "sponsors-section-list")
	if section.Error != nil {
		return `{"sponsors": {"count":"Error: GitHub Sponsors aren't setup with this user.", "error": "` + section.Error.Error() + `"}}`
	}

	soup.SetDebug(true)

	sponsorSection := section.Find("div")
	currentCountElement := sponsorSection.Find("span")
	pastCountElement := sponsorSection.FindNextElementSibling().Find("span")

	if pastCountElement.Error != nil {
		return `{"sponsors": {"count":"Error: GitHub Sponsors aren't setup with this user.", "error": "` + pastCountElement.Error.Error() + `"}}`
	}

	if currentCountElement.Error != nil {
		return `{"sponsors": {"count":"Error: GitHub Sponsors aren't setup with this user.", "error": "` + currentCountElement.Error.Error() + `"}}`
	}

	currentCount, err := strconv.Atoi(currentCountElement.Text())
	if err != nil {
		return `{"sponsors": {"message": ` + err.Error() + `"}}`
	}

	pastCount, err := strconv.Atoi(pastCountElement.Text())
	if err != nil {
		return `{"sponsors": {"message": ` + err.Error() + `"}}`
	}

	resp := Response{
		Sponsors: Sponsors{
			Current: currentCount,
			Past:    pastCount,
		},
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		return `{"sponsors": {"message":"failed to marshal response"}}`
	}

	return string(jsonData)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sponsors := getSponsorCount(r.RequestURI)
	fmt.Fprint(w, sponsors)
}
