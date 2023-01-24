package main

import (
	"fmt"
)

type event int32

const (
	event_2x2x2 event = iota
	event_3x3x3
	event_4x4x4
	event_5x5x5
	event_6x6x6
	event_7x7x7
	event_8x8x8
	event_9x9x9
	event_3x3x3_blind
	event_4x4x4_blind
	event_5x5x5_blind
	event_megamix
	event_pyrammid
	event_skewb
	event_square1
	event_one_hand
)

func (e event) String() string {
	return []string{
		"event_2x2x2",
		"event_3x3x3",
		"event_4x4x4",
		"event_5x5x5",
		"event_6x6x6",
		"event_7x7x7",
		"event_8x8x8",
		"event_9x9x9",
		"event_3x3x3_blind",
		"event_4x4x4_blind",
		"event_5x5x5_blind",
		"event_megamix",
		"event_pyrammid",
		"event_skewb",
		"event_square1",
		"event_one_hand"}[e]
}

type tournament struct {
	tournament_id int64
	player_map    map[int64]player
	event_id      event
}

func (t *tournament) getEventStr() string {
	return t.event_id.String()
}

func (t *tournament) printState() {
	fmt.Printf("Tournament: %d\n", t.tournament_id)
	fmt.Printf("Event: %d", t.event_id)
	if len(t.player_map) == 0 {
		fmt.Printf("Player map is empty\n")
	} else {
		for id, player := range t.player_map {
			fmt.Printf("\tPlayer tuid:%d , guid: %d\n", id, player.id)
		}
	}
}

type player struct {
	id uint32
}

type tournamentB interface {
	get_stage() (uint32, error)
	get_players() ([]player, error)
	receive_results() error
	decide_stage() error
}
