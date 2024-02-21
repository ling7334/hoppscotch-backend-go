# hoppscotch-backend-go

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ling7334/hoppscotch-backend-go)
[![GitHub License](https://img.shields.io/github/license/ling7334/hoppscotch-backend-go)](https://github.com/ling7334/hoppscotch-backend-go/blob/master/LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/ling7334/hoppscotch)](https://hub.docker.com/r/ling7334/hoppscotch)


> A Go port of the popular API development tool, Hoppscotch.

## 🚀 Why 'hoppscotch-backend-go'?

With the aim of reduce the image size and resource usage, I've taken up the challenge to create a backend implementation in Go. This not only leverages the power and performance of Go but also contributes to the diversity of tools available for developers.

## 🛠️ Key Features:

- Go Power: Benefit from the speed and efficiency of Go.
- Compatibility: Designed to seamlessly integrate with the Hoppscotch frontend.
- Easy Setup: Simple installation and configuration for a Hoppscotch-like experience.

## Demo
[hoppscotch.ling.sh](hoppscotch.ling.sh)

## Developing
This project is build on `gqlgen` and `gorm`. please see related [documents](https://gqlgen.com/) for more information.

### Project structure
```
hoppscotch-backend-go
  ├── api
  │   ├── graphql               # graphql schema
  ├── cmd                       # main entrypoints
  │   ├── dummper               # database schema dumpper
  │   ├── server                # main server
  │   ├── import-meta-env       # go port of nodejs import-meta-env
  ├── internal
  │   ├── graphql               # graphql endpoints
  │   │     ├── dto             # dto models
  │   ├── rest                  # rest endpoints
  ├── pkg           
  │   ├── exception
  │   ├── mail
  │   ├── middleware
  │   ├── model                 # database models
  │   ├── models
  │   ├── session               # jwt token related
  ├── template                  # email templates
```

### Tips:
> `OwnersCount`, `EditorsCount` and `ViewersCount` alterd to return `int64` instead of `int`, because gorm return `int64` type in `count` function.


## Continuous Integration
We use GitHub Actions for continuous integration. Check out our build workflows.

### Build & Push
`import-meta-env` and main `server` are built using `go build` in `go_builder` image.

`nginx.conf`, `healthcheck.sh`, `template` from project root folder, `import-meta-env`, `server` from `go_builder` and `site` from `hoppscotch/hoppscotch:latest` are copied to `nginx:latest` image.

Port `3000/tcp`, `3100/tcp`, `3170/tcp` are exposed.

The product image are pushed to [DockerHub](https://hub.docker.com/r/ling7334/hoppscotch).


## 👩‍💻 How to Contribute:
Excited to contribute or explore further? Dive into the GitHub repository, check out the issues, and feel free to submit pull requests. This project is open to collaboration and welcomes your valuable contributions.

Please contribute using GitHub Flow. Create a branch, add commits, and open a pull request.

Please read CONTRIBUTING for details on our CODE OF CONDUCT, and the process for submitting pull requests to us.


## 🙏 Acknowledgments:
A big shoutout to the creators of Hoppscotch for inspiring this project, and to the open-source community for fostering a collaborative environment.

Let's build, test, and explore together! 🌟

Feel free to star the repository, share your thoughts, and let's make API development even more exciting and efficient!
