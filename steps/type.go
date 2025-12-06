package steps

import "strings"

type PageType string

const (
	PageHome        PageType = "home"
	PageProduct     PageType = "product"
	PageCollection  PageType = "collection"
	PagePage        PageType = "page"
	PageAccount     PageType = "account"
	PageSearch      PageType = "search"
	PageCart        PageType = "cart"
	PageCheckout    PageType = "checkout"
	PageOrderStatus PageType = "order_status"
	PageOther       PageType = "other"
)

func Classify(path string) PageType {
	if path == "" || path == "/" {
		return PageHome
	}

	p := path
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	if strings.Contains(p, "/products/") {
		return PageProduct
	}

	if strings.HasPrefix(p, "/collections/") {
		return PageCollection
	}

	if strings.HasPrefix(p, "/pages/") {
		return PagePage
	}

	if strings.HasPrefix(p, "/account") {
		return PageAccount
	}

	if strings.HasPrefix(p, "/search") {
		return PageSearch
	}

	if strings.HasPrefix(p, "/cart") {
		return PageCart
	}

	if strings.HasPrefix(p, "/checkouts") {
		return PageCheckout
	}

	if strings.HasPrefix(p, "/orders/") || strings.Contains(p, "/order/") {
		return PageOrderStatus
	}

	return PageOther
}
