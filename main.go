package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func ErrorMessage(s string) {
	log.Fatal("ERROR: invalid data format", s)
}

func PrintResults(ls []string) {
	for _, l := range ls {
		fmt.Println(l)
	}

	fmt.Printf("\n")
}

func GetFileLines(p string) []string {
	file, err := os.Open(p)
	if err != nil {
		ErrorMessage(", invalid file")
	}

	// Use bufio to scan lines and store them in slice lines
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	file.Close()

	return lines
}

// Error handling
func FormatCheck(rooms []string, coords []Vector2) {
	if len(rooms) < 2 {
		ErrorMessage(", not enough rooms")
	}
	var roomsCopy []string
	for _, r := range rooms {
		roomsCopy = append(roomsCopy, r)
	}
	slices.Sort(roomsCopy)
	roomsCopy = slices.Compact(roomsCopy)
	if len(rooms) > len(roomsCopy) {
		ErrorMessage(", duplicate rooms")
	}
	for i, xy := range coords {
		if xy.x > 1000 || xy.y > 1000 {
			ErrorMessage(", invalid room co-ordinates")
		}

		for j, xy2 := range coords {
			if i == j {
				continue
			}

			if xy.x == xy2.x && xy.y == xy2.y {
				ErrorMessage(", rooms share same co-ordinates")
			}
		}
	}
}

func countLine(s string) int {
	count := 0

	for _, c := range s {
		if c == '\n' {
			count++
		}
	}

	return count
}

// Split the file into rooms, co-ordinates and tunnels
func GetData(lines []string) ([]string, []Vector2, int) {
	var rooms []string
	var end string
	var connections int
	var coords []Vector2
	var endCoords Vector2

	// Regex for how rooms, tunnels and commands should be formatted
	roomFormat, _ := regexp.Compile(`([\S]+) ([\d]+) ([\d]+)`)
	connectFormat, _ := regexp.Compile(`([\S]+)-([\S]+)`)
	commandFormat, _ := regexp.Compile(`#([\S]+)`)

	// Start by obtaining the starting room
	for i := 1; i < len(lines); i++ {
		if lines[i] == "##start" {
			if !roomFormat.MatchString(lines[i+1]) || string(lines[i+1][0]) == "L" || string(lines[i+1][0]) == "#" {
				ErrorMessage(", invalid start room - " + lines[i+1])
			}

			room := strings.Split(lines[i+1], " ")
			x, _ := strconv.Atoi(room[1])
			y, _ := strconv.Atoi(room[2])

			rooms = append(rooms, room[0])
			coords = append(coords, Vector2{x, y})
			break
		}
	}

	// Find every other room
	for i := 1; i < len(lines); i++ {
		// Save the index for where tunnels start for later
		if connectFormat.MatchString(lines[i]) {
			if i > 4 {
				connections = i
			} else {
				ErrorMessage(", not enough rooms")
			}
			break
		}

		// Ignoring commands, only looking for rooms
		if lines[i-1] == "##start" || commandFormat.MatchString(lines[i]) {
			continue
		}

		if !roomFormat.MatchString(lines[i]) || string(lines[i][0]) == "L" || string(lines[i][0]) == "#" {
			ErrorMessage(", invalid room - \"" + lines[i] + "\"")
		}

		// For valid rooms, split into room name and co-ordinates
		room := strings.Split(lines[i], " ")
		x, _ := strconv.Atoi(room[1])
		y, _ := strconv.Atoi(room[2])

		if lines[i-1] == "##end" {
			end = room[0]
			endCoords.x = x
			endCoords.y = y
			continue
		}

		rooms = append(rooms, strings.Split(lines[i], " ")[0])
		coords = append(coords, Vector2{x, y})
	}

	// Append end room after every other room
	if end != "" {
		rooms = append(rooms, end)
		coords = append(coords, endCoords)
	}

	return rooms, coords, connections
}

// Add rooms as keys to a map, with a slice of rooms they're connected to as the value
func MapConnections(lines []string, connections int) map[string][]string {
	connects := make(map[string][]string)
	connectFormat, _ := regexp.Compile(`([\S]+)-([\S]+)`)

	for i := connections; i < len(lines); i++ {
		// Error handling
		if string(lines[i][0]) == "#" {
			continue
		}
		connect := strings.Split(lines[i], "-")
		if len(connect) != 2 || !connectFormat.MatchString(lines[i]) {
			ErrorMessage(", tunnels formatted incorrectly")
		}
		if connect[0] == connect[1] {
			ErrorMessage(", room has path to itself")
		}

		// Checking if rooms are already in the map and appending rooms they're connected to
		if _, ok := connects[connect[0]]; ok {
			connects[connect[0]] = append(connects[connect[0]], connect[1])
		} else {
			connects[connect[0]] = append(connects[connect[0]], connect[1])
		}

		if _, ok := connects[connect[1]]; ok {
			connects[connect[1]] = append(connects[connect[1]], connect[0])
		} else {
			connects[connect[1]] = append(connects[connect[1]], connect[0])
		}
	}

	return connects
}

