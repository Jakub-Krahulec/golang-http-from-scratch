package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type TokenType string

const (
	String     TokenType = "string"
	Int        TokenType = "int"
	Double     TokenType = "double"
	Structural TokenType = "structural"
	Keyword    TokenType = "keyword"
)

type Token struct {
	Type  TokenType
	Value string
}

type TokenizerState int

const (
	Normal             TokenizerState = 0
	ReadingNumber      TokenizerState = 1
	ReadingString      TokenizerState = 2
	ReadingKeyword     TokenizerState = 3
	ReadingInsideArray TokenizerState = 4
)

type Tokenizer struct {
	State  TokenizerState
	Tokens []Token
}

func testJsonTokenize() {
	testJson := "{\"name\": \"Jakub\", \"age\": 25, \"isGod\": true}"
	// testJson := "{\"name\": \"Jakub\", \"age\": 28}"

	tokenizer := Tokenizer{State: Normal, Tokens: []Token{}}
	err := tokenizer.checkEmptyJson(testJson)
	if err != nil {
		fmt.Print(err)
		return
	}
	err = tokenizer.checkStartingAndEndingChar(testJson)
	if err != nil {
		fmt.Print(err)
		return
	}

	err = tokenizer.tokenizeInput(testJson)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print("TOKENY: ")
	for _, token := range tokenizer.Tokens {
		fmt.Print(token.Value + " (" + string(token.Type) + ")")
		fmt.Print("\n")
	}
}

func (t *Tokenizer) checkEmptyJson(s string) error {
	if s == "" {
		return errors.New("json input is empty")
	}
	return nil
}

func (t *Tokenizer) checkStartingAndEndingChar(s string) error {
	strLen := len(s)
	if s[0] != '{' {
		return errors.New("json must start with {")
	}
	if s[strLen-1] != '}' {
		return errors.New("json must end with }")
	}
	return nil
}

func (t *Tokenizer) tokenizeInput(s string) error {
	structuralTokens := []string{"{", "}", ",", ":"}
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	keywords := []string{"null", "true", "false"}

	stringBuffer := ""
	numberBuffer := ""
	keywordBuffer := ""

	for _, char := range s {
		currentCharAsString := string(char)

		switch t.State {
		case ReadingString:
			if currentCharAsString == "\"" {
				t.State = Normal
				t.Tokens = append(t.Tokens, Token{Type: String, Value: stringBuffer})
				stringBuffer = ""
			} else {
				stringBuffer += string(char)
			}
		case ReadingNumber:
			if slices.Contains(numbers, currentCharAsString) {
				numberBuffer += currentCharAsString
			} else if currentCharAsString == "." {
				if strings.Contains(numberBuffer, ".") {
					return errors.New("invalid number - it contains 2 dots")
				} else {
					numberBuffer += currentCharAsString
				}
			} else if currentCharAsString == "," || currentCharAsString == "}" {
				t.State = Normal
				if strings.Contains(numberBuffer, ".") {
					t.Tokens = append(t.Tokens, Token{Type: Double, Value: numberBuffer})
				} else {
					t.Tokens = append(t.Tokens, Token{Type: Int, Value: numberBuffer})
				}
				numberBuffer = ""
				t.Tokens = append(t.Tokens, Token{Type: Structural, Value: currentCharAsString})
			} else {
				return errors.New("invalid character in number")
			}
		case ReadingKeyword:
			if currentCharAsString == "," || currentCharAsString == "}" {
				keywordBuffer = strings.Trim(keywordBuffer, " ")
				if slices.Contains(keywords, keywordBuffer) {
					t.Tokens = append(t.Tokens, Token{Type: Keyword, Value: keywordBuffer})
				} else {
					return errors.New("invalid keyword")
				}

				keywordBuffer = ""
				t.State = Normal
				t.Tokens = append(t.Tokens, Token{Type: Structural, Value: currentCharAsString})
			} else {
				currentCharLowered := strings.ToLower(currentCharAsString)
				keywordBuffer += currentCharLowered
			}
		case Normal:
			currentCharLowered := strings.ToLower(currentCharAsString)
			if currentCharAsString == " " {
				continue
			} else if currentCharAsString == "\"" {
				t.State = ReadingString
			} else if slices.Contains(structuralTokens, currentCharAsString) {
				t.Tokens = append(t.Tokens, Token{Type: Structural, Value: currentCharAsString})
			} else if slices.Contains(numbers, currentCharAsString) {
				numberBuffer += currentCharAsString
				t.State = ReadingNumber
			} else if currentCharLowered == "f" || currentCharLowered == "t" || currentCharLowered == "n" {
				keywordBuffer += currentCharLowered
				t.State = ReadingKeyword
			} else {
				return errors.New("invalid character " + string(char))
			}
		}
	}

	if t.State != Normal {
		// I might add better error with which state it ended
		return errors.New("invalid json.")
	}

	return nil
}
