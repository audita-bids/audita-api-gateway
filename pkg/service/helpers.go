package service

import (
	model "audita-api-gateway/request"

	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/client"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func BuildPatchBusinessClient(user *client.ClientComplete, r *model.BusinessClientRequest) *client.PatchBusinessClientRequest {
	merged := &client.ClientComplete{
		Id:       r.Id,
		ClientId: r.ClientId,
		Status:   r.Status,
		OwnerId:  user.GetId(),
		Name:     r.Name,
		Email:    r.Email,
	}

	return &client.PatchBusinessClientRequest{
		Client:     merged,
		OwnerId:    user.GetId(),
		UpdateMask: &fieldmaskpb.FieldMask{Paths: r.UpdateMask},
	}
}
