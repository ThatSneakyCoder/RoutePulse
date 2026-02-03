# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### K8s Config ###

# Uncomment to use secrets
k8s_yaml('./infra/development/kubernetes/secrets.yaml')
k8s_yaml('./infra/development/kubernetes/app-config.yaml')

### End of K8s Config ###

### RabbitMQ ###

# k8s_yaml('./infra/development/k8s/rabbitmq-deployment.yaml')
# k8s_resource('rabbitmq', port_forwards=['5672', '15672'], labels='tooling')

### rmq end ###

### Postgresql start ###

k8s_yaml('./infra/development/kubernetes/pg.yaml')
k8s_resource('postgres', port_forwards=['15432:5432'], labels='tooling')

### Postgresql end ###

# ### Migration Job Start ###

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


# ### Migration Job End ###

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


### Trip Service ###

# Uncomment once we have a trip service

# trip_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/trip-service ./services/trip-service/cmd/main.go'
# if os.name == 'nt':
#  trip_compile_cmd = './infra/development/docker/trip-build.bat'

# local_resource(
#   'trip-service-compile',
#   trip_compile_cmd,
#   deps=['./services/trip-service', './shared'], labels="compiles")

# docker_build_with_restart(
#   'ride-sharing/trip-service',
#   '.',
#   entrypoint=['/app/build/trip-service'],
#   dockerfile='./infra/development/docker/trip-service.Dockerfile',
#   only=[
#     './build/trip-service',
#     './shared',
#   ],
#   live_update=[
#     sync('./build', '/app/build'),
#     sync('./shared', '/app/shared'),
#   ],
# )

# k8s_yaml('./infra/development/k8s/trip-service-deployment.yaml')
# k8s_resource('trip-service', resource_deps=['trip-service-compile', 'rabbitmq'], labels="services")

# ### End of Trip Service ###

# ### Driver Service ###

# # Uncomment once we have a Driver service

# driver_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/driver-service ./services/driver-service'
# if os.name == 'nt':
#  driver_compile_cmd = './infra/development/docker/driver-build.bat'

# local_resource(
#   'driver-service-compile',
#   driver_compile_cmd,
#   deps=['./services/driver-service', './shared'], labels="compiles")

# docker_build_with_restart(
#   'ride-sharing/driver-service',
#   '.',
#   entrypoint=['/app/build/driver-service'],
#   dockerfile='./infra/development/docker/driver-service.Dockerfile',
#   only=[
#     './build/driver-service',
#     './shared',
#   ],
#   live_update=[
#     sync('./build', '/app/build'),
#     sync('./shared', '/app/shared'),
#   ],
# )

# k8s_yaml('./infra/development/k8s/driver-service-deployment.yaml')
# k8s_resource('driver-service', resource_deps=['driver-service-compile', 'rabbitmq'], labels="services")

# ### End of driver Service ###

# ### Web Frontend ###

# docker_build(
#   'ride-sharing/web',
#   '.',
#   dockerfile='./infra/development/docker/web.Dockerfile',
# )

# k8s_yaml('./infra/development/k8s/web-deployment.yaml')
# k8s_resource('web', port_forwards=3000, labels="frontend")

# ### End of Web Frontend ###