package gofast

import (
	"io"
)

//Decode method read FAST message from in and write decoded message to out.
func Decode(in io.Reader, out io.Writer) error {
	var msg string
	var err error
	for msg, err = DecodeMessage(in); err != nil; {
		out.Write([]byte(msg))
	}

	if err != io.EOF {
		return err
	}
	return nil
}

//DecodeMessage method decode one FAST message from input in.
func DecodeMessage(in io.Reader) (string, error) {
	_, err := DecodePMap(in)
	if err != nil {
		return "", err
	}
	return "", nil
}

//DecodePMap decode the PMap from input.
func DecodePMap(in io.Reader) ([]byte, error) {
	return nil, nil
}

//ReadBinaryByStopbit read binary data from input until stop bit of the byte is set.
func ReadBinaryByStopbit(in io.Reader) []byte {
	return nil
}
