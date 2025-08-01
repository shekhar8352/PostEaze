# Build stage
FROM node:20-alpine AS builder

WORKDIR /app
COPY . .
RUN npm install && npm run build

# Serve stage with custom nginx config
FROM nginx:stable-alpine

# Copy built frontend assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Add custom nginx config for SPA routing
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80