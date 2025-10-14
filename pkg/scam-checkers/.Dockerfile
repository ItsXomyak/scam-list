# ========================
# Stage 1: Build
# ========================
FROM node:20-slim AS builder

# Устанавливаем зависимости для Puppeteer
RUN apt-get update && apt-get install -y \
    ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcups2 \
    libdbus-1-3 \
    libdrm2 \
    libgbm1 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libx11-xcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    xdg-utils \
    wget \
    curl \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Создаем рабочую директорию
WORKDIR /usr/src/app

# Копируем package.json и package-lock.json (для быстрого кэширования npm install)
COPY pkg/scamcheck-parser/package*.json ./

# Устанавливаем зависимости
RUN npm install --production

# Копируем остальной код
COPY pkg/scamcheck-parser/ .

# ========================
# Stage 2: Run
# ========================
FROM node:20-slim

# Создаем рабочую директорию
WORKDIR /usr/src/app

# Копируем зависимости и код из builder
COPY --from=builder /usr/src/app .

# Устанавливаем необходимые библиотеки для Puppeteer
RUN apt-get update && apt-get install -y \
    ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcups2 \
    libdbus-1-3 \
    libdrm2 \
    libgbm1 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libx11-xcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    xdg-utils \
    wget \
    curl \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Puppeteer запускается без sandbox
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=false
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/google-chrome
ENV NODE_ENV=production

# Проброс порта (берется из .env)
EXPOSE ${PORT:-3000}

# Запуск сервера
CMD ["node", "server.js"]
