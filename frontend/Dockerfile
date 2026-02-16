ARG NODE_VERSION=25.5

FROM node:${NODE_VERSION}-alpine AS base

FROM base AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build

FROM base AS runner
WORKDIR /app
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/.next/standalone ./
# COPY --from=builder /app/public ./public

ENV HOSTNAME="0.0.0.0"
ENV PORT=3000
EXPOSE 3000

CMD ["node", "server.js"]