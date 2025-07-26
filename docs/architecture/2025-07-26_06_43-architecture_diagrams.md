# Architecture Diagrams - Code Quality Assessment

Date: 2025-07-26T06:43:02+02:00

## Current State Architecture

### Dual Application Problem

```mermaid
graph TD
    A[User Request] --> B{Which System?}
    B -->|Portfolio| C[cmd/server/main.go]
    B -->|Contact Form| D[main.go]
    
    C --> E[internal/container/container.go]
    D --> F[internal/infrastructure/container.go]
    
    E --> G[Memory Repositories]
    F --> H[Domain Services]
    
    G --> I[No Database]
    H --> J[No Database]
    
    K[sqlc Generated Code] -.->|Unused| I
    K -.->|Unused| J
    
    style K fill:#ff9999
    style I fill:#ffcc99
    style J fill:#ffcc99
```

### Current Dependency Graph

```mermaid
graph LR
    subgraph "main.go Application"
        A[main.go] --> B[infrastructure/container.go]
        B --> C[domain.ContactRepository]
        B --> D[domain.EmailService]
        B --> E[domain.LoggingService]
        B --> F[application.ContactService]
    end
    
    subgraph "cmd/server Application"
        G[cmd/server/main.go] --> H[container/container.go]
        H --> I[Memory Repositories]
        H --> J[No Services]
    end
    
    subgraph "Ghost Systems"
        K[sqlc Generated Code] -.->|Disconnected| L[Database Layer]
        M[Config System] -.->|Partial| N[Both Applications]
    end
    
    style K fill:#ff9999
    style L fill:#ff9999
    style M fill:#ffcc99
    style N fill:#ffcc99
```

## Missing Library Dependencies

```mermaid
mindmap
    root((Missing Libraries))
        spf13/viper
            Configuration Management
            Environment Variables
            Config File Parsing
        gin-gonic/gin
            HTTP Framework
            Middleware Support
            Route Handling
        charmbracelet/log
            Structured Logging
            Log Levels
            Pretty Printing
        onsi/ginkgo
            BDD Testing
            Test Suites
            Matcher Library
        LarsArtmann/uniflow
            Error Handling
            Flow Control
            Result Types
        OpenTelemetry
            Distributed Tracing
            Metrics Collection
            Observability
```

## Code Quality Issues Categorization

```mermaid
pie title Issue Distribution by Impact
    "Architectural (51%)" : 51
    "Security (20%)" : 20
    "Performance (15%)" : 15
    "Maintainability (10%)" : 10
    "Cosmetic (4%)" : 4
```

## Quality Assessment Flow

```mermaid
flowchart TD
    A[Code Quality Request] --> B[Surface Analysis]
    B --> C{Linting Clean?}
    C -->|No| D[Fix Lint Issues]
    C -->|Yes| E[Duplication Analysis]
    
    D --> F[Static Errors]
    D --> G[Magic Numbers]
    D --> H[Formatting]
    
    E --> I[Architectural Analysis]
    I --> J{Systems Connected?}
    J -->|No| K[Ghost System Detection]
    J -->|Yes| L[Performance Analysis]
    
    K --> M[Dual Applications Found]
    M --> N[Missing Dependencies]
    N --> O[Integration Strategy]
    
    L --> P[Success]
    O --> Q[Comprehensive Plan]
    
    style M fill:#ff9999
    style N fill:#ff9999
    style K fill:#ffcc99
```

## Proposed Target Architecture

### Unified Application Architecture

```mermaid
graph TB
    subgraph "Presentation Layer"
        A[HTTP Server] --> B[Gin Router]
        B --> C[HTMX Handlers]
        B --> D[API Handlers]
    end
    
    subgraph "Application Layer"
        C --> E[Contact Service]
        C --> F[Portfolio Service]
        D --> G[API Service]
    end
    
    subgraph "Domain Layer"
        E --> H[Contact Aggregate]
        F --> I[Portfolio Aggregate]
        G --> J[User Aggregate]
        
        H --> K[Contact Events]
        I --> L[Portfolio Events]
        J --> M[User Events]
    end
    
    subgraph "Infrastructure Layer"
        K --> N[Event Store]
        L --> N
        M --> N
        
        E --> O[Contact Repository]
        F --> P[Portfolio Repository]
        G --> Q[User Repository]
        
        O --> R[SQLite Database]
        P --> R
        Q --> R
    end
    
    subgraph "Cross-Cutting Concerns"
        S[Configuration - Viper]
        T[Logging - LarsArtmann/uniflow]
        U[Observability - OpenTelemetry]
        V[Testing - Ginkgo]
    end
    
    style N fill:#99ff99
    style R fill:#99ff99
    style S fill:#99ccff
    style T fill:#99ccff
    style U fill:#99ccff
    style V fill:#99ccff
```

### Event Sourcing Architecture

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Handler
    participant A as Aggregate
    participant E as Event Store
    participant P as Projection
    participant R as Read Model
    
    C->>H: Command
    H->>A: Load from Events
    E->>A: Event Stream
    A->>A: Business Logic
    A->>E: New Events
    E->>P: Event Notification
    P->>R: Update Read Model
    H->>C: Response
