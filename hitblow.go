package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
)

/**
Hit & Blow

Colors:
- 0: None
- 1: blue
- 2: red
- 3: green
- 4: yellow
- 5: pink
- 6: white

Response:
- nhit: number of hits
- nblow: number of blowns

**/

type Color uint8

const (
	Blue   Color = 0
	Red          = 1
	Green        = 2
	Yellow       = 3
	Pink         = 4
	White        = 5
)

const NUM_COLORS = 6
const LINE_LENGTH = 4
const NUM_ALL_LINES = 6 * 6 * 6 * 6
const LINE_ID_MIN = 0
const LINE_ID_MAX = NUM_ALL_LINES - 1

func (color Color) String() string {
	switch color {
	case Blue:
		return "Blue"
	case Red:
		return "Red"
	case Green:
		return "Green"
	case Yellow:
		return "Yellow"
	case Pink:
		return "Pink"
	case White:
		return "White"
	}
	return "Unknown"
}

func (color Color) ShortString() string {
	switch color {
	case Blue:
		return "B"
	case Red:
		return "R"
	case Green:
		return "G"
	case Yellow:
		return "Y"
	case Pink:
		return "P"
	case White:
		return "W"
	}
	return "U"
}

type Line [4]Color
type LineId uint

func (line *Line) String() string {
	return fmt.Sprintf(
		"%s%s%s%s",
		line[0].ShortString(),
		line[1].ShortString(),
		line[2].ShortString(),
		line[3].ShortString())
}

func (line *Line) GetLineId() LineId {
	return LineId((int(line[0]) * 6 * 6 * 6) + (int(line[1]) * 6 * 6) + (int(line[2]) * 6) + int(line[3]))
}

func CreateLineForId(lineId uint) Line {
	return Line{
		Color(lineId / (6 * 6 * 6)),
		Color((lineId / (6 * 6)) % 6),
		Color((lineId / 6) % 6),
		Color(lineId % 6),
	}
}

func GetLineForId(lineId LineId) *Line {
	return &AllLines[lineId]
}

var AllLines []Line = GenerateAllLines()

func GenerateAllLines() []Line {
	lines := make([]Line, NUM_ALL_LINES)
	for i := 0; i < NUM_ALL_LINES; i++ {
		lines[i] = CreateLineForId(uint(i))
	}
	return lines
}

type HitBlow struct {
	nhit  int
	nblow int
}

func (hb *HitBlow) String() string {
	return fmt.Sprintf("{%d hits, %d blows}", hb.nhit, hb.nblow)
}

func calcHB(answer *Line, sub *Line) *HitBlow {
	var nh_ans [6]int
	var nh_sub [6]int
	var nhit int
	var nblow int
	for i := 0; i < 4; i++ {
		if answer[i] == sub[i] {
			nhit++
		} else {
			nh_ans[answer[i]]++
			nh_sub[sub[i]]++
		}
	}
	for i := 0; i < 6; i++ {
		nblow += min(nh_ans[i], nh_sub[i])
	}
	return &HitBlow{nhit, nblow}
}

type Restriction struct {
	line *Line
	hb   *HitBlow
}

func IsAccepted(attemptLine *Line, restLine *Line, restHb *HitBlow) bool {
	var hb = calcHB(restLine, attemptLine)
	return restHb.nhit == hb.nhit && restHb.nblow == hb.nblow
}

func IsAcceptedByAllRestrictions(line *Line, restrictions []Restriction) bool {
	for _, rest := range restrictions {
		if !IsAccepted(line, rest.line, rest.hb) {
			return false
		}
	}
	return true
}

type Game struct {
	restrictions    []Restriction
	possibleLineIds []LineId
}

func NewGame(restrictions []Restriction) *Game {
	var possibleLineIds = make([]LineId, 0, NUM_ALL_LINES)
	for i := 0; i < NUM_ALL_LINES; i++ {
		var lineId = LineId(i)
		var line = GetLineForId(lineId)
		if IsAcceptedByAllRestrictions(line, restrictions) {
			possibleLineIds = append(possibleLineIds, lineId)
		}
	}
	return &Game{
		restrictions:    restrictions,
		possibleLineIds: possibleLineIds,
	}
}

func (game *Game) AddRestriction(restrictions ...Restriction) {
	game.restrictions = append(game.restrictions, restrictions...)
	var newPossibleLineIds = make([]LineId, 0, len(game.possibleLineIds))
	for _, lineId := range game.possibleLineIds {
		if IsAcceptedByAllRestrictions(GetLineForId(lineId), restrictions) {
			newPossibleLineIds = append(newPossibleLineIds, lineId)
		}
	}
	game.possibleLineIds = newPossibleLineIds
}

func (game *Game) DescribeStatus() string {
	var text = ""

	text += fmt.Sprintf("Restrictions: %d\n", len(game.restrictions))
	for _, rest := range game.restrictions {
		text += fmt.Sprintf("  [%s] -> %s\n", rest.line, rest.hb)
	}
	text += "\n"

	text += fmt.Sprintf("Possible lines: %d\n", len(game.possibleLineIds))
	for _, lineId := range game.possibleLineIds[:min(10, len(game.possibleLineIds))] {
		text += fmt.Sprintf("  [%s]\n", GetLineForId(lineId))
	}
	text += "\n"

	return text

}

type Suggestion struct {
	line     *Line
	score    float64
	variance float64
}

var firstSuggestions []Suggestion = nil
var firstTopSuggestions []Suggestion = nil

func (game Game) SuggestNextAttempts() (topSuggestions []Suggestion, allSuggestions []Suggestion) {
	if len(game.possibleLineIds) == 0 {
		return []Suggestion{}, []Suggestion{}
	}
	if len(game.possibleLineIds) == NUM_ALL_LINES {
		if firstTopSuggestions != nil && firstSuggestions != nil {
			return firstTopSuggestions, firstSuggestions
		}
	}

	// var nextAttemptsCandidates = game.possibleLineIds
	var nextAttemptsCandidates = make([]LineId, NUM_ALL_LINES)
	for i := 0; i < NUM_ALL_LINES; i++ {
		nextAttemptsCandidates[i] = LineId(i)
	}

	// answerとして可能なlineすべてのうち、最も情報量の期待値が多くなるものを選ぶ.
	var suggestions = make([]Suggestion, len(nextAttemptsCandidates))

	var wg = &sync.WaitGroup{}
	wg.Add(len(nextAttemptsCandidates))

	for i, nextLineId := range nextAttemptsCandidates {
		go func(i int, nextLineId LineId) {
			var nextLine = GetLineForId(nextLineId)
			var mean, variance = game.estimateInformationValue(nextLine)
			suggestions[i] = Suggestion{nextLine, mean, variance}
			wg.Add(-1)
		}(i, nextLineId)
	}
	wg.Wait()

	sort.SliceStable(suggestions, func(i, j int) bool {
		return suggestions[i].score > suggestions[j].score || (suggestions[i].score == suggestions[j].score && suggestions[i].variance < suggestions[j].variance)
	})

	var bestMean = suggestions[0].score
	var bestVariance = suggestions[0].variance

	var filteredSuggestions = make([]Suggestion, 0, len(suggestions))
	for _, sug := range suggestions {
		if sug.score == bestMean && sug.variance == bestVariance {
			filteredSuggestions = append(filteredSuggestions, sug)
		} else {
			break
		}
	}

	if len(game.possibleLineIds) == NUM_ALL_LINES {
		if firstSuggestions == nil || firstTopSuggestions == nil {
			firstSuggestions = suggestions
			firstTopSuggestions = filteredSuggestions
		}
	}

	return filteredSuggestions, suggestions
}

func (game Game) countNextPossibleLines(nextAttemptLine *Line, hb *HitBlow) int {
	// lineを次の試行として結果hbが得られた場合の、次の試行で可能なlineの数を計算する.
	var count = 0
	for _, possibleLineId := range game.possibleLineIds {
		var possibleLine = GetLineForId(possibleLineId)
		if IsAccepted(possibleLine, nextAttemptLine, hb) {
			count++
		}
	}
	return count
}

