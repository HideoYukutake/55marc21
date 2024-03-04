package marc21

import (
	"io"
	"strings"
)

type Record struct {
    Leader *Leader
    Fields []Field
}

func  ReadRecord(reader io.Reader) (record *Record, err error)  {
    record := &Record{}
    record.Fields = make([]Field, 0, 8)
    if record.Leader, err = readLeader(reader); err != nil {
        return
    }
    return
}

func (record *Record) AddField(f Field)  {
    record.Fields = append(record.Fields, f)
}

func (record *Record) String() string {
    estrings := make([]string, len(record.Fields))
    for i, entry := range record.Fields {
        estrings[i] = entry.String()
    }
    return strings.Join(estrings, "\n")
}

func (record *Record) GetField(tag string)  (fields []Field){
    fields = make([]Field, 0, 4)
    return
}

func (record *Record) GetSubField(tag string, code byte)  (subfields []SubField){
    subfields = make([]SubField, 0, 4)
    return
}
