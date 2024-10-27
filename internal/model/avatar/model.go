package avatar

import "time"

type Info struct {
	ID        int       `db:"id"`
	UserUUID  string    `db:"user_uuid"`
	Link      string    `db:"link"`
	CreatedAt time.Time `db:"create_at"`
}
