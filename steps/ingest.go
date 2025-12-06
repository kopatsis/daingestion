package steps

func CheckEvent(param string) bool {

	valid := map[string]struct{}{
		"cart_viewed":                      {},
		"checkout_address_info_submitted":  {},
		"checkout_completed":               {},
		"checkout_contact_info_submitted":  {},
		"checkout_shipping_info_submitted": {},
		"checkout_started":                 {},
		"collection_viewed":                {},
		"page_viewed":                      {},
		"payment_info_submitted":           {},
		"product_added_to_cart":            {},
		"product_removed_from_cart":        {},
		"product_viewed":                   {},
		"search_submitted":                 {},
	}

	_, ok := valid[param]
	return ok
}
