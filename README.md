
# ☕ Coffee Shop API

Une API REST pour gérer un coffee shop, développée en Go avec MySQL.

## 🚀 Démarrage rapide

### Interface

https://hellodamien.github.io/drink-ordering-app/

### Prérequis
- Docker & Docker Compose
- Go 1.23+ (optionnel, pour dev local)

### Installation

1. **Cloner le projet**
```bash
git clone <votre-repo>
cd coffee-shop-api
```

2. **Configurer les variables d'environnement**
```bash
cp .env.example .env
```

Modifiez `.env` avec vos valeurs :
```env
DB_HOST=mysql
DB_PORT=3306
DB_USER=coffee_user
DB_PASSWORD=your_password
DB_NAME=coffee_shop
DB_ROOT_PASSWORD=your_root_password
```

3. **Lancer l'application**
```bash
docker-compose up -d
```

L'API est accessible sur `http://localhost:8080`

## 📡 Endpoints

### Menu
- `GET /menu` - Récupérer toutes les boissons

### Boissons
- `GET /drinks/{id}` - Récupérer une boisson par ID

### Commandes
- `GET /orders` - Récupérer toutes les commandes
- `GET /orders/{id}` - Récupérer une commande par ID
- `POST /orders` - Créer une nouvelle commande
- `PATCH /orders/{id}/status` - Mettre à jour le statut d'une commande
- `DELETE /orders/{id}` - Supprimer une commande

## 🛠️ Développement

### Avec Docker (recommandé)
```bash
# Démarrer en mode dev avec hot reload
docker-compose up

# Voir les logs
docker-compose logs -f api

# Arrêter
docker-compose down
```

### En local
```bash
# Lancer juste MySQL en Docker
docker-compose up mysql -d

# Lancer l'API avec Air (hot reload)
air
```

## 🗄️ Base de données

### Accéder à MySQL
```bash
docker exec -it coffee-shop-db mysql -u coffee_user -p
```

### Tables
- `drinks` - Liste des boissons disponibles
- `orders` - Commandes des clients

## 📦 Structure du projet

```
coffee-shop-api/
├── database/          # Configuration base de données
├── handlers/          # Handlers des routes API
├── Dockerfile.api     # Dockerfile pour l'API
├── Dockerfile.mysql   # Dockerfile pour MySQL
├── docker-compose.yml # Orchestration des services
├── init.sql          # Script d'initialisation de la BDD
├── main.go           # Point d'entrée de l'application
└── .env              # Variables d'environnement (à créer)
```

## 🧪 Tests

```powershell
# Tester le menu
Invoke-RestMethod -Uri http://localhost:8080/menu

# Créer une commande
$body = @{
    drink_id = "1"
    customer_name = "John"
    size = "large"
    extras = @("milk", "sugar")
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/orders -Method POST -Body $body -ContentType "application/json"
```

## 🔧 Commandes utiles

```bash
# Rebuild les images
docker-compose up --build

# Voir les containers actifs
docker-compose ps

# Supprimer les volumes (reset BDD)
docker-compose down -v

# Logs d'un service spécifique
docker-compose logs -f mysql
docker-compose logs -f api
```

## 📝 License

MIT