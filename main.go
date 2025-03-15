package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// File for logging session rolls
const logFileName = "session_log.txt"

// Keeps track of roll history in memory
var rollHistory []string

// rollDice rolls 'count' dice with 'sides' sides and returns a slice of individual rolls and the total sum.
func rollDice(count, sides int) ([]int, int) {
	rolls := []int{}
	total := 0
	for i := 0; i < count; i++ {
		roll := rand.Intn(sides) + 1
		rolls = append(rolls, roll)
		total += roll
	}
	return rolls, total
}

// rollWithFortuneOrMisfortune handles PF2e "Fortune" (2d20kh1) and "Misfortune" (2d20kl1)
func rollWithFortuneOrMisfortune(rollType string) (string, int) {
	rolls, _ := rollDice(2, 20) // Always rolls 2d20
	rollsStr := fmt.Sprintf("%d, %d", rolls[0], rolls[1])

	if rollType == "fortune" { // Keep highest (Fortune)
		return rollsStr, max(rolls[0], rolls[1])
	} else if rollType == "misfortune" { // Keep lowest (Misfortune)
		return rollsStr, min(rolls[0], rolls[1])
	}
	return "", 0 // Should never reach here
}

// parseRollInput extracts count, sides, and modifier from a roll string like "2d6+3" or "2d20kh1"
func parseRollInput(input string) (int, int, int, string, error) {
	pattern := `^(\d*)d(\d+)(kh1|kl1)?([+-]\d+)?$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(strings.TrimSpace(input))
	if matches == nil {
		return 0, 0, 0, "", fmt.Errorf("invalid format: use XdY+Z or 2d20kh1/kl1")
	}

	// Extract dice count (default to 1 if empty)
	count := 1
	if matches[1] != "" {
		count, _ = strconv.Atoi(matches[1])
	}

	// Extract dice sides
	sides, _ := strconv.Atoi(matches[2])

	// Detect Fortune (kh1) or Misfortune (kl1)
	rollType := ""
	if matches[3] == "kh1" {
		rollType = "fortune"
	} else if matches[3] == "kl1" {
		rollType = "misfortune"
	}

	// Extract modifier (default to 0 if empty)
	modifier := 0
	if matches[4] != "" {
		modifier, _ = strconv.Atoi(matches[4])
	}

	return count, sides, modifier, rollType, nil
}

// logRoll saves a roll to the history and writes it to the log file.
func logRoll(entry string) {
	rollHistory = append(rollHistory, entry) // Store in memory

	// Append to log file
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error writing to log file:", err)
		return
	}
	defer file.Close()

	_, _ = file.WriteString(entry + "\n")
}

// processMultipleRolls handles multiple rolls entered as a comma-separated string
func processMultipleRolls(input string) {
	rolls := strings.Split(input, ",") // Split input into separate rolls

	for _, roll := range rolls {
		roll = strings.TrimSpace(roll) // Remove extra spaces
		count, sides, modifier, rollType, err := parseRollInput(roll)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		var resultStr string

		// Handle Fortune (2d20kh1) and Misfortune (2d20kl1)
		if rollType == "fortune" || rollType == "misfortune" {
			rollsStr, result := rollWithFortuneOrMisfortune(rollType)
			finalResult := result + modifier
			resultStr = fmt.Sprintf("Rolling %s... (%s) %+d = %d [%s]", roll, rollsStr, modifier, finalResult, rollType)
		} else {
			// Regular roll
			rollResults, rollTotal := rollDice(count, sides)
			rollsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(rollResults)), ", "), "[]")
			finalResult := rollTotal + modifier
			resultStr = fmt.Sprintf("Rolling %s... (%s) %+d = %d", roll, rollsStr, modifier, finalResult)
		}

		// Print and log the roll
		fmt.Println(resultStr)
		logRoll(resultStr)
	}
}

// showHistory displays all past rolls from the session
func showHistory() {
	if len(rollHistory) == 0 {
		fmt.Println("No rolls have been made yet.")
		return
	}

	fmt.Println("\n--- Roll History ---")
	for _, entry := range rollHistory {
		fmt.Println(entry)
	}
	fmt.Println("--------------------")
}

// Helper functions to get max and min values
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

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed once at program start

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nEnter your rolls (e.g., 1d20+5, 2d6+3, 2d20kh1), or type 'history' to view previous rolls, or 'exit' to quit: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Exiting dice roller. Goodbye!")
			break
		} else if input == "history" {
			showHistory()
		} else {
			processMultipleRolls(input)
		}
	}
}
