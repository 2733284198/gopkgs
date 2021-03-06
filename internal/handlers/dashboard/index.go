package dashboard

import (
	"context"
	"net/http"
	"time"

	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/core"
	"pkg.razonyang.com/gopkgs/internal/helper"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
)

type Handler struct {
	core.Handler
}

func (h *Handler) Register(router clevergo.Router) {
	router.Get("/dashboard", h.index)
}

func (h *Handler) index(c *clevergo.Context) error {
	ctx := c.Context()
	user := h.User(ctx)
	userID := user.ID
	domainCard, err := h.getDomainCard(ctx, userID)
	if err != nil {
		return err
	}
	packageCard, err := h.getPackageCard(ctx, userID)
	if err != nil {
		return err
	}
	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		return err
	}
	dailyReportCard, err := h.getDailyReportCard(ctx, userID, loc)
	if err != nil {
		return err
	}
	monthlyReportCard, err := h.getMonthlyReportCard(ctx, userID, loc)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dashboard/index.tmpl", clevergo.Map{
		"page": web.NewPage("Dashboard"),
		"cards": []Card{
			domainCard,
			packageCard,
			dailyReportCard,
			monthlyReportCard,
		},
	})
}

func (h *Handler) getDomainCard(ctx context.Context, userID int64) (card Card, err error) {
	card = NewCard("DOMAINS", "globe", "primary", "/domain")
	err = models.CountDomainsByUser(ctx, h.DB, &card.Count, userID)
	return
}

func (h *Handler) getPackageCard(ctx context.Context, userID int64) (card Card, err error) {
	card = NewCard("PACKAGES", "cubes", "success", "/package")
	err = models.CountPackagesByUser(ctx, h.DB, &card.Count, userID)
	return
}

func (h *Handler) getDailyReportCard(ctx context.Context, userID int64, loc *time.Location) (card Card, err error) {
	card = NewCard("Daily Report", "download", "secondary", "/report")
	return card, h.getReport(ctx, &card.Count, userID, helper.CurrentUTC().In(loc))
}

func (h *Handler) getMonthlyReportCard(ctx context.Context, userID int64, loc *time.Location) (card Card, err error) {
	card = NewCard("Monthly Report", "download", "info", "/report")
	return card, h.getReport(ctx, &card.Count, userID, helper.CurrentUTC().In(loc).AddDate(0, 0, -29))
}

func (h *Handler) getReport(ctx context.Context, count *int64, userID int64, fromDate time.Time) error {
	query := `
SELECT COUNT(1) FROM actions 
LEFT JOIN packages ON packages.id = actions.package_id
LEFT JOIN domains ON domains.id = packages.domain_id
WHERE domains.user_id = ?
	AND actions.created_at >= ?
`
	return h.DB.GetContext(ctx, count, query, userID, fromDate.Format("2006-01-02"))
}

type Card struct {
	Background string
	Icon       string
	Title      string
	Count      int64
	Link       string
}

func NewCard(title, icon, background, link string) Card {
	return Card{
		Background: background,
		Icon:       icon,
		Title:      title,
		Link:       link,
	}
}
