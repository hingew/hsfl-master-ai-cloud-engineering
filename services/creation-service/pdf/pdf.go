package pdf

import (
	"bytes"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
)

type Pdf interface {
	Render(*model.PdfTemplate, map[string]string)
	Out() (*bytes.Buffer, error)
}
