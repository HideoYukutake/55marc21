package marc21

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type dirent struct {
    tag           string
    length        int
    startCharPos  int
}

const (
    // Record Terminator
    RT = 0x1D
    // Record Separator
    RS = 0x1E
    // Subfield Delimiter
    DELIM = 0x1F
)

var ErrFieldSeparator = errors.New("Record Separator (field terminator)")

type Leader struct {
    Length                              int
    Status, Type                        byte
    ImplementationDefined               [5]byte
    CharacterEncoding                   byte
    BaseAddress                         int
    IndicattorCount, SubfieldCodeLength int
    LengthOfLength, LengthOfStartPos    int
}

func (leader Leader) Bytes() (buf []byte) {
    buf = make([]byte, 24)
    copy(buf[0:5], []byte(fmt.Sprintf("%05d", leader.Length)))
    buf[5] = leader.Status
    buf[6] = leader.Type
    copy(buf[7:9], leader.ImplementationDefined[0:2])
    buf[9] = leader.CharacterEncoding
    copy(buf[10:11], fmt.Sprintf("%d", leader.IndicattorCount))
    copy(buf[11:12], fmt.Sprintf("%d", leader.SubfieldCodeLength))
    copy(buf[12:17], fmt.Sprintf("%05d", leader.BaseAddress))
    copy(buf[17:20], leader.ImplementationDefined[2:5])
    copy(buf[20:21], fmt.Sprintf("%d", leader.LengthOfLength))
    copy(buf[21:22], fmt.Sprintf("%d", leader.LengthOfStartPos))
    buf[22] = '0'
    buf[23] = '0'
    return
}

func (leader Leader) String() string {
    return string(leader.Bytes())
}

func ParseLeader(r io.Reader) (leader *Leader, err error) {
    return readReader(r)
}

func readReader(reader io.Reader) (leader *Leader, err error) {
    data := make([]byte, 24)
    n, err := io.ReadFull(reader, data)
    if err != nil {
        return nil, err
    }
    if n < 23 {
        err = fmt.Errorf("invalid leader: expected 24 bytes, read %d", n)
        return
    }
    leader = &Leader{}
    leader.Length, err = strconv.Atoi(string(data[0:5]))
    if err != nil {
        err = fmt.Errorf("invalid record length: %s", err)
        return
    }
    leader.Status = data[5]
    leader.Type = data[6]
    copy(leader.ImplementationDefined[0:2], data[7:9])
    leader.CharacterEncoding = data[9]
    leader.IndicattorCount, err = strconv.Atoi(string(data[10:11]))
    if err != nil || leader.IndicattorCount != 2{
        err = fmt.Errorf("erronous indicator count, expected '2', got %v", err)
        return
    }

    leader.SubfieldCodeLength, err = strconv.Atoi(string(data[11:12]))
    if err != nil || leader.SubfieldCodeLength != 2{
        err = fmt.Errorf("erronous subfield code length, expected '2', got %v", err)
        return
    }

    leader.BaseAddress, err = strconv.Atoi(string(data[12:17]))
    if err != nil {
        err = fmt.Errorf("invalid base address: %s", err)
        return
    }
    copy(leader.ImplementationDefined[2:5], data[17:20])
    leader.LengthOfLength, err = strconv.Atoi(string(data[20:21]))
    if err != nil {
        return
    }
    leader.LengthOfStartPos, err = strconv.Atoi(string(data[21:22]))
    if err != nil {
        return
    }
    return
}

func readDirEnt(reader io.Reader) (dent *dirent, err error) {
    data := make([]byte, 12)
    if _, err = reader.Read(data[0:1]); err != nil {
        return
    }
    if data[0] == RS {
        err = ErrFieldSeparator
    }
    n, err := io.ReadFull(reader, data[1:])
    if err != nil {
        return
    }
    if n != 11 {
        err = fmt.Errorf("invalid directory entry, expected 12 bytes, got %d", n)
    }
    dent = &dirent{}
    dent.tag = string(data[0:3])
    if dent.length, err = strconv.Atoi(string(data[3:7])); err != nil {
        return
    }
    if dent.startCharPos, err = strconv.Atoi(string(data[7:12])); err != nil {
        return
    }
    return

}

