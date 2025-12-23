package live

import (
	"dmd/logging"
	"encoding/json"
)

type ProductIDs struct {
	ProductVariant struct {
		ID      string `json:"id"`
		Product struct {
			ID string `json:"id"`
		} `json:"product"`
	} `json:"productVariant"`
}

// variant ID, product ID
func ExtractProductIDs(raw json.RawMessage, store, eventType, requestID string) (string, string, error) {
	var d ProductIDs
	err := json.Unmarshal(raw, &d)
	if err != nil {
		logging.LogError(
			"ERROR",
			"extraction_failure",
			"json",
			store,
			eventType,
			requestID,
			true,
			"unable to unmarshal product ids",
		)
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
func ExtractCollectionID(raw json.RawMessage, store, eventType, requestID string) (string, error) {
	var d DataCollection
	err := json.Unmarshal(raw, &d)
	if err != nil {
		logging.LogError(
			"ERROR",
			"extraction_failure",
			"json",
			store,
			eventType,
			requestID,
			true,
			"unable to unmarshal collection",
		)
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

func ExtractCheckoutLineIDs(raw json.RawMessage, store, eventType, requestID string) ([]LineIDs, error) {
	var d DataCheckout
	err := json.Unmarshal(raw, &d)
	if err != nil {
		logging.LogError(
			"ERROR",
			"extraction_failure",
			"json",
			store,
			eventType,
			requestID,
			true,
			"unable to unmarshal checkout",
		)
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

type DataProductVariantTitle struct {
	ProductVariant struct {
		Title string `json:"title"`
	} `json:"productVariant"`
}

func ExtractProductVariantTitle(raw json.RawMessage, store, eventType, requestID string) (string, error) {
	var d DataProductVariantTitle
	err := json.Unmarshal(raw, &d)
	if err != nil {
		logging.LogError(
			"ERROR",
			"extraction_failure",
			"json",
			store,
			eventType,
			requestID,
			true,
			"unable to unmarshal variant title",
		)
		return "", err
	}
	return d.ProductVariant.Title, nil
}

type DataCollectionTitle struct {
	Collection struct {
		Title string `json:"title"`
	} `json:"collection"`
}

func ExtractCollectionTitle(raw json.RawMessage, store, eventType, requestID string) (string, error) {
	var d DataCollectionTitle
	err := json.Unmarshal(raw, &d)
	if err != nil {
		logging.LogError(
			"ERROR",
			"extraction_failure",
			"json",
			store,
			eventType,
			requestID,
			true,
			"unable to unmarshal collection title",
		)
		return "", err
	}
	return d.Collection.Title, nil
}

type DataSearchQuery struct {
	SearchResult struct {
		Query string `json:"query"`
	} `json:"searchResult"`
}

func ExtractSearchQuery(raw json.RawMessage, store, eventType, requestID string) (string, error) {
	var d DataSearchQuery
	err := json.Unmarshal(raw, &d)
	if err != nil {
		logging.LogError(
			"ERROR",
			"extraction_failure",
			"json",
			store,
			eventType,
			requestID,
			true,
			"unable to unmarshal search query",
		)
		return "", err
	}
	return d.SearchResult.Query, nil
}

type DataCheckoutOrderID struct {
	Checkout struct {
		Order struct {
			ID string `json:"id"`
		} `json:"order"`
	} `json:"checkout"`
}

func ExtractCheckoutOrderID(raw json.RawMessage, store, eventType, requestID string) (string, error) {
	var d DataCheckoutOrderID
	err := json.Unmarshal(raw, &d)
	if err != nil {
		logging.LogError(
			"ERROR",
			"extraction_failure",
			"json",
			store,
			eventType,
			requestID,
			true,
			"unable to unmarshal order ID",
		)
		return "", err
	}
	return d.Checkout.Order.ID, nil
}
