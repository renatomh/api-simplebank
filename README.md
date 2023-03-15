<h1 align="center"><img alt="API Simple Bank" title="API Simple Bank" src="https://go.dev/images/go-logo-blue.svg" width="250" /></h1>

# API Simple Bank

## 💡 Project's Idea

This project was developed while studying Go, Docker and Kubernetes. It aims to create a simple bank API.

## 🔍 Features

* Create new accounts with different currencies;
* Transfer money from your accounts to another ones;

## 🛠 Technologies

During the development of this project, the following techologies were used:

- [Go](https://go.dev/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Viper](https://github.com/spf13/viper)
- [gomock](https://github.com/golang/mock)
- [Docker](https://www.docker.com/)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [sqlc](https://sqlc.dev/)
- [gorilla/mux](https://github.com/gorilla/mux)

## 💻 Project Configuration

### First, you must install these tools:

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go](https://go.dev/dl/)
- [scoop](https://scoop.sh/) for Windows or [Homebrew](https://brew.sh/) for Mac
- [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [choco](https://chocolatey.org/install) and [choco install make](https://stackoverflow.com/questions/2532234/how-to-run-a-makefile-in-windows) for Windows
- [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html) (for WIndows, we'll use it on a Docker image)

### Setting up infrastructure:

Start postgres container, creating database and running database migration up all versions:

```bash
$ make postgres
$ make createdb
$ make migrateup
```

## 🌐 Setting up config files

In order to build docker iamges from a Dockerfile, we can do it with the following command (if we already have a *Dockerfile* on the directory where the commnd is being executed):

```bash
$ docker build -t simplebank:latest .
```

We can create docker networks with the following command:

```bash
$ docker network create bank-network
```

## ⏯️ Running

To run the application locally, you can use the following command:

```bash
$ make server
```

## 🔨 Project's *Deploy*

In order to deploy the application to the kubernetes cluster ...

### Documentation:
* [A Tour of Go](https://go.dev/tour/welcome/1)
* [postgres | Docker Hub](https://hub.docker.com/_/postgres)
* [How to run a makefile in Windows?](https://stackoverflow.com/questions/2532234/how-to-run-a-makefile-in-windows)
* [Docker tries to mkdir the folder that I mount](https://stackoverflow.com/questions/50817985/docker-tries-to-mkdir-the-folder-that-i-mount)
* [NGINX Ingress Controller - Installation Guide # AWS](https://kubernetes.github.io/ingress-nginx/deploy/#aws)

## 📄 License

This project is under the **MIT** license. For more information, access [LICENSE](./LICENSE).
