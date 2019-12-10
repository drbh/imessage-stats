package messagestats

import (
	"encoding/json"
	"fmt"
	"github.com/cdipaolo/sentiment"
	gomojicount "github.com/drbh/gomoji-counter"
	imessagehooks "github.com/drbh/imessage-stats/go-imessage-hooks"
	_ "github.com/mattn/go-sqlite3"
	"math"
	"strconv"
	"strings"
	"time"
)

type MessageStats struct {
	Messages               []Msg
	EmojiMap               map[string]int
	WkHr                   map[string]int
	SentimentScore         uint8
	MessageCount           int
	FirstSeen              time.Time
	AverageResponseSeconds float64
	ResponseTimes          []ResponseTime
}

type Msg struct {
	Year      int
	Month     int
	Day       int
	Wkday     string
	Hour      int
	Len       int
	Positve   uint8
	Timestamp string
}

func (ms *MessageStats) getCountSince() {
	//
}

type Key struct {
	Weekday   string
	HourOfDay int
}

type ResponseTime struct {
	IsentTime       time.Time
	TheyRespondTime time.Time
	Diff            time.Duration
}

func getStats(handle string) MessageStats {
	myStats := MessageStats{}
	allMsgs := imessagehooks.Fetch(handle, "0")
	model, err := sentiment.Restore()
	if err != nil {
		panic(fmt.Sprintf("Could not restore model!\n\t%v\n", err))
	}
	var days = [...]string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}

	msf := map[Key]int{}

	for _, w := range days {
		for i := 0; i <= 23; i++ {
			msf[Key{w, i}] = 0
		}
	}
	count := 0
	minTime := math.Inf(1)

	responseTimes := []ResponseTime{}
	var allText string
	var lastTm time.Time
	var lastFrom string

	for _, v := range allMsgs {
		count++
		s, err := strconv.ParseFloat(v.Date, 64)
		if err != nil {

		}

		tm := imessagehooks.AppleTimestampToTime(v.Date)

		if v.IsFromMe == "1" && lastFrom == "0" {
			lastFrom = "1"
			diff := tm.Sub(lastTm)
			responseTimes = append(responseTimes, ResponseTime{tm, lastTm, diff})
			lastTm = tm
			// fmt.Println(v.IsFromMe, tm, lastTm, diff)
		} else {
			lastFrom = "0"
			// fmt.Println(v.IsFromMe, tm, lastTm)
		}

		// time.Sleep(1 * time.Second)

		if s < minTime {
			minTime = s
		}

		m := Msg{
			Year:      tm.Year(),
			Month:     int(tm.Month()),
			Day:       tm.Day(),
			Wkday:     tm.Weekday().String(),
			Hour:      tm.Hour(),
			Len:       len(v.Text),
			Positve:   0, //analysis.Score,
			Timestamp: v.Date,
		}
		key := Key{tm.Weekday().String(), tm.Hour()}
		val := msf[key]
		msf[key] = val + 1
		myStats.Messages = append(myStats.Messages, m)
		allText = allText + " " + v.Text
	}

	analysis := model.SentimentAnalysis(allText, sentiment.English) // 0
	total := gomojicount.GetEmojiFrequencyCount(allText)

	mojis := map[string]int{}
	for _, v := range total {
		mojis[v.Emoji] = v.Count
	}

	// myStats.Messages = []Msg{}

	_, responseTimes = responseTimes[0], responseTimes[1:]
	myStats.ResponseTimes = responseTimes

	diffs := int64(0)
	for _, rtm := range responseTimes {
		diffs = diffs + int64(rtm.Diff)
	}

	myStats.EmojiMap = mojis
	myStats.SentimentScore = analysis.Score

	myStats.MessageCount = count
	myStats.FirstSeen = imessagehooks.AppleTimestampToTime(fmt.Sprint(int(minTime)))
	myStats.AverageResponseSeconds = (float64(diffs) / float64(len(responseTimes))) / 1000000000

	msfS := map[string]int{}

	for k, v := range msf {
		msfS[k.Weekday+"_"+fmt.Sprint(k.HourOfDay)] = v
	}

	myStats.WkHr = msfS

	return myStats
}

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	m := make(map[string]int)
	for _, word := range words {
		m[word] += 1
	}
	return m
}

