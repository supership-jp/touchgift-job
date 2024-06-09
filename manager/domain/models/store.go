package models

type Store struct {
	ID               int    `db:"id" json:"id"`
	OrganizationCode string `db:"organization_code" json:"organization_code"`
	Name             string `db:"name" json:"name"`
}
