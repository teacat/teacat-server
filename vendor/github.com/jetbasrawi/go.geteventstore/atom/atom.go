package atom

import (
	"encoding/xml"
	"strings"
	"time"
)

// Feed represents an atom feed page from the eventstore.
type Feed struct {
	XMLName      xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title        string   `xml:"title"`
	ID           string   `xml:"id"`
	StreamID     string   `xml:"streamId"`
	HeadOfStream bool     `xml:"headOfStream"`
	Link         []Link   `xml:"link"`
	Updated      TimeStr  `xml:"updated"`
	Author       *Person  `xml:"author"`
	Entry        []*Entry `xml:"entry"`
}

// GetLink gets the link with the name specified by the link argument.
func (f *Feed) GetLink(name string) *Link {
	if f == nil {
		return nil
	}

	for _, v := range f.Link {
		if v.Rel == name {
			return &v
		}
	}
	return nil
}

// PrettyPrint returns an indented string representation of the feed.
func (f *Feed) PrettyPrint() string {
	b, err := xml.MarshalIndent(f, "", "	")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// GetEventURLs extracts a slice of event urls from the feed object.
func (f *Feed) GetEventURLs() ([]string, error) {
	s := make([]string, len(f.Entry))
	for i := 0; i < len(f.Entry); i++ {
		e := f.Entry[i]
		s[i] = strings.TrimRight(e.Link[1].Href, "/")
	}
	return s, nil
}

// Entry represents a feed entry.
type Entry struct {
	Title     string  `xml:"title"`
	ID        string  `xml:"id"`
	Link      []Link  `xml:"link"`
	Published TimeStr `xml:"published"`
	Updated   TimeStr `xml:"updated"`
	Author    *Person `xml:"author"`
	Summary   *Text   `xml:"summary"`
	Content   *Text   `xml:"content"`
}

// Link represents a Link entry in the feed.
type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

// Person represents a person
type Person struct {
	Name string `xml:"name"`
}

// Text represents a text entry
type Text struct {
	Type string `xml:"type,attr,omitempty"`
	Body string `xml:",chardata"`
}

// TimeStr is a formatted time string
type TimeStr string

// Time returns a TimeStr
func Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}
