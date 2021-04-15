package musicxml

import (
	"io"
	"math"
	"encoding/xml"
	"strconv"
	"fmt"
	"os"
	"io/ioutil"
)

type Builder struct {
	framesPerQuarter int
	divisions int
	measures []Measure
	beats int
	beatType int

	keyNotes []keyNotes
}

type keyNotes struct {
	notes []Note
	index int
}

func NewBuilder(framesPerQuarter, divisions, beats, beatType int) *Builder {
	return &Builder{
		framesPerQuarter: framesPerQuarter,
		divisions: divisions,
		measures: make([]Measure, 1),
		beats: beats,
		beatType: beatType,
	}
}

func (b *Builder) BuildXML(w io.Writer) {
	r, _ := os.Open("header.xml")
	header, _ := ioutil.ReadAll(r)
	w.Write(header)
	r.Close()

	b.keyNotes = make([]keyNotes, 52)

	c := make([]chan bool, 2)
	c[0] = make(chan bool)
	c[1] = make(chan bool)

	go b.processNote(25, "E", &b.keyNotes[25], c[0])
	go b.processNote(28, "A", &b.keyNotes[28], c[1])

	<- c[0]
	<- c[1]

	const maxVal = 999999999
	//Find next division that plays a note and all notes that play on that division
	for {
		min := maxVal
		var values []int
		for i, n := range b.keyNotes{
			if n.index >= len(n.notes) {
				continue
			}

			if n.notes[n.index].StartsAtDivision < min {
				values = []int{i}
				min = n.notes[n.index].StartsAtDivision
			} else if n.notes[n.index].StartsAtDivision == min {
				//insert in order of duration
				for j, index := range values {
					if n.notes[n.index].Duration > b.keyNotes[index].notes[b.keyNotes[index].index].Duration {
						valueCopy := make([]int, 0)
						valueCopy = append(valueCopy, values[:j]...)
						valueCopy = append(valueCopy, i)
						valueCopy = append(valueCopy, values[j:]...)
						values = valueCopy
						break
					}
				}
			}
		}
		if min == maxVal {
			break
		} else {
			b.addToMeasure(values)
		}
	}

	b.measures[0].Attributes = NewAttributes(b.divisions, b.beats, b.beatType, 0, "major")

	for i, _ := range b.measures {
		b.writeMeasureToXML(w, i)
	}

	r, _ = os.Open("footer.xml")
	footer, _ := ioutil.ReadAll(r)
	w.Write(footer)
	r.Close()
}