```

### CQRS Implementation

```mermaid
graph LR
    subgraph "Command Side"
        A[Commands] --> B[Command Handlers]
        B --> C[Aggregates]
        C --> D[Event Store]
    end
    
    subgraph "Query Side"
        E[Queries] --> F[Query Handlers]
        F --> G[Read Models]
        G --> H[Database Views]
    end
    
    subgraph "Event Processing"
        D --> I[Event Bus]
        I --> J[Projections]
        J --> G
        I --> K[Process Managers]
        K --> A
    end
    
    style D fill:#99ff99
    style I fill:#99ccff
    style G fill:#ffcc99
```

## Migration Strategy

### Phase 1: Foundation

```mermaid
gantt
    title Migration Timeline
    dateFormat  YYYY-MM-DD
    section Foundation
    Merge Applications         :done, merge, 2025-07-26, 1d
    Add Missing Libraries      :done, libs, after merge, 2d
    Setup Configuration        :active, config, after libs, 1d
    
    section Architecture
    Implement DI Container     :di, after config, 2d
    Add Structured Logging     :logging, after di, 1d
    Setup Database Layer       :db, after logging, 2d
    
    section Quality
    Add Comprehensive Tests    :tests, after db, 3d
    Setup CI/CD Pipeline       :cicd, after tests, 2d
    Add Observability          :otel, after cicd, 2d
    
    section Advanced
    Implement Event Sourcing   :es, after otel, 5d
    Add CQRS Separation        :cqrs, after es, 3d
    Performance Optimization   :perf, after cqrs, 2d
```

### Dependency Resolution Order

```mermaid
graph TD
    A[1. Merge Applications] --> B[2. Add Core Libraries]
    B --> C[3. Setup Configuration]
    C --> D[4. Implement DI]
    D --> E[5. Add Database Layer]
    E --> F[6. Implement Services]
    F --> G[7. Add Event System]
    G --> H[8. Setup CQRS]
    H --> I[9. Add Observability]
    I --> J[10. Performance Tuning]
    
    style A fill:#ff9999,color:#000
    style B fill:#ffcc99,color:#000
    style C fill:#ffff99,color:#000
    style D fill:#ccff99,color:#000
    style E fill:#99ff99,color:#000
    style F fill:#99ffcc,color:#000
    style G fill:#99ccff,color:#000
    style H fill:#9999ff,color:#fff
    style I fill:#cc99ff,color:#fff
    style J fill:#ff99cc,color:#000
```

## System Integration Points

### External Services Integration

```mermaid
graph TB
    subgraph "Holger Hahn Website"
        A[Web Application]
        B[API Layer]
        C[Event System]
    end
    
    subgraph "External Services"
        D[Email Service - SMTP]
        E[Analytics - Google]
        F[Monitoring - OTEL]
        G[Database - SQLite/Turso]
    end
    
    subgraph "Development Tools"
        H[GitHub Actions]
        I[Docker Registry]
        J[Testing - Ginkgo]
        K[Linting - golangci-lint]
    end
    
    A --> D
    B --> E
    C --> F
    A --> G
    
    H --> I
    J --> A
    K --> A
    
    style D fill:#99ff99
    style E fill:#99ff99
    style F fill:#99ff99
    style G fill:#99ff99
```

## Testing Architecture

### Test Pyramid

```mermaid
graph TD
    subgraph "E2E Tests (10%)"
        A[Browser Tests]
        B[API Integration Tests]
    end
    
    subgraph "Integration Tests (30%)"
        C[Service Integration]
        D[Database Tests]
        E[Event System Tests]
    end
    
    subgraph "Unit Tests (60%)"
        F[Domain Logic Tests]
        G[Handler Tests]
        H[Repository Tests]
        I[Utility Function Tests]
    end
    
    style A fill:#ff9999
    style B fill:#ff9999
    style C fill:#ffcc99
    style D fill:#ffcc99
    style E fill:#ffcc99
    style F fill:#99ff99
    style G fill:#99ff99
    style H fill:#99ff99
    style I fill:#99ff99
```

## Deployment Architecture

### Production Topology

```mermaid
graph TB
    subgraph "Load Balancer"
        A[Nginx/Cloudflare]
    end
    
    subgraph "Application Tier"
        B[Go Application Pod 1]
        C[Go Application Pod 2]
        D[Go Application Pod 3]
    end
    
    subgraph "Data Tier"
        E[SQLite/Turso Primary]
        F[SQLite/Turso Replica]
    end
    
    subgraph "Observability"
        G[Metrics - Prometheus]
        H[Logs - Loki]
        I[Traces - Jaeger]
    end
    
    A --> B
    A --> C
    A --> D
    
    B --> E
    C --> E
    D --> E
    
    E --> F
    
    B --> G
    C --> G
    D --> G
    
    B --> H
    C --> H
    D --> H
    
    B --> I
    C --> I
    D --> I
    
    style E fill:#99ff99
    style F fill:#99ccff
```

---
Generated during comprehensive code quality improvement project
Claude (Sonnet 4) - 2025-07-26