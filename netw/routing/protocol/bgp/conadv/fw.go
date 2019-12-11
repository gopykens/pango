package conadv

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
)

// FwConAdv is the client.Network.BgpConditionalAdv namespace.
type FwConAdv struct {
	con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwConAdv) Initialize(con util.XapiClient) {
	c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *FwConAdv) ShowList(vr string) ([]string, error) {
	c.con.LogQuery("(show) list of %s", plural)
	path := c.xpath(vr, nil)
	return c.con.EntryListUsing(c.con.Show, path[:len(path)-1])
}

// GetList performs GET to retrieve a list of values.
func (c *FwConAdv) GetList(vr string) ([]string, error) {
	c.con.LogQuery("(get) list of %s", plural)
	path := c.xpath(vr, nil)
	return c.con.EntryListUsing(c.con.Get, path[:len(path)-1])
}

// Get performs GET to retrieve information for the given uid.
func (c *FwConAdv) Get(vr, name string) (Entry, error) {
	c.con.LogQuery("(get) %s %q", singular, name)
	return c.details(c.con.Get, vr, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *FwConAdv) Show(vr, name string) (Entry, error) {
	c.con.LogQuery("(show) %s %q", singular, name)
	return c.details(c.con.Show, vr, name)
}

// Set performs SET to create / update one or more objects.
func (c *FwConAdv) Set(vr string, e ...Entry) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if vr == "" {
		return fmt.Errorf("vr must be specified")
	}

	_, fn := c.versioning()
	names := make([]string, len(e))

	// Build up the struct.
	d := util.BulkElement{XMLName: xml.Name{Local: "policy"}}
	for i := range e {
		d.Data = append(d.Data, fn(e[i]))
		names[i] = e[i].Name
	}
	c.con.LogAction("(set) %s: %v", plural, names)

	// Set xpath.
	path := c.xpath(vr, names)
	if len(e) == 1 {
		path = path[:len(path)-1]
	} else {
		path = path[:len(path)-2]
	}

	// Create the objects.
	_, err = c.con.Set(path, d.Config(), nil, nil)
	return err
}

// Edit performs EDIT to create / update one object.
func (c *FwConAdv) Edit(vr string, e Entry) error {
	var err error

	if vr == "" {
		return fmt.Errorf("vr must be specified")
	}

	_, fn := c.versioning()

	c.con.LogAction("(edit) %s %q", singular, e.Name)

	// Set xpath.
	path := c.xpath(vr, []string{e.Name})

	// Edit the object.
	_, err = c.con.Edit(path, fn(e), nil, nil)
	return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *FwConAdv) Delete(vr string, e ...interface{}) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if vr == "" {
		return fmt.Errorf("vr must be specified")
	}

	names := make([]string, len(e))
	for i := range e {
		switch v := e[i].(type) {
		case string:
			names[i] = v
		case Entry:
			names[i] = v.Name
		default:
			return fmt.Errorf("Unknown type sent to delete: %s", v)
		}
	}
	c.con.LogAction("(delete) %s: %v", plural, names)

	// Remove the objects.
	path := c.xpath(vr, names)
	_, err = c.con.Delete(path, nil, nil)
	return err
}

/** Internal functions for this namespace struct **/

func (c *FwConAdv) versioning() (normalizer, func(Entry) interface{}) {
	return &container_v1{}, specify_v1
}

func (c *FwConAdv) details(fn util.Retriever, vr, name string) (Entry, error) {
	path := c.xpath(vr, []string{name})
	obj, _ := c.versioning()
	if _, err := fn(path, nil, obj); err != nil {
		return Entry{}, err
	}
	ans := obj.Normalize()

	return ans, nil
}

func (c *FwConAdv) xpath(vr string, vals []string) []string {
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"virtual-router",
		util.AsEntryXpath([]string{vr}),
		"protocol",
		"bgp",
		"policy",
		"conditional-advertisement",
		"policy",
		util.AsEntryXpath(vals),
	}
}
