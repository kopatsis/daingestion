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

    async function SetSessionCookie(browser, id) {
        await browser.cookies.set('px_sess_id', id, {
            maxAge: 1800,
            sameSite: 'None',
            secure: true,
            path: '/'
        })
    }

    function CheckParsed(parsed) {
        if (!parsed) return false
        if (typeof parsed.t !== 'number' || !Number.isFinite(parsed.t)) return false
        if (typeof parsed.id !== 'string') return false

        if (!parsed.id.startsWith('PXID-')) return false

        const uuid = parsed.id.slice(5)
        const parts = uuid.split('-')
        if (parts.length !== 5) return false
        if (parts[0].length !== 8) return false
        if (parts[1].length !== 4) return false
        if (parts[2].length !== 4) return false
        if (parts[3].length !== 4) return false
        if (parts[4].length !== 12) return false

        return true
    }


    async function GetSessionId(browser) {
        const storage = browser.localStorage
        const now = Date.now()
        const raw = await storage.getItem('px_sess')

        let expired = null

        if (raw) {
            const parsed = JSON.parse(raw)
            if (CheckParsed(parsed)) {
                if ((now - parsed.t) < 30 * 60 * 1000) {
                    const updated = { id: parsed.id, t: now }
                    await storage.setItem('px_sess', JSON.stringify(updated))
                    return [parsed.id, null]
                } else {
                    expired = parsed.id
                }
            }
        }

        const id = 'PXID-' + crypto.randomUUID()
        const obj = { id: id, t: now }
        await storage.setItem('px_sess', JSON.stringify(obj))
        return [id, expired]
    }


    for (const ev of events) {
        analytics.subscribe(ev, async (event) => {
            const [sid, old] = await GetSessionId(browser)
            await SetSessionCookie(browser, sid)

            const sessionObj = old ? { current: sid, previous: old } : { current: sid }

            const payload = {
                event: event,
                init: baseInit,
                session: sessionObj
            }

            fetch('https://YOUR_ENDPOINT_HERE', {
                method: 'POST',
                credentials: 'include',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            })
        })
    }
})
