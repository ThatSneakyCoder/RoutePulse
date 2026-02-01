# Core services & data ownership
Who talks to whom synchronously, and who owns which data?
```mermaid
flowchart LR
    %% =====================
    %% Actors
    %% =====================
    User["End User
(Owner / Operations Manager / Dispatcher)
- Uses Web / Frontend
- Makes HTTP requests"]

    %% =====================
    %% Edge Layer
    %% =====================
    APIGW["API Gateway Service (HTTP)
- JWT authentication
- RBAC permission checks
- Routes requests to backend services
- No database access"]

    %% =====================
    %% Core Backend Services
    %% =====================
    Identity["Identity Service (gRPC)
- User signup / login
- Password hashing
- JWT generation
- Emits auth-related events"]

    Org["Organization Service (gRPC)
- Organizations (companies)
- Memberships
- Role assignments
- Invites users to orgs"]

    Fleet["Fleet Service (gRPC)
- Vehicles
- Drivers
- Vehicle lifecycle (active / inactive)"]

    %% =====================
    %% Databases (owned per service)
    %% =====================
    IdentityDB["PostgreSQL
Identity Schema
- users
- credentials"]

    OrgDB["PostgreSQL
Organization Schema
- organizations
- memberships
- roles"]

    FleetDB["PostgreSQL
Fleet Schema
- vehicles
- drivers"]

    %% =====================
    %% Request Flow
    %% =====================
    User -->|HTTP| APIGW

    APIGW -->|gRPC| Identity
    APIGW -->|gRPC| Org
    APIGW -->|gRPC| Fleet

    %% =====================
    %% Data Ownership
    %% =====================
    Identity -->|owns & writes| IdentityDB
    Org -->|owns & writes| OrgDB
    Fleet -->|owns & writes| FleetDB

```

# Authorization & RBAC (OpenFGA only)
How permissions are written and checked
```mermaid
flowchart LR
    APIGW["API Gateway
(RBAC Enforcement)"]

    Identity["Identity Service"]
    Org["Organization Service"]

    OpenFGA["OpenFGA Authorization Engine
(RBAC Policy Store)"]

    Identity -->|write user-role relations| OpenFGA
    Org -->|write org-membership relations| OpenFGA
    APIGW -->|check permission| OpenFGA
```

# Async events & RabbitMQ
Who publishes events, who consumes them, and why
```mermaid
flowchart LR
    Identity["Identity Service"]
    Fleet["Fleet Service"]
    Tracking["Tracking Service"]
    Trip["Trip Service"]

    MQ["RabbitMQ<br/>Event Broker"]

    Alert["Alert Service<br/>(Async)"]
    Analytics["Analytics Service<br/>(Async)"]
    Notification["Notification Service<br/>(Async)"]

    Identity -->|UserCreated| MQ
    Fleet -->|VehicleRegistered| MQ
    Tracking -->|LocationUpdated| MQ
    Trip -->|TripCompleted| MQ

    MQ --> Alert
    MQ --> Analytics
    MQ --> Notification

```