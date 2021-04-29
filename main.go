package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

type SafeCounter struct {
	mu        sync.Mutex
	sharedMap map[string]int
}

func main() {

	stringsArray := readInputFile("test.txt")
	sC := SafeCounter{sharedMap: make(map[string]int)}
	mapper(&sC, stringsArray)
	opt := reducer(&sC)
	writeOutput("WordCountOutput.txt", opt)
	// printMap(sC.sharedMap)

}

func mapper(sC *SafeCounter, stringsArray []string) {
	arrayLength := len(stringsArray)
	partLength := int(arrayLength / 5)

	for i := 0; i < arrayLength; i++ {

		stringsArray[i] = strings.ToLower(stringsArray[i])
		sC.sharedMap[stringsArray[i]] = 0
	}
	doneBuffer := make(chan int, 5)

	if partLength < 1 {
		go wordCounter(sC, stringsArray[:], doneBuffer)
		doneBuffer <- 1
		doneBuffer <- 1
		doneBuffer <- 1
		doneBuffer <- 1
	} else {
		go wordCounter(sC, stringsArray[0:partLength], doneBuffer)
		go wordCounter(sC, stringsArray[partLength:2*partLength], doneBuffer)
		go wordCounter(sC, stringsArray[2*partLength:3*partLength], doneBuffer)
		go wordCounter(sC, stringsArray[3*partLength:4*partLength], doneBuffer)
		go wordCounter(sC, stringsArray[4*partLength:], doneBuffer)
	}
	for i := 0; i < 5; i++ {
		<-doneBuffer
	}
}

func wordCounter(sC *SafeCounter, stringsArray []string, doneBuffer chan int) {
	for i := 0; i < len(stringsArray); i++ {
		// stringsArray[i] = strings.ToLower(stringsArray[i])
		sC.mu.Lock()
		sC.sharedMap[stringsArray[i]] += 1
		sC.mu.Unlock()
	}
	doneBuffer <- 1
}

func readInputFile(fileName string) []string {
	file, _ := os.Open(fileName)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	words := []string{}

	for scanner.Scan() {
		newLine := scanner.Text()

		words = append(words, strings.Split(newLine, " ")...)
	}

	return words
}

type WordFrequency struct {
	word      string
	frequency int
}

type AlphabetSorter struct {
	wordsList []WordFrequency
}

func (sorter AlphabetSorter) Len() int {
	return len(sorter.wordsList)
}

func (sorter AlphabetSorter) Less(i, j int) bool {
	return sorter.wordsList[i].word < sorter.wordsList[j].word
}

func (sorter AlphabetSorter) Swap(i, j int) {
	temp := sorter.wordsList[i]
	sorter.wordsList[i] = sorter.wordsList[j]
	sorter.wordsList[j] = temp
}

type FrequencySorter struct {
	wordsList []WordFrequency
}

func (sorter FrequencySorter) Len() int {
	return len(sorter.wordsList)
}

func (sorter FrequencySorter) Less(i, j int) bool {
	return sorter.wordsList[i].frequency > sorter.wordsList[j].frequency
}

func (sorter FrequencySorter) Swap(i, j int) {
	temp := sorter.wordsList[i]
	sorter.wordsList[i] = sorter.wordsList[j]
	sorter.wordsList[j] = temp
}

func turnMapToList(wordsMap map[string]int) []WordFrequency {
	wordsList := []WordFrequency{}

	for word, frequency := range wordsMap {
		wordsList = append(wordsList, WordFrequency{word, frequency})
	}

	return wordsList
}

func reducer(sC *SafeCounter) []WordFrequency {
	wordsMap := sC.sharedMap

	wordsList := turnMapToList(wordsMap)

	sorterAlpha := AlphabetSorter{wordsList}

	sort.Stable(sorterAlpha)

	sorterFreq := FrequencySorter{sorterAlpha.wordsList}

	sort.Stable(sorterFreq)

	return sorterFreq.wordsList
}

func writeOutput(fileName string, words []WordFrequency) {
	file, _ := os.Create(fileName)

	defer file.Close()

	for _, wordFreq := range words {
		newLine := wordFreq.word + " : " + fmt.Sprint(wordFreq.frequency) + " \n"

		file.WriteString(newLine)
	}
}

// FOR testing
//--------------------------------------

// func printMap(m map[string]int) {
// 	var maxLenKey int
// 	for k, _ := range m {
// 		if len(k) > maxLenKey {
// 			maxLenKey = len(k)
// 		}
// 	}

// 	for k, v := range m {
// 		fmt.Println(k + ": " + strings.Repeat(" ", maxLenKey-len(k)) + fmt.Sprint(v))
// 	}
// }

// to compare two files
// cmp --silent ExampleOut.txt ../ExampleOut.txt || echo "files are different"
