// This code was originally based on the Digital Ocean provider from
// https://github.com/terraform-providers/terraform-provider-digitalocean.

package main

import (
	"fmt"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVultrDNSDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrDNSDomainCreate,
		Read:   resourceVultrDNSDomainRead,
		Delete: resourceVultrDNSDomainDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVultrDNSDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	name := d.Get("name").(string)
	ipv4Address := d.Get("ipv4_address").(string)

	err := client.CreateDNSDomain(name, ipv4Address)
	if err != nil {
		return fmt.Errorf("Error creating domain: %s", err)
	}

	d.SetId(name)

	return resourceVultrDNSDomainRead(d, meta)
}

func resourceVultrDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	domains, err := client.GetDNSDomains()
	if err != nil {
		return fmt.Errorf("Error retrieving domain: %s", err)
	}

	var domain *lib.DNSDomain

	for _, c := range domains {
		if c.Domain == d.Id() {
			domain = &c
			break
		}
	}

	// if the domain is somehow already destroyed mark as succesfully gone.
	if domain == nil {
		d.SetId("")
		return nil
	}

	// find the ipv4 address record associated with the domain.
	records, err := client.GetDNSRecords(domain.Domain)
	if err != nil {
		return fmt.Errorf("Error retrieving domain records: %s", err)
	}
	var record *lib.DNSRecord
	for _, r := range records {
		if r.Type == "A" && r.Name == "" {
			record = &r
			break
		}
	}

	// if we cannot find the default ipv4 record for the domain, mark the entire domain as succesfully gone.
	if record == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", domain.Domain)
	d.Set("ipv4_address", record.Data)

	return nil
}

func resourceVultrDNSDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	err := client.DeleteDNSDomain(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting domain: %s", err)
	}

	d.SetId("")
	return nil
}
