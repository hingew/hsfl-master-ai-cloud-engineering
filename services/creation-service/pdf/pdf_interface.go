package pdf

import "bytes"

type Pdf interface {
	AddElement()
	Out() (*bytes.Buffer, error)
}