func (game Game) estimateInformationValue(nextAttemptLine *Line) (mean float64, variance float64) {
	// lineを次の試行とした場合に除外できるlineの割合の期待値を計算する.
	var currentLineCount = len(game.possibleLineIds)
	var filteredCounts = make([]uint64, len(game.possibleLineIds))
	var wg = &sync.WaitGroup{}
	wg.Add(len(game.possibleLineIds))

	for i, possibleLineId := range game.possibleLineIds {
		go func(i int, possibleLineId LineId) {
			var possibleLine = GetLineForId(possibleLineId)
			var hb = calcHB(nextAttemptLine, possibleLine)
			var nextLineCount = game.countNextPossibleLines(nextAttemptLine, hb)
			filteredCounts[i] = uint64(currentLineCount - nextLineCount)
			wg.Add(-1)
		}(i, possibleLineId)
	}

	wg.Wait()

	var totalFilteredCount uint64 = 0
	var totalFilteredCountSquared uint64 = 0
	for _, count := range filteredCounts {
		totalFilteredCount += count
		totalFilteredCountSquared += count * count
	}

	mean = float64(totalFilteredCount) / float64(currentLineCount) / float64(currentLineCount)
	variance = float64(totalFilteredCountSquared)/float64(currentLineCount)/float64(currentLineCount) - mean*mean

	return mean, variance
}

func main1() {
	var game = NewGame([]Restriction{
		// ここに制約を追加
		//{&Line{Blue, Blue, Red, Red}, &HitBlow{0, 1}},
		//{&Line{Red, Pink, Green, Green}, &HitBlow{0, 1}},
		//{&Line{Green, Yellow, Pink, White}, &HitBlow{3, 0}},
	})

	fmt.Println(game.DescribeStatus())

	var topSuggestions, allSuggestions = game.SuggestNextAttempts()

	if len(topSuggestions) > 0 {
		// ファイル suggestions.txt に出力
		var file, _ = os.Create("suggestions.txt")
		defer file.Close()
		for i, sug := range allSuggestions {
			var istop = ""
			if i < len(topSuggestions) {
				istop = "(top)"
			}
			file.WriteString(fmt.Sprintf(
				"[%s]\t%.3f\t%.2f\t%s\n",
				sug.line,
				sug.score,
				sug.variance,
				istop,
			))
		}
		fmt.Printf("Suggestions: %d (see suggestions.txt for detail)\n", len(topSuggestions))
		for _, sug := range topSuggestions[:min(10, len(topSuggestions))] {
			fmt.Printf("  [%s] %.3f (%.2f)\n", sug.line, sug.score, sug.variance)
		}
		fmt.Println()

		// ランダムに選択
		var sugChoosen = topSuggestions[rand.Intn(len(topSuggestions))]
		fmt.Printf("Next: [%s]\n", sugChoosen.line)
		fmt.Println()
	}
}

func main() {

	for i := 0; i < 10000; i++ {
		var answer = AllLines[rand.Intn(len(AllLines))]
		var game = NewGame([]Restriction{})
		fmt.Printf("Trial #%d: Answer=[%s]\n", i+1, &answer)

		for attempt_i := 1; attempt_i <= 10; attempt_i++ {
			var numPossibleLinesBefore = len(game.possibleLineIds)
			var topSuggestions, _ = game.SuggestNextAttempts()
			if len(topSuggestions) == 0 {
				break
			}
			var sugChoosen = topSuggestions[rand.Intn(len(topSuggestions))]
			var hb = calcHB(&answer, sugChoosen.line)
			game.AddRestriction(Restriction{sugChoosen.line, hb})

			fmt.Printf("  Attempt #%d: [%s] => %s (%d => %d)\n", attempt_i, sugChoosen.line, hb, numPossibleLinesBefore, len(game.possibleLineIds))

			if len(game.possibleLineIds) == 1 {
				fmt.Printf("  Answer Found! [%s] (attempts: %d)\n", &AllLines[game.possibleLineIds[0]], attempt_i)
				break
			}
		}
	}

}
