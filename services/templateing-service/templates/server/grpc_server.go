package server

import (
	"context"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/constants"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/proto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/repository"
)

type GrpcServer struct {
	proto.UnimplementedTemplateServiceServer
	repo repository.Repository
}

func NewGrpcServer(repo repository.Repository) *GrpcServer {
	return &GrpcServer{repo: repo}
}

func (s *GrpcServer) GetTemplate(ctx context.Context, r *proto.TemplateRequest) (*proto.Template, error) {
	templateId := uint(r.GetId())

	template, err := s.repo.GetTemplateById(templateId)
	if err != nil {
		return nil, err
	}

	return toGrpcTemplate(template), nil
}

func toGrpcTemplate(template *model.PdfTemplate) *proto.Template {
	var elements []*proto.Element

	for _, element := range template.Elements {
		elements = append(elements, toGrpcElement(&element))
	}

	return &proto.Template{
		Id:        uint32(template.ID),
		CreatedAt: template.CreatedAt.Format(constants.GRPC_DATE_FORMAT),
		UpdatedAt: template.UpdatedAt.Format(constants.GRPC_DATE_FORMAT),
		Name:      template.Name,
		Elements:  elements,
	}
}

func toGrpcElement(element *model.Element) *proto.Element {
	return &proto.Element{
		Id:        uint32(element.ID),
		Type:      element.Type,
		X:         int32(element.X),
		Y:         int32(element.Y),
		Width:     int32(element.Width),
		Height:    int32(element.Height),
		ValueFrom: element.ValueFrom,
		Font:      element.Font,
		FontSize:  int32(element.FontSize),
	}
}
