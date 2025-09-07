FROM node:20-alpine AS builder
WORKDIR /app
# Install required system packages
RUN apk add --no-cache libc6-compat
# Install dependencies based on lockfile
COPY package.json pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
# Copy source code and build
COPY . .
# Build-time env for public API base URL
# NOTE: NEXT_PUBLIC_ で始まる値はクライアントに埋め込まれるため、
# 本番適用にはビルド時に値を渡す必要があります
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
RUN pnpm build && pnpm prune --prod

FROM node:20-alpine AS runner
WORKDIR /app
RUN apk add --no-cache libc6-compat && corepack enable
ENV NODE_ENV=production
# ランタイムでも参照できるように環境変数を保持（SSR 等の将来対応を考慮）
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
COPY --from=builder /app ./
EXPOSE 3000
CMD ["pnpm", "start"]
