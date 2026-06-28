package model

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/bids"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/whitelabel"
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
	UserId string         `json:"user_id"`
	BidId  string         `json:"bid_id"`
	Origin bids.BidOrigin `json:"origin"`
}

type WhitelabelRequest struct {
	Id            string           `json:"id"`
	OwnerId       string           `json:"owner_id"`
	CompanyName   string           `json:"company_name"`
	LogoUri       string           `json:"logo_uri"`
	MobileLogoUri string           `json:"mobile_logo_uri"`
	BackgroundUri string           `json:"background_uri"`
	Theme         whitelabel.Theme `json:"theme"`
	FontFamily    string           `json:"font_family"`
	FontUri       string           `json:"font_uri"`

	PrimaryColor    string `json:"primary_color"`
	SecondaryColor  string `json:"secondary_color"`
	AccentColor     string `json:"accent_color"`
	BackgroundColor string `json:"background_color"`
	SurfaceColor    string `json:"surface_color"`
	TextColor       string `json:"text_color"`
	BorderColor     string `json:"border_color"`
	SuccessColor    string `json:"success_color"`
	ErrorColor      string `json:"error_color"`
	WarningColor    string `json:"warning_color"`

	// specific image cases
	LogoImage       io.Reader `json:"logo_image"`
	MobileLogoImage io.Reader `json:"mobile_logo_image"`
	BackgroundImage io.Reader `json:"background_image"`
}

type BidRequest struct {
}

type CdnResponse struct {
	URL string `json:"url"`
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

func (w *WhitelabelRequest) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(w)
}

func (c *CdnResponse) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(c)
}
