const express = require('express')
const puppeteer = require('puppeteer')
const cluster = require('cluster')
const os = require('os')
require('dotenv').config()

const app = express()
app.use(express.json())

// Настройки для оптимизации
const BROWSER_POOL_SIZE = process.env.BROWSER_POOL_SIZE || 2
const MAX_CONCURRENT_REQUESTS = process.env.MAX_CONCURRENT_REQUESTS || 5
const REQUEST_TIMEOUT = process.env.REQUEST_TIMEOUT || 30000

class BrowserPool {
	constructor(size = BROWSER_POOL_SIZE) {
		this.pool = []
		this.size = size
		this.activeConnections = 0
		this.maxConnections = MAX_CONCURRENT_REQUESTS
	}

	async initialize() {
		console.log(`Initializing browser pool with ${this.size} browsers...`)
		for (let i = 0; i < this.size; i++) {
			try {
				const browser = await this.createBrowser()
				this.pool.push(browser)
			} catch (error) {
				console.error(`Failed to create browser ${i + 1}:`, error.message)
			}
		}
		console.log(`Browser pool initialized with ${this.pool.length} browsers`)
	}

	async createBrowser() {
		return await puppeteer.launch({
			headless: process.env.NODE_ENV === 'production' ? 'new' : false,
			args: [
				'--no-sandbox',
				'--disable-setuid-sandbox',
				'--disable-dev-shm-usage',
				'--disable-accelerated-2d-canvas',
				'--disable-gpu',
				'--window-size=1920,1080',
				'--disable-background-networking',
				'--disable-background-timer-throttling',
				'--disable-renderer-backgrounding',
				'--disable-backgrounding-occluded-windows',
				'--disable-client-side-phishing-detection',
				'--disable-ipc-flooding-protection',
				'--disable-hang-monitor',
				'--disable-popup-blocking',
				'--disable-prompt-on-repost',
				'--disable-sync',
				'--disable-translate',
				'--disable-features=TranslateUI',
				'--metrics-recording-only',
				'--no-first-run',
			],
			defaultViewport: { width: 1920, height: 1080 },
			ignoreHTTPSErrors: true,
		})
	}

	async getBrowser() {
		if (this.activeConnections >= this.maxConnections) {
			throw new Error('Maximum concurrent requests reached')
		}

		this.activeConnections++

		if (this.pool.length === 0) {
			// Создаем временный браузер если пул пуст
			return await this.createBrowser()
		}

		return this.pool.shift()
	}

	async returnBrowser(browser) {
		this.activeConnections--

		try {
			// Проверяем состояние браузера
			if (browser && browser.process() && !browser.process().killed) {
				// Закрываем все страницы кроме первой
				const pages = await browser.pages()
				for (let i = 1; i < pages.length; i++) {
					await pages[i].close()
				}
				this.pool.push(browser)
			} else {
				// Создаем новый браузер если старый неисправен
				const newBrowser = await this.createBrowser()
				this.pool.push(newBrowser)
			}
		} catch (error) {
			console.error('Error returning browser to pool:', error.message)
			try {
				await browser.close()
			} catch (closeError) {
				console.error('Error closing browser:', closeError.message)
			}
		}
	}

	async destroy() {
		console.log('Destroying browser pool...')
		const browsers = [...this.pool]
		this.pool = []

		await Promise.allSettled(browsers.map(browser => browser.close()))
		console.log('Browser pool destroyed')
	}
}

// Создаем пул браузеров
const browserPool = new BrowserPool()

// Кэш для результатов
const cache = new Map()
const CACHE_TTL = process.env.CACHE_TTL || 0 // 1 час

function getCacheKey(domain) {
	return `domain_${domain}`
}

function isCacheValid(timestamp) {
	return Date.now() - timestamp < CACHE_TTL
}

// Валидация домена
function validateDomain(domain) {
	if (!domain || typeof domain !== 'string') {
		return false
	}

	// Базовая проверка формата домена
	const domainRegex =
		/^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
	return domainRegex.test(domain.replace(/-/g, '.'))
}

