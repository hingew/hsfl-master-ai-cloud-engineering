package pdf

import (
	"bytes"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
)

type Pdf interface {
	Render(*model.PdfTemplate, map[string]interface{})
	Out() (*bytes.Buffer, error)
}
