# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instala dependências necessárias para build
RUN apk add --no-cache git

# Copia os arquivos go mod e soma para cache eficiente
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app

# Final stage: imagem mínima
FROM scratch

WORKDIR /app

# Copia o binário do estágio de build
COPY --from=builder /app/app .

# Expõe a porta padrão (ajuste conforme necessário)
EXPOSE 8080

# Comando de inicialização
ENTRYPOINT ["./app"]