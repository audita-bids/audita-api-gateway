package model

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/bids"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/billings"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/client"
	"github.com/newdesksoftwares/private-kit/pkg/pb/protocols/whitelabel"
)

type FavoriteBidRequest struct {
	Id        string `json:"id"`
	BidId     string `json:"bid_id"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type AnalysisRequest struct {
	Id     string `json:"id"`
	BidId  string `json:"bid_id"`
	Base64 string `json:"base64"`
	UserId string `json:"user_id"`
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
	UserId           string         `json:"user_id"`
	BidId            string         `json:"bid_id"`
	PublicationMonth int32          `json:"publication_month"`
	Origin           bids.BidOrigin `json:"origin"`
}

type WhitelabelRequest struct {
	Id            string           `json:"id"`
	ManagerId     string           `json:"manager_id"`
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
	Id string
}

type ProposalRequest struct {
	Id          string          `json:"id"`
	BidId       string          `json:"bid_id"`
	Files       []*bids.BidFile `json:"files"`
	Value       float32         `json:"value"`
	Observation string          `json:"observation"`
}

type PaymentRequest struct {
	Id                      string                   `json:"id"`
	PlanId                  string                   `json:"plan_id"`
	PaymentMethod           billings.PaymentMethod   `json:"payment_method"`
	ProviderPaymentMethodId string                   `json:"provider_payment_method_id"`
	Token                   string                   `json:"token"`
	Installments            int32                    `json:"installments"`
	IssuerId                string                   `json:"issuer_id"`
	Payer                   *billings.PayerComplete  `json:"payer"`
	IdempotencyKey          string                   `json:"idempotency_key"`
	DeviceId                string                   `json:"device_id"`
	IpAddress               string                   `json:"-"`
	CallbackUrl             string                   `json:"callback_url"`
	CouponCode              string                   `json:"coupon_code"`
	Metadata                map[string]string        `json:"metadata"`
	PaymentProvider         billings.PaymentProvider `json:"payment_provider"`
}

type PlanRequest struct {
	billings.PlanComplete
}

type ClientRequest struct {
}

type ScopeRequest struct{}

type BusinessClientRequest struct {
	Id         string              `json:"id"`
	Status     client.ClientStatus `json:"status"`
	UpdateMask []string            `json:"update_mask"`

	client.CreateBusinessClientRequest
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
	if r.Body == nil {
		return nil
	}
	if err := json.NewDecoder(r.Body).Decode(w); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func (c *CdnResponse) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(c)
}

func (p *ProposalRequest) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(p)
}

func (p *PaymentRequest) Decode(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return err
	}

	p.IpAddress = clientIP(r)

	return nil
}

func (p *PlanRequest) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(p)
}

func (b *BusinessClientRequest) Decode(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(b)
}

func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if first, _, found := strings.Cut(forwarded, ","); found {
			return strings.TrimSpace(first)
		}
		return strings.TrimSpace(forwarded)
	}

	if real := r.Header.Get("X-Real-Ip"); real != "" {
		return strings.TrimSpace(real)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