// Оптимизированная функция парсинга
async function parseDomainData(page, url) {
	// Настройка перехватчика запросов для блокировки ненужных ресурсов
	await page.setRequestInterception(true)
	page.on('request', request => {
		const resourceType = request.resourceType()
		if (['image', 'stylesheet', 'font', 'media'].includes(resourceType)) {
			request.abort()
		} else {
			request.continue()
		}
	})

	// Переход на страницу с оптимизированными параметрами
	await page.goto(url, {
		waitUntil: 'domcontentloaded',
		timeout: REQUEST_TIMEOUT,
	})

	// Ожидание загрузки ключевых элементов
	try {
		await page.waitForSelector('div.factcesAccordion', { timeout: 10000 })
	} catch (error) {
		console.warn('factcesAccordion not found, continuing...')
	}

	// Извлечение данных одним запросом evaluate
	const scrapedData = await page.evaluate(() => {
		const result = {
			technicalAnalysis: {},
			mainInfo: {
				totalPercent: null,
				domainAge: null,
				domainDate: null,
				allData: null,
			},
		}

		// Technical Analysis
		try {
			const accordion = document.querySelector(
				'div.factcesAccordion .AccordionWrapper .accordion'
			)
			if (accordion) {
				const panels = accordion.querySelectorAll('.panel')

				panels.forEach(panel => {
					const headerEl = panel.querySelector('.panel-heading h4')
					if (!headerEl) return

					const header = headerEl.innerText.trim()
					const bodyEl = panel.querySelector('.panel-body .content-wrapper')

					if (!bodyEl || !bodyEl.innerText.trim()) {
						result.technicalAnalysis[header] = null
						return
					}

					const items = Array.from(bodyEl.querySelectorAll('p')).map(p => {
						const strong = p.querySelector('strong')
						if (strong) {
							const br = strong.nextElementSibling
							if (br && br.nodeName === 'BR') {
								const value = br.nextSibling?.textContent?.trim() || ''
								return { [strong.innerText.trim()]: value }
							}
							const text = p.innerText.replace(strong.innerText, '').trim()
							return { [strong.innerText.trim()]: text }
						}
						return p.innerText.trim()
					})

					result.technicalAnalysis[header] = items.every(
						i => typeof i === 'object'
					)
						? Object.assign({}, ...items)
						: items
				})
			}
		} catch (error) {
			console.warn('Error parsing technical analysis:', error.message)
		}

		// Main Info
		try {
			result.mainInfo.totalPercent =
				document.querySelector('p.totalPercent strong')?.innerText || null
			result.mainInfo.domainAge =
				document.querySelector('div.panel-body p')?.innerText || null
			result.mainInfo.domainDate =
				document.querySelector('ul.WOTDetailsList p.orange')?.innerText || null
			result.mainInfo.allData =
				document.querySelector('div.onlinePaymentsSec')?.innerText || null
		} catch (error) {
			console.warn('Error parsing main info:', error.message)
		}

		return result
	})

	// Обработка дополнительных данных
	let blackList = null
	let httpsConnection = null
	let siteDescription = null

	if (scrapedData.mainInfo.allData) {
		try {
			const filteredData = scrapedData.mainInfo.allData
				.split(/\n\s*\n/)
				.map(s => s.trim())
				.filter(Boolean)
				.slice(1, -1)

			blackList = filteredData[3] || null
			httpsConnection = filteredData[5] || null
			siteDescription = filteredData.slice(7).join(' ') || null
		} catch (error) {
			console.warn('Error processing additional data:', error.message)
		}
	}

	return {
		technicalAnalysis: scrapedData.technicalAnalysis,
		summary: {
			totalPercent: scrapedData.mainInfo.totalPercent,
			domainAge: scrapedData.mainInfo.domainAge,
			domainDate: scrapedData.mainInfo.domainDate,
			blackList,
			httpsConnection,
			siteDescription,
		},
	}
}

// Middleware для ограничения количества запросов
const rateLimiter = (req, res, next) => {
	const now = Date.now()
	const windowMs = 60000 // 1 минута
	const maxRequests = 30

	if (!rateLimiter.requests) {
		rateLimiter.requests = new Map()
	}

	const clientId = req.ip || 'unknown'
	const clientRequests = rateLimiter.requests.get(clientId) || []

	// Очищаем старые запросы
	const validRequests = clientRequests.filter(time => now - time < windowMs)

	if (validRequests.length >= maxRequests) {
		return res.status(429).json({
			error: 'Too many requests',
			retryAfter: Math.ceil((validRequests[0] + windowMs - now) / 1000),
		})
	}

	validRequests.push(now)
	rateLimiter.requests.set(clientId, validRequests)
	next()
}

