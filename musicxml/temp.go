package musicxml

import (
	"bufio"
	"os"
	"strings"
)

type key struct {
	start  int
	finish int
	Pixels []bool
	pitch  string
}

func LoadFromFile(file string, index int, pitch string) key {
	f, _ := os.Open(file)

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

		k.Pixels = append(k.Pixels, arr[index] == "1")
	}

	return k
}
