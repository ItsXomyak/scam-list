const express = require('express')
const puppeteer = require('puppeteer')
require('dotenv').config()

const app = express()
app.use(express.json())

// POST /parse-domain
// body: { "domain": "vk-com" }
app.post('/parse-domain', async (req, res) => {
	const { domain } = req.body
	if (!domain) {
		return res.status(400).json({ error: 'Domain is required' })
	}

	const url = `https://scam-detector.com/validator/${domain}-review`

	let browser
	try {
		browser = await puppeteer.launch({
			headless: false,
			args: ['--no-sandbox', '--disable-setuid-sandbox'],
			defaultViewport: null,
		})

		const page = await browser.newPage()
		await page.goto(url, { waitUntil: 'domcontentloaded' })
		await page.waitForSelector('div.factcesAccordion', { timeout: 30000 })

		// --- Technical Analysis ---
		const technicalAnalysis = await page.evaluate(() => {
			const accordion = document.querySelector(
				'div.factcesAccordion .AccordionWrapper .accordion'
			)
			if (!accordion) return {}

			const panels = accordion.querySelectorAll('.panel')
			const data = {}

			panels.forEach(panel => {
				const headerEl = panel.querySelector('.panel-heading h4')
				if (!headerEl) return

				const header = headerEl.innerText.trim()
				const bodyEl = panel.querySelector('.panel-body .content-wrapper')

				if (!bodyEl || !bodyEl.innerText.trim()) {
					data[header] = null
					return
				}

				const items = Array.from(bodyEl.querySelectorAll('p')).map(p => {
					const strong = p.querySelector('strong')
					if (strong) {
						const br = strong.nextElementSibling
						if (br && br.nodeName === 'BR') {
							const value = br.nextSibling
								? br.nextSibling.textContent.trim()
								: ''
							return { [strong.innerText.trim()]: value }
						}
						const text = p.innerText.replace(strong.innerText, '').trim()
						return { [strong.innerText.trim()]: text }
					}
					return p.innerText.trim()
				})

				if (items.every(i => typeof i === 'object')) {
					data[header] = Object.assign({}, ...items)
				} else {
					data[header] = items
				}
			})

			return data
		})

		// --- Основная информация ---
		const mainInfo = await page.evaluate(() => {
			const totalPercent =
				document.querySelector('p.totalPercent strong')?.innerText || null
			const domainAge =
				document.querySelector('div.panel-body p')?.innerText || null
			const domainDate =
				document.querySelector('ul.WOTDetailsList p.orange')?.innerText || null
			const allData =
				document.querySelector('div.onlinePaymentsSec')?.innerText || null

			return { totalPercent, allData, domainAge, domainDate }
		})

		let blackList = null
		let httpsConnection = null
		let siteDescription = null
		if (mainInfo.allData) {
			const filteredData = mainInfo.allData
				.split(/\n\s*\n/)
				.map(s => s.trim())
				.filter(Boolean)
				.slice(1, -1)

			blackList = filteredData[3] || null
			httpsConnection = filteredData[5] || null
			siteDescription = filteredData.slice(7).join(' ') || null
		}

		const result = {
			technicalAnalysis,
			summary: {
				totalPercent: mainInfo.totalPercent,
				domainAge: mainInfo.domainAge,
				domainDate: mainInfo.domainDate,
				blackList,
				httpsConnection,
				siteDescription,
			},
		}

		res.json(result)
	} catch (err) {
		console.error(err)
		res.status(500).json({ error: 'Parsing failed', details: err.message })
	} finally {
		if (browser) await browser.close()
	}
})

const PORT = process.env.PORT || 4000
app.listen(PORT, () => {
	console.log(`Server running on port ${PORT}`)
})
