package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func HttpGet(url string) (string, error) {
	//TODO - implement http get logic
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
		return "", fmt.Errorf("Error making GET request: %v", err)
	}
	// Defer the closing of the response body
	// The body is of type io.ReadCloser and must be closed to prevent resource leaks
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Print the response body as a string
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body:\n%s\n", string(body))
	return string(body), nil
}

func HttpPost(url string, contentType string, data io.Reader) (string, error) {
	//TODO - implement http post logic
	resp, err := http.Post(url, contentType, data)
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
		return "", fmt.Errorf("Error making POST request: %v", err)
	}
	// Defer the closing of the response body
	// The body is of type io.ReadCloser and must be closed to prevent resource leaks
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Print the response body as a string
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body:\n%s\n", string(body))
	return string(body), nil
}
