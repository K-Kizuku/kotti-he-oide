FROM node:20-alpine AS builder
WORKDIR /app
# Install required system packages
RUN apk add --no-cache libc6-compat
# Install dependencies based on lockfile
COPY package.json pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
# Copy source code and build
COPY . .
RUN pnpm build && pnpm prune --prod

FROM node:20-alpine AS runner
WORKDIR /app
RUN apk add --no-cache libc6-compat && corepack enable
ENV NODE_ENV=production
COPY --from=builder /app ./
EXPOSE 3000
CMD ["pnpm", "start"]
