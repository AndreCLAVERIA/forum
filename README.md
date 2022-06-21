# Forum in golang

## How to start the project ?

**If you use Dockerfile run those following commands in your terminal** :

    1. docker build -t forum .
    2. docker run -p 8081:8080 -it forum

8081 is the port which your PC will listen
8080 is the port where docker will serve

**If you want to use just commands line in your terminal** :

    - cd forum
    - go run main.go
    - ctrl + click on http://localhost:8080

PS: Even if the directory code say it has problem, it doesnt just go run or go build it will works

## What is the subject of the project ?

![picture](https://raw.githubusercontent.com/AndreClaveria/forum/master/img/accueil.png)
This project consists in creating a web forum that allows :

    - communication between users.
    - associating categories to posts.
    - liking and disliking posts and comments.
    - filtering posts.

Which also have 5 additional modules :

    - image-upload
    - security
    - authentication
    - advanced-features
    - moderation

### Competences that we need

    **SQLite - Database**
    **Encrypting**
    **Dockerize**
    **HTML and CSS advanced features**

# Groupe 9 :  4Head Forum - By Sacha PERRIN, Damien COURTEAUX, André CLAVERIA, Noé MOYEN, Maxime ROADLEY-BATTIN
