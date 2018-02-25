// This code was originally based on the Digital Ocean provider from
// https://github.com/terraform-providers/terraform-provider-digitalocean.

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVultrDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrDNSRecordCreate,
		Read:   resourceVultrDNSRecordRead,
		Update: resourceVultrDNSRecordUpdate,
		Delete: resourceVultrDNSRecordDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"data": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	domain := d.Get("domain").(string)
	rtype := d.Get("type").(string)
	name := d.Get("name").(string)
	data := d.Get("data").(string)
	priority := 0 // only used in MX and SRV resource records.
	ttl := 0

	if attr, ok := d.GetOk("priority"); ok {
		priority = attr.(int)
	}

	if attr, ok := d.GetOk("ttl"); ok {
		ttl = attr.(int)
	}

	err := client.CreateDNSRecord(domain, name, rtype, data, priority, ttl)
	if err != nil {
		return fmt.Errorf("Failed to create dns record: %s", err)
	}

	records, err := client.GetDNSRecords(domain)
	if err != nil {
		return fmt.Errorf("Failed to get dns records: %s", err)
	}
	for _, r := range records {
		if r.Name == name && r.Type == rtype && r.Data == data {
			d.SetId(strconv.Itoa(r.RecordID))
			return resourceVultrDNSRecordRead(d, meta)
		}
	}

	return fmt.Errorf("Failed to get the just created dns record: %s", err)
}

func resourceVultrDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to parse resource id: %s", err)
	}
	domain := d.Get("domain").(string)

	records, err := client.GetDNSRecords(domain)
	if err != nil {
		return fmt.Errorf("Failed to get dns records: %s", err)
	}

	var record *lib.DNSRecord
	for _, r := range records {
		if r.RecordID == id {
			record = &r
			break
		}
	}

	// if it is somehow already destroyed mark as succesfully gone.
	if record == nil {
		d.SetId("")
		return nil
	}

	d.Set("type", record.Type)
	d.Set("name", record.Name)
	d.Set("data", record.Data)
	d.Set("priority", strconv.Itoa(record.Priority))
	d.Set("ttl", strconv.Itoa(record.TTL))
	d.Set("fqdn", constructFqdn(record.Name, domain))

	return nil
}

func resourceVultrDNSRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to parse resource id: %s", err)
	}
	domain := d.Get("domain").(string)

	record := lib.DNSRecord{
		RecordID: id,
		Type:     d.Get("type").(string),
		Name:     d.Get("name").(string),
		Data:     d.Get("data").(string),
		Priority: d.Get("priority").(int),
		TTL:      d.Get("ttl").(int),
	}
	if err := client.UpdateDNSRecord(domain, record); err != nil {
		return fmt.Errorf("Failed to update record: %v", err)
	}

	return resourceVultrDNSRecordRead(d, meta)
}

func resourceVultrDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	domain := d.Get("domain").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to parse resource id: %s", err)
	}

	err = client.DeleteDNSRecord(domain, id)
	if err != nil {
		return fmt.Errorf("Failed to delete record: %s", err)
	}

	return nil
}

func constructFqdn(name, domain string) string {
	rn := strings.ToLower(strings.TrimSuffix(name, "."))
	domain = strings.TrimSuffix(domain, ".")
	if !strings.HasSuffix(rn, domain) {
		rn = strings.Join([]string{name, domain}, ".")
	}
	return rn
}
