package tests

import (
	"database/sql"
	"time"

	"github.com/driver005/database"
)

// User has one `Account` (has one), many `Pets` (has many) and `Toys` (has many - polymorphic)
// He works in a Company (belongs to), he has a Manager (belongs to - single-table), and also managed a Team (has many - single-table)
// He speaks many languages (many to many) and has many friends (many to many - single-table)
// His pet also has one Toy (has one - polymorphic)
// NamedPet is a reference to a Named `Pets` (has many)
type User struct {
	database.Model
	Name      string
	Age       uint
	Birthday  *time.Time
	Account   Account
	Pets      []*Pet
	NamedPet  *Pet
	Toys      []Toy `database:"polymorphic:Owner"`
	CompanyID *int
	Company   Company
	ManagerID *uint
	Manager   *User
	Team      []User     `database:"foreignkey:ManagerID"`
	Languages []Language `database:"many2many:UserSpeak;"`
	Friends   []*User    `database:"many2many:user_friends;"`
	Active    bool
}

type Account struct {
	database.Model
	UserID sql.NullInt64
	Number string
}

type Pet struct {
	database.Model
	UserID *uint
	Name   string
	Toy    Toy `database:"polymorphic:Owner;"`
}

type Toy struct {
	database.Model
	Name      string
	OwnerID   string
	OwnerType string
}

type Company struct {
	ID   int
	Name string
}

type Language struct {
	Code string `database:"primarykey"`
	Name string
}

type Coupon struct {
	ID               int              `database:"primarykey; size:255"`
	AppliesToProduct []*CouponProduct `database:"foreignKey:CouponId;constraint:OnDelete:CASCADE"`
	AmountOff        uint32           `database:"amount_off"`
	PercentOff       float32          `database:"percent_off"`
}

type CouponProduct struct {
	CouponId  int    `database:"primarykey;size:255"`
	ProductId string `database:"primarykey;size:255"`
	Desc      string
}

type Order struct {
	database.Model
	Num      string
	Coupon   *Coupon
	CouponID string
}

type Parent struct {
	database.Model
	FavChildID uint
	FavChild   *Child
	Children   []*Child
}

type Child struct {
	database.Model
	Name     string
	ParentID *uint
	Parent   *Parent
}
