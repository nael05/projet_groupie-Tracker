# Groupie Tracker

**Projet scolaire collaboratif**
Ce projet a été réalisé en groupe dans le cadre de nos études en informatique à Ynov Campus.

## Présentation du projet
Groupie Tracker est une application de bureau que nous avons développée en Go (Golang) à l'aide du toolkit graphique Fyne. Elle interagit avec une API REST externe pour afficher des informations sur divers artistes et groupes de musique (historique, membres, dates et lieux de concerts).

L'objectif de ce projet était de nous familiariser avec la manipulation de données, la conception d'interface utilisateur en Go, et la gestion des requêtes API asynchrones (y compris le géocodage et la visualisation de cartes).

## Fonctionnalités
- **Annuaire d'artistes** : Affichage des artistes sous forme de grille avec leurs noms et images.
- **Recherche avancée** : Nous avons implémenté une barre de recherche permettant de filtrer par nom d'artiste, nom de membre, date de création, année du premier album ou lieu de concert.
- **Système de filtres** : Il est possible d'affiner la liste selon l'année de création, le nombre de membres, etc.
- **Vue détaillée** : Affiche les informations complètes d'un artiste sélectionné (membres, albums, concerts).
- **Géolocalisation & Carte** :
  - L'application convertit les noms de lieux en coordonnées géographiques via l'API Nominatim.
  - Nous avons intégré une carte interactive (OpenStreetMap) pour visualiser les lieux des concerts.

## Stack Technique
- **Langage** : Go (Golang)
- **Framework GUI** : Fyne v2
- **Format de données** : JSON
- **APIs externes** :
  - Données des artistes (API fournie dans le sujet)
  - Géocodage : Nominatim (OpenStreetMap)
  - Tuiles de carte : OpenStreetMap

## Arborescence du Projet
```
projet_groupie-Tracker/
├── main.go            # Point d'entrée de l'application
├── appli/             # Logique centrale et composants UI
│   ├── api.go         # Gestion des requêtes HTTP et structures de données
│   └── page.go        # Layout Fyne, événements et rendu de la carte
├── go.mod / go.sum    # Dépendances du module Go
└── README.md          # Documentation
```

## Prérequis
Pour lancer ce projet, vous aurez besoin de :
1. **Go** (version 1.25 ou compatible).
2. **Un compilateur C** : Fyne nécessite un compilateur C (GCC) pour les bindings CGO liés au rendu graphique.
   - *Windows* : TDM-GCC ou MinGW-w64.
   - *macOS* : Xcode Command Line Tools.
   - *Linux* : GCC.

## Installation et Lancement
1. Clonez ce dépôt.
2. Ouvrez un terminal dans le dossier racine.
3. Téléchargez les dépendances :
   ```bash
   go mod tidy
   ```
4. Lancez l'application :
   ```bash
   go run main.go
   ```

### 👥 Contributeurs
- Groupe d'étudiants Ynov Campus (Projet Go).
