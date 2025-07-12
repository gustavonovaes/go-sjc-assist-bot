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
Usage: go run ./cmd/cli/main.go <sspsp|sspsp:detailed [year]|cetesb|news|news:filtered|model:train|model:test [text]>


# example commands

> go run ./cmd/cli/main.go cetesb

=======================================================
| N1 - BOA        | N2 - MODERADA   | N3 - RUIM       |
| N4 - MUITO RUIM | N5 - PÉSSIMA    |                 |
=======================================================
Índice atual: N1 - BOA (23.00)
=======================================================

> go run ./cmd/cli/main.go sspsp:detailed 2022

 Jan | Fev | Mar | Abr | Mai | Jun | Jul | Ago | Set | Out | Nov | Dez | Delito 
-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+--------
 369 | 378 | 482 | 453 | 539 | 482 | 345 | 568 | 552 | 563 | 456 | 453 | FURTO - OUTROS 
  76 |  98 | 109 | 102 |  98 | 105 | 117 | 116 | 126 | 116 | 126 |  97 | FURTO DE VEÍCULO 
   0 |   0 |   0 |   1 |   0 |   0 |   0 |   0 |   0 |   0 |   0 |   0 | HOMICÍDIO CULPOSO OUTROS...
...

> go run ./cmd/cli/main.go model:train

Training model with:
 - good subjects: [emergencial povo apoio conscientização projeto infraestrutura conquista governo prefeitura municipio sjc são josé dos campos são josé estado de sp sp registra paulista investimento economia atinge meta operação acontece inaugura festa do fim de semana feira]
 - bad subjects: [acidente violencia mort familia corpo assassinato roubo furt incendio atropel apreensão cachorr sexual mutilad agredid sem vida asfixiou em coma desaparecid quadrilha trafico confusão agressão polícia suspeita drogas armas tiroteio assalto sequestro explosão criminoso bolsonaro lula taubaté jacareí rio de janeiro tecnico de copa américa brasileirão copa de clubes fifa]


> go run ./cmd/cli/main.go model:test "A prefeitura de São José dos Campos anunciou um novo projeto de infraestrutura que promete melhorar a qualidade de vida na cidade."
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