package main

type Slot struct {
  Time    string `json:"time"`
  Court   int `json:"court"`
  Booked  bool `json:"booked"`
  Link    string `json:"link"`
  Slot    int `json:"slot"`
}

type Slots []Slot
