package marc21

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Field interface {
    String() string
    GetTag() string
}

type ControlField struct {
    Tag     string
    Data    string
}

func (cf *ControlField) String() string {
    return fmt.Sprintf("%s %s", cf.Tag, cf.Data)
}

func (cf *ControlField) GetTag() string {
    return cf.Tag
}

func readControl(reader io.Reader, dent *dirent) (field Field, err error) {
    data := make([]byte, dent.length)
    n, err := io.ReadFull(reader, data)
    if err != nil {
        return
    }
    if n != dent.length {
        err = fmt.Errorf("invalid control entry, expected %d bytes, read %d", dent.length, n)
        return
    }
    if data[dent.length - 1] != RS {
        err = fmt.Errorf("invalid control entry, does not end with a field terminator")
        return
    }
    field = &ControlField{Tag: dent.tag, Data: string(data[:dent.length-1])}
    return
}

type SubField struct {
    Code  byte
    Value string
}

func (sf SubField) String() string {
    return fmt.Sprintf("%s %s", sf.Code, sf.Value)
}

type DataField struct {
    Tag         string
    Indicator1  byte
    Indicator2  byte
    SubFields   []*SubField
}

func (df *DataField) GetTag() string {
    return df.Tag
}


func (df *DataField) String() string {
    subfields := make([]string, 0, len(df.SubFields))
    for _, sf := range df.SubFields {
        subfields = append(subfields, "["+sf.String()+"]")
    }
    return fmt.Sprintf("%s [%c%c] %s", df.Tag, df.Indicator1, df.Indicator2, strings.Join(subfields, ", "))
}

func readData(reader io.Reader, dent *dirent) (field Field, err error) {
    data := make([]byte, dent.length)
    n, err := io.ReadFull(reader, data)
    if err != nil {
        return
    }
    if n != dent.length {
        err = fmt.Errorf("invalid data entry, expected %d bytes, read %d", dent.length, n)
        return
    }
    if data[dent.length - 1] != RS {
        err = fmt.Errorf("invalid data entry, does not end with a field terminator")
        return
    }
    df := &DataField{Tag: dent.tag}
    df.Indicator1, df.Indicator2 = data[0], data[1]
    df.SubFields = make([]*SubField, 0, 1)
    for _, sfBytes := range bytes.Split(data[2:dent.length-1], []byte{DELIM}) {
        if len(sfBytes) == 0 {
            continue
        }
        sf := &SubField{Code: sfBytes[0], Value: string(sfBytes[1:])}
        df.SubFields = append(df.SubFields, sf)
    }
    field = df
    return
}
