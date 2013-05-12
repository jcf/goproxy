package main

import (
  "fmt"
  "html"
  "io/ioutil"
  "log"
  "net/http"
  "time"
)

/* func benchmark(description string, f func()) interface{} { */
/*   t := time.Now() */
/*   f() */
/*   delta := time.Since(t) */
/*   log.Printf("%s took %v", description, delta) */
/* } */

func fetchRemote(url string, c chan *http.Response) {
  client := &http.Client{
  /* CheckRedirect: redirectPolicyFunc, */
  }

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    log.Fatal("Failed to create new request for %s", url)
  }

  req.Header.Add("User-Agent", "GoProxy")

  t := time.Now()
  resp, err := client.Do(req)

  if err == nil {
    log.Printf("Fetched %s in %v", url, time.Since(t))
    c <- resp
  } else {
    log.Fatal("Request to %s failed", url)
  }
}

func readResponse(c chan *http.Response) {
  resp := <- c

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  if err == nil {
    log.Printf("Got response!\n%s", body[:80])
    fmt.Printf("%s\n", body)
  } else {
    log.Fatal("Failed to read response from %v", resp)
  }
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
  url := html.EscapeString(r.FormValue("url"))
  // log.Printf("Proxying %s\n", url)

  c := make (chan *http.Response)
  go fetchRemote(url, c)
  go readResponse(c)
}

func main() {
  // resp, err := http.Get("http://example.com/")

  http.HandleFunc("/", handleRequest)

  // http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
  //   fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
  // })

  log.Printf("Starting up... http://localhost:8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
