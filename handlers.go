package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "net/url"
  "strconv"
  "strings"
  "time"

  "github.com/PuerkitoBio/goquery"
)

func Scrape(w http.ResponseWriter, r *http.Request) {

  slots := Slots{}

  doc := getDocument("http://tynemouth-squash.herokuapp.com/?day=1")
  //slots = parseAvailableSlots(*doc)
  slots = parseBookedSlots(*doc)

  enc := json.NewEncoder(w)
  enc.SetEscapeHTML(false)
  err := enc.Encode(slots)
  if err != nil {
    panic(err)
  }
}

func getDocument(url string) *goquery.Document {

  resp, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }

  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
  }

  doc, err := goquery.NewDocumentFromReader(resp.Body)
  if err != nil {
    log.Fatal(err)
  }

  return doc
}

func parseAvailableSlots(doc goquery.Document) []Slot {

  slots := Slots{}

  doc.Find(".booking div.book a.booking_link").Each(func(i int,s *goquery.Selection) {
    link, found := s.Attr("href")
    if found {
      slot := parseLink(link)
      slots = append(slots, slot)
    }
  })

  return slots
}

func parseBookedSlots(doc goquery.Document) []Slot {

  slots := Slots{}

  href := ""

  doc.Find(".booking div.booked a").Each(func(i int, s *goquery.Selection) {
    link, found := s.Attr("href")
    if found {
      if href != link {
        href = link
        detailsDoc := getDocument("http://tynemouth-squash.herokuapp.com" + link)
        parseSlotDetails(*detailsDoc, link)
      }
    }
  })

  return slots
}

func parseSlotDetails(doc goquery.Document, link string) {
  s := doc.Find("body h1")
  c := s.Text()[6:7]
  court, _ := strconv.Atoi(c)

  s = doc.Find("body h2")
  w := strings.Fields(s.Text())

  const shortForm = "3:04pm on Monday 02 January 2006"

  // Mess about to get the time and date of the court.
  year := time.Now().Year()
  y := strconv.Itoa(year)
  tim := w[0] + " " +w[1] + " " +w[2] + " " +w[3][0:2] + " " +w[4] + " " + y
  t, e := time.Parse(shortForm, tim)
  if e != nil {
    log.Fatal(e)
  }

  slot := Slot {
    Court: court,
    Time: t.Format("2006-01-02 15:04:05.000"),
    Booked: true,
    Link: "http://tynemouth-squash.herokuapp.com" + link,
  }
  fmt.Println(slot)

}

func parseLink(link string) Slot {
  u, err := url.Parse(link)
  if err != nil {
    log.Fatal(err)
  }
  q := u.Query()

  court, _    := strconv.Atoi(q.Get("court"))
  timeslot, _ := strconv.Atoi(q.Get("timeSlot"))
  days, _     := strconv.Atoi(q.Get("days"))
  hour, _     := strconv.Atoi(q.Get("hour"))
  min, _      := strconv.Atoi(q.Get("min"))

  // Generate a time for the court based on the days parameter
  ti := time.Now()
  t  := time.Date(ti.Year(), ti.Month(), ti.Day(), hour, min, 0, 0, time.UTC)
  t   = t.AddDate(0, 0, days)

  slot := Slot {
    Court: court,
    Time: t.Format("2006-01-02 15:04:05.000"),
    Slot: timeslot,
    Booked: false,
    Link: "http://tynemouth-squash.herokuapp.com" + link,
  }

  return slot
}
