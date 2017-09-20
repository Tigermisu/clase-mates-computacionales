package utilities

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// AssertError checks an error and panics if it is not null
func AssertError(e error) {
	if e != nil {
		fmt.Println("Error!\n", e)
		panic(e)
	}
}

// GetConsoleInput reads a string from os.Stdin and returns it
func GetConsoleInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// PrettyPrint formats a variable as an indented JSON and prints it
func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

// LoadFile reads a file and returns a string with its contents
func LoadFile(path string) string {
	data, err := ioutil.ReadFile(path)
	AssertError(err)
	return string(data)
}

// GetData requests an HTTP URL for arbitrary data and returns the response's body
func GetData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// PostJSONData sends a variable encoded as a JSON to a given URL using the HTTP verb POST
func PostJSONData(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

// PutJSONData sends a variable encoded as a JSON to a given URL using the HTTP verb PUT
func PutJSONData(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}
