package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDoQuery(t *testing.T) {
	query := "的"
	entry, err := doQuery(query)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(entry.Forms)
}

func TestParseArgs(t *testing.T) {
	want := []string{"木", "尸", "腂"}
	got := parseArgs([]string{"木尸腂"})

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestAtohan(t *testing.T) {
	want := "日月金木水火土竹戈十大中一弓人心手口尸廿山女田難卜符"
	got := atohan("abcdefghijklmnopqrstuvwxyz")
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
