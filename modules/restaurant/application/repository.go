package application

import (
	"context"

	"github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

// Repository, application katmanının domain repository ile
// çalışabilmesi için kullanılan arayüzdür.
// İstersen doğrudan domain.Repository de kullanılabilir,
// ancak bu katman ayrımı ileride esneklik sağlar.
type Repository interface {
	Save(ctx context.Context, restaurant domain.Restaurant) error
	List(ctx context.Context) ([]domain.Restaurant, error)
}
