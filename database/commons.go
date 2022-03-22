package database

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

//Contains (string in []string)
func Contains(s []string, e string) bool {
	//Check inputs
	if s == nil || e == "" {
		return false
	}

	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsOnType(s []OnType, e int) bool {
	for _, a := range s {
		if int(a) == e {
			return true
		}
	}
	return false
}

func CheckTableRule(table Table, rule OnType) bool {
	//Check if the rule is on the table
	for _, r := range table.Rule.RuleTypes {
		if ContainsOnType(r.Can, int(rule)) {
			return true
		}
	}

	return false
}

func Remove(s []string, scr string) []string {
	for i, a := range s {
		if a == scr {
			s = append(s[:i], s[i+1:]...)
			return s
		}
	}

	return s
}

func RemoveTabsFromLines(lines []string) []string {
	//Remove the first group of tabs (before the real characters)
	for i := 0; i < len(lines); i++ {
	redo_tabs:
		if strings.HasPrefix(lines[i], "\t") {
			lines[i] = strings.TrimPrefix(lines[i], "\t")
			goto redo_tabs
		}

	redo_spaces:
		if strings.HasPrefix(lines[i], "    ") {
			lines[i] = strings.TrimPrefix(lines[i], "    ")
			goto redo_spaces
		}
	}

	return lines
}

func md5sum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func AddHashToDB(hash string, path string) {
	//The hash will be saved in "./db/SCRIPTS_HASH" like that:
	//[hash], [path]

	//Create the file if it doesn't exist
	if _, err := os.Stat("./db/SCRIPTS_HASH"); os.IsNotExist(err) {
		file, err := os.Create("./db/SCRIPTS_HASH")
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Close()
	}

	//Open the file
	file, err := os.OpenFile("./db/SCRIPTS_HASH", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Write the data
	fmt.Fprintf(file, "%s, %s\n", hash, path)

	//Close the file
	file.Close()

}

func HashAlreadyContained(hash string) bool {
	//Open the file
	file, err := os.OpenFile("./db/SCRIPTS_HASH", os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//Read the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, hash) {
			return true
		}
	}

	return false
}

func GenerateUserAccessKey() string {
	rand.Seed(int64(os.Getpid()) * time.Now().UnixNano())

	//Generate a random string
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
