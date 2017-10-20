package slashgordon

import (
	"fmt"
	"net/http"
	"github.com/gorilla/schema"
	"encoding/json"
	"bytes"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

type Request struct {
	Token			string `schema:"token"`
	TeamID			string `schema:"team_id"`
	TeamDomain		string `schema:"team_domain"`
	EnterpriseID	string `schema:"enterprise_id"`
	EnterpriseName	string `schema:"enterprise_name"`
	ChannelID		string `schema:"channel_id"`
	ChannelName		string `schema:"channel_name"`
	UserID			string `schema:"user_id"`
	UserName		string `schema:"user_name"`
	Command			string `schema:"command"`
	Text			string `schema:"text"`
	ResponseURL		string `schema:"response_url"`
	TriggerID		string `schema:"trigger_id"`
}

var decoder = schema.NewDecoder()

type Response struct {
	ResponseType	string `json:"response_type,omitifempty"`
	Text			string `json:"text"`
	Attachments		[]Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Text			string `json:"text"`
	ImageURL		string `json:"image_url"`
	CallbackID		string `json:"callback_id,omitempty"`
	Actions			[]Action `json:"actions,omitempty"`
}

type Action struct {
	Name			string `json:"name"`
	Text			string `json:"text,omitempty"`
	Type			string `json:"type"`
	Value			string `json:"value"`
	Confirmation	Confirmation `json:"confirm,omitempty"`
}

type Confirmation struct {
	Title			string `json:"title"`
	Text			string `json:"text"`
	OKText			string `json:"ok_text"`
	DismissText		string `json:"dismiss_text"`
}

type Payload struct {
	Actions			[]Action `json:"actions"`
	CallbackID		string `json:"callback_id"`
	User			User `json:"user"`
	ResponseURL		string `json:"response_url"`
}

type User struct {
	ID				string `json:"id"`
	Name			string `json:"name"`
}

func init() {
	http.HandleFunc("/doge", DogeHandler)
	http.HandleFunc("/goodboy", GoodBoyHandler)
}

func DogeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			// Error 500
		}

		var request Request
		err = decoder.Decode(&request, r.PostForm)
		if err != nil {
			// Error 500
		}

		// Verify Token
		// Check if it's an action

		response := Response{
			ResponseType: "in_channel",
			Text: "Här kommer varsågod:",
			Attachments: []Attachment{
				{
					Text: "rolig_hund.gif",
					ImageURL: "http://i0.kym-cdn.com/photos/images/original/000/581/296/c09.jpg",
				},
				{
					Text: "Duktig vovve?",
					CallbackID: "good-boy-callback",
					Actions: []Action{
						{
							Name: "good_boy",
							Text: "Ja! :heart_eyes:",
							Type: "button",
							Value: "true",
							Confirmation: Confirmation{
								Title: "Är du säker?",
								Text: "Alltså helt säker?",
								OKText: "Ja!",
								DismissText: "Nej!",
							},
						},
						{
							Name: "good_boy",
							Text: "Nej! :disappointed:",
							Type: "button",
							Value: "false",
							Confirmation: Confirmation{
								Title: "Är du säker?",
								Text: "Alltså helt säker?",
								OKText: "Ja!",
								DismissText: "Nej!",
							},
						},
					},
				},
			},
		}

		json, err := json.Marshal(response)
		if err != nil {
			// Error 500
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", json)
	}
}

func GoodBoyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			// Error 500
		}

		var p Payload
		err = json.Unmarshal([]byte(r.FormValue("payload")), &p)
		if err != nil {
			// Error 500
		}

		var msg string
		if p.CallbackID == "good-boy-callback" {
			// Make this better
			if p.Actions[0].Value == "true" {
				msg = fmt.Sprintf("Woho, <@%s> tycker om mig! :heart_eyes:", p.User.ID)
			} else {
				msg = fmt.Sprintf("Vad har jag _någonsin_ gjort dig, <@%s>? :cry:", p.User.ID)
			}
		}
		response := Response{
			ResponseType: "in_channel",
			Text: msg,
		}

		json, err := json.Marshal(response)
		if err != nil {
			// Error 500
		}

		jsonBlob := bytes.NewBuffer([]byte(json))

		ctx := appengine.NewContext(r)
		client := urlfetch.Client(ctx)
		_, err = client.Post(p.ResponseURL, "application/json", jsonBlob)
		if err != nil {
			// Error 500
		}
	}
}