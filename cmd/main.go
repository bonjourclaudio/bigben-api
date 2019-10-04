package main

import (
	"bytes"
	"encoding/json"
	"github.com/claudioontheweb/bigben-api/models"
	"github.com/gorilla/mux"
	"github.com/ledongthuc/pdf"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", GetMenuHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func GetMenuHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	menu := fetchBigBenSite(w)

	err := json.NewEncoder(w).Encode(menu)
	if err != nil {
		panic(err)
	}
}

func fetchBigBenSite(w http.ResponseWriter) models.Menu {

	dt := time.Now()
	today := dt.Format("02.01.06")
	url := "https://bigbenwestside.ch/uploads/Mittagsmen%C3%BCs/" + today + ".pdf"

	menu, err := downloadFile(url, w, today)

	if err != nil {
		panic(err)
	}

	return menu
}

// Get Menu as PDF from URL and download
func downloadFile(url string, w http.ResponseWriter, fn string) (models.Menu, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	date := strings.Replace(fn, ".", "_", -1)
	filename := date + ".pdf"

	out, err := os.Create("./assets/" + filename)
	if err != nil {
		panic(err)
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	menu := printFile(filename)

	return menu, err
}

// Open PDF and return content as string
func printFile(filename string) models.Menu {
	f, r, err := pdf.Open("./assets/" + filename)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()

	if err != nil {
		panic(err)
	}

	_, err = buf.ReadFrom(b)
	if err != nil {
		panic(err)
	}

	ms := buf.String()

	menuString := strings.TrimSpace(ms)

	re := regexp.MustCompile("DAILY DISH|STEAK OF THE WEEK|BURGER OF THE WEEK")

	res := re.Split(menuString, -1)

	var menu models.Menu

	for k, _ := range res {
		menu.Date = res[0]

		res[k] = strings.TrimSpace(res[k])

		dailyDishString := []rune(res[1])
		steakString := []rune(res[2])
		burgerString := []rune(res[3])

		menu.DailyDish.Price = string(dailyDishString[len(dailyDishString)-5:])
		menu.DailyDish.Content = strings.Trim(string(dailyDishString), menu.DailyDish.Price)

		menu.SteakOfTheWeek.Price = string(steakString[len(steakString)-5:])
		menu.SteakOfTheWeek.Content = strings.Trim(string(steakString), menu.SteakOfTheWeek.Price)

		menu.BurgerOfTheWeek.Price = string(burgerString[len(burgerString)-5:])
		menu.BurgerOfTheWeek.Content = strings.Trim(string(burgerString), menu.BurgerOfTheWeek.Price)

	}

	return menu


}
