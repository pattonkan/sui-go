package isc

type NormalizedID struct {
	ID string `json:"id"`
}

type NormalizedAssets struct {
	Type   string `json:"type"`
	Fields struct {
		Coins struct {
			Type   string `json:"type"`
			Fields struct {
				NormalizedID `json:"id"`
				// FIXME this should be int
				Size string `json:"size"`
			} `json:"fields"`
		} `json:"coins"`
		NormalizedID `json:"id"`
		Nfts         []interface{} `json:"nfts"`
	} `json:"fields"`
}
