package rebrickable

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type User struct {
	*Client
	token string
}

func (u *User) Do(req *http.Request) (*http.Response, error) {
	fmt.Println(req.URL.Path)
	fmt.Println(req.URL.RawPath)
	return u.Client.Do(req)
}

func (u *User) buildEndpoint(endpoint string, a ...interface{}) string {
	return fmt.Sprintf("users/%v/%v", u.token, fmt.Sprintf(endpoint, a...))
}

// AllParts get a list of all the Part/s in all the user's PartList/s as well as the Part/s
// inside Set/s in the user's SetList/s.
//
// WARNING: This call is very resource intensive, do not overuse it!
func (u *User) AllParts(opts ...RequestOption) (parts []Part, err error) {
	err = u.get("allparts", true, &parts, opts...)
	return
}

// Build find out how many parts the user needs to build the specified Set.
func (u *User) Build(setNumber string) (build Build, err error) {
	err = u.get(u.buildEndpoint("build/%v", setNumber), false, &build)
	return
}

type Build struct {
	User                  int         `json:"user"`
	Inventory             int         `json:"inventory"`
	UserList              interface{} `json:"user_list"`
	PctOwned              int         `json:"pct_owned"`
	NumMissing            int         `json:"num_missing"`
	NumIgnored            int         `json:"num_ignored"`
	NumOwnedLessIgnored   int         `json:"num_owned_less_ignored"`
	TotalParts            int         `json:"total_parts"`
	TotalPartsLessIgnored int         `json:"total_parts_less_ignored"`
	BuildOptions          struct {
		IgnorePrint    bool        `json:"ignore_print"`
		IgnoreMold     bool        `json:"ignore_mold"`
		IgnoreAltp     bool        `json:"ignore_altp"`
		IgnoreMinifigs bool        `json:"ignore_minifigs"`
		IgnoreNonLego  bool        `json:"ignore_non_lego"`
		SortBy         int         `json:"sort_by"`
		Color          int         `json:"color"`
		Theme          interface{} `json:"theme"`
		MinParts       int         `json:"min_parts"`
		MaxParts       int         `json:"max_parts"`
		MinYear        int         `json:"min_year"`
		MaxYear        int         `json:"max_year"`
		AddedDaysAgo   int         `json:"added_days_ago"`
		IncOfficial    bool        `json:"inc_official"`
		IncCustom      bool        `json:"inc_custom"`
		IncBmodels     bool        `json:"inc_bmodels"`
		IncAccessory   bool        `json:"inc_accessory"`
		IncPremium     bool        `json:"inc_premium"`
		IncAlts        bool        `json:"inc_alts"`
		IncOwned       bool        `json:"inc_owned"`
	} `json:"build_options"`
}

// LostParts get a list of all the LostPart/s from the user's LEGO collection.
func (u *User) LostParts(opts ...RequestOption) (parts []LostPart, err error) {
	err = u.get(u.buildEndpoint("lost_parts"), true, &parts, opts...)
	return
}

// NewLostPart add a LostPart to the user.
//
// TODO handle multiple lost parts
func (u *User) NewLostPart(invPartId, quantity int) (part LostPart, err error) {
	form := url.Values{}
	form.Set("lost_quantity", fmt.Sprint(quantity))
	form.Set("inv_part_id", fmt.Sprint(invPartId))
	err = u.post(u.buildEndpoint("lost_parts"), form, &part)
	return
}

// DeleteLostPart remove the LostPart from the user.
func (u *User) DeleteLostPart(id string) error {
	return u.delete(u.buildEndpoint("lost_parts/%v", id))
}

type LostPart struct {
	LostPartID   int `json:"lost_part_id"`
	LostQuantity int `json:"lost_quantity"`
	InvPart      struct {
		ID        int `json:"id"`
		InvPartID int `json:"inv_part_id"`
		Part      `json:"part"`
		Color     `json:"color"`
		SetNum    string `json:"set_num"`
		Quantity  int    `json:"quantity"`
		IsSpare   bool   `json:"is_spare"`
		ElementID string `json:"element_id"`
		NumSets   int    `json:"num_sets"`
	} `json:"inv_part"`
}

