package musicxml

import (
	"io"
	"math"
	"encoding/xml"
	"strconv"
	"fmt"
)

type Builder struct {
	framesPerQuarter int
	divisions int
	measures []Measure
	beats int
	beatType int
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
	/*notes := map[int]string{
		0: "A",
		1: "B",
		2: "C",
		3: "D",
		4: "E",
		5: "F",
		6: "G",
	}

	for i := 0; i < 52; i++ {
		b.processNote(i, notes[i % len(notes)])
	}*/
	b.processNote(25, "E")
	b.processNote(28, "A")

	b.measures[0].Attributes = NewAttributes(b.divisions, b.beats, b.beatType, 0, "major")

	for i, _ := range b.measures {
		b.writeMeasureToXML(w, i)
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
func (b *Builder) processNote(index int, note string) error {
	k := LoadFromFile("out1.txt", index, note)

	octave := int((index + 5) / 8)
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

					staff := 1
					if octave < 4 {
						staff = 2
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
					}

					b.addToMeasure(snappedI, n)

					i = j
					break
				}
			}
		}
	}

	return nil
}

func (b *Builder) addToMeasure(noteStart int, note Note) {
	framesPerDivision := b.framesPerQuarter / b.divisions

	measureIndex := noteStart / (b.framesPerQuarter * b.beats)

	noteStartRelative := (noteStart - b.framesPerQuarter * b.beats * measureIndex) / framesPerDivision

	if measureIndex >= len(b.measures) {
		newMeasures := make([]Measure, measureIndex + 1)
		copy(newMeasures, b.measures)
		b.measures = newMeasures
	}

	if b.measures[measureIndex].Counter < noteStartRelative {
		forward :=Forward {
				XMLName: xml.Name{Local: "forward"},
				Duration: noteStartRelative - b.measures[measureIndex].Counter,
			}
		b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, forward)
		b.measures[measureIndex].Counter = b.measures[measureIndex].Counter + (noteStartRelative - b.measures[measureIndex].Counter)
	} else if b.measures[measureIndex].Counter > noteStartRelative {
		backward := Backup{
				XMLName: xml.Name{Local: "backup"},
				Duration: b.measures[measureIndex].Counter - noteStartRelative,
			}
		b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, backward)
		b.measures[measureIndex].Counter = b.measures[measureIndex].Counter - (b.measures[measureIndex].Counter - noteStartRelative)
	}

	//Check if the note carries into next measure
	measureDuration := (b.beats * b.framesPerQuarter / framesPerDivision - b.measures[measureIndex].Counter)
	if measureDuration - note.Duration < 0 {
		remainingDuration := note.Duration - measureDuration
		newNote := note
		newNote.Duration = remainingDuration

		b.addToMeasure(noteStart + b.beats * b.framesPerQuarter, newNote)
		note.Duration = measureDuration
	}

	b.measures[measureIndex].Notes = append(b.measures[measureIndex].Notes, note)
	b.measures[measureIndex].Counter = b.measures[measureIndex].Counter + note.Duration
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
