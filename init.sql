-- Créer la table des boissons
CREATE TABLE IF NOT EXISTS drinks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    base_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Créer la table des commandes
CREATE TABLE IF NOT EXISTS orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    drink_id VARCHAR(50) NOT NULL,
    drink_name VARCHAR(100) NOT NULL,
    size VARCHAR(20) DEFAULT 'medium',
    extras TEXT,
    customer_name VARCHAR(100),
    status VARCHAR(50) DEFAULT 'pending',
    total_price DECIMAL(10, 2) NOT NULL,
    ordered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insérer des données de test pour les boissons
INSERT INTO drinks (name, category, base_price) VALUES 
    ('Espresso', 'coffee', 2.0),
    ('Cappuccino', 'coffee', 3.0),
    ('Latte', 'coffee', 3.5),
    ('Black Tea', 'tea', 2.5),
    ('Green Tea', 'tea', 2.5),
    ('Iced Coffee', 'cold', 3.0),
    ('Iced Tea', 'cold', 2.5);
