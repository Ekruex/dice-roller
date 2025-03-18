package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RollResult struct {
	Input       string `json:"input"`
	Rolls       []int  `json:"rolls"`
	Modifier    int    `json:"modifier"`
	Total       int    `json:"total"`
	FortuneType string `json:"fortuneType,omitempty"`
}

// Roll XdY dice
func rollDice(count, sides int) ([]int, int) {
	rolls := make([]int, count)
	total := 0
	for i := 0; i < count; i++ {
		rolls[i] = rand.Intn(sides) + 1
		total += rolls[i]
	}
	return rolls, total
}

// Roll 2d20 and pick highest (fortune) or lowest (misfortune)
func rollWithFortuneOrMisfortune(rollType string) ([]int, int) {
	rolls, _ := rollDice(2, 20)
	if rollType == "fortune" {
		return rolls, max(rolls[0], rolls[1])
	} else if rollType == "misfortune" {
		return rolls, min(rolls[0], rolls[1])
	}
	return nil, 0
}

// Parse roll input (e.g., "2d6+3", "1d20kh1", "3d8-1")
func parseRollInput(input string) (int, int, int, string, error) {
	input = strings.ToLower(strings.ReplaceAll(input, "D", "d")) // Normalize "D" to "d"

	// Updated regex pattern
	pattern := `^(\d*)d(\d+)(kh1|kl1)?(?:([+-]?)\s*(\d+))?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(strings.TrimSpace(input))

	if matches == nil {
		return 0, 0, 0, "", fmt.Errorf("invalid format: use XdY+Z")
	}

	// Dice count (default = 1)
	count := 1
	if matches[1] != "" {
		count, _ = strconv.Atoi(matches[1])
	}

	// Dice sides
	sides, _ := strconv.Atoi(matches[2])

	// Fortune/Misfortune type
	rollType := ""
	if matches[3] == "kh1" {
		rollType = "fortune"
	} else if matches[3] == "kl1" {
		rollType = "misfortune"
	}

	// ❌ Prevent Fortune/Misfortune on non-d20 rolls
	if (rollType == "fortune" || rollType == "misfortune") && sides != 20 {
		return 0, 0, 0, "", fmt.Errorf("Fortune/Misfortune only applies to d20 rolls")
	}

	// Modifier (default = 0)
	modifier := 0
	if matches[5] != "" {
		mod, _ := strconv.Atoi(matches[5])
		if matches[4] == "-" {
			modifier = -mod
		} else {
			modifier = mod
		}
	}

	return count, sides, modifier, rollType, nil
}

// Handle dice roll requests
func rollHandler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	input = strings.ReplaceAll(input, " ", "+")

	if input == "" {
		http.Error(w, "Missing input parameter", http.StatusBadRequest)
		return
	}

	// Log request
	fmt.Println("Received request: /roll?input=", input)

	// Normalize and parse input
	input, _ = url.QueryUnescape(input)
	count, sides, modifier, rollType, err := parseRollInput(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var rolls []int
	var total int
	var resultStr string

	if rollType == "fortune" || rollType == "misfortune" {
		rolls, total = rollWithFortuneOrMisfortune(rollType)
		total += modifier
		resultStr = fmt.Sprintf("Rolling %s... (%v) %+d = %d [%s]", input, rolls, modifier, total, rollType)
	} else {
		rolls, total = rollDice(count, sides)
		total += modifier
		rollsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(rolls)), ", "), "[]")
		resultStr = fmt.Sprintf("Rolling %s... (%s) %+d = %d", input, rollsStr, modifier, total)
	}

	fmt.Println("Result:", resultStr)
	w.Write([]byte(resultStr))
}

// Utility functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Dice Roller API! Try /roll?input=1d20+5"))
}

// Main function
func main() {
	rand.Seed(time.Now().UnixNano())

	// Define routes
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/roll", rollHandler)

	// Start server
	fmt.Println("Server running on http://localhost:5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
