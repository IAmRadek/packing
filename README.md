# Packing Optimization System

A web application for optimizing the allocation of items into containers of different sizes. This system helps determine the most efficient way to pack items, minimizing the number of containers used while satisfying demand.

## Overview

The Packing Optimization System is designed to solve the bin packing problem, specifically for scenarios where you need to distribute items into containers of different sizes. The system uses a dynamic programming algorithm to find the optimal allocation of items to minimize waste.

Key features:

- Create and manage inventories with different container sizes
- Calculate optimal packing distributions for a given demand
- Web interface for easy interaction with the system
- RESTful API for programmatic access

Live Demo:

https://packing.fly.dev

## How to Run

### Prerequisites

- Go 1.24 or higher
- Make

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/IAmRadek/packing.git
   cd packing
   ```

2. Create a `.development.env` file with your configuration:
   ```
   ADDR=:8080
   READ_TIMEOUT=10s
   READ_HEADER_TIMEOUT=10s
   WRITE_TIMEOUT=10s
   IDLE_TIMEOUT=10s
   MAX_HEADER_BYTES=1024
   GRACEFUL_SHUTDOWN_DURATION=5s
   ```

### Testing the Application

Using Make:

```
make tests
```

### Running the Application

Using Make:

```
make run
```

### Building the Application

Using Make:

```
make build
```

The built binary will be in the `build` directory.

## Project Structure

```
.
├── cmd/                  # Application entry points
│   └── webd/             # Web server
├── internal/             # Internal packages
│   ├── algorithms/       # Packing algorithms
│   │   └── dp/           # Dynamic programming implementation
│   ├── app/              # Application services
│   │   ├── allocation/   # Allocation service
│   │   └── inventory/    # Inventory service
│   ├── domain/           # Domain models
│   │   └── pack/         # Packing domain models
│   ├── handlers/         # HTTP handlers
│   ├── infra/            # Infrastructure code
│   └── templates/        # HTML templates
│       ├── layouts/      # Layout templates
│       └── pages/        # Page templates
├── Makefile              # Build and run commands
└── README.md             # This file
```

## Architecture

The application follows a clean architecture approach with a clear separation of concerns:

1. **Domain Layer**: Contains the core business logic and entities (pack models, inventory)
2. **Application Layer**: Implements use cases using domain entities (allocation and inventory services)
3. **Infrastructure Layer**: Provides implementations for external dependencies (repositories)
4. **Interface Layer**: Handles external interactions (HTTP handlers, templates)

## API Endpoints

- `GET /`: Home page
- `GET /inventory`: List all inventories
- `GET/POST /inventory/create`: Create a new inventory
- `GET/POST /inventory/{sku}`: View inventory details and calculate allocations
- `POST /inventory/{sku}/update`: Update inventory sizes
- `POST /inventory/{sku}/delete`: Deletes inventory
- `GET /api/allocate`: API endpoint for allocation calculation


## Possible improvements

- Error handling: export more sentinel errors and translate them to human errors when returning to UI.
- Tracing and measuring algorythm performance.
- Middleware for tracing and logging.
- Persistence using a real database, not in memory one.

## License

© 2025 Radosław Dejnek. All rights reserved.
