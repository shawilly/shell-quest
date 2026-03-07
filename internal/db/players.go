package db

import "time"

type Player struct {
	ID        int64
	Name      string
	Tier      string
	Avatar    string
	CreatedAt time.Time
}

func (d *DB) CreatePlayer(name, tier string) (*Player, error) {
	res, err := d.DB.Exec(
		`INSERT INTO players (name, tier) VALUES (?, ?)`, name, tier,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Player{ID: id, Name: name, Tier: tier, Avatar: "default"}, nil
}

func (d *DB) ListPlayers() ([]*Player, error) {
	rows, err := d.DB.Query(`SELECT id, name, tier, avatar, created_at FROM players ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []*Player
	for rows.Next() {
		p := &Player{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Tier, &p.Avatar, &p.CreatedAt); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, rows.Err()
}
