package musicxml

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type key struct {
	start  int
	finish int
	frames []bool
	pitch  string
}

func LoadFromFile(file string, index int, pitch string) key {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	k := key{
		pitch: pitch,
	}

	r := bufio.NewReader(f)

	for {

		val, err := r.ReadString('\n')
		if err != nil {
			break
		}

		arr := strings.Split(val, " ")

		k.frames = append(k.frames, arr[index] == "1")
	}

	return k
}
