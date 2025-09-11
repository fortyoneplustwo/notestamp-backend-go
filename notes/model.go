package notes

const (
	ContentType = "application/json"
	Prefix      = "notes"
)

type Notes struct {
	Id   string `json:"id"`
	Data []byte `json:"content"`
}
