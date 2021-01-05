package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// endpoint urls
const (
	NEW  = "https://hacker-news.firebaseio.com/v0/newstories.json"
	TOP  = "https://hacker-news.firebaseio.com/v0/topstories.json"
	BEST = "https://hacker-news.firebaseio.com/v0/beststories.json"
	ITEM = "https://hacker-news.firebaseio.com/v0/item/%d.json"
	PAGE = "https://news.ycombinator.com/item?id=%d"
)

// Item is the news item
type Item struct {
	sno   int    // The serial number of the item as listed on hn
	ID    int    `json:"id"`    // The item's unique id.
	Type  string `json:"type"`  // The type of item. One of "job", "story", "comment", "poll", or "pollopt".
	By    string `json:"by"`    // The username of the item's author.
	Time  int    `json:"time"`  // Creation date of the item, in Unix Time.
	URL   string `json:"url"`   // The URL of the story.
	Score int    `json:"score"` // The story's score, or the votes for a pollopt.
	Title string `json:"title"` // The title of the story, poll or job. HTML.
}

var items map[int]Item
var context = TOP
var startIndex int
var countPerPage = 30
var layout = "01/02 15:04" // reference "2006-01-02 15:04:05.999999999 -0700 MST"

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

type listItem struct {
	sno    int
	itemID int
}

// refresh loads and lists stories from the context and provides an option to reset startIndex
func refresh(resetStartIndex bool) {
	if resetStartIndex {
		startIndex = 0
	}
	loadItems(context)
	listStories()
}

// loadItems uses `context` to load items from the webservice
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
	idChan := make(chan listItem, countPerPage)
	itemChan := make(chan Item, countPerPage)

	for i, item := range itemIds[startIndex : startIndex+countPerPage] {
		idChan <- listItem{i + startIndex + 1, item}
		go getItem(idChan, itemChan)
	}

	for i := 0; i < countPerPage; i++ {
		item := <-itemChan
		items[item.sno] = item
	}
}

// listStories formats the prints the contents of items map and should be called after the map is populated by loadItems
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
	fmt.Printf("item %d to %d of %s at %s\n", startIndex+1, startIndex+countPerPage, contextName, time.Now().Format(layout))
	for i := startIndex + 1; i <= startIndex+countPerPage; i++ {
		ts := time.Unix(int64(items[i].Time), 0).Format(layout)
		sc := items[i].Score
		by := items[i].By
		title := items[i].Title + domain(items[i].URL)
		fmt.Printf("[%s %4d %15s] %02d. %s\n", ts, sc, by, i, title)
	}
}

// getItem gets the details of an item based on item id and is expected to be run in a go routine
func getItem(listItemChan chan listItem, itemChan chan Item) {
	li := <-listItemChan
	res, err := http.Get(fmt.Sprintf(ITEM, li.itemID))
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
	item.sno = li.sno
	itemChan <- item
}

// getInput presents available options to user and returns the option provided by user
func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter your choice [<sno> | (m)ore | (t)op | (b)est | (n)ew | (q)uit | (r)efresh]: ")
	text, _ := reader.ReadString('\n')
	return strings.Trim(text, "\n \t")
}

// openItemInBrowser opens the hn page for a given item, and it assumes the application is running on MacOS and also
// Firefox is installed. Change the command to suite your environment.
func openItemInBrowser(ch string) {
	if choice, err := strconv.Atoi(ch); err == nil {
		if item, ok := items[choice]; ok {
			fmt.Printf("opening item %d: %s\n", choice, item.Title)
			cmd := exec.Command("open", "-a", "/Applications/Firefox.app", fmt.Sprintf(PAGE, item.ID))
			if err := cmd.Run(); err != nil {
				panic(err.Error())
			}
			return
		}
	}
	fmt.Println("unknown option:", ch)
}

func domain(uri string) (ul string) {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	ul = u.Hostname()
	if len(ul) > 0 {
		ul = fmt.Sprintf(" (%s)", ul)
	}
	return
}
