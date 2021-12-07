package rebrickable

import (
	"fmt"
	"time"
)

type lClient struct {
	*Client
}

func NewLEGOClient(apiKey string, opts ...ClientOption) *lClient {
	return &lClient{
		NewClient(apiKey, opts...),
	}
}

func (c *lClient) endpoint(endpoint string, a ...interface{}) string {
	return fmt.Sprintf("lego/%v", fmt.Sprintf(endpoint, a...))
}

// Colors get a list of all Color.
func (c *lClient) Colors(opts ...RequestOption) (colors []Color, err error) {
	err = c.GetDecode("lego/colors/", true, &colors, opts...)
	return
}

// Color get details about a specific Color.
func (c *lClient) Color(id int, opts ...RequestOption) (color Color, err error) {
	err = c.GetDecode(c.endpoint("colors/%v", id), false, &color, opts...)
	return
}

type Color struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Rgb         string `json:"rgb"`
	IsTrans     bool   `json:"is_trans"`
	ExternalIds struct {
		BrickLink struct {
			ExtIds    []int      `json:"ext_ids"`
			ExtDescrs [][]string `json:"ext_descrs"`
		} `json:"BrickLink"`
		BrickOwl struct {
			ExtIds    []int      `json:"ext_ids"`
			ExtDescrs [][]string `json:"ext_descrs"`
		} `json:"BrickOwl"`
		Lego struct {
			ExtIds    []int      `json:"ext_ids"`
			ExtDescrs [][]string `json:"ext_descrs"`
		} `json:"LEGO"`
		Peeron struct {
			ExtIds    []interface{} `json:"ext_ids"`
			ExtDescrs [][]string    `json:"ext_descrs"`
		} `json:"Peeron"`
		LDraw struct {
			ExtIds    []int      `json:"ext_ids"`
			ExtDescrs [][]string `json:"ext_descrs"`
		} `json:"LDraw"`
	} `json:"external_ids"`
}

// Element get details about a specific Element ID.
func (c *lClient) Element(id string) (element Element, err error) {
	err = c.GetDecode(c.endpoint("elements/%v", id), false, &element)
	return
}

type Element struct {
	Part struct {
		PartNum     string        `json:"part_num"`
		Name        string        `json:"name"`
		PartCatID   int           `json:"part_cat_id"`
		YearFrom    int           `json:"year_from"`
		YearTo      int           `json:"year_to"`
		PartURL     string        `json:"part_url"`
		PartImgURL  string        `json:"part_img_url"`
		Prints      []interface{} `json:"prints"`
		Molds       []interface{} `json:"molds"`
		Alternates  []interface{} `json:"alternates"`
		ExternalIds struct {
			BrickLink []string `json:"BrickLink"`
			BrickOwl  []string `json:"BrickOwl"`
			Brickset  []string `json:"Brickset"`
			Lego      []string `json:"LEGO"`
		} `json:"external_ids"`
		PrintOf string `json:"print_of"`
	} `json:"part"`
	Color struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Rgb         string `json:"rgb"`
		IsTrans     bool   `json:"is_trans"`
		ExternalIds struct {
			BrickLink struct {
				ExtIds    []int      `json:"ext_ids"`
				ExtDescrs [][]string `json:"ext_descrs"`
			} `json:"BrickLink"`
			BrickOwl struct {
				ExtIds    []int      `json:"ext_ids"`
				ExtDescrs [][]string `json:"ext_descrs"`
			} `json:"BrickOwl"`
			Lego struct {
				ExtIds    []int      `json:"ext_ids"`
				ExtDescrs [][]string `json:"ext_descrs"`
			} `json:"LEGO"`
			Peeron struct {
				ExtIds    []interface{} `json:"ext_ids"`
				ExtDescrs [][]string    `json:"ext_descrs"`
			} `json:"Peeron"`
			LDraw struct {
				ExtIds    []int      `json:"ext_ids"`
				ExtDescrs [][]string `json:"ext_descrs"`
			} `json:"LDraw"`
		} `json:"external_ids"`
	} `json:"color"`
	ElementID     string `json:"element_id"`
	DesignID      string `json:"design_id"`
	ElementImgURL string `json:"element_img_url"`
	PartImgURL    string `json:"part_img_url"`
}

// Minifigs get a list of Minifig.
func (c *lClient) Minifigs(opts ...RequestOption) (minifigs []Minifig, err error) {
	err = c.GetDecode("lego/minifigs", true, &minifigs, opts...)
	return
}

// Minifig get details for a specific Minifig.
func (c *lClient) Minifig(setNumber string) (minifig Minifig, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/minifigs/%v", setNumber), false, &minifig)
	return
}

// MinifigParts get a list of all inventory Part\s in this Minifig.
func (c *lClient) MinifigParts(setNumber string, opts ...RequestOption) (parts []Part, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/minifigs/%v/parts", setNumber), true, &parts, opts...)
	return
}

// MinifigSets get a list of Set a Minifig has appeared in.
func (c *lClient) MinifigSets(setNumber string, opts ...RequestOption) (sets []Set, err error) {
	err = c.GetDecode(c.endpoint("minifigs/%v/sets", setNumber), true, &sets, opts...)
	return
}

