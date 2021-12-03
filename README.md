# What is PPVMIO?

When I was first introduced to AWS in Fall of 2019, I was determined to setup some basic infrastructure. To better keep track of all the infrastructure I began brainstorming a potential application or namespace that I could use to label everything. I settled on PPVMIO which you can break down into three parts

1. PP - my initials
2. VM - virtual machine; I was initially worried about labeling EC2 instances
3. IO - input/ouput, also a comain top level domain for tech companies 

Eventually the set of characters become a moniker that I use frequently across the internet and across all of my various projects. I combine both my passions for art/design/music with my software engineer skills. From a creative perspective, my projects fall under the category of [digital brutalism](https://brutalistwebsites.com/). From an engineering perspective, I build modern applications and use various tools to help me deploy and manage my projects in the cloud. 

Below you'll find the true README instructions for my personal website [ppvm.io](https://ppvm.io/)

I hope you enjoy!

# README

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

## Build & Deploy

### Docker Hub

Link to [docker image](https://hub.docker.com/repository/docker/ppdocx/ppvmio)

**Note**: dockerhub no longer supports automated docker builds on pushes to github. To build a new image

```
docker build -t ppdocx/ppvmio:<VERSION TAG> .
```

and to push

```
docker push ppdocx/ppvmio:<VERSION TAG>
```

### AWS

**Update Task**

To deploy the new tag to the [fargate cluster](https://console.aws.amazon.com/ecs/home?region=us-east-1#/clusters/ppvmio-cluster)

1. Navigate to the UI and [create a new task revision](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/update-task-definition.html).
2. Scroll down to where the image name is and click on it
3. Edit the image to have the proper tag and click `Update`
4. Click `Create`
5. Go back to the cluster and Update the service to have the new revision.

*Note*: you may have to stop the previous task before the new one will be allowed to run.

To refresh the cloudfront catche run

```
aws cloudfront create-invalidation --distribution-id E20PNAIG7SWPHU --paths "/*"
```

**Force New Deployment**

To force a new deployment run

```
aws ecs update-service --cluster ppvmio-cluster --service ppvmio-prod --force-new-deployment
```

### Photo Resizing
You can use ffmpeg to resize photos before uploading to s3 in order to compress size of photos + speed up page load time. Run the following command to convert all `*.jpg` photos in a folder to size 768x515 (about 3:2 ratio).

```
for i in *.jpg; do ffmpeg -i "$i" -vf scale=768:515 "output_$i"; done
```


## Improvements
* Refactor to be a static frontend web application, too expensive to run as fargate
* Maybe use something like github actions to automate builds
* Figure out easy way to deploy new task to service