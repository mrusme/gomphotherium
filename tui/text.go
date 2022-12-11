package tui

import (
	"container/list"
	"math"
	"math/rand"
	"strings"

	"github.com/mattn/go-runewidth"
)

type Indent struct {
	Width int
	Lines []string
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (x *Indent) InitializeWithString(width int, line string) {
	var lines = make([]string, 1)
	lines[0] = line
	x.InitializeWithArray(width, lines)
}

func (x *Indent) InitializeWithArray(width int, lines []string) {
	x.Width = width
	x.Lines = lines // TODO: validation?
}

func (x *Indent) GetLine(line int) string {
	if line < len(x.Lines) {
		return x.Lines[line]
	} else {
		return strings.Repeat(" ", x.Width)
	}
}

func (x *Indent) IndentSlice(lines *[]string) *[]string {
	numberOfLines := Max(len(x.Lines), len(*lines))
	var newLines []string = make([]string, 0)
	r := 0
	for r < numberOfLines {
		line := ""
		if len(x.Lines) <= r {
			line = x.Lines[len(x.Lines)-1]
		} else {
			line = x.Lines[r]
		}

		if len(*lines) <= r {
			line += strings.Repeat(" ", 1)
		} else {
			line += (*lines)[r]
		}
		newLines = append(newLines, line)
		r++
	}
	return &newLines
}

func (x *Indent) ExtendWithIndent(with Indent) *Indent {
	// this method will access the variable
	// food in Animal class
	numberOfLines := Max(len(x.Lines), len(with.Lines))
	indent := &Indent{
		Width: x.Width + with.Width,
		Lines: make([]string, numberOfLines),
	}
	r := 0
	for r < numberOfLines {
		line := ""
		if len(x.Lines) <= r {
			line = strings.Repeat(" ", x.Width)
		} else {
			line = x.Lines[r]
		}

		if len(with.Lines) <= r {
			line += strings.Repeat(" ", with.Width)
		} else {
			line += with.Lines[r]
		}

		indent.Lines[r] = line
		r++
	}
	return indent
}

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
	rand.Seed(112358132134) // consistent seed so that results are the same between reloads
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
	justifyText bool,
) *[]string {

	// list of words that will be written to the current line
	wordList := list.New()
	// current line we're on
	lineCount := 0
	// number of pre-justified spaces that will be present in the current line
	spaceCount := 0
	// number of characters that have already been committed to the current line
	committedCharacterCount := 0
	characterCountOfCurrentWord := 0

	var outList []string // = make([]string, 1)
	word := ""
	for _, character := range stringToWrap {
		characterWidth := runewidth.RuneWidth(character)
		if character == '\n' || character == '\r' {
			if characterCountOfCurrentWord+committedCharacterCount >= maximumWidth {
				line := CreateLine(*wordList, 0, 0)
				outList = append(outList, line)
				lineCount++
				wordList = list.New()
			}
			wordList.PushBack(word)

			// lines that end with a newline don't get justified
			line := CreateLine(*wordList, 0, 0)
			outList = append(outList, line)
			lineCount++

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
			outList = append(outList, line)
			lineCount++

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
	wordList.PushBack(word)
	line := CreateLine(*wordList, 0, 0)
	outList = append(outList, line)

	return &outList
}
