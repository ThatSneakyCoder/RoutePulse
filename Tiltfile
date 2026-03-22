# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### K8s Config ###

# SECRETS and CONFIGMAP
k8s_yaml('./infra/development/kubernetes/secrets/secrets-identity.yaml')
k8s_yaml('./infra/development/kubernetes/secrets/secrets-organization.yaml')
k8s_yaml('./infra/development/kubernetes/secrets/secrets-fleet.yaml')
k8s_yaml('./infra/development/kubernetes/secrets/secrets-rabbitmq.yaml')
k8s_yaml('./infra/development/kubernetes/secrets/secrets-clickhouse.yaml')
k8s_yaml('./infra/development/kubernetes/configmaps/configmap.yaml')
k8s_yaml('./infra/development/kubernetes/configmaps/configmap-ch.yaml')
k8s_yaml('./infra/development/kubernetes/configmaps/configmap-org.yaml')
k8s_yaml('./infra/development/kubernetes/configmaps/configmap-fleet.yaml')
k8s_yaml('./infra/development/kubernetes/configmaps/configmap-tracking.yaml')

### End of K8s Config ###

### RabbitMQ ###

k8s_yaml('./infra/development/kubernetes/rabbitmq/rabbitmq.yaml')
k8s_resource('rabbitmq', port_forwards=['5672', '15672'], labels='tooling')

### RabbitMQ end ###

### Identity Postgresql start ###

k8s_yaml('./infra/development/kubernetes/pg.yaml')
k8s_resource('postgres', port_forwards=['15432:5432'], labels='tooling')

### Identity Postgresql end ###

### Organization Postgresql start ###

k8s_yaml('./infra/development/kubernetes/organization/postgres-organization.yaml')
k8s_resource('postgres-org', port_forwards=['15433:5432'], labels='tooling')

### Organization Postgresql end ###

### Fleet Postgresql start ###

k8s_yaml('./infra/development/kubernetes/fleet/fleet-pg.yaml')
k8s_resource('postgres-fleet', port_forwards=['15434:5432'], labels='tooling')

### Fleet Postgresql end ###

### clickhouse start ###

k8s_yaml('./infra/development/kubernetes/clickhouse.yaml')
k8s_resource('clickhouse', port_forwards=['18123:8123'], labels='tooling')

### clickhouse end ###

### tracking service postgres Migration Job Start ###

# k8s_yaml('./infra/development/kubernetes/tracking/tracking-migrate-job.yaml')

# k8s_resource(
#   'tracking-db-migrate',
#   resource_deps=['postgres-org', 'tracking-service-compile'],
#   labels='migrations'
# )

# docker_build(
#   'routepulse/tracking-migrate',
#   '.',
#   dockerfile='./infra/development/docker/tracking-migrate.Dockerfile',
#   only=[
#     './services/tracking-service/internal/migrate/migrations',
#   ],
# )

### tracking service postgres Migration Job End ###

### organization service postgres Migration Job Start ###

k8s_yaml('./infra/development/kubernetes/organization/organization-migrate-job.yaml')

k8s_resource(
  'organization-db-migrate',
  resource_deps=['postgres-org', 'organization-service-compile'],
  labels='migrations'
)

docker_build(
  'routepulse/organization-migrate',
  '.',
  dockerfile='./infra/development/docker/organization-migrate.Dockerfile',
  only=[
    './services/organization-service/internal/migrate/migrations',
  ],
)

### organization service postgres Migration Job End ###

### fleet service postgres Migration Job Start ###

k8s_yaml('./infra/development/kubernetes/fleet/fleet-migrate-job.yaml')

k8s_resource(
  'fleet-db-migrate',
  resource_deps=['postgres-fleet', 'fleet-service-compile'],
  labels='migrations'
)

docker_build(
  'routepulse/fleet-migrate',
  '.',
  dockerfile='./infra/development/docker/fleet-migrate.Dockerfile',
  only=[
    './services/fleet-service/internal/migrate/migrations',
  ],
)

### fleet service postgres Migration Job End ###

### identity service postgres Migration Job Start ###

k8s_yaml('./infra/development/kubernetes/identity-migrate-job.yaml')

k8s_resource(
  'identity-db-migrate',
  resource_deps=['postgres', 'identity-service-compile'],
  labels='migrations'
)

docker_build(
  'routepulse/identity-migrate',
  '.',
  dockerfile='./infra/development/docker/identity-migrate.Dockerfile',
  only=[
    './services/identity-service/internal/migrate/migrations',
    # './tools/migrate',
  ],
)

### identity service postgres Migration Job End ###

### clickhouse Migration Job Start ###

k8s_yaml('./infra/development/kubernetes/analytics-migrate-job.yaml')

k8s_resource(
  'analytics-db-migrate',
  resource_deps=['clickhouse', 'analytics-service-compile'],
  labels='migrations'
)

docker_build(
  'routepulse/analytics-migrate',
  '.',
  dockerfile='./infra/development/docker/analytics-migrate.Dockerfile',
  only=[
    './services/analytics-service/internal/migrate/migrations',
  ],
)

### clickhouse Migration Job End ###

### API Gateway ###

gateway_compile_cmd = '''
cd services/api-gateway && \
swag init \
  -g main.go \
  -d cmd \
  -o docs && \
cd ../.. && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o build/api-gateway ./services/api-gateway/cmd
'''

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway', './shared'], ignore=['./services/api-gateway/docs'], labels="compiles")


docker_build_with_restart(
  'routepulse/api-gateway',
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./infra/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/kubernetes/api-gateway-deployment.yaml')
k8s_resource(
  'api-gateway', 
  port_forwards=8080,
  resource_deps=['api-gateway-compile'], labels="services",
  links="http://localhost:8080/v1/health"
)

### End of API Gateway ###

### Tracking Service start ###

tracking_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/tracking-service ./services/tracking-service/cmd'

local_resource(
  'tracking-service-compile',
  tracking_compile_cmd,
  deps=['./services/tracking-service', './shared'], labels="compiles")


docker_build_with_restart(
  'routepulse/tracking-service',
  '.',
  entrypoint=['/app/build/tracking-service'],
  dockerfile='./infra/development/docker/tracking-service.Dockerfile',
  only=[
    './build/tracking-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/kubernetes/tracking/tracking-service-deployment.yaml')
k8s_resource(
  'tracking-service', 
  port_forwards=9093,
  resource_deps=[
    # 'postgres', 
    'tracking-service-compile'], 
    labels="services"
)

### Tracking Service end ###

### Identity Service start ###

identity_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/identity-service ./services/identity-service/cmd'

local_resource(
  'identity-service-compile',
  identity_compile_cmd,
  deps=['./services/identity-service', './shared'], labels="compiles")


docker_build_with_restart(
  'routepulse/identity-service',
  '.',
  entrypoint=['/app/build/identity-service'],
  dockerfile='./infra/development/docker/identity-service.Dockerfile',
  only=[
    './build/identity-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/kubernetes/identity-service-deployment.yaml')
k8s_resource(
  'identity-service', 
  port_forwards=9090,
  resource_deps=['postgres', 'identity-service-compile'], labels="services"
)

### Identity Service end ###

### Organization Service start ###

identity_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/organization-service ./services/organization-service/cmd'

local_resource(
  'organization-service-compile',
  identity_compile_cmd,
  deps=['./services/organization-service', './shared'], labels="compiles")


docker_build_with_restart(
  'routepulse/organization-service',
  '.',
  entrypoint=['/app/build/organization-service'],
  dockerfile='./infra/development/docker/organization-service.Dockerfile',
  only=[
    './build/organization-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/kubernetes/organization/organization-service-deployment.yaml')
k8s_resource(
  'organization-service', 
  port_forwards=9091,
  resource_deps=['postgres-org', 'organization-service-compile'], labels="services"
)

### Organization Service end ###

### Fleet Service start ###

identity_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/fleet-service ./services/fleet-service/cmd'

local_resource(
  'fleet-service-compile',
  identity_compile_cmd,
  deps=['./services/fleet-service', './shared'], labels="compiles")


docker_build_with_restart(
  'routepulse/fleet-service',
  '.',
  entrypoint=['/app/build/fleet-service'],
  dockerfile='./infra/development/docker/fleet-service.Dockerfile',
  only=[
    './build/fleet-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/kubernetes/fleet/fleet-service-deployment.yaml')
k8s_resource(
  'fleet-service', 
  port_forwards=9092,
  resource_deps=['postgres-fleet', 'fleet-service-compile'], labels="services"
)

### Fleet Service end ###

### Analytics Service start ###

identity_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/analytics-service ./services/analytics-service/cmd'

local_resource(
  'analytics-service-compile',
  identity_compile_cmd,
  deps=['./services/analytics-service', './shared'], labels="compiles")


docker_build_with_restart(
  'routepulse/analytics-service',
  '.',
  entrypoint=['/app/build/analytics-service'],
  dockerfile='./infra/development/docker/analytics-service.Dockerfile',
  only=[
    './build/analytics-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/kubernetes/analytics-service-deployment.yaml')
k8s_resource(
  'analytics-service', 
  port_forwards=9096,
  resource_deps=['clickhouse', 'analytics-service-compile'], labels="services"
)

### Analytics Service end ###

### Web Frontend ###

docker_build(
  'routepulse/frontend/web',
  './frontend/web',
  dockerfile='./infra/development/docker/web.Dockerfile',
)

k8s_yaml('./infra/development/kubernetes/web-deployment.yaml')
k8s_resource('web', port_forwards=['5173:5173'], labels="frontend")

### End of Web Frontend ###
