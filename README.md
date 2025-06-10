# 🏄🏼‍♀️ SurfSpot Backend - API REST en Go avec Gin 

Ce dépôt contient le backend de l'application **SurfSpot**, développé en **Go** avec le framework **Gin**, connecté à une base de données PostgreSQL hébergée sur **Neon**.

## Fonctionnalités pricipales 
- API REST pour gérer les spots de surf
- Gestion des erreurs et des statuts HTTP

## Technologies utilisées

- [Go](https://golang.org/)
- [Gin Gonic](https://github.com/gin-gonic/gin)
- [PostgreSQL (Neon)](https://neon.tech/)

## Configugation (.env)
Avant de lancer l'application, créer un fichier `.env`à la racine du dossier back-end avec les variables suivantes:

- DB="postgres://user:password@localhost:5432/surfspot"
- PORT=4000
- SECRET=exemple1234
