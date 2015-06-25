// This code was originally based on the Digital Ocean provider from
// https://github.com/hashicorp/terraform/tree/master/builtin/providers/digitalocean.

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVultrServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrServerCreate,
		Read:   resourceVultrServerRead,
		Update: resourceVultrServerUpdate,
		Delete: resourceVultrServerDelete,

		Schema: map[string]*schema.Schema{
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"power_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"plan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"os_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"ipxe_chain_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			// if you are using an iso make sure you set `os_id` to `159` (Custom).
			"iso_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"script_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"ssh_key_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ipv4_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_private_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv6": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"ipv6_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_networking": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"auto_backups": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceVultrServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	name := d.Get("name").(string)
	regionId := d.Get("region_id").(int)
	planId := d.Get("plan_id").(int)
	osId := d.Get("os_id").(int)

	options := &lib.ServerOptions{
		IPXEChainURL: d.Get("ipxe_chain_url").(string),
		ISO:          d.Get("iso_id").(int),
		Script:       d.Get("script_id").(int),
		UserData:     d.Get("user_data").(string),
		Snapshot:     d.Get("snapshot_id").(string),
	}

	if attr, ok := d.GetOk("ipv6"); ok {
		options.IPV6 = attr.(bool)
	}

	if attr, ok := d.GetOk("private_networking"); ok {
		options.PrivateNetworking = attr.(bool)
	}

	if attr, ok := d.GetOk("auto_backups"); ok {
		options.AutoBackups = attr.(bool)
	}

	sshKeyIdsLen := d.Get("ssh_key_ids.#").(int)
	if sshKeyIdsLen > 0 {
		sshKeyIds := make([]string, 0, sshKeyIdsLen)
		for i := 0; i < sshKeyIdsLen; i++ {
			key := fmt.Sprintf("ssh_key_ids.%d", i)
			sshKeyIds = append(sshKeyIds, d.Get(key).(string))
		}
		options.SSHKey = strings.Join(sshKeyIds, ",")
	}

	log.Printf("[DEBUG] Server create configuration: %#v", options)

	server, err := client.CreateServer(name, regionId, planId, osId, options)

	if err != nil {
		return fmt.Errorf("Error creating server: %s", err)
	}

	d.SetId(server.ID)

	log.Printf("[INFO] Server ID: %s", d.Id())

	// wait for the server to be "ready". we have to wait for status=active and power_status=running.

	_, err = WaitForServerAttribute(d, "active", []string{"pending"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) to become active: %s", d.Id(), err)
	}

	_, err = WaitForServerAttribute(d, "running", []string{"stopped"}, "power_status", meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) to become running: %s", d.Id(), err)
	}

	return resourceVultrServerRead(d, meta)
}

func resourceVultrServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	server, err := client.GetServer(d.Id())
	if err != nil {
		// check if the server not longer exists.
		if err.Error() == "Invalid server." {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving server: %s", err)
	}

	d.Set("name", server.Name)
	d.Set("region_id", server.RegionID)
	d.Set("plan_id", server.PlanID)
	d.Set("status", server.Status)
	d.Set("power_status", server.PowerStatus)
	d.Set("default_password", server.DefaultPassword)
	d.Set("ipv4_address", server.MainIP)
	d.Set("ipv6_address", server.MainIPV6)
	d.Set("ipv4_private_address", server.InternalIP)

	d.SetConnInfo(map[string]string{
		"type":     "ssh",
		"host":     server.MainIP,
		"password": server.DefaultPassword,
	})

	return nil
}

func resourceVultrServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")

		err := client.RenameServer(d.Id(), newName.(string))

		if err != nil {
			return fmt.Errorf("Error renaming server (%s): %s", d.Id(), err)
		}

		_, err = WaitForServerAttribute(d, newName.(string), []string{"", oldName.(string)}, "name", meta)

		if err != nil {
			return fmt.Errorf("Error waiting for rename server (%s) to finish: %s", d.Id(), err)
		}

		d.SetPartial("name")
	}

	return resourceVultrServerRead(d, meta)
}

func resourceVultrServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*lib.Client)

	log.Printf("[INFO] Deleting server: %s", d.Id())

	err := client.DeleteServer(d.Id())

	if err != nil && strings.Contains(err.Error(), "404 Not Found") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error deleting server: %s", err)
	}

	return nil
}

func WaitForServerAttribute(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for server (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    newServerStateRefreshFunc(d, attribute, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForState()
}

// TODO This function still needs a little more refactoring to make it
// cleaner and more efficient
func newServerStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*lib.Client)
	return func() (interface{}, string, error) {
		err := resourceVultrServerRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		// See if we can access our attribute
		if attr, ok := d.GetOk(attribute); ok {
			// Retrieve the server properties
			server, err := client.GetServer(d.Id())
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving server: %s", err)
			}

			return &server, attr.(string), nil
		}

		return nil, "", nil
	}
}