// Основной эндпоинт
app.post('/parse-domain', rateLimiter, async (req, res) => {
	const startTime = Date.now()
	let browser = null
	let page = null

	try {
		const { domain } = req.body

		// Валидация
		if (!validateDomain(domain)) {
			return res.status(400).json({
				error: 'Invalid domain format',
				domain: domain,
			})
		}

		// Проверка кэша
		const cacheKey = getCacheKey(domain)
		const cached = cache.get(cacheKey)

		if (cached && isCacheValid(cached.timestamp)) {
			return res.json({
				...cached.data,
				cached: true,
				processingTime: Date.now() - startTime,
			})
		}

		const url = `https://scam-detector.com/validator/${domain}-review`

		// Получение браузера из пула
		browser = await browserPool.getBrowser()
		page = await browser.newPage()

		// Настройка страницы для оптимальной производительности
		await page.setUserAgent(
			'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
		)
		await page.setViewport({ width: 1920, height: 1080 })
		await page.setDefaultTimeout(REQUEST_TIMEOUT)

		// Парсинг данных
		const result = await parseDomainData(page, url)

		// Сохранение в кэш
		cache.set(cacheKey, {
			data: result,
			timestamp: Date.now(),
		})

		// Очистка старых записей кэша
		if (cache.size > 1000) {
			const entries = Array.from(cache.entries())
			entries.slice(0, 500).forEach(([key]) => cache.delete(key))
		}

		await page.close()
		await browserPool.returnBrowser(browser)

		res.json({
			...result,
			cached: false,
			processingTime: Date.now() - startTime,
		})
	} catch (error) {
		console.error('Parsing error:', error)

		try {
			if (page) await page.close()
			if (browser) await browserPool.returnBrowser(browser)
		} catch (cleanupError) {
			console.error('Cleanup error:', cleanupError.message)
		}

		const statusCode = error.message.includes('timeout')
			? 504
			: error.message.includes('Maximum concurrent')
			? 503
			: 500

		res.status(statusCode).json({
			error: 'Parsing failed',
			details: error.message,
			domain: req.body?.domain,
			processingTime: Date.now() - startTime,
		})
	}
})

// Эндпоинт для проверки здоровья сервиса
app.get('/health', (req, res) => {
	res.json({
		status: 'ok',
		uptime: process.uptime(),
		memory: process.memoryUsage(),
		browserPool: {
			available: browserPool.pool.length,
			active: browserPool.activeConnections,
			max: browserPool.maxConnections,
		},
		cache: {
			size: cache.size,
			maxSize: 1000,
		},
	})
})

// Эндпоинт для очистки кэша
app.delete('/cache', (req, res) => {
	const beforeSize = cache.size
	cache.clear()
	res.json({
		message: 'Cache cleared',
		clearedEntries: beforeSize,
	})
})

// Graceful shutdown
const gracefulShutdown = async () => {
	console.log('Received shutdown signal, closing server...')

	try {
		await browserPool.destroy()
		process.exit(0)
	} catch (error) {
		console.error('Error during shutdown:', error)
		process.exit(1)
	}
}

process.on('SIGTERM', gracefulShutdown)
process.on('SIGINT', gracefulShutdown)

// Инициализация и запуск сервера
const PORT = process.env.PORT || 4000

async function startServer() {
	try {
		await browserPool.initialize()

		app.listen(PORT, () => {
			console.log(`Server running on port ${PORT}`)
			console.log(`Browser pool size: ${browserPool.pool.length}`)
			console.log(`Max concurrent requests: ${MAX_CONCURRENT_REQUESTS}`)
			console.log(`Request timeout: ${REQUEST_TIMEOUT}ms`)
			console.log(`Cache TTL: ${CACHE_TTL}ms`)
		})
	} catch (error) {
		console.error('Failed to start server:', error)
		process.exit(1)
	}
}

startServer()
