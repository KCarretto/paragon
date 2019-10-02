package mocks

//go:generate mockgen -destination=io.gen.go -package=mocks io Writer,WriteCloser
//go:generate mockgen -destination=codec.gen.go -package=mocks github.com/kcarretto/paragon/transport Encoder,Decoder
