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
	sixteenth    NoteType = "sixteenth"
	thirtysecond NoteType = "32nd"
)

var durationToNoteType = map[int]NoteType{
	16: whole,
	8:  half,
	4:  quarter,
	2:  eighth,
	1:  sixteenth,
}

type Pitch struct {
	Step   string `xml:"step"`
	Octave int    `xml:"octave"`
}

type Note struct {
	XMLName xml.Name
	Pitch    *Pitch    `xml:"pitch,omitempty"`
	Duration int      `xml:"duration,omitempty"`
	NoteType NoteType `xml:"type,omitempty"`
	Staff    int      `xml:"staff,omitempty"`
	Tie string `xml:",innerxml"`
	Tied string `xml:",innerxml"`
	StartsAtDivision int `xml:"-"`
}

type ChordedNote struct {
	Note
	Chord string `xml:"allowempty>doesnotmatter,omitempty"`
}

func NewChordedNoteFromNote(note Note) ChordedNote {
	return ChordedNote {
		Note {
			XMLName: note.XMLName,
			Pitch: note.Pitch,
			Duration: note.Duration,
			NoteType: note.NoteType,
			Staff: note.Staff,
			StartsAtDivision: note.StartsAtDivision,
			Tie: note.Tie,
			Tied: note.Tied,
		},
		"",
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
	Number     string     `xml:"number,attr"`
	Notes      []interface {}     //TODO: May have to delete notes elements
	Attributes *Attributes `xml:"attributes,omitempty"`
	Counter int `xml:"-"`
}

type Key struct {
	Fifths int `xml:"fifths"`
}

type Time struct {
	Beats    int `xml:"beats"`
	BeatType int `xml:"beat-type"`
}

type Attributes struct {
	Divisions int  `xml:"divisions,omitempty"`
	Key       *Key  `xml:"key,omitempty"`
	Time      *Time `xml:"time,omitempty"`
}
