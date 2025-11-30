package meta_service

import (
	"github.com/shekhar8352/PostEaze/provider/meta"
)

type MetaService interface {
	GetPagesFromCode(code string, redirectURI string) ([]meta.Page, error)
}

type MetaServiceImpl struct {
	provider meta.MetaProvider
}

func NewMetaService() *MetaServiceImpl {
	return &MetaServiceImpl{
		provider: meta.NewMetaProvider(),
	}
}

func (s *MetaServiceImpl) GetPagesFromCode(code string, redirectURI string) ([]meta.Page, error) {
	// 1. Exchange code for short-lived token
	tokenResp, err := s.provider.ExchangeCodeForToken(code, redirectURI)
	if err != nil {
		return nil, err
	}

	// 2. Exchange short-lived token for long-lived token
	longLivedTokenResp, err := s.provider.GetLongLivedToken(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	// 3. Fetch pages using long-lived token
	pages, err := s.provider.GetPages(longLivedTokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	return pages, nil
}
