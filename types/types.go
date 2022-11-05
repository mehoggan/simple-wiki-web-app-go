package types

type Page struct {
	Title string
	Body  []byte
}

type Server struct {
	DocRoot string `yaml:"doc_root"`
}

type Config struct {
	Server Server `yaml:"server"`
}
