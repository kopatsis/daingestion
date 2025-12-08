package live

import "encoding/json"

type ProductIDs struct {
	ProductVariant struct {
		ID      string `json:"id"`
		Product struct {
			ID string `json:"id"`
		} `json:"product"`
	} `json:"productVariant"`
}

// variant ID, product ID
func ExtractProductIDs(raw json.RawMessage) (string, string, error) {
	var d ProductIDs
	err := json.Unmarshal(raw, &d)
	if err != nil {
		return "", "", err
	}
	return d.ProductVariant.ID, d.ProductVariant.Product.ID, nil
}

type DataCollection struct {
	Collection struct {
		ID string `json:"id"`
	} `json:"collection"`
}

// collection ID
func ExtractCollectionID(raw json.RawMessage) (string, error) {
	var d DataCollection
	err := json.Unmarshal(raw, &d)
	if err != nil {
		return "", err
	}
	return d.Collection.ID, nil
}