// Search for all possible paths from start to end
func SearchPaths(start string, end string, m map[string][]string, path []string, paths *string) {
	BackTrack := false

	if len(path) == 0 {
		path = append(path, start)
	}

	// For all the rooms linked to the current room, add to current path and do so until reaching the end
	for _, p := range m[start] {
		for i := 0; i < len(path); i++ {
			if p == path[i] {
				BackTrack = true
			}
		}

		// If room has already been added to the path, ignore
		if BackTrack {
			BackTrack = false
			continue
		}

		path = append(path, p)
		if path[len(path)-1] == end {
			for _, v := range path {
				*paths = *paths + v + " "
			}
		}

		if len(m[p]) > 1 {
			SearchPaths(p, end, m, path, paths)
			path = path[0 : len(path)-1]
		}
	}
}

// Find the max number of nodes connected to the end room that also have a path to the start room
func FindMaxPaths(paths [][]string) int {
	var exits []string
	var zeroRoomPath int

	// Accounting for if the shortest path is direct from start to end
	if len(paths[0]) > 2 {
		exits = append(exits, paths[0][len(paths[0])-2])
	} else {
		zeroRoomPath = 1
	}

	for i := 1; i < len(paths); i++ {
		if !slices.Contains(exits, paths[i][len(paths[i])-2]) {
			exits = append(exits, paths[i][len(paths[i])-2])
		}
	}

	return len(exits) + zeroRoomPath
}

// Finding paths with unique rooms from each other
func OptimisePaths(MaxPaths int, paths [][]string, temppaths [][]string, allPaths *[][][]string, found *bool) {
	if len(temppaths) == MaxPaths {
		*allPaths = append(*allPaths, temppaths)
		*found = true
		return
	}

	if len(temppaths) == 0 {
		temppaths = append(temppaths, paths[0])
	}

	var dupecheck []string
	var dupechecker []string
	for _, p := range temppaths {
		for _, n := range p[1 : len(p)-1] {
			dupecheck = append(dupecheck, n)
		}
	}

	for i := 1; i < len(paths); i++ {
		if *found {
			return
		}

		redundant := false
		for _, p := range temppaths {
			if p[1] == paths[i][1] {
				redundant = true
			}
		}

		if redundant {
			continue
		}

		for _, n := range paths[i][1 : len(paths[i])-1] {
			dupecheck = append(dupecheck, n)
		}

		for _, n := range dupecheck {
			dupechecker = append(dupechecker, n)
		}

		slices.Sort(dupechecker)
		dupechecker = slices.Compact(dupechecker)

		if len(dupecheck) == len(dupechecker) {
			temppaths = append(temppaths, paths[i])
			OptimisePaths(MaxPaths, paths, temppaths, allPaths, found)

			if *found {
				return
			}

			temppaths = temppaths[:len(temppaths)-1]
		}

		dupecheck = nil
		for _, p := range temppaths {
			for _, n := range p[1 : len(p)-1] {
				dupecheck = append(dupecheck, n)
			}
		}
		dupechecker = nil
	}
}

// Convert string returned from SearchPaths into a slice of paths
func GetPaths(ps string, end string) [][]string {
	var paths [][]string
	var path []string

	prepaths := strings.Split(ps, " ")

	for _, v := range prepaths {
		path = append(path, v)
		if v == end {
			paths = append(paths, path)
			path = nil
		}
	}

	return paths
}

type Ant struct {
	Id       int
	Path     []string
	Pos      int
	LeftHome bool
}

type Vector2 struct {
	x int
	y int
}

func LemIn(ants []Ant, rooms []string, paths [][]string, Occupied map[string]bool, result *string) bool {
	var steps string
	var antsRemaining int
	Complete := true
	startToEnd := false

	// Check how many ants are still at home
	for _, ant := range ants {
		if !ant.LeftHome {
			antsRemaining++
		}
	}

	// For all ants, move along/change paths so long as they haven't reached the end
	for i := range ants {
		if ants[i].Pos != len(ants[i].Path)-1 {
			newPath := false

			// If the next room in the path is occupied, move to the next free path available, else wait
			if Occupied[ants[i].Path[ants[i].Pos+1]] || startToEnd {
				for _, p := range paths {
					if ants[i].Path[1] != p[1] && !Occupied[p[ants[i].Pos+1]] {
						// Don't go to the next path if it's quicker to wait for current path to be free
						if antsRemaining <= len(paths) && len(p)-len(ants[i].Path) > 1 {
							break
						}

						ants[i].Path = p
						newPath = true
						break
					}

				}

				if !newPath {
					continue
				}
			}

			// Moving ants forward
			ants[i].Pos = ants[i].Pos + 1
			Occupied[ants[i].Path[ants[i].Pos-1]] = false
			ants[i].LeftHome = true
			steps = steps + "L" + strconv.Itoa(ants[i].Id) + "-" + ants[i].Path[ants[i].Pos] + " "

			// If ant hasn't reached end, set their new room to occupied, else check if 2 room path
			if ants[i].Pos == len(ants[i].Path)-1 {
				if len(ants[i].Path) == 2 {
					startToEnd = true
				}
			} else {
				Occupied[ants[i].Path[ants[i].Pos]] = true
			}
		}
	}

	// Turn has finished, append to final result
	steps = steps[:len(steps)-1]
	*result = *result + steps + "\n"

	// If all ants have reached the end, return, else run another Lem-in step
	for _, ant := range ants {
		if ant.Pos != len(ant.Path)-1 {
			Complete = false
		}
	}

	if Complete {
		return Complete
	} else {
		return LemIn(ants, rooms, paths, Occupied, result)
	}
}

