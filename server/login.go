package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type authForm struct {
	Username string
	Password string
}

type TimesCollect struct {
	time       string
	scramble   string
	parsedTime string
	cookie     string
}

type TimesStore struct {
	times map[string][]TimesCollect
}

func (srp *sessionResolutionPoints) registerUserCA(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Register method not supported", http.StatusBadRequest)
		return
	}

	Debugf("req: %s\n", string(b))
	var plf authForm
	json.Unmarshal(b, &plf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Debugf("%s\n", b)
	username := plf.Username
	password := plf.Password
	Debugf("\tpass %s\n", password)
	Debugf("\tuser %s\n", username)
	_, ok := srp.users[username]

	if ok {
		http.Error(w, "User exists Credentials", http.StatusConflict)
		return
	}
	srp.users[username] = password
	w.Write([]byte("registration succesful"))
	for k, v := range srp.users {
		fmt.Printf("%s, %s \n", k, v)
	}
}

func (srp *sessionResolutionPoints) loginHandlerCA(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	fmt.Println("xd")
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Login method not supported", http.StatusBadRequest)
		return
	}
	Debugf("req: %s\n", string(b))
	var plf authForm
	err = json.Unmarshal(b, &plf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Debugf("%s\n", b)
	name := plf.Username
	pass := plf.Password
	//databese connection here???
	storedPassword, ok := srp.users[name]

	if !ok {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	//origin := r.Header.Get("Origin")

	session, _ := srp.store.Get(r, "session.id")
	if storedPassword == pass {
		Debugf("[INFO] user validated")
		session.Values["authenticatred"] = "true"
		session.Save(r, w)
		http.Redirect(w, r, "timer", 302)
	} else {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
	}
}

func (srp *sessionResolutionPoints) logoutHandlerCA(w http.ResponseWriter, r *http.Request) {
	session, _ := srp.store.Get(r, "session.id")
	session.Values["authenticated"] = false
	session.Save(r, w)
	w.Write([]byte("Logout Succesfull"))
}

func (srp *sessionResolutionPoints) healtcheckCA(w http.ResponseWriter, r *http.Request) {
	session, _ := srp.store.Get(r, "session.id")
	authenticated := session.Values["authenticated"]
	if authenticated != nil && authenticated != false {
		w.Write([]byte("Welcome!"))
		return
	}
	http.Error(w, "forbidden", http.StatusForbidden)
	return
}

func (srp *sessionResolutionPoints) timerCollectCA(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Collection method not supported", http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Register method not supported", http.StatusBadRequest)
		return
	}

	Debugf("req: %s\n", string(b))
	var tc TimesCollect
	json.Unmarshal(b, &tc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(tc)
	//srp.times.times[tc.cookie].
}
