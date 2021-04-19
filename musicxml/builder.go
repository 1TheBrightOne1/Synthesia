package musicxml

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

var PitchIndices = map[int]string{
	0: "A",
	1: "B",
	2: "C",
	3: "D",
	4: "E",
	5: "F",
	6: "G",
}

var NoteTypes = []NoteType{
	"whole",
	"half",
	"quarter",
	"eighth",
	"16th",
	"32nd",
	"64th",
}

type Builder struct {
	framesPerQuarter int
	divisions        int
	measures         []Measure
	beats            int
	beatType         int

	keyNotes []keyNotes
	durationToNoteType map[int]NoteType
}

type keyNotes struct {
	notes []Note
	index int
}

func NewBuilder(framesPerQuarter, divisions, beats, beatType int) *Builder {
	durationToNoteType := make(map[int]NoteType)

	//TODO: validate division is a power of 2 and we don't overflow
	index := 0
	for i := divisions * 4 /*whole note*/; i >= 1; i /= 2 {
		durationToNoteType[i] = NoteTypes[index]
		index++
	}

	return &Builder{
		framesPerQuarter: framesPerQuarter,
		divisions:        divisions,
		measures:         make([]Measure, 1),
		beats:            beats,
		beatType:         beatType,
		durationToNoteType: durationToNoteType,
	}
}

func (b *Builder) BuildXML(w io.Writer) {
	r, _ := os.Open("header.xml")
	header, _ := ioutil.ReadAll(r)
	w.Write(header)
	r.Close()

	b.keyNotes = make([]keyNotes, 52)

	c := make([]chan bool, 52)

	for i, _ := range c {
		c[i] = make(chan bool)
		go b.processNote(i, PitchIndices[i%len(PitchIndices)], &b.keyNotes[i], c[i])
	}

	for _, channel := range c {
		<-channel
	}

	const maxVal = 999999999
	//Find next division that plays a note and all notes that play on that division
	for {
		min := maxVal
		var values []int
		for i, n := range b.keyNotes {
			if n.index >= len(n.notes) {
				continue
			}

			if n.notes[n.index].StartsAtDivision < min {
				values = []int{i}
				min = n.notes[n.index].StartsAtDivision
			} else if n.notes[n.index].StartsAtDivision == min {
				//insert in order of duration
				inserted := false
				for j, index := range values {
					if n.notes[n.index].Duration > b.keyNotes[index].notes[b.keyNotes[index].index].Duration {
						valueCopy := make([]int, 0)
						valueCopy = append(valueCopy, values[:j]...)
						valueCopy = append(valueCopy, i)
						valueCopy = append(valueCopy, values[j:]...)
						values = valueCopy
						inserted = true
						break
					}
				}
				if !inserted {
					values = append(values, i)
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
	relativeNoteStart := noteStart - measureIndex*b.beats*b.divisions

	//TODO: Remove debug code here
	if noteStart == 41 {
		noteStart = 41
	}

	if measureIndex >= len(b.measures) {
		newMeasures := make([]Measure, measureIndex+1)
		copy(newMeasures, b.measures)
		b.measures = newMeasures
	}

	//Forward to the next starting note
	if b.measures[measureIndex].Counter < relativeNoteStart {
		b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, NewForward(relativeNoteStart-b.measures[measureIndex].Counter))
		b.measures[measureIndex].Counter = relativeNoteStart
	}

	isChord := false //Only the first note in each chord should be false
	voice := [...]int{1, 1}
	var previousNote *Note
	for _, currentIndex := range noteIndices {
		currentNote := b.keyNotes[currentIndex].notes[b.keyNotes[currentIndex].index]

		//Backup to the previous starting position is multiple durations are detected
		//Or if the measure counter has surpassed the next notes start
		if previousNote != nil && previousNote.Duration != currentNote.Duration || relativeNoteStart < b.measures[measureIndex].Counter {
			backupDuration := 0
			if relativeNoteStart < b.measures[measureIndex].Counter {
				backupDuration = b.measures[measureIndex].Counter - relativeNoteStart
			} else {
				backupDuration = previousNote.Duration
			}
			b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, NewBackup(backupDuration))
			isChord = false
			voice[currentNote.Staff-1]++
		}

		currentNote.Voice = voice[currentNote.Staff-1] + 1

		//Check if note extends into next measure
		if noteStart+currentNote.Duration > (measureIndex+1)*b.beats*b.divisions {
			newDuration := (measureIndex+1)*b.beats*b.divisions - noteStart
			remainingDuration := currentNote.Duration - newDuration
			newNote := currentNote

			currentNote.Duration = newDuration
			currentNote.Tie = currentNote.Tie + "<tie type=\"start\"/>"
			currentNote.Notations = append(currentNote.Notations, Tied{T: "start"})

			newNote.StartsAtDivision = (measureIndex + 1) * b.beats * b.divisions
			newNote.Duration = remainingDuration
			newNote.Tie = "<tie type=\"stop\"/>"
			newNote.Notations = make([]Tied, 0)
			newNote.Notations = append(newNote.Notations, Tied{T: "stop"})

			if b.keyNotes[currentIndex].index >= len(b.keyNotes[currentIndex].notes) {
				b.keyNotes[currentIndex].notes = append(b.keyNotes[currentIndex].notes, newNote)
			} else {
				newNotes := make([]Note, 0)
				newNotes = append(newNotes, b.keyNotes[currentIndex].notes[:b.keyNotes[currentIndex].index+1]...)
				newNotes = append(newNotes, newNote)
				newNotes = append(newNotes, b.keyNotes[currentIndex].notes[b.keyNotes[currentIndex].index+1:]...)
				b.keyNotes[currentIndex].notes = newNotes
			}
		}

		if isChord {
			b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, NewChordedNoteFromNote(currentNote))
		} else {
			b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, currentNote)
		}
		b.measures[measureIndex].Counter = (currentNote.StartsAtDivision - measureIndex*b.beats*b.divisions) + currentNote.Duration

		isChord = true
		previousNote = &currentNote

		b.keyNotes[currentIndex].index++
	}
}

func (b *Builder) writeMeasureToXML(w io.Writer, i int) error {
	measure := b.measures[i]
	measure.Number = strconv.Itoa(i)
	measure.XMLName = xml.Name{Local: "measure"}

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
	k := LoadFromFile("out3.txt", index, note)

	octave := int((index + 5) / 7) //TODO: dynamically pick the octave offset

	staff := 1
	if octave < 4 {
		staff = 2
	}

	framesPerDivision := b.framesPerQuarter / b.divisions
	for i, frame := range k.frames {
		if frame && (i == 0 || !k.frames[i-1]) {
			for j := i + 1; i < len(k.frames); j++ {
				if !k.frames[j] {
					snappedI := int(math.Round(float64(i)/float64(framesPerDivision))) * framesPerDivision
					snappedJ := int(math.Round(float64(j)/float64(framesPerDivision))) * framesPerDivision
					duration := (snappedJ - snappedI) / framesPerDivision
					if duration == 0 {
						duration = 1
					}

					n := Note{
						XMLName: xml.Name{Local: "note"},
						Pitch: &Pitch{
							Step:   note,
							Octave: octave,
						},
						Duration:         duration,
						NoteType:         b.durationToNoteType[duration],
						Staff:            staff,
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
		Staves: 2,
		Clefs: []interface{}{
			Clef{
				XMLName: xml.Name{Local: "clef"},
				Number: 1,
				Sign: "G",
				Line: 2,
		},
			Clef{
				XMLName: xml.Name{Local: "clef"},
				Number: 2,
				Sign: "F",
				Line: 4,
			},
		},
	}
}
