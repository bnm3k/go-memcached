package main

import (
	"encoding/json"
	"net/http"
)

type stdReply struct {
	Reply string `json:"reply"`
}

type getReply struct {
	Reply string `json:"reply"`
	Val   string `json:"val"`
}

type getReplyToken struct {
	Reply string `json:"reply"`
	Val   string `json:"val"`
	Token string `json:"token"`
}

func (api *httpAPI) home(w http.ResponseWriter, r *http.Request) {
	reply := "Hello go-memcached"
	jsonString, _ := json.Marshal(
		stdReply{reply})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func getStdParams(r *http.Request) (string, string, string) {
	key := r.URL.Query().Get(":key")
	val := r.URL.Query().Get("val")
	exptimeStr := r.URL.Query().Get("exp")
	if exptimeStr == "" {
		exptimeStr = "-1"
	}
	return key, val, exptimeStr
}

func (api *httpAPI) handleSet(w http.ResponseWriter, r *http.Request) {
	key, val, exptimeStr := getStdParams(r)
	reply := api.cache.Set(key, val, exptimeStr)
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleAdd(w http.ResponseWriter, r *http.Request) {
	key, val, exptimeStr := getStdParams(r)
	reply := api.cache.Add(key, val, exptimeStr)
	//return only reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleReplace(w http.ResponseWriter, r *http.Request) {
	key, val, exptimeStr := getStdParams(r)
	reply := api.cache.Replace(key, val, exptimeStr)
	//return only reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleAppend(w http.ResponseWriter, r *http.Request) {
	key, val, exptimeStr := getStdParams(r)
	reply := api.cache.Append(key, val, exptimeStr)
	//return only reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handlePrepend(w http.ResponseWriter, r *http.Request) {
	key, val, exptimeStr := getStdParams(r)
	reply := api.cache.Prepend(key, val, exptimeStr)
	//return only reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleIncrement(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	numStr := r.URL.Query().Get("num")

	reply := api.cache.Increment(key, numStr)
	//return only reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleDecrement(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	numStr := r.URL.Query().Get("num")
	reply := api.cache.Decrement(key, numStr)
	//return only reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleCompareAndSwap(w http.ResponseWriter, r *http.Request) {
	key, val, exptimeStr := getStdParams(r)
	reply := api.cache.CompareAndSwap(key, val, exptimeStr, "")
	//return only reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	reply, val := api.cache.Get(key)
	//return reply & val

	jsonString, _ := json.Marshal(
		getReply{string(reply), val})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleGetEntryPlusToken(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	reply, val, token := api.cache.GetEntryPlusToken(key)
	//return reply & val & token
	jsonString, _ := json.Marshal(
		getReplyToken{string(reply), val, string(token)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleDelete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	reply := api.cache.Delete(key)
	//return reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleClear(w http.ResponseWriter, r *http.Request) {
	reply := api.cache.Clear()
	//return reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func (api *httpAPI) handleStats(w http.ResponseWriter, r *http.Request) {
	reply := api.cache.Stats()
	//return reply
	jsonString, _ := json.Marshal(
		stdReply{string(reply)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}
