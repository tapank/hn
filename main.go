package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	NEW  = "https://hacker-news.firebaseio.com/v0/newstories.json"
	TOP  = "https://hacker-news.firebaseio.com/v0/topstories.json"
	BEST = "https://hacker-news.firebaseio.com/v0/beststories.json"
	ITEM = "https://hacker-news.firebaseio.com/v0/item/%d.json"
	PAGE = "https://news.ycombinator.com/item?id=%d"
)

type Item struct {
	Id    int    `json:"id"`    // The item's unique id.
	Type  string `json:"type"`  // The type of item. One of "job", "story", "comment", "poll", or "pollopt".
	By    string `json:"by"`    // The username of the item's author.
	Time  int    `json:"time"`  // Creation date of the item, in Unix Time.
	Url   string `json:"url"`   // The URL of the story.
	Score int    `json:"score"` // The story's score, or the votes for a pollopt.
	Title string `json:"title"` // The title of the story, poll or job. HTML.
}

var items map[int]Item
var context = TOP
var startIndex int
var countPerPage = 30

func main() {
	refresh(true)
	for {
		switch choice := getInput(); choice {
		case "m", "more":
			startIndex += countPerPage
			refresh(false)
		case "n", "new":
			context = NEW
			refresh(true)
		case "t", "top":
			context = TOP
			refresh(true)
		case "b", "best":
			context = BEST
			refresh(true)
		case "q", "quit":
			return
		case "r", "refresh":
			refresh(false)
		default:
			openItemInBrowser(choice)
		}
	}

}

func refresh(resetStartIndex bool) {
	if resetStartIndex {
		startIndex = 0
	}
	loadItems(context)
	listStories()
}

func loadItems(endpoint string) {
	res, err := http.Get(endpoint)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	itemIds := []int{}
	if err := json.Unmarshal(body, &itemIds); err != nil {
		log.Fatal(err)
	}
	items = map[int]Item{}
	for i, item := range itemIds[startIndex : startIndex+countPerPage] {
		item := getItem(item)
		items[i+1+startIndex] = item
	}
}

func listStories() {
	contextName := "unknown context"
	switch context {
	case TOP:
		contextName = "top"
	case NEW:
		contextName = "new"
	case BEST:
		contextName = "best"
	}
	fmt.Printf("item %d to %d of %s\n", startIndex+1, startIndex+countPerPage, contextName)
	for i := startIndex + 1; i <= startIndex+countPerPage; i++ {
		fmt.Printf("%02d. [%5d %15s] %s\n", i, items[i].Score, items[i].By, items[i].Title)
	}
}

func getItem(id int) Item {
	res, err := http.Get(fmt.Sprintf(ITEM, id))
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var item Item
	err = json.Unmarshal(body, &item)
	if err != nil {
		panic(err.Error())
	}
	return item
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter choice [<sno>|(m)ore|(t)op|(b)est|(n)ew|(q)uit|(r)efresh]: ")
	text, _ := reader.ReadString('\n')
	return strings.Trim(text, "\n \t")
}

func openItemInBrowser(ch string) {
	if choice, err := strconv.Atoi(ch); err == nil {
		if item, ok := items[choice]; ok {
			fmt.Printf("opening item %d: %s\n", choice, item.Title)
			cmd := exec.Command("open", "-a", "/Applications/Firefox.app", fmt.Sprintf(PAGE, item.Id))
			if err := cmd.Run(); err != nil {
				panic(err.Error())
			}
			return
		}
	}
	fmt.Println("unknown option:", ch)
}
