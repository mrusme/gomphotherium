package tui

import (
	"container/list"
	"math"
	"math/rand"

	"github.com/mattn/go-runewidth"
)

func StringPtr(s string) *string {
	return &s
}

func CreateLine(
	wordList list.List,
	numberOfSpaces int,
	leftoverSpaces int) string {

	out := ""
	spaceList := make([]*string, 0)

	// we're gonna add every space to a list so that we can pad them out
	for e := wordList.Front(); e != nil; e = e.Next() {
		value, ok := e.Value.(*string)
		if ok {
			spaceList = append(spaceList, value)
		}
	}

	// allocate each justified space to one of the spaces in the list at random
	for i := leftoverSpaces; i > 0; i-- {
		selection := rand.Intn(numberOfSpaces)
		*spaceList[selection] += " "
	}

	// write the line
	for e := wordList.Front(); e != nil; e = e.Next() {
		value, ok := e.Value.(string)
		if ok {
			out += value
		} else {
			value, ok := e.Value.(*string)
			if ok {
				out += *value
			}
		}
	}
	return out
}

func WrapWithIndent(
	stringToWrap string,
	maximumWidth int,
	indentString string,
	justifyText bool,
) string {

	// list of words that will be written to the current line
	wordList := list.New()
	// number of pre-justified spaces that will be present in the current line
	spaceCount := 0
	// number of characters that have already been committed to the current line
	committedCharacterCount := 0
	characterCountOfCurrentWord := 0

	out := indentString
	word := ""
	for _, character := range stringToWrap {
		characterWidth := runewidth.RuneWidth(character)
		if character == '\n' || character == '\r' {
			if characterCountOfCurrentWord+committedCharacterCount >= maximumWidth {
				out += "\n"
				out += indentString
			}
			wordList.PushBack(word)

			// lines that end with a newline don't get justified
			line := CreateLine(*wordList, 0, 0)
			out += line
			out += "\n"
			out += indentString

			// reset all vars to new line
			word = ""
			characterCountOfCurrentWord = 0
			wordList = list.New()
			spaceCount = 0
			committedCharacterCount = 0
		} else if character == ' ' {
			// commit the current word and a space
			wordList.PushBack(word)
			wordList.PushBack(StringPtr(" "))
			spaceCount++
			committedCharacterCount += characterCountOfCurrentWord + 1

			// reset vars to new word
			characterCountOfCurrentWord = 0
			word = ""
		} else if committedCharacterCount+runewidth.StringWidth(word)+characterWidth >= maximumWidth {
			// by including this word we would exceed the maximum line. Time to wrap.

			// if this is a giant word (longer than line length) - we need to print
			// part of it otherwise we'll never write it.
			if characterCountOfCurrentWord+characterWidth >= maximumWidth {
				wordList.PushBack(word)

				//reset word vars
				word = ""
				characterCountOfCurrentWord = 0
			}

			word += string(character)
			characterCountOfCurrentWord += characterWidth

			leftoverSpace := 0
			if justifyText {
				// the amount of spaces to justify the text. 1.25 x the number of
				// spaces in the committed line looks nice.
				leftoverSpace = int(
					math.Min(
						float64(maximumWidth-(committedCharacterCount)),
						float64(spaceCount)*1.25))
			}

			// commit the line
			line := CreateLine(*wordList, spaceCount, leftoverSpace)
			out += line
			out += "\n"
			out += indentString

			// reset line vars
			spaceCount = 0
			committedCharacterCount = 0
			wordList = list.New()
		} else {
			// regular char? just add it to the word
			characterCountOfCurrentWord += characterWidth
			word += string(character)
		}
	}

	// dump what's left
	line := CreateLine(*wordList, 0, 0)
	out += line
	out += word
	return out
}
