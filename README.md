# Intructions for running this really clunky poorly built but cool ass website

Someday I will refactor this to a proper frontend website but for now this is what it is :)

`it's like an avant-garde website - "kevin chou"`

## Local

### Raw Execution

Note: this app uses `dep` which has been deprecated. 

1. Clone the project into `$GOPATH/src`. Then run

```
dep ensure
```

2. Then run the following command and it should magically work!
```
go run app.go
```

### Docker Container

1. To build the container
```
docker build -t <name of container> .
```

2. You'll need to mount the aws folder into the container and pass an environment variable to load from the default aws config directory.

```
docker run -p 3000:3000 -v $HOME/.aws:/root/.aws -e AWS_SDK_LOAD_CONFIG=1 <name of container>
```

## Build Process

### Docker Hub

* Link to [docker image](https://hub.docker.com/repository/docker/ppdocx/ppvmio)

Pushes to git and tagging a release will trigger a docker build

### AWS

* [AWS Console Cluster](https://console.aws.amazon.com/ecs/home?region=us-east-1#/clusters/ppvmio-cluster)
