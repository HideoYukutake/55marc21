package marc21

type Field interface {
    String() string
    GetTag() string
}

type ControlField struct {
    Tag     string
    Data    string
}

type SubField struct {
    Code  byte
    Value byte
}

type DataField struct {
    Tag         string
    Indicator1  byte
    Indicator2  byte
    SubFields   []*SubField
}
