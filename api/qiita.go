package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
	"github.com/olekukonko/tablewriter"
	"strconv"
)

// jsonのパース用に構造体を定義
type Data struct {
	Url            string `json:"url"`
	Title          string `json:"title"`
	LikesCount     int    `json:"likes_count"`
	ReactionsCount int    `json:"reactions_count"`
	PageViewsCount int    `json:"page_views_count"`
}

// Qiitaからデータを取得
func fetchQiitaData(accessToken string, targetDate time.Time) []Data {

	baseUrl := "https://qiita.com/api/v2/"
	action := "items"

	// 1リクエストあたり最大30件取得
	baseParam := "?page=1&per_page=30"

	// monthだけintではなくMonth型のため型変換が必要
	year, month, day := targetDate.Date()
	targetDay := dateNum2String(year, int(month), day)

	year, month, day = targetDate.AddDate(0, 0, 1).Date()
	nextDay := dateNum2String(year, int(month), day)

	// 投稿の検索クエリを作成
	// 検索クエリ stocks:>NUM created:<YYYY-MM-DD created:>YYYY-MM-DD
	// 指定日に投稿されたストック数30以上の記事を取得
	varParam := "&query=stocks:>30+created:>=" + targetDay + "+created:<" + nextDay

	endpointURL, err := url.Parse(baseUrl + action + baseParam + varParam)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(Data{})
	if err != nil {
		panic(err)
	}

	var resp = &http.Response{}
	// qiitaのアクセストークンがない場合はAuthorizationを付与しない
	if len(accessToken) > 0 {
		resp, err = http.DefaultClient.Do(&http.Request{
			URL:    endpointURL,
			Method: "GET",
			Header: http.Header{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer " + accessToken},
			},
		})
	} else {
		resp, err = http.DefaultClient.Do(&http.Request{
			URL:    endpointURL,
			Method: "GET",
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
		})
	}
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var datas []Data

	if err := json.Unmarshal(b, &datas); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil
	}
	return datas
}

// データの出力
func outputQiitaData(datas []Data) {

	for _, val := range datas {
		fmt.Println(val.LikesCount, val.Title, val.Url)
	}

	fmt.Println()
}

// 年月日の数値を文字列に変換
func dateNum2String(year int, month int, day int) string {
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func RunQiitaSearch() {

	// アクセストークン取得
	qiitaToken := os.Getenv("QIITA_TOKEN")

	var baseDate time.Time

	//fmt.Println("いいね数		タイトル		URL")

	// コマンドライン引数がある場合 YYYY-MM-DD の形で指定 不適切な方の場合エラー
	//var err error
	//if len(os.Args) >= 2 {
	//	arg := os.Args[1]
	//	layout := "2006-01-02"
	//	baseDate, err = time.Parse(layout, arg)
	//	if err != nil {
	//		panic(err)
	//	}
	//} else {
		// 引数がない場合
	baseDate = time.Now()
	//}
	// 一週間前から取得
	startDate := baseDate.AddDate(0, 0, -6)
	numOfDates := 7

	// 対象の日付から一週間分のデータを取得
	for i := 0; i < numOfDates; i++ {
		targetDate := startDate.AddDate(0, 0, i)
		fmt.Printf("%d-%d-%d\n", targetDate.Year(), int(targetDate.Month()), targetDate.Day())

		var output [][]string
		// データ取得
		datas := fetchQiitaData(qiitaToken, targetDate)

		// 出力用データ整形
		for _, val := range datas {
			title := ""
			splitlen := 40
			runes := []rune(val.Title)

			for i := 0; i< len(runes); i += splitlen {
				if i+splitlen < len(runes) {
					title += string(runes[i:(i+splitlen)]) + "\n"
				} else {
					title += string(runes[i:])
				}
			}

			output = append(output, []string{strconv.Itoa(val.LikesCount), title, val.Url})
		}

		// テーブル描画
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"いいね数", "タイトル", "URL"})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(output)
		table.Render()
	}
}
