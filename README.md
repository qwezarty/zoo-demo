# zoo-demo

*A zoo-animal demo written in Go(gin/gorm).*

## introductions

I've implemented basic CRUD in two ways.

1. a native gin/gorm way as it described in their docs
2. a reusable pattern with a little Object-Oriented tricks

## native way

* apps/zoo/configure.go: router table situated here
* apps/zoo/handlers.go: Create, Get, Gets, Update, Remove handler entrance

## reusable pattern

* apps/apps.go: a parent CRUD handler entrance with a more general form
* apps/animal/configure.go: router table still situated here
* apps/animal/configure.go: you can override parent CRUD here

