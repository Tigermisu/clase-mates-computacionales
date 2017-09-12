package main

import (
	"bufio"
	"clase-sistemas-inteligentes/utilities"
	"fmt"
	"os"
	"regexp"
)

func main() {
	fmt.Println("String Identifier via Regex.")
	fmt.Print("Team Members:\nChristopher Jáquez\nHomero González\n\n")

	for {
		fmt.Print("Enter a string: ")
		input := getConsoleInput()
		fmt.Println()
		identifyString(input)
		fmt.Println()
		if !askForYesOrNo() {
			break
		}
	}
}

func askForYesOrNo() bool {
	for {
		fmt.Print("Repeat? (y/n): ")
		answer := utilities.GetConsoleInput()
		if answer != "y" && answer != "n" {
			fmt.Print("Sorry, invalid response. ")
		} else {
			return answer == "y"
		}
	}
}

func identifyString(str string) {
	// IP Match
	re := regexp.MustCompile(`(?is)^(((25[0-5])|(2[0-4][0-9])|([0-1]?[0-9]{1,2}))(\.|$)){4}$`)

	if len(re.FindStringIndex(str)) > 0 {
		fmt.Println("You have entered an I.P. Address.")
		return
	}

	// CURP Match
	re = regexp.MustCompile(`(?is)^[A-Z]{1}[AEIOU]{1}[A-Z]{2}[0-9]{2}(0[1-9]|1[0-2])(0[1-9]|1[0-9]|2[0-9]|3[0-1])[HM]{1}(AS|BC|BS|CC|CS|CH|CL|CM|DF|DG|GT|GR|HG|JC|MC|MN|MS|NT|NL|OC|PL|QT|QR|SP|SL|SR|TC|TS|TL|VZ|YN|ZS|NE)[B-DF-HJ-NP-TV-Z]{3}[0-9A-Z]{1}[0-9]{1}$`)

	if len(re.FindStringIndex(str)) > 0 {
		fmt.Println("You have entered a CURP.")
		return
	}

	// URL Match
	re = regexp.MustCompile(`(?is)^(?:[a-z]+:\/\/)?(:?[a-z0-9]*\.)+[a-z]+(?:\/[a-z0-9\/\?_\(\)\%\#\&\.\=\-]*)?`)

	if len(re.FindStringIndex(str)) > 0 {
		fmt.Println("You have entered a URL.")
		return
	}

	// Valid Date Match
	re = regexp.MustCompile(`(?is)^(?:(?:(?:(?:[0-2]?[0-9])|30)\/(?:(?:0?4)|(?:0?6)|(?:0?9)|(?:11)))|(?:(?:(?:[0-2]?[0-9])|(?:(?:30)|(?:31)))\/(?:(?:0?1)|(?:0?3)|(?:0?5)|(?:0?7)|(?:0?8)|(?:10)|(?:12)))|(?:(?:[0-1]?[0-9])|(?:2[0-8]))\/0?2)\/[0-9]+$`)

	if len(re.FindStringIndex(str)) > 0 {
		fmt.Println("You have entered a valid date.")
		return
	}

	// Credit Card
	re = regexp.MustCompile(`(?is)^(?:4[0-9]{12}(?:[0-9]{3})?|(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}|3[47][0-9]{13})$`)

	if len(re.FindStringIndex(str)) > 0 {
		fmt.Println("You have entered a Visa, MasterCard or AmericanExpress Credit Card.")
		return
	}

	// Positive Integers
	re = regexp.MustCompile(`(?is)^\+?0*[1-9]+[0-9]*(\.0*)?$`)

	if len(re.FindStringIndex(str)) > 0 {
		fmt.Println("You have entered a Positive Integer.")
		return
	}

	// Valid Hexadecimal Color in HTML
	re = regexp.MustCompile(`(?is)^#([a-f0-9]{3}){1,2}$`)

	if len(re.FindStringIndex(str)) > 0 {
		fmt.Println("You have entered a Valid Hex Color in HTML.")
		return
	}

	// HTML Tag
	re = regexp.MustCompile(`(?is)^(<[a-z]+.*)((>.*?<\/[a-z]+[a-z0-9]*>)|(\/ *>))`)

	if len(re.FindStringIndex(str)) > 0 {
		fmt.Println("You have entered an HTML tag.")
		return
	}

	fmt.Println("Your string cannot be identified.")
}

func getConsoleInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
