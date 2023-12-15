package client

import (
	"context"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/proto"
	"google.golang.org/grpc"
	"time"
)

const DATE_FORMAT = "2006-01-02 15:04:05.999999999 -0700 MST"

type TemplatingGrpcClient struct {
	proto.TemplateServiceClient
}

func NewGrpcClient(conn grpc.ClientConnInterface) TemplatingGrpcClient {
	client := proto.NewTemplateServiceClient(conn)
	return TemplatingGrpcClient{client}
}

func (c *TemplatingGrpcClient) FetchTemplate(templateID uint) (*model.PdfTemplate, error) {
	request := &proto.TemplateRequest{Id: uint32(templateID)}

	res, err := c.GetTemplate(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return fromGrpcTemplate(res)
}

func fromGrpcTemplate(template *proto.Template) (*model.PdfTemplate, error) {
	var elements []model.Element

	for _, element := range template.Elements {
		elements = append(elements, fromGrpcElement(element))
	}

	createdAt, err := time.Parse(template.CreatedAt, DATE_FORMAT)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(template.UpdatedAt, DATE_FORMAT)
	if err != nil {
		return nil, err
	}

	return &model.PdfTemplate{
		ID:        uint(template.Id),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      template.Name,
		Elements:  elements,
	}, nil
}

func fromGrpcElement(element *proto.Element) model.Element {
	return model.Element{
		ID:        uint(element.Id),
		Type:      element.Type,
		X:         int(element.X),
		Y:         int(element.Y),
		Width:     int(element.Width),
		Height:    int(element.Height),
		ValueFrom: element.ValueFrom,
		Font:      element.Font,
		FontSize:  int(element.FontSize),
	}
}
