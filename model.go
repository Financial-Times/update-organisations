package main

type organisation struct {
	UUID                   string                 `json:"uuid"`
	Type                   OrgType                `json:"type"`
	ProperName             string                 `json:"properName"`
	PrefLabel              string                 `json:"prefLabel"`
	LegalName              string                 `json:"legalName,omitempty"`
	ShortName              string                 `json:"shortName,omitempty"`
	HiddenLabel            string                 `json:"hiddenLabel,omitempty"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers"`
	TradeNames             []string               `json:"tradeNames,omitempty"`
	LocalNames             []string               `json:"localNames,omitempty"`
	FormerNames            []string               `json:"formerNames,omitempty"`
	Aliases                []string               `json:"aliases,omitempty"`
	IndustryClassification string                 `json:"industryClassification,omitempty"`
	ParentOrganisation     string                 `json:"parentOrganisation,omitempty"`
}

type alternativeIdentifiers struct {
	TME               []string `json:"TME,omitempty"`
	UUIDS             []string `json:"uuids"`
	FactsetIdentifier string   `json:"factsetIdentifier,omitempty"`
	LeiCode           string   `json:"leiCode,omitempty"`
}

//OrgType is the type of an Organisation
type OrgType string