func main() {
	if len(os.Args) != 2 {
		ErrorMessage(", argument should be 1 file")
	}

	// Open the file
	lines := GetFileLines(os.Args[1])
	if lines == nil {
		ErrorMessage(", input file is empty")
	}

	for _, l := range lines {
		fmt.Println(l)
	}
	fmt.Println()

	// Error handling
	numberOfAnts, err := strconv.Atoi(lines[0])
	if err != nil || numberOfAnts < 1 || numberOfAnts > 9223372036854775807 {
		ErrorMessage(", invalid number of Ants")
	}

	if !slices.Contains(lines, "##end") || !slices.Contains(lines, "##start") {
		ErrorMessage(", no start/end room found")
	}

	var linesCopy []string
	for _, l := range lines {
		linesCopy = append(linesCopy, l)
	}
	slices.Sort(linesCopy)
	linesCopy = slices.Compact(linesCopy)
	if len(linesCopy) < len(lines) {
		ErrorMessage(", duplicate data found in file")
	}

	// Extracting data from file to separate rooms, room co-ordinates and tunnels
	rooms, coords, connections := GetData(lines)
	FormatCheck(rooms, coords)

	// Find room connections
	connects := MapConnections(lines, connections)
	for _, c := range connects {
		for _, v := range c {
			if !slices.Contains(rooms, v) {
				ErrorMessage(", unknown room in tunnels")

			}
		}
	}

	var allPaths [][][]string
	var tempPaths [][]string
	var pathBuild []string
	var pathsBuild string

	SearchPaths(rooms[0], rooms[len(rooms)-1], connects, pathBuild, &pathsBuild)

	if pathsBuild == "" {
		ErrorMessage(", no viable paths")
	}

	paths := GetPaths(pathsBuild, rooms[len(rooms)-1])

	// Sort paths from shortest to longest
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// Finding the optimal paths (unique rooms) for 1 path up to if there are 'maxpaths' paths
	maxpaths := FindMaxPaths(paths)
	for i := maxpaths; i > 0; i-- {
		found := false
		if i == 1 {
			allPaths = append(allPaths, [][]string{paths[0]})
			break
		}

		if len(paths) == i {
			allPaths = append(allPaths, paths)
			continue
		}

		// From the shortest path, try to find optimal paths, if not possible start again from next shortest
		for j := 0; j < len(paths)-i+1; j++ {
			OptimisePaths(i, paths[j:], tempPaths, &allPaths, &found)
			if found {
				break
			}
		}
	}

	// Map to determine if a room is currently occupied by an ant
	roomChecker := make(map[string]bool)
	for _, v := range rooms {
		roomChecker[v] = false
	}

	var ants []Ant
	var allSteps []string
	var steps string

	for i := 1; i < numberOfAnts+1; i++ {
		ants = append(ants, Ant{i, nil, 0, false})
	}

	for _, paths := range allPaths {
		for i := 0; i < len(ants); i++ {
			ants[i].Path = paths[0]
		}
		LemIn(ants, rooms[0:len(rooms)-1], paths, roomChecker, &steps)
		allSteps = append(allSteps, steps)
		steps = ""
		for i := 0; i < numberOfAnts; i++ {
			ants[i].Path, ants[i].Pos, ants[i].LeftHome = nil, 0, false
		}
	}

	var fastest int

	for i := 0; i < len(allSteps); i++ {
		// fmt.Printf("%v\n\n", allSteps[i])

		if countLine(allSteps[i]) < countLine(allSteps[fastest]) {
			fastest = i
			continue
		}

		if countLine(allSteps[i]) == countLine(allSteps[fastest]) {
			if len(strings.Split(allSteps[i], " ")) < len(strings.Split(allSteps[fastest], " ")) {
				fastest = i
			}
		}
	}

	fmt.Println(allSteps[fastest])
}
