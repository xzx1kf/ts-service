package main

import (
  "encoding/json"
  "log"
  "net/http"
  "net/url"
  "strconv"
  "time"

  "github.com/PuerkitoBio/goquery"
)

func Scrape(w http.ResponseWriter, r *http.Request) {

  slots := Slots{}

  doc := getDocument("1")
  slots = parseAvailableSlots(*doc)

  enc := json.NewEncoder(w)
  enc.SetEscapeHTML(false)
  err := enc.Encode(slots)
  if err != nil {
    panic(err)
  }
}

func getDocument(days string) *goquery.Document {

  // TODO: for testing purposes set the day to 1. Day should be a parameter
  resp, err := http.Get("http://tynemouth-squash.herokuapp.com/?day=" + days)
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
