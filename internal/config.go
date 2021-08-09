package internal

const (
	ClientName  = "spa"
	ProductName = "Secure password assistant"
)

var (
	ClientVersion = "0.0.0"
	ClientOS      = "unknown"
)

type Configuration struct {
	BaseDir string
	// Server  ServerConfig
	Domain string
	User   string
	Secret SecretDigest
	Salt   SaltDigest

	Pin      int
	Pwd      PwdConfig
	Question QuestionConfig
	Count    int
}

type SaltDigest struct {
	Raw  string
	Hash []byte
}

type SecretDigest struct {
	Raw    string
	NewRaw string

	// secret is split into two parts
	Stock []byte
	Foil  []byte
}

type PwdConfig struct {
	Pepper string `json:"pepper"`
	Mask   string `json:"mask"`
	Length int    `json:"length"`
	Note   string `json:"note"`
}

type QuestionConfig struct {
	Question string `json:"-"`
}

type ServerConfig struct {
	Port int
}
