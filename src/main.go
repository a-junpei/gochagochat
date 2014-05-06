package gochagochat

import (
	"appengine"
	"appengine/channel"
	"appengine/user"
	"html/template"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/receive", receive)
	http.HandleFunc("/logout", logout)
}

var mainTemplate = template.Must(template.ParseFiles("tpl/index.html"))
var key_suffix = "12345"
var keys = map[string]bool{}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	c.Infof("keys = %v", keys)

	key := u.ID + key_suffix
	tok, err := channel.Create(c, key)
	if err != nil {
		http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
		c.Errorf("channel.Create: %v", err)
		return
	}
	_, ok := keys[key]
	if ok {
		c.Infof("Already")
	} else {
		keys[key] = true
	}

	err = mainTemplate.Execute(w, map[string]string{
		"token": tok,
		"me":    u.ID,
		"key":   key,
	})
	if err != nil {
		c.Errorf("mainTemplate: %v", err)
	}
}

func receive(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// g := r.FormValue("g")
	msg := r.FormValue("msg")
	c.Infof("keys = %v", keys)

	for k := range keys {
		channel.Send(c, k, "go receive!["+time.Now().String()+"] "+msg)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	g := r.FormValue("g")
	c.Infof("g = " + g)

	delete(keys, g)
	c.Infof("keys = %v", keys)

	url, err := user.LogoutURL(c, "/")
	if err != nil {
		return
	}

	c.Infof("url = " + url)

	http.Redirect(w, r, url, 302)
}
