package musicxml

import (
	"encoding/xml"
)

type NoteType string

const (
	whole        NoteType = "whole"
	half         NoteType = "half"
	quarter      NoteType = "quarter"
	eighth       NoteType = "eighth"
	sixteenth    NoteType = "16th"
	thirtysecond NoteType = "32nd"
)

type Pitch struct {
	Step   string `xml:"step"`
	Alter int `xml:"alter,omitempty"`
	Octave int    `xml:"octave"`
}

type Tied struct {
	T string `xml:"type,attr"`
}

type Note struct {
	XMLName xml.Name
	Pitch    *Pitch    `xml:"pitch,omitempty"`
	Duration int      `xml:"duration,omitempty"`
	Tie string `xml:",innerxml"`
	Voice int `xml:"voice,omitempty"`
	NoteType NoteType `xml:"type,omitempty"`
	Staff    int      `xml:"staff,omitempty"`
	Notations []Tied `xml:"notations>tied,omitempty"`
	StartsAtDivision int `xml:"-"`
}

type ChordedNote struct {
	XMLName xml.Name
	Chord string `xml:"chord,allowempty"`
	Pitch    *Pitch    `xml:"pitch,omitempty"`
	Duration int      `xml:"duration,omitempty"`
	Tie string `xml:",innerxml"`
	Voice int `xml:"voice,omitempty"`
	NoteType NoteType `xml:"type,omitempty"`
	Staff    int      `xml:"staff,omitempty"`
	Notations []Tied `xml:"notations>tied,omitempty"`
	StartsAtDivision int `xml:"-"`
}

func NewChordedNoteFromNote(note Note) ChordedNote {
	return ChordedNote {
		XMLName: xml.Name{Local:"note"},
		Pitch: note.Pitch,
		Duration: note.Duration,
		NoteType: note.NoteType,
		Staff: note.Staff,
		StartsAtDivision: note.StartsAtDivision,
		Tie: note.Tie,
		Voice: note.Voice,
		Notations: note.Notations,
		Chord: "",
	}
}

type Backup struct {
	XMLName xml.Name
	Duration int `xml:"duration,omitempty"`
}

func NewBackup(duration int) Backup {
	return Backup{
		XMLName: xml.Name{Local: "backup"},
		Duration: duration,
	}
}

type Forward struct {
	XMLName xml.Name
	Duration int `xml:"duration,omitempty"`
}

func NewForward(duration int) Forward {
	return Forward{
		XMLName: xml.Name{Local: "forward"},
		Duration: duration,
	}
}

type Measure struct {
	XMLName xml.Name
	Attributes *Attributes `xml:"attributes,omitempty"`
	Number     string     `xml:"number,attr"`
	Notes      []interface {}
	Counter int `xml:"-"`
}

type Key struct {
	Fifths int `xml:"fifths"`
}

type Time struct {
	Beats    int `xml:"beats"`
	BeatType int `xml:"beat-type"`
}

type Clef struct {
	XMLName xml.Name
	Number int `xml:"number,attr"`
	Sign string `xml:"sign"`
	Line int `xml:"line"`
}

type Attributes struct {
	Divisions int  `xml:"divisions,omitempty"`
	Key       *Key  `xml:"key,omitempty"`
	Time      *Time `xml:"time,omitempty"`
	Staves int `xml:"staves,omitempty"`
	Clefs []interface{}
}
