package server

type Server struct {
	Start func() error
}