//addToMeasure expects a slice of indices for notes that start at the same time with their durations sorted in decreasing order
func (b *Builder) addToMeasure(noteIndices []int) {
	if len(noteIndices) == 0 {
		return
	}

	noteStart := b.keyNotes[noteIndices[0]].notes[b.keyNotes[noteIndices[0]].index].StartsAtDivision
	measureIndex := noteStart / (b.beats * b.divisions)
	relativeNoteStart := noteStart - measureIndex * b.beats * b.divisions

	if measureIndex >= len(b.measures) {
		newMeasures := make([]Measure, measureIndex + 1)
		copy(newMeasures, b.measures)
		b.measures = newMeasures
	}

	//Forward to the next starting note
	if b.measures[measureIndex].Counter < relativeNoteStart {
		b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, NewForward(relativeNoteStart))
		b.measures[measureIndex].Counter += relativeNoteStart - b.measures[measureIndex].Counter
	}

	isChord := false //Only the first note in each chord should be false
	var previousNote *Note
	for _, currentIndex := range noteIndices {
		currentNote := b.keyNotes[currentIndex].notes[b.keyNotes[currentIndex].index]

		//Backup to the previous starting position is multiple durations are detected
		if previousNote != nil && previousNote.Duration != currentNote.Duration {
			b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, NewBackup(previousNote.Duration))
			isChord = false
		}

		//Check if note extends into next measure
		if noteStart + currentNote.Duration > (measureIndex + 1) * b.beats * b.divisions {
			newDuration := (measureIndex + 1) * b.beats * b.divisions - noteStart
			remainingDuration := currentNote.Duration - newDuration
			currentNote.Duration = newDuration
			currentNote.Tie = currentNote.Tie + "<tie>start</tie>"
			currentNote.Tied = currentNote.Tied + "<tied>start</tied>"

			newNote := currentNote
			newNote.StartsAtDivision = (measureIndex + 1) * b.beats * b.divisions
			newNote.Duration = remainingDuration
			newNote.Tie = "<tie>stop</tie>"
			newNote.Tied = "<tied>stop</tied>"

			if b.keyNotes[currentIndex].index >= len(b.keyNotes[currentIndex].notes) {
				b.keyNotes[currentIndex].notes = append(b.keyNotes[currentIndex].notes, newNote)
			} else {
				newNotes := make([]Note, 0)
				newNotes = append(newNotes, b.keyNotes[currentIndex].notes[:b.keyNotes[currentIndex].index + 1]...)
				newNotes = append(newNotes, newNote)
				newNotes = append(newNotes, b.keyNotes[currentIndex].notes[b.keyNotes[currentIndex].index + 1:]...)
				b.keyNotes[currentIndex].notes = newNotes
			}
		}

		if isChord {
			b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, NewChordedNoteFromNote(currentNote))
		} else {
			b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, currentNote)
		}
		b.measures[measureIndex].Counter = (currentNote.StartsAtDivision - measureIndex * b.beats * b.divisions) + currentNote.Duration

		isChord = true
		previousNote = &currentNote

		b.keyNotes[currentIndex].index++
	}
}

func (b *Builder) writeMeasureToXML(w io.Writer, i int) error {
	measure := b.measures[i]
	measure.Number = strconv.Itoa(i + 1)
	measure.XMLName = xml.Name{ Local: "measure" }

	expectedLength := b.divisions * b.beats
	actualLength := 0
	for _, note := range measure.Notes {
		switch n := note.(type) {
		case Forward:
			actualLength += n.Duration
		case Backup:
			actualLength -= n.Duration
		case Note:
			actualLength += n.Duration
		}
	}

	if expectedLength > actualLength {
		remainder := expectedLength - actualLength
		measure.Notes = append(measure.Notes, NewForward(remainder))
	}
	
	bytes, err := xml.Marshal(measure)
	if err != nil {
		fmt.Println(i)
		fmt.Println(err)
		return err
	}
	w.Write(bytes)
	return nil
}

//processNote goes through all frames and writes them to the measures
func (b *Builder) processNote(index int, note string, keyNote *keyNotes, c chan bool) {
	k := LoadFromFile("out1.txt", index, note)

	octave := int((index + 5) / 8)
	
	staff := 1
	if octave < 4 {
		staff = 2
	}

	framesPerDivision := b.framesPerQuarter / b.divisions
	for i, frame := range k.frames {
		if frame && (i == 0 || !k.frames[i-1]) {
			for j := i + 1; i < len(k.frames); j++ {
				if !k.frames[j] {
					snappedI := int(math.Round(float64(i) / float64(framesPerDivision))) * framesPerDivision
					snappedJ := int(math.Round(float64(j) / float64(framesPerDivision))) * framesPerDivision
					duration := (snappedJ - snappedI) / framesPerDivision
					if duration == 0 {
						duration = 1
					}

					n := Note{
						XMLName: xml.Name{Local: "note"},
						Pitch: &Pitch{
							Step:   note,
							Octave: 4,
						},
						Duration: duration,
						NoteType: durationToNoteType[duration],
						Staff:    staff,
						StartsAtDivision: snappedI / framesPerDivision,
					}

					keyNote.notes = append(keyNote.notes, n)

					i = j
					break
				}
			}
		}
	}

	c <- true
}

func NewAttributes(divisions, beats, beatType, key int, mode string) *Attributes {
	return &Attributes{
		Divisions: divisions,
		Key: &Key{
			Fifths: key,
		},
		Time: &Time{
			Beats:    beats,
			BeatType: beatType,
		},
	}
}
