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

type LineIDs struct {
	ProductID string
	VariantID string
}

type DataCart struct {
	Cart struct {
		Lines []struct {
			Merchandise struct {
				ID      string `json:"id"`
				Product struct {
					ID string `json:"id"`
				} `json:"product"`
			} `json:"merchandise"`
		} `json:"lines"`
	} `json:"cart"`
}

func ExtractLineIDs(raw json.RawMessage) ([]LineIDs, error) {
	var d DataCart
	err := json.Unmarshal(raw, &d)
	if err != nil {
		return nil, err
	}

	out := make([]LineIDs, len(d.Cart.Lines))
	for i, l := range d.Cart.Lines {
		out[i] = LineIDs{
			ProductID: l.Merchandise.Product.ID,
			VariantID: l.Merchandise.ID,
		}
	}

	return out, nil
}

type DataCheckout struct {
	Checkout struct {
		LineItems []struct {
			Variant struct {
				ID      string `json:"id"`
				Product struct {
					ID string `json:"id"`
				} `json:"product"`
			} `json:"variant"`
		} `json:"lineItems"`
	} `json:"checkout"`
}

func ExtractCheckoutLineIDs(raw json.RawMessage) ([]LineIDs, error) {
	var d DataCheckout
	err := json.Unmarshal(raw, &d)
	if err != nil {
		return nil, err
	}

	out := make([]LineIDs, len(d.Checkout.LineItems))
	for i, l := range d.Checkout.LineItems {
		out[i] = LineIDs{
			ProductID: l.Variant.Product.ID,
			VariantID: l.Variant.ID,
		}
	}

	return out, nil
}
