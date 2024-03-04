package marc21

type Record struct {
    Leader *Leader `xml:"leader"`
    Fields []Field
}
