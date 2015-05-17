// +build ignore

// This is a small chat client for Podio
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/andreas/podio-go"
	"github.com/andreas/podio-go/conversation"
)

var (
	podioClient string
	podioSecret string

	defaultCacheFile = filepath.Join(os.Getenv("HOME"), ".chat-cli-request-token")

	cacheFile = flag.String("cachefile", defaultCacheFile, "Authentication token cache file")
)

func main() {
	flag.Parse()

	podioClient = envDefault("PODIO_CLIENT", "chatcli")
	podioSecret = envDefault("PODIO_SECRET", "4qCaud5yZTt56w6WWbsKp1ldoq0egbEqzTuq7kIU6X6IKy9f9Gjp4K9M9zttXJul")

	token, err := readToken(*cacheFile)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Err reading token cache: %s\n", err)
	}
	if token == nil {
		authcode, err := getOauthToken()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting auth from podio:", err)
			os.Exit(1)
		}
		token, err = podio.AuthWithAuthCode(
			podioClient, podioSecret,
			authcode, "http://127.0.0.1/",
		)
	}
	err = writeToken(*cacheFile, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err writing token file: %s\n", err)
	}

	client := &conversation.Client{podio.NewClient(token)}

	id, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Bad or no conversation ID given. Listing conversations\n")
		listConversations(client)
	} else {
		talkTo(client, uint(id))
	}

}

func prompt(q string) string {
	fmt.Printf("%s: ", q)
	defer fmt.Println("")

	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	out := s.Text()
	return out
}

func envDefault(key, deflt string) string {
	val := os.Getenv(key)
	if val == "" {
		return deflt
	}
	return val
}

func listConversations(client *conversation.Client) {
	convs, err := client.GetConversations(conversation.WithLimit(200))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting conversation list:", err)
		return
	}
	for _, conv := range convs {
		if conv.Type == "direct" {
			fmt.Println("Id:", conv.ConversationId, "direct", conv.Participants[0].Name)
		} else {
			fmt.Println("Id:", conv.ConversationId, "group", len(conv.Participants), "colleagues on", conv.Subject)
		}
	}
}

func talkTo(client *conversation.Client, convId uint) {
	var (
		eventChan = make(chan conversation.Event, 1)
		inputChan = make(chan string)
	)

	go func() {
		last := conversation.Event{}
		for {

			events, err := client.GetEvents(convId, conversation.WithLimit(1))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Err getting update:", err)
				time.Sleep(800 * time.Millisecond)
				continue
			}

			if len(events) == 0 || last.EventID == events[0].EventID {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			last = events[0]
			eventChan <- last
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			inputChan <- s.Text()
		}

		if s.Err() != nil {
			fmt.Fprintln(os.Stderr, "Error scanning for input:", s.Err())
		}
	}()

	lastTalker := uint(0)

	for {
		select {
		case t := <-inputChan:
			_, err := client.Reply(convId, t)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error replying to Podio:", err)
			}

		case event := <-eventChan:
			if event.CreatedBy.Id != lastTalker {
				fmt.Println(event.CreatedBy.Name, "said:")
			}
			lastTalker = event.CreatedBy.Id
			fmt.Println(" > ", event.Data.Text)
		}
	}

}

// Inspired by how github.com/nf/streak uses Google Oauth to avoid having user type in password
// Requires that the redirect_url of the client is set to 127.0.0.1
// TODO: part of podio-go?

func getOauthToken() (string, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	defer l.Close()

	code := make(chan string)
	go http.Serve(l, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "You can close this window now")
		fmt.Println(req)
		code <- req.FormValue("code")
	}))

	u, _ := url.Parse("https://podio.com/oauth/authorize")
	params := url.Values{}

	params.Add("client_id", podioClient)
	params.Add("redirect_uri", fmt.Sprintf("http://%s/", l.Addr()))
	u.RawQuery = params.Encode()
	openURL(u.String())

	return <-code, nil
}

func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err

}

func readToken(f string) (*podio.AuthToken, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var t podio.AuthToken
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func writeToken(f string, t *podio.AuthToken) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f, b, 0600)
}