func getCountStrings(period string) map[int]map[string]int {
	allMsgs := imessagehooks.FetchFullDatabase("0")

	count := 0
	minTime := math.Inf(1)

	var allText string
	// var messageTimeMap map[int]string
	messageTimeMapYear := make(map[int]string)
	messageTimeMapMonth := make(map[int]string)

	for _, v := range allMsgs {
		tm := imessagehooks.AppleTimestampToTime(v.Date)
		count++
		s, err := strconv.ParseFloat(v.Date, 64)
		if err != nil {
		}
		if s < minTime {
			minTime = s
		}
		messageTimeMapYear[tm.Year()] = messageTimeMapYear[tm.Year()] + " " + v.Text
		messageTimeMapMonth[int(tm.Month())] = messageTimeMapMonth[int(tm.Month())] + " " + v.Text
		allText = allText + " " + v.Text
	}

	switch os := period; os {
	case "everything":
		resultAll := make(map[int]map[string]int)
		resultAll[0] = WordCount(allText)
		return resultAll

	case "years":
		resultsYear := make(map[int]map[string]int)
		for k, v := range messageTimeMapYear {
			resultsYear[k] = WordCount(v)
		}
		return resultsYear

	case "months":
		resultsMonth := make(map[int]map[string]int)
		for k, v := range messageTimeMapMonth {
			resultsMonth[k] = WordCount(v)
		}
		return resultsMonth

	default:
		resultAll := make(map[int]map[string]int)
		resultAll[0] = WordCount(allText)
		return resultAll
	}

}

func getStatsFullDatabase() MessageStats {
	myStats := MessageStats{}
	allMsgs := imessagehooks.FetchFullDatabase("0")
	model, err := sentiment.Restore()
	if err != nil {
		panic(fmt.Sprintf("Could not restore model!\n\t%v\n", err))
	}
	var days = [...]string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}

	msf := map[Key]int{}

	for _, w := range days {
		for i := 0; i <= 23; i++ {
			msf[Key{w, i}] = 0
		}
	}
	count := 0
	minTime := math.Inf(1)

	responseTimes := []ResponseTime{}
	var allText string
	var lastTm time.Time
	var lastFrom string

	for _, v := range allMsgs {
		count++
		s, err := strconv.ParseFloat(v.Date, 64)
		if err != nil {

		}

		tm := imessagehooks.AppleTimestampToTime(v.Date)

		if v.IsFromMe == "1" && lastFrom == "0" {
			lastFrom = "1"
			diff := tm.Sub(lastTm)
			responseTimes = append(responseTimes, ResponseTime{tm, lastTm, diff})
			lastTm = tm
			// fmt.Println(v.IsFromMe, tm, lastTm, diff)
		} else {
			lastFrom = "0"
			// fmt.Println(v.IsFromMe, tm, lastTm)
		}

		// time.Sleep(1 * time.Second)

		if s < minTime {
			minTime = s
		}

		m := Msg{
			Year:      tm.Year(),
			Month:     int(tm.Month()),
			Day:       tm.Day(),
			Wkday:     tm.Weekday().String(),
			Hour:      tm.Hour(),
			Len:       len(v.Text),
			Positve:   0, //analysis.Score,
			Timestamp: v.Date,
		}
		key := Key{tm.Weekday().String(), tm.Hour()}
		val := msf[key]
		msf[key] = val + 1
		myStats.Messages = append(myStats.Messages, m)
		allText = allText + " " + v.Text
	}

	analysis := model.SentimentAnalysis(allText, sentiment.English) // 0
	total := gomojicount.GetEmojiFrequencyCount(allText)

	mojis := map[string]int{}
	for _, v := range total {
		mojis[v.Emoji] = v.Count
	}

	// myStats.Messages = []Msg{}

	_, responseTimes = responseTimes[0], responseTimes[1:]
	myStats.ResponseTimes = responseTimes

	diffs := int64(0)
	for _, rtm := range responseTimes {
		diffs = diffs + int64(rtm.Diff)
	}

	myStats.EmojiMap = mojis
	myStats.SentimentScore = analysis.Score

	myStats.MessageCount = count
	myStats.FirstSeen = imessagehooks.AppleTimestampToTime(fmt.Sprint(int(minTime)))
	myStats.AverageResponseSeconds = (float64(diffs) / float64(len(responseTimes))) / 1000000000

	msfS := map[string]int{}

	for k, v := range msf {
		msfS[k.Weekday+"_"+fmt.Sprint(k.HourOfDay)] = v
	}

	myStats.WkHr = msfS

	return myStats
}

func GetFullProfileStats(handle string) []byte {
	msf := getStats(handle)

	js, _ := json.Marshal(msf)

	if string(js) == "" {

	}
	return js
	// fmt.Println(len(js))
}

func GetStringCountsFullDatabase(period string) []byte {
	msf := getCountStrings(period)

	js, _ := json.Marshal(msf)

	if string(js) == "" {

	}
	return js
	// fmt.Println(len(js))
}

func GetFullProfileStatsFullDatabase() []byte {
	msf := getStatsFullDatabase()
	js, _ := json.Marshal(msf)
	if string(js) == "" {

	}
	return js
	// fmt.Println(len(js))
}
