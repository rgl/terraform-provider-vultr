package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVultrSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrSSHKeyCreate,
		Read:   resourceVultrSSHKeyRead,
		Update: resourceVultrSSHKeyUpdate,
		Delete: resourceVultrSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					// Vultr always trims the key, so to get a correct diff we
					// also need to trim.
					return strings.TrimSpace(v.(string))
				},
			},
		},
	}
}

func resourceVultrSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	name := d.Get("name").(string)
	publicKey := d.Get("public_key").(string)

	log.Printf("[DEBUG] SSH Key create: %s", name)

	key, err := client.CreateSSHKey(name, publicKey)
	if err != nil {
		return fmt.Errorf("Error creating SSH Key: %s", err)
	}

	d.SetId(key.ID)

	log.Printf("[INFO] SSH Key: %s", key.ID)

	return resourceVultrSSHKeyRead(d, meta)
}

func resourceVultrSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	// NB the API has no support for getting a key by id, so we have to get'em all.

	keys, err := client.GetSSHKeys()
	if err != nil {
		return fmt.Errorf("Error retrieving SSH key: %s", err)
	}

	var key *lib.SSHKey

	for _, k := range keys {
		if k.ID == d.Id() {
			key = &k
			break
		}
	}

	// if the key is somehow already destroyed mark as succesfully gone.
	if key == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", key.Name)
	d.Set("public_key", key.Key)

	return nil
}

func resourceVultrSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	key := lib.SSHKey{
		ID:   d.Id(),
		Name: d.Get("name").(string),
		Key:  d.Get("public_key").(string),
	}

	log.Printf("[DEBUG] SSH key update: %s", key.ID)
	err := client.UpdateSSHKey(key)
	if err != nil {
		return fmt.Errorf("Failed to update SSH key: %s", err)
	}

	return resourceVultrSSHKeyRead(d, meta)
}

func resourceVultrSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	log.Printf("[INFO] Deleting SSH key: %s", d.Id())
	err := client.DeleteSSHKey(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting SSH key: %s", err)
	}

	return nil
}
