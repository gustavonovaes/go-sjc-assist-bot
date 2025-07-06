# go-sjc-assist-bot

A lightweight bot designed to provide useful information about São José dos Campos, Brazil, integrating with government APIs to deliver real-time data about air quality, crime statistics, and more.

## Supported Platforms
- [x] **Telegram** - [Try it now!](https://t.me/sjc_assist_bot)
- [ ] **Discord** - Coming soon

## Available Commands

| Command | Description |
|---------|-------------|
| `/about` | Exibe informações sobre o bot e como contribuir |
| `/qualidadeAr` | Exibe o índice de qualidade do ar da cidade via CETESB |
| `/crimes` | Exibe o total de crimes registrados na cidade nos últimos anos |
| `/mapaCrimes` | Exibe link para o mapa com a marcações dos crimes registrados no último semestre |

## Running as CLI
```bash
go run cmd/cli/main.go -service cetesb -city_id 49

# Output
Nome: S.José Campos
Indice qualidade do Ar: 22.00000
```

## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, please open an issue or submit a pull request.

## Data Sources

- **CETESB**: Air quality monitoring data
- **SSP-SP**: Crime statistics from São Paulo State Public Security
- **Local News**: Aggregated from regional news sources

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- CETESB for providing air quality data
- SSP-SP for crime statistics
- São José dos Campos discord community [DEVs de SJC](https://discord.gg/ENXmaanGFD)