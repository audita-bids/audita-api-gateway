package model

import (
	"encoding/json"
	"net/http"

	"github.com/project-pncp/private-kit/pkg/pb/protocols/bids"
)

type FavoriteBidRequest struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Sequence  int32  `json:"sequence"`
	ProcessId string `json:"process_id"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type AnalysisRequest struct {
	Id        string `json:"id"`
	UserID    string `json:"user_id"`
	ProcessID string `json:"process_id"`
	Sequence  int32  `json:"sequence"`
	Content   string `json:"content"`
	Base64    string `json:"base64"`
}

type IntegrationRequest struct {
	StartDate    string `json:"start_date"`
	Uf           string `json:"uf"`
	Cnpj         string `json:"cnpj"`
	ModalityCode string `json:"modality_code"`
	FinalDate    string `json:"final_date"`
	Page         int32  `json:"page"`
	PageSize     int32  `json:"page_size"`
	CityCode     string `json:"city_code"`
	Sequence     int32  `json:"sequence"`
}

type HoldingRequest struct {
	ProcessId             string         `json:"process_id"`
	Sequence              int32          `json:"sequence"`
	UserId                string         `json:"user_id"`
	Origin                bids.BidOrigin `json:"origin"`
	PublicationDate       string         `json:"publication_date"`
	ProposalOpeningDate   string         `json:"proposal_opening_date"`
	ProposalClosingDate   string         `json:"proposal_closing_date"`
	DisputeDate           string         `json:"dispute_date"`
	HomologationDate      string         `json:"homologation_date"`
	ContractSignDate      string         `json:"contract_sign_date"`
	ContractStartDate     string         `json:"contract_start_date"`
	ContractEndDate       string         `json:"contract_end_date"`
	ClarificationDeadline string         `json:"clarification_deadline"`
	AppealDeadline        string         `json:"appeal_deadline"`
}

func (f *FavoriteBidRequest) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(f)
}

func (i *IntegrationRequest) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(i)
}

func (a *AnalysisRequest) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(a)
}

func (h *HoldingRequest) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(h)
}
