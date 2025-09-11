package cleanup

type Options struct {
	HasMedia bool
}

type Log struct {
	Uid string `json:"uid"`
	Id  string `json:"id"`
	Err string  `json:"err"`
}
