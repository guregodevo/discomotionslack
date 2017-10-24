package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"strings"

	"github.com/dailymotion-leo/discomotion/models"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

type ServerStatus struct {
	ServerInfo

	Uptime string `json:"uptime"`
}

// Responds OK if service is running properly
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	var status = ServerStatus{
		ServerInfo: ServerInfo{
			Hostname:  s.Info.Hostname,
			Server:    s.Info.Server,
			Version:   s.Info.Version,
			Build:     s.Info.Build,
			BuildTime: s.Info.BuildTime,
		},
		Uptime: time.Now().Sub(s.Uptime).String(),
	}

	WriteJson(w, &status, http.StatusOK)
}

func sendCoreEvent(idx string, now bool, baseUrl string, channelId string) {

	v := models.PlayVideo{VideoId: idx, Now: now}

	b, _ := json.Marshal(v)

	log.Debug("json", string(b))

	url := fmt.Sprintf("%s/playlist/%s/music", baseUrl, channelId)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

// NOTE: if in the future we start having problems with non-received messages
// we must start using action.ResponseURL to send instead of writing directly.

func (s *Server) Interactive(w http.ResponseWriter, r *http.Request) {
	action := &slack.AttachmentActionCallback{}

	err := json.Unmarshal([]byte(r.PostFormValue("payload")), action)
	if err != nil {
		writeText(w, "ERROR Cannot get payload", http.StatusForbidden)
		return
	}

	//TODO if action.Token != c.VerificationToken {
	//	return 0, nil
	//}

	//if action.CallbackID != "accept_channel" {
	//	writeJson(w, "ERROR Channel Not accept", http.StatusForbidden)
	//	return
	//}

	if action.Actions[0].Value == "" {
		writeText(w, "You've rejected this website approval request.", http.StatusBadRequest)
		return
	}

	if len(action.Actions) == 0 {
		writeText(w, "missing action", http.StatusBadRequest)
		return
	}

	names := strings.Split(action.Actions[0].Name, "::::")
	actionname := names[0]
	username := names[1]
	title := names[2]

	idx := action.Actions[0].Value

	channelId := action.Channel.ID

	log.Debug("channelID", channelId)
	log.Debug("videoID", idx)
	log.Debug("action", actionname)
	log.Debug("username", username)

	//usernames := strings.Split(username, ".")
	//username = fmt.Sprintf("%s%s", usernames[0][:1], usernames[1])

	if actionname == "queue" {
		actionname = "queu"
		go sendCoreEvent(idx, false, s.BaseURL, channelId)
	} else if actionname == "play" {
		go sendCoreEvent(idx, true, s.BaseURL, channelId)
	}

	writeText(w, fmt.Sprintf("%s %sed \"%s\"", username, actionname, title), http.StatusOK)

	return
}

// Responds to slack command /play
func (s *Server) Play(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("r.Form", r.Form)

	values := r.Form

	channelId := values.Get("channel_id")

	log.Debug("channel_id found", channelId)

	params := slack.PostMessageParameters{}

	text := values.Get("text")
	username := values.Get("user_name")

	res, rerr := http.Get(fmt.Sprintf("%s/search/%s", s.BaseURL, text))
	defer res.Body.Close()

	if rerr != nil {
		log.Fatal(err)
	}
	searchresp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("searchresp %s", searchresp)

	var searchResult []models.PlaySearchItem

	err = json.Unmarshal(searchresp, &searchResult)
	if err != nil {
		log.Fatalln("can't parse response:", err)
	}

	attachments := []slack.Attachment{}

	for _, item := range searchResult {
		attachments = append(attachments, slack.Attachment{
			Fallback:   "fallback",
			CallbackID: "interactive_fdqsdfjqsldkfj",
			Fields: []slack.AttachmentField{slack.AttachmentField{
				Title: item.Title,
				Value: "",
			}},
			Actions: []slack.AttachmentAction{slack.AttachmentAction{
				Text:  "queue",
				Type:  "button",
				Name:  strings.Join([]string{"queue", username, item.Title}, "::::"),
				Value: item.Id,
			}, slack.AttachmentAction{
				Text:  "play now",
				Type:  "button",
				Name:  strings.Join([]string{"play", username, item.Title}, "::::"),
				Value: item.Id,
			}},
		})
	}

	params.Attachments = attachments
	channelID, timestamp, e := s.Api.PostMessage(channelId, fmt.Sprintf("You searched for \"%s\"", text), params)
	if e != nil {
		log.Error("%s\n", e.Error())
		writeText(w, e.Error(), http.StatusOK)
		return
	}
	respMsg := fmt.Sprintf("Message successfully sent to channel %s at %s", channelID, timestamp)
	log.Debug(respMsg)

	w.WriteHeader(http.StatusOK)
	//writeJson(w, "", http.StatusOK)
}

/*func writeText(w http.ResponseWriter, text string, status int) {

	data := url.Values{}
	data.Add("text", text)
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Cache-Control", "no-store")

	w.WriteHeader(status)
	dataEncoded := data.Encode()
	w.Write([]byte(dataEncoded))
}*/

func writeText(w http.ResponseWriter, data string, status int) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")

	w.WriteHeader(status)
	w.Write([]byte(data))
}

func WriteJson(w http.ResponseWriter, info interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")

	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(info); err != nil {
		log.WithField("error", err).Info(fmt.Sprintf("Failed to write JSON"))
	}
}
