## go-sjc-assist-bot

A lightweight bot designed to provide useful information about São José dos Campos, Brazil.

# Supported Commands

- `/qualidadeAr` - Exibe o índice de qualidade do ar da cidade via CETESB
- `/crimes` - Exibe o total de crimes registrados na cidade nos últimos anos
- `/mapaCrimes` - Exibe link para o mapa com a marcações dos crimes registrados no último semestre
- `/about` - Exibe informações sobre o bot e como contribuir

# Supported Platforms
- [x] Telegram ([Telegram Bot](https://t.me/sjc_assist_bot))
- [ ] Discord (Em breve)


# How to run as CLI
```bash
go run cmd/cli/main.go -service cetesb -city_id 49

# Output
Nome: S.José Campos
Indice qualidade do Ar: 22.00000
```
