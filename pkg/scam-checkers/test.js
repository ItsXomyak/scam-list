const express = require('express')
const puppeteer = require('puppeteer')
require('dotenv').config()

const app = express()
app.use(express.json())

// POST /parse-domain
// body: { "domain": "youtube.com" }
app.post('/parse-domain', async (req, res) => {
	const { domain } = req.body
	if (!domain) {
		return res.status(400).json({ error: 'Domain is required' })
	}
	const url = `https://www.phishtank.org`

	let browser
	try {
		browser = await puppeteer.launch({
			headless: false,
			args: ['--no-sandbox', '--disable-setuid-sandbox'],
			defaultViewport: null,
		})

		const page = await browser.newPage()

		// Добавляем обработку различных Cloudflare проверок
		page.on('response', response => {
			if (
				response.url().includes('captcha') ||
				response.url().includes('challenge')
			) {
				console.log('Cloudflare challenge detected:', response.url())
			}
		})

		await page.goto(url, {
			waitUntil: 'domcontentloaded',
			timeout: 30000,
		})

		// Ждем и проверяем разные селекторы Cloudflare
		const cloudflareSelectors = [
			'input[type="checkbox"]',
			'#challenge-form',
			'.cf-challenge',
			'[data-sitekey]',
			'#cf-challenge-running',
		]

		let cloudflareDetected = false

		// Проверяем все возможные селекторы Cloudflare
		for (const selector of cloudflareSelectors) {
			const element = await page.$(selector)
			if (element) {
				console.log('Found Cloudflare element:', selector)
				cloudflareDetected = true

				// Если это чекбокс - кликаем
				if (selector === 'input[type="checkbox"]') {
					await page.click(selector)
					await new Promise(resolve => setTimeout(resolve, 2000))
				}

				// Ждем исчезновения Cloudflare
				await page
					.waitForFunction(() => !document.querySelector('${selector}'), {
						timeout: 15000,
					})
					.catch(() => console.log('Cloudflare might still be present'))

				break
			}
		}

		// Ждем основную страницу после Cloudflare
		await page
			.waitForSelector('input[name="isaphishurl"]', {
				timeout: 20000,
			})
			.catch(() => {
				throw new Error('Main page not loaded after Cloudflare')
			})

		const fullUrl = `http://${domain}`

		// Очищаем поле и вводим URL
		await page.evaluate(() => {
			const input = document.querySelector('input[name="isaphishurl"]')
			if (input) input.value = ''
		})

		await page.type('input[name="isaphishurl"]', fullUrl)
		await new Promise(resolve => setTimeout(resolve, 1000))

		// Нажимаем кнопку проверки
		await page.click('input[type="submit"][value="Is it a phish?"]')

		// Ждем навигации с таймаутом
		await Promise.race([
			page.waitForNavigation({
				waitUntil: 'domcontentloaded',
				timeout: 15000,
			}),
			new Promise(resolve => setTimeout(resolve, 10000)),
		])

		// Проверяем снова на Cloudflare после отправки
		for (const selector of cloudflareSelectors) {
			const element = await page.$(selector)
			if (element) {
				console.log('Cloudflare after submit:', selector)
				if (selector === 'input[type="checkbox"]') {
					await page.click(selector)
				}
				await new Promise(resolve => setTimeout(resolve, 5000))
				break
			}
		}

		// Получаем текущий URL для отладки
		const currentUrl = page.url()
		console.log('Current URL:', currentUrl)

		res.json({
			domain,
			checkedUrl: fullUrl,
			currentUrl: currentUrl,
			status: 'page_loaded',
			message: 'Navigation completed',
			timestamp: new Date().toISOString(),
		})
	} catch (error) {
		if (browser) await browser.close()
		console.error('Error:', error)
		res.status(500).json({
			error: 'Processing failed',
			details: error.message,
		})
	}
})

// Запуск сервера
const PORT = process.env.PORT || 4001
app.listen(PORT, () => {
	console.log(`Server running on port ${PORT}`)
})
