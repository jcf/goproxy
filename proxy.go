package main

import (
  "fmt"
  "html"
  "io/ioutil"
  "log"
  "net/http"
  "time"
)

func fetchRemote(url string) (*http.Response, error) {
  client := &http.Client{
  /* CheckRedirect: redirectPolicyFunc, */
  }

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    log.Fatal("Failed to create new request for %s", url)
  }

  req.Header.Add("User-Agent", "GoProxy")

  return client.Do(req)
}

func readResponse(resp *http.Response) ([]byte, error) {
  defer resp.Body.Close()
  return ioutil.ReadAll(resp.Body)
}

// TODO DRY
func handleRequest(w http.ResponseWriter, r *http.Request) {
  url := html.EscapeString(r.FormValue("url"))

  t := time.Now()
  resp, err := fetchRemote(url)
  duration := time.Since(t)

  if err == nil {
    log.Printf("Fetched %s in %v", url, duration)
  } else {
    log.Fatal("Request to %s failed", url)
  }

  body, err := readResponse(resp)

  if err == nil {
    fmt.Fprintf(w, "%s\n", body)
  } else {
    log.Fatal("Request to %s failed", url)
  }
}

func main() {
  http.HandleFunc("/", handleRequest)

  log.Printf("Starting up... http://localhost:8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
