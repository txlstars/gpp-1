package gppio

// import "io"

type EmptyWriter struct{}

func (EmptyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
