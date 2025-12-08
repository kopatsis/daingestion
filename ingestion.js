import { register } from '@shopify/web-pixels-extension';

register(({ analytics, browser, init }) => {
    const privacy = init.customerPrivacy
    const allowed = privacy.analyticsProcessingAllowed
    const baseInit = allowed ? init.data : { cart: init.data.cart, shop: init.data.shop }

    const events = [
        'cart_viewed',
        'checkout_address_info_submitted',
        'checkout_completed',
        'checkout_contact_info_submitted',
        'checkout_shipping_info_submitted',
        'checkout_started',
        'collection_viewed',
        'page_viewed',
        'payment_info_submitted',
        'product_added_to_cart',
        'product_removed_from_cart',
        'product_viewed',
        'search_submitted'
    ]


    for (const ev of events) {
        analytics.subscribe(ev, async (event) => {
            const privacy = init.customerPrivacy
            const allowed = privacy.analyticsProcessingAllowed
            const baseInit = allowed ? init.data : { cart: init.data.cart, shop: init.data.shop }        

            const payload = {
                event: event,
                init: baseInit,
                time: time.Now()
            }

            fetch('https://YOUR_ENDPOINT_HERE/' + ev, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            })
        })
    }
})
