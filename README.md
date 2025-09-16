# UTTE Universe

UTTE Universe is a backend service for simulating a resource-based universe with buildings, resources, and game ticks. It is written in Go and designed for extensibility and integration.

> **Note:** This project is a case study for AI-assisted software development and experimentation. It was designed and developed using ChatGPT and GitHub Copilot.

## Features

- Resource and building type management
- Configurable game tick duration and universe seed
- YAML-based configuration
- Extensible core entities

## Getting Started

### Prerequisites

- Go 1.20+
- [Optional] Docker for containerized deployment

### Installation

```sh
git clone https://github.com/your-org/utte-universe.git
cd utte-universe
go mod tidy
```

### Configuration

Create a `config.yml` file in the project root. Example:

```yaml
tickDuration: 1s
universeSeed: 42
resources:
  - type: "Energy"
    initial: 1000
buildings:
  - type: "SolarPlant"
    count: 5
```

### Running

```sh
go run backend/main.go
```

### Testing

```sh
go test ./...
```

## Project Structure

- `core/` — Core entities and configuration
- `backend/` — Main backend service
- `config.yml` — Game configuration

## Contributing

Pull requests are welcome. For major changes, please open an issue first.

## License

MIT
