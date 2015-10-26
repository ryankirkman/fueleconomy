package handlers

import (
	"bytes"
	"net/url"
	"strconv"
)

var (
	PageLengthDefault int = 10
	PageLengthMax     int = 100
	PageNoDefault     int = 1
)

func getPageFromQueryVals(queryVals url.Values, URL *url.URL) (p *PageInfo) {
	p = &PageInfo{
		BaseUrl:    URL,
		PageLength: PageLengthDefault,
		PageNo:     PageNoDefault,
	}
	pageNo := getIntFromQueryVals(queryVals, "page")
	pageLength := getIntFromQueryVals(queryVals, "pageLength")
	if pageNo > 0 {
		p.PageNo = pageNo
	}
	if pageLength > 0 {
		p.PageLength = maxInt(pageLength, PageLengthMax)
	}

	return p
}

type PageInfo struct {
	BaseUrl      *url.URL `json:"-"`
	NextPage     string   `json:"nextPage,omitEmpty"`
	PageLength   int      `json:"pageLength"`
	PageNo       int      `json:"page"`
	PrevPage     string   `json:"prevPage,omitEmpty"`
	TotalResults int      `json:"totalResults,omitEmpty"`
	TotalPages   int      `json:"totalPages,omitEmpty"`
}

func (p *PageInfo) Fill(queryVals url.Values, resultCount int) {
	p.TotalResults = resultCount
	p.TotalPages = (resultCount / p.PageLength) + 1
	if (resultCount % p.PageLength) == 0 {
		p.TotalPages--
	}
	p.generateUrls(queryVals)
}

func (p *PageInfo) generateUrls(queryVals url.Values) {
	var thisPageNo int = p.PageNo
	if thisPageNo > 1 {
		p.PrevPage = generateUrlForPage(p.BaseUrl, queryVals, thisPageNo-1)
	}
	if thisPageNo < p.TotalPages {
		p.NextPage = generateUrlForPage(p.BaseUrl, queryVals, thisPageNo+1)
	}
}

func generateUrlForPage(URL *url.URL, queryVals url.Values, pageNo int) string {
	urlBuff := bytes.Buffer{}
	urlBuff.WriteString("http://fueleconomy.io")
	urlBuff.WriteString(URL.Path)
	urlBuff.WriteString("?")
	queryVals.Set("page", strconv.Itoa(pageNo))
	urlBuff.WriteString(queryVals.Encode())
	return urlBuff.String()
}
