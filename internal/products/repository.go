package products

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, p *Product) error
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, p *Product) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const productQuery = `
		INSERT INTO products (id, sku, name, slug, description, price, stock, category_id, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(ctx, productQuery,
		p.ID, p.SKU, p.Name, p.Slug, p.Description, p.Price, p.Stock, p.CategoryID, p.IsFeatured,
	)
	if err != nil {
		return err
	}

	if len(p.Images) > 0 {
		const imageQuery = `
			INSERT INTO products_images (id, product_id, path, original_url, thumbnail_url, is_primary)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

		for _, img := range p.Images {
			_, err = tx.Exec(ctx, imageQuery,
				img.ID, p.ID, img.Path, img.Original, img.Thumbnail, img.IsPrimary,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]*Product, error) {
	const query = `
		SELECT p.id, p.sku, p.name, p.slug, p.description, p.price, p.stock,
			p.category_id, c.name as category_name, p.is_featured, p.created_at
		FROM products p
		JOIN categories c ON p.category_id = c.id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		p := &Product{}
		err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Stock,
			&p.CategoryID, &p.CategoryName, &p.IsFeatured, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	return nil, nil
}

func (r *PostgresRepository) Update(ctx context.Context, p *Product) error {
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