type Minifig struct {
	SetNum         string    `json:"set_num"`
	Name           string    `json:"name"`
	NumParts       int       `json:"num_parts"`
	SetImgURL      string    `json:"set_img_url"`
	SetURL         string    `json:"set_url"`
	LastModifiedDt time.Time `json:"last_modified_dt"`
}

// PartCategories get a list of all PartCategory.
func (c *lClient) PartCategories(opts ...RequestOption) (partCategories []PartCategory, err error) {
	err = c.GetDecode(c.endpoint("part_categories"), true, &partCategories, opts...)
	return
}

// PartCategory get details about a specific PartCategory.
func (c *lClient) PartCategory(id int, opts ...RequestOption) (partCategory PartCategory, err error) {
	err = c.GetDecode(c.endpoint("part_categories/%v", id), false, &partCategory, opts...)
	return
}

type PartCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	PartCount int    `json:"part_count"`
}

// Parts get a list of Part.
func (c *lClient) Parts(opts ...RequestOption) (parts []Part, err error) {
	err = c.GetDecode("lego/parts", true, &parts, opts...)
	return
}

// Part get details about a specific Part.
func (c *lClient) Part(partNumber string) (part Part, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/parts/%v", partNumber), false, &part)
	return
}

// PartColors get a list of all Color a Part has appeared in.
func (c *lClient) PartColors(partNumber string, opts ...RequestOption) (colors []Color, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/parts/%v/colors", partNumber), true, &colors, opts...)
	return
}

// PartColor get details about a specific Part Color combination.
func (c *lClient) PartColor(partNumber string, colorId int) (partColor PartColor, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/parts/%v/colors/%v", partNumber, colorId), false, &partColor)
	return
}

type PartColor struct {
	PartImgURL  string   `json:"part_img_url"`
	YearFrom    int      `json:"year_from"`
	YearTo      int      `json:"year_to"`
	NumSets     int      `json:"num_sets"`
	NumSetParts int      `json:"num_set_parts"`
	Elements    []string `json:"elements"`
}

// PartColorSets get a list of all Set the Part Color combination has appeared in.
func (c *lClient) PartColorSets(partNumber string, colorId int, opts ...RequestOption) (sets []Set, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/parts/%v/colors/%v/sets", partNumber, colorId), true, &sets, opts...)
	return
}

type Part struct {
	PartNum     string        `json:"part_num"`
	Name        string        `json:"name"`
	PartCatID   int           `json:"part_cat_id"`
	YearFrom    int           `json:"year_from"`
	YearTo      int           `json:"year_to"`
	PartURL     string        `json:"part_url"`
	PartImgURL  string        `json:"part_img_url"`
	Prints      []string      `json:"prints"`
	Molds       []interface{} `json:"molds"`
	Alternates  []string      `json:"alternates"`
	ExternalIds struct {
		BrickOwl []string `json:"BrickOwl"`
		Brickset []string `json:"Brickset"`
		LDraw    []string `json:"LDraw"`
		LEGO     []string `json:"LEGO"`
	} `json:"external_ids"`
	PrintOf interface{} `json:"print_of"`
}

// Sets get a list of Set.
func (c *lClient) Sets(opts ...RequestOption) (sets []Set, err error) {
	err = c.GetDecode("lego/sets", true, &sets, opts...)
	return
}

// Set get details for a specific Set.
func (c *lClient) Set(setNumber string) (set Set, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/sets/%v", setNumber), false, &set)
	return
}

// SetAlternates get a list of MOCs which are alternate builds of a specific Set,
// i.e. all parts in the MOC can be found in the set.
func (c *lClient) SetAlternates(setNumber string, opts ...RequestOption) (sets []Set, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/sets/%v/alternates", setNumber), true, &sets, opts...)
	return
}

// SetMinifigs get a list of all inventory Minifig in this Set.
func (c *lClient) SetMinifigs(setNumber string, opts ...RequestOption) (minifigs []Minifig, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/sets/%v/minifigs", setNumber), true, &minifigs, opts...)
	return
}

// SetParts get a list of all inventory Part in this Set.
func (c *lClient) SetParts(setNumber string, opts ...RequestOption) (parts []Part, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/sets/%v/parts", setNumber), true, &parts, opts...)
	return
}

// SetSets get a list of all inventory Set in this Set.
func (c *lClient) SetSets(setNumber string, opts ...RequestOption) (sets []Set, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/sets/%v/sets", setNumber), true, &sets, opts...)
	return
}

type Set struct {
	SetNum         string    `json:"set_num"`
	Name           string    `json:"name"`
	Year           int       `json:"year"`
	ThemeID        int       `json:"theme_id"`
	NumParts       int       `json:"num_parts"`
	SetImgURL      string    `json:"set_img_url"`
	SetURL         string    `json:"set_url"`
	LastModifiedDt time.Time `json:"last_modified_dt"`
}

// Themes return all themes
func (c *lClient) Themes(opts ...RequestOption) (themes []Theme, err error) {
	err = c.GetDecode("lego/themes", true, &themes, opts...)
	return
}

// Theme get details for a specific Theme.
func (c *lClient) Theme(id int, opts ...RequestOption) (theme Theme, err error) {
	err = c.GetDecode(fmt.Sprintf("lego/themes/%v", id), false, &theme, opts...)
	return
}

type Theme struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
}
