package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type RoomFunc func(r *Room)

type Item struct {
	name        string
	description string
}

type ItemStatus struct {
	uses   int
	broken bool
}

type Monster struct {
	name        string
	description string
	// movement...
}

type Room struct {
	short_description string
	long_description  string
	exits             map[string]*Room
	entrances         map[string]*Room
	enterFunc         RoomFunc
	exitFunc          RoomFunc
	visits            int
	items             map[string]Item
}

type Game struct {
	current_room *Room
	inventory    map[string]Item
	moves        int
}

var directions map[string]string = map[string]string{
	"n": "north",
	"s": "south",
	"e": "east",
	"w": "west"}

var opposites map[string]string = map[string]string{
	"n": "s",
	"e": "w",
	"u": "d"}

func get_opposite(dir string) string {
	var opposite string
	for k, v := range opposites {
		if k == dir {
			opposite = v
		}
		if v == dir {
			opposite = k
		}
	}
	return opposite
}

func look_dir(dir string, r *Room) string {
	switch dir {
	case "n":
	case "s":
	case "e":
	case "w":
		dirname := directions[dir]
		return "To the " + dirname + " you see "
	case "u":
		return "Above you, you see " + r.short_description
	case "d":
		return "Below you, you see " + r.short_description
	default:
		return "You can't quite make out what you see"
	}
	return "What?"
}

func look(g *Game, line []string) {
	fmt.Println(g.current_room.long_description)
	for direction, e := range g.current_room.exits {
		fmt.Println(look_dir(direction, e))
	}

	if len(g.current_room.items) > 0 {
		fmt.Println("You see the following items before you:")
		for _, item := range g.current_room.items {
			_, item_in_inventory := g.inventory[item.name]
			if !item_in_inventory {
				fmt.Println("  " + item.name)
			}
		}
	}
}

func take(g *Game, line []string) {
	item_name := line[1]
	item, item_in_room := g.current_room.items[item_name]
	_, item_in_inventory := g.inventory[item_name]
	if !item_in_room {
		fmt.Println("You can't take the " + item.name)
	} else {
		if item_in_inventory {
			fmt.Println("You already have the " + item.name)
		} else {
			fmt.Println("You take the " + line[1])
			g.inventory[item.name] = item
			// remove the item from the current room
			delete(g.current_room.items, item.name)
		}
	}
}

func move(g *Game, line []string) {
	direction := line[0]
	next_room, exists := g.current_room.exits[direction]
	if exists {
		g.current_room = next_room
		look(g, nil)
	} else {
		fmt.Println("You can't go that way")
	}
}

// ground, sky, u
func link_rooms(room_a *Room, room_b *Room, dir string) {
	opposite := get_opposite(dir)
	if room_a.exits == nil {
		room_a.exits = make(map[string]*Room)
	}
	if room_b.exits == nil {
		room_b.exits = make(map[string]*Room)
	}
	room_a.exits[dir] = room_b
	room_b.exits[opposite] = room_a
}

//func room_links(links ...Link) map[string]Room {
//	roomlinks := make(map[string]Room)
//	for _, link := range links {
//		roomlinks[link.direction] = link.room
//	}
//	return roomlinks
//}

var hammer = Item{"hammer", "This hammer has seen better days"}

var verbs = map[string]func(g *Game, line []string){
	"look":    look,
	"examine": nil,
	"n":       move,
	"e":       move,
	"s":       move,
	"w":       move,
	"u":       move,
	"d":       move,
	"jump":    nil,
	"run":     nil,
	"take":    take}

var manhole = Room{
	"a manhole with a ladder leads into darkness",
	"Graffiti lines the wall around the manhole ladder.",
	nil,
	nil,
	nil,
	nil,
	0,
	nil}

var large_field = Room{
	"large field",
	"You are standing in a large field. \nRows of tall trees " +
		"block your view to the north, east, and south.",
	nil,
	nil,
	nil,
	nil,
	0,
	items(hammer)}

func items(is ...Item) map[string]Item {
	m := make(map[string]Item)
	for _, v := range is {
		m[v.name] = v
	}
	return m
}

func build_rooms() {
	link_rooms(&large_field, &manhole, "d")
}

func main() {
	fmt.Println("Welcome to Dork\n")
	bio := bufio.NewReader(os.Stdin)

	game := Game{&large_field, nil, 0}
	game.inventory = make(map[string]Item)
	build_rooms()
	look(&game, nil)
	for {
		fmt.Print(string(game.moves) + "> ")
		lineBytes, _, _ := bio.ReadLine()
		line := string(lineBytes)
		chunks := strings.Split(line, " ")
		if len(chunks) > 0 {
			f, found := verbs[chunks[0]]
			if found {
				if f != nil {
					f(&game, chunks)
				} else {
					fmt.Println("Nothing happens")
				}
				game.moves++
			} else {
				fmt.Println("Hm?")
			}
		}
	}

}
