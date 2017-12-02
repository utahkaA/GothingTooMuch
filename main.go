package main

import (
  "fmt"
  "time"
  "strings"
  "regexp"
  "math/rand"
  "github.com/PuerkitoBio/goquery"
)

type Word interface {
  Format()
}

type KatakanaWord struct {
  body string
  parts []string
}

func (this *KatakanaWord) Format() string {
  return fmt.Sprintf("%s %s", this.body, this.parts)
}

type TwoCharIdiomaticWord struct {
  body string
  parts []string
}

func (this *TwoCharIdiomaticWord) Format() string {
  return fmt.Sprintf("%s", this.body)
}

func getKatakanaWordList() []*KatakanaWord {
  doc, _ := goquery.NewDocument("http://web.ydu.edu.tw/~uchiyama/data/kata_2010.html")

  var words []*KatakanaWord
  doc.Find("tbody").Each(func(_ int, s *goquery.Selection) {
    s.Find("tr").Each(func(i int, s *goquery.Selection) {
      body := s.Find(".kata").Text()

      part := s.Find(".part").Text()
      parts := []string{"Noun", }
      if part == "-な" {
        part = "Adj"
        parts = append(parts, part)
      } else if part == "-する" {
        part = "Verb"
        parts = append(parts, part)
      }

      words = append(words, &KatakanaWord{body, parts})
    })
  })

  return words
}

func getTwoCharIdiomaticWordList() []*TwoCharIdiomaticWord {
  var words []*TwoCharIdiomaticWord
  re := regexp.MustCompile(`>>[0-9]+`)

  doc, _ := goquery.NewDocument("http://nantuka.blog119.fc2.com/blog-entry-28981.html")

  doc.Find(".article").Each(func(_ int, s *goquery.Selection) {
    s.Find("div[style='margin: 10px 0px 60px 20px; font-weight:bold;']").Each(func(_ int, s *goquery.Selection) {
      msg := strings.TrimLeft(s.Text(), "\n ")
      msg = strings.TrimRight(msg, "\n")
      if re.Match([]byte(msg)) {
        msg = re.ReplaceAllString(msg, "")
      }
      rune_msg := []rune(msg)

      if len(rune_msg) != 2 {
        ws := strings.Fields(msg)

        // check whether each elements' length is two or not.
        cnt := 0
        for _, w := range ws {
          if len([]rune(w)) == 2 {
            cnt++
          }
        }

        // append target word to dictionary of Twocharidiomaticword
        if cnt == len(ws) {
          for _, w := range ws {
            words = append(words, &TwoCharIdiomaticWord{w, []string{"Noun",}})
          }
        }
      } else {
        words = append(words, &TwoCharIdiomaticWord{msg, []string{"Noun",}})
      }
    })
  })

  return words
}

func containsAdj(parts []string) bool {
  for _, part := range parts {
    if part == "Adj" {
      return true
    }
  }
  return false
}

func selectAdjKatakanaWord(katakanaWordList []*KatakanaWord) []*KatakanaWord {
  var selectedWords []*KatakanaWord
  for _, word := range katakanaWordList {
    if containsAdj(word.parts) {
      selectedWords = append(selectedWords, word)
    }
  }
  return selectedWords
}

func main() {
  katakanaWords := getKatakanaWordList()
  adjKatakanaWords := selectAdjKatakanaWord(katakanaWords)
  twoCharIdiomaticWords := getTwoCharIdiomaticWordList()

  rand.Seed(time.Now().UnixNano())
  i := rand.Intn(len(adjKatakanaWords))
  j := rand.Intn(len(katakanaWords))
  k := rand.Intn(len(twoCharIdiomaticWords))

  fmt.Println("ナッシング トゥー マッチ！")
  fmt.Println("オーマイ ゴッド ファーザー降臨！")
  fmt.Printf("%s %s %s！\n", adjKatakanaWords[i].body,
                             katakanaWords[j].body,
                             twoCharIdiomaticWords[k].body)
  fmt.Println("ヨイショ！")
}