// Minifigs get a list of all the Minifig/s in all the user's Set/s. Note that this is a
// read-only list as Minifig/s are automatically determined by the Set/s in the user's SetList/s
func (u *User) Minifigs(opts ...RequestOption) (minifigs []Minifig, err error) {
	err = u.get(u.buildEndpoint("minifigs"), true, &minifigs, opts...)
	return
}

// PartLists get a list of all the user's PartList/s.
func (u *User) PartLists(opts ...RequestOption) (partLists []PartList, err error) {
	err = u.get(u.buildEndpoint("partlists"), true, &partLists, opts...)
	return
}

// NewPartList add a new PartList.
func (u *User) NewPartList(name string, numParts int, buildable bool) (partList PartList, err error) {
	form := url.Values{}
	form.Set("is_buildable", strconv.FormatBool(buildable))
	form.Set("name", name)
	form.Set("num_parts", fmt.Sprint(numParts))

	err = u.post(u.buildEndpoint("partlists"), form, &partList)
	return
}

// DeletePartList delete a PartList and all it's Part/s.
func (u *User) DeletePartList(listId string) error {
	return u.delete(u.buildEndpoint("partlists/%v", listId))
}

// PartList get details about a specific PartList.
func (u *User) PartList(listId string) (partList PartList, err error) {
	err = u.get(u.buildEndpoint("partlists/%v", listId), false, &partList)
	return
}

// UpdatePartList update an existing PartList's details.
//
// TODO handle replace part list?
func (u *User) UpdatePartList(pList PartList, listId string) (partList PartList, err error) {
	err = u.patch(u.buildEndpoint("partlists/%v", listId), pList.Form(), &partList)
	return
}

// PartListParts get a list of all the Part/s in a specified PartList.
func (u *User) PartListParts(listId string, opts ...RequestOption) (parts []Part, err error) {
	err = u.get(u.buildEndpoint("partlists/%v/parts", listId), true, &parts, opts...)
	return
}

// NewPartListPart add a Part to the PartList.
//
// TODO handle multiple new parts
func (u *User) NewPartListPart(listId string, partNum string, quantity, colorId int) (part Part, err error) {
	form := url.Values{}
	form.Set("part_num", partNum)
	form.Set("quantity", fmt.Sprint(quantity))
	form.Set("color_id", fmt.Sprint(colorId))

	err = u.post(u.buildEndpoint("partlists/%v/parts", listId), form, &part)
	return
}

// DeletePartListPart delete a Part from the PartList
func (u *User) DeletePartListPart(listId, colorId, partNum string, opts ...RequestOption) error {
	return u.delete(u.buildEndpoint("partlists/%v/parts/%v/%v", listId, partNum, colorId), opts...)
}

// PartListPart get details about a specific Part in the PartList.
func (u *User) PartListPart(listId, partNum, colorId string) (part Part, err error) {
	err = u.get(u.buildEndpoint("partlists/%v/parts/%v/%v", listId, partNum, colorId), false, &part)
	return
}

// UpdatePartListPart replace an existing Part's details in the PartList.
func (u *User) UpdatePartListPart(listId, partNum, colorId string, quantity int) (part Part, err error) {
	form := url.Values{}
	form.Set("quantity", fmt.Sprint(quantity))

	err = u.put(u.buildEndpoint("partlists/%v/parts/%v/%v", listId, partNum, colorId), nil, &part)
	return
}

type PartList struct {
	ID          int    `json:"id"`
	IsBuildable bool   `json:"is_buildable"`
	Name        string `json:"name"`
	NumParts    int    `json:"num_parts"`
}

func (p PartList) Form() url.Values {
	form := url.Values{}
	form.Set("is_buildable", strconv.FormatBool(p.IsBuildable))
	form.Set("name", p.Name)
	form.Set("num_parts", fmt.Sprint(p.NumParts))

	return form
}

