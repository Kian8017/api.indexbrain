package main

import ()

type ClientMessageType string

const (
	QueryType     ClientMessageType = "query"
	GetPlacesType ClientMessageType = "getplaces"
	GetBannerType ClientMessageType = "getbanner"
)

type ClientMessage struct {
	Type   ClientMessageType `json:"type"`
	Params string            `json:"params"`
}
