package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type WordStat struct {
	Word  string
	Count int
}

func TopN(text string, n int) []WordStat {
	var topsz int = n
	all := make(map[string]int)
	var pairs, top []WordStat
	for _, word := range strings.Fields(text) {
		all[word]++
	}
	for word, cnt := range all {
		pairs = append(pairs, WordStat{Word: word, Count: cnt})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})
	if len(pairs) < topsz {
		topsz = len(pairs)
	}
	top = make([]WordStat, topsz)
	copy(top, pairs[:topsz])
	return top
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %s", err)
		os.Exit(1)
	}
	text := string(data)
	top := TopN(text, 10)
	for i, stat := range top {
		fmt.Printf("% 2d: %v\t%d\n", i, stat.Word, stat.Count)
	}
}
