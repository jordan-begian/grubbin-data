FROM oven/bun:alpine

WORKDIR /app

COPY package.json bun.lockb* ./
RUN bun install

COPY . .

EXPOSE 5173

CMD ["sh", "-c", "bun install && bun run dev"]
