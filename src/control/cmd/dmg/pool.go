//
// (C) Copyright 2019-2020 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//

package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"

	"github.com/daos-stack/daos/src/control/client"
	"github.com/daos-stack/daos/src/control/common"
)

const (
	maxNumSvcReps = 13
)

// PoolCmd is the struct representing the top-level pool subcommand.
type PoolCmd struct {
	Create       PoolCreateCmd       `command:"create" alias:"c" description:"Create a DAOS pool"`
	Destroy      PoolDestroyCmd      `command:"destroy" alias:"d" description:"Destroy a DAOS pool"`
	List         systemListPoolsCmd  `command:"list" alias:"l" description:"List DAOS pools"`
	Query        PoolQueryCmd        `command:"query" alias:"q" description:"Query a DAOS pool"`
	GetACL       PoolGetACLCmd       `command:"get-acl" alias:"ga" description:"Get a DAOS pool's Access Control List"`
	OverwriteACL PoolOverwriteACLCmd `command:"overwrite-acl" alias:"oa" description:"Overwrite a DAOS pool's Access Control List"`
	UpdateACL    PoolUpdateACLCmd    `command:"update-acl" alias:"ua" description:"Update entries in a DAOS pool's Access Control List"`
	DeleteACL    PoolDeleteACLCmd    `command:"delete-acl" alias:"da" description:"Delete an entry from a DAOS pool's Access Control List"`
	SetProp      PoolSetPropCmd      `command:"set-prop" alias:"sp" description:"Set pool property"`
}

// PoolCreateCmd is the struct representing the command to create a DAOS pool.
type PoolCreateCmd struct {
	logCmd
	connectedCmd
	GroupName  string `short:"g" long:"group" description:"DAOS pool to be owned by given group, format name@domain"`
	UserName   string `short:"u" long:"user" description:"DAOS pool to be owned by given user, format name@domain"`
	ACLFile    string `short:"a" long:"acl-file" description:"Access Control List file path for DAOS pool"`
	ScmSize    string `short:"s" long:"scm-size" required:"1" description:"Size of SCM component of DAOS pool"`
	NVMeSize   string `short:"n" long:"nvme-size" description:"Size of NVMe component of DAOS pool"`
	RankList   string `short:"r" long:"ranks" description:"Storage server unique identifiers (ranks) for DAOS pool"`
	NumSvcReps uint32 `short:"v" long:"nsvc" default:"1" description:"Number of pool service replicas"`
	Sys        string `short:"S" long:"sys" default:"daos_server" description:"DAOS system that pool is to be a part of"`
	UUID       string `short:"p" long:"pool" description:"UUID to be used when creating the pool, randomly generated if not specified"`
}

// Execute is run when PoolCreateCmd subcommand is activated
func (c *PoolCreateCmd) Execute(args []string) error {
	msg := "SUCCEEDED: "

	scmBytes, err := humanize.ParseBytes(c.ScmSize)
	if err != nil {
		return errors.Wrap(err, "pool SCM size")
	}

	var nvmeBytes uint64
	if c.NVMeSize != "" {
		nvmeBytes, err = humanize.ParseBytes(c.NVMeSize)
		if err != nil {
			return errors.Wrap(err, "pool NVMe size")
		}
	}

	var acl *client.AccessControlList
	if c.ACLFile != "" {
		acl, err = readACLFile(c.ACLFile)
		if err != nil {
			return err
		}
	}

	if c.NumSvcReps > maxNumSvcReps {
		return errors.Errorf("max number of service replicas is %d, got %d",
			maxNumSvcReps, c.NumSvcReps)
	}

	usr, grp, err := formatNameGroup(c.UserName, c.GroupName)
	if err != nil {
		return errors.WithMessage(err, "formatting user/group strings")
	}

	var ranks []uint32
	if err := common.ParseNumberList(c.RankList, &ranks); err != nil {
		return errors.WithMessage(err, "parsing rank list")
	}

	req := &client.PoolCreateReq{
		ScmBytes: scmBytes, NvmeBytes: nvmeBytes, RankList: ranks,
		NumSvcReps: c.NumSvcReps, Sys: c.Sys, Usr: usr, Grp: grp, ACL: acl,
		UUID: c.UUID,
	}

	resp, err := c.conns.PoolCreate(req)
	if err != nil {
		msg = errors.WithMessage(err, "FAILED").Error()
	} else {
		msg += fmt.Sprintf("UUID: %s, Service replicas: %s",
			resp.UUID, formatPoolSvcReps(resp.SvcReps))
	}

	c.log.Infof("Pool-create command %s\n", msg)

	return err
}

// formatNameGroup converts system names to principal and if both user and group
// are unspecified, takes effective user name and that user's primary group.
func formatNameGroup(usr string, grp string) (string, string, error) {
	if usr == "" && grp == "" {
		eUsr, err := user.Current()
		if err != nil {
			return "", "", err
		}

		eGrp, err := user.LookupGroupId(eUsr.Gid)
		if err != nil {
			return "", "", err
		}

		usr, grp = eUsr.Username, eGrp.Name
	}

	if usr != "" && !strings.Contains(usr, "@") {
		usr += "@"
	}

	if grp != "" && !strings.Contains(grp, "@") {
		grp += "@"
	}

	return usr, grp, nil
}

// PoolDestroyCmd is the struct representing the command to destroy a DAOS pool.
type PoolDestroyCmd struct {
	logCmd
	connectedCmd
	// TODO: implement --sys & --svc options (currently unsupported server side)
	UUID  string `long:"pool" required:"1" description:"UUID of DAOS pool to destroy"`
	Force bool   `short:"f" long:"force" description:"Force removal of DAOS pool"`
}

// Execute is run when PoolDestroyCmd subcommand is activated
func (d *PoolDestroyCmd) Execute(args []string) error {
	msg := "succeeded"

	req := &client.PoolDestroyReq{UUID: d.UUID, Force: d.Force}

	err := d.conns.PoolDestroy(req)
	if err != nil {
		msg = errors.WithMessage(err, "failed").Error()
	}

	d.log.Infof("Pool-destroy command %s\n", msg)

	return err
}

// PoolQueryCmd is the struct representing the command to destroy a DAOS pool.
type PoolQueryCmd struct {
	logCmd
	connectedCmd
	UUID string `long:"pool" required:"1" description:"UUID of DAOS pool to query"`
}

// Execute is run when PoolQueryCmd subcommand is activated
func (c *PoolQueryCmd) Execute(args []string) error {
	req := client.PoolQueryReq{
		UUID: c.UUID,
	}

	resp, err := c.conns.PoolQuery(req)
	if err != nil {
		return errors.Wrap(err, "pool query failed")
	}

	// Maintain output compability with the `daos pool query` output.
	var bld strings.Builder
	fmt.Fprintf(&bld, "Pool %s, ntarget=%d, disabled=%d\n",
		resp.UUID, resp.TotalTargets, resp.DisabledTargets)
	bld.WriteString("Pool space info:\n")
	fmt.Fprintf(&bld, "- Target(VOS) count:%d\n", resp.ActiveTargets)
	if resp.Scm != nil {
		bld.WriteString("- SCM:\n")
		fmt.Fprintf(&bld, "  Total size: %s\n", humanize.Bytes(resp.Scm.Total))
		fmt.Fprintf(&bld, "  Free: %s, min:%s, max:%s, mean:%s\n",
			humanize.Bytes(resp.Scm.Free), humanize.Bytes(resp.Scm.Min),
			humanize.Bytes(resp.Scm.Max), humanize.Bytes(resp.Scm.Mean))
	}
	if resp.Nvme != nil {
		bld.WriteString("- NVMe:\n")
		fmt.Fprintf(&bld, "  Total size: %s\n", humanize.Bytes(resp.Nvme.Total))
		fmt.Fprintf(&bld, "  Free: %s, min:%s, max:%s, mean:%s\n",
			humanize.Bytes(resp.Nvme.Free), humanize.Bytes(resp.Nvme.Min),
			humanize.Bytes(resp.Nvme.Max), humanize.Bytes(resp.Nvme.Mean))
	}
	if resp.Rebuild != nil {
		if resp.Rebuild.Status == 0 {
			fmt.Fprintf(&bld, "Rebuild %s, %d objs, %d recs\n",
				resp.Rebuild.State, resp.Rebuild.Objects, resp.Rebuild.Records)
		} else {
			fmt.Fprintf(&bld, "Rebuild failed, rc=%d, status=%d", resp.Status, resp.Rebuild.Status)
		}
	}

	c.log.Info(bld.String())
	return nil
}

// PoolSetPropCmd represents the command to set a property on a pool.
type PoolSetPropCmd struct {
	logCmd
	connectedCmd
	UUID     string `long:"pool" required:"1" description:"UUID of DAOS pool"`
	Property string `short:"n" long:"name" required:"1" description:"Name of property to be set"`
	Value    string `short:"v" long:"value" required:"1" description:"Value of property to be set"`
}

// Execute is run when PoolSetPropCmd subcommand is activated.
func (c *PoolSetPropCmd) Execute(_ []string) error {
	req := client.PoolSetPropReq{
		UUID:     c.UUID,
		Property: c.Property,
	}

	req.SetString(c.Value)
	if numVal, err := strconv.ParseUint(c.Value, 10, 64); err == nil {
		req.SetNumber(numVal)
	}

	resp, err := c.conns.PoolSetProp(req)
	if err != nil {
		return errors.Wrap(err, "pool set-prop failed")
	}

	c.log.Infof("pool set-prop succeeded (%s=%q)", resp.Property, resp.Value)
	return nil
}

// PoolGetACLCmd represents the command to fetch an Access Control List of a
// DAOS pool.
type PoolGetACLCmd struct {
	logCmd
	connectedCmd
	UUID    string `long:"pool" required:"1" description:"UUID of DAOS pool"`
	File    string `short:"o" long:"outfile" required:"0" description:"Output ACL to file"`
	Force   bool   `short:"f" long:"force" required:"0" description:"Allow to clobber output file"`
	Verbose bool   `short:"v" long:"verbose" required:"0" description:"Add descriptive comments to ACL entries"`
}

// Execute is run when the PoolGetACLCmd subcommand is activated
func (d *PoolGetACLCmd) Execute(args []string) error {
	req := client.PoolGetACLReq{UUID: d.UUID}

	resp, err := d.conns.PoolGetACL(req)
	if err != nil {
		d.log.Infof("Pool-get-ACL command failed: %s\n", err.Error())
		return err
	}

	d.log.Debugf("Pool-get-ACL command succeeded, UUID: %s\n", d.UUID)

	acl := formatACL(resp.ACL, d.Verbose)

	if d.File != "" {
		err = d.writeACLToFile(acl)
		if err != nil {
			return err
		}
		d.log.Infof("Wrote ACL to output file: %s", d.File)
	} else {
		d.log.Info(acl)
	}

	return nil
}

func (d *PoolGetACLCmd) writeACLToFile(acl string) error {
	if !d.Force {
		// Keep the user from clobbering existing files
		_, err := os.Stat(d.File)
		if err == nil {
			return errors.New(fmt.Sprintf("file already exists: %s", d.File))
		}
	}

	f, err := os.Create(d.File)
	if err != nil {
		d.log.Errorf("Unable to create file: %s", d.File)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(acl)
	if err != nil {
		d.log.Errorf("Failed to write to file: %s", d.File)
		return err
	}

	return nil
}

// PoolOverwriteACLCmd represents the command to overwrite the Access Control
// List of a DAOS pool.
type PoolOverwriteACLCmd struct {
	logCmd
	connectedCmd
	UUID    string `long:"pool" required:"1" description:"UUID of DAOS pool"`
	ACLFile string `short:"a" long:"acl-file" required:"1" description:"Path for new Access Control List file"`
}

// Execute is run when the PoolOverwriteACLCmd subcommand is activated
func (d *PoolOverwriteACLCmd) Execute(args []string) error {
	acl, err := readACLFile(d.ACLFile)
	if err != nil {
		return err
	}

	req := client.PoolOverwriteACLReq{
		UUID: d.UUID,
		ACL:  acl,
	}

	resp, err := d.conns.PoolOverwriteACL(req)
	if err != nil {
		d.log.Infof("Pool-overwrite-ACL command failed: %s\n", err.Error())
		return err
	}

	d.log.Infof("Pool-overwrite-ACL command succeeded, UUID: %s\n", d.UUID)
	d.log.Info(formatACLDefault(resp.ACL))

	return nil
}

// PoolUpdateACLCmd represents the command to update the Access Control List of
// a DAOS pool.
type PoolUpdateACLCmd struct {
	logCmd
	connectedCmd
	UUID    string `long:"pool" required:"1" description:"UUID of DAOS pool"`
	ACLFile string `short:"a" long:"acl-file" required:"0" description:"Path for new Access Control List file"`
	Entry   string `short:"e" long:"entry" required:"0" description:"Single Access Control Entry to add or update"`
}

// Execute is run when the PoolUpdateACLCmd subcommand is activated
func (d *PoolUpdateACLCmd) Execute(args []string) error {
	if (d.ACLFile == "" && d.Entry == "") || (d.ACLFile != "" && d.Entry != "") {
		return errors.New("either ACL file or entry parameter is required")
	}

	var acl *client.AccessControlList
	if d.ACLFile != "" {
		aclFileResult, err := readACLFile(d.ACLFile)
		if err != nil {
			return err
		}
		acl = aclFileResult
	} else {
		acl = &client.AccessControlList{
			Entries: []string{d.Entry},
		}
	}

	req := client.PoolUpdateACLReq{
		UUID: d.UUID,
		ACL:  acl,
	}

	resp, err := d.conns.PoolUpdateACL(req)
	if err != nil {
		d.log.Infof("Pool-update-ACL command failed: %s\n", err.Error())
		return err
	}

	d.log.Infof("Pool-update-ACL command succeeded, UUID: %s\n", d.UUID)
	d.log.Info(formatACLDefault(resp.ACL))

	return nil
}

// PoolDeleteACLCmd represents the command to delete an entry from the Access
// Control List of a DAOS pool.
type PoolDeleteACLCmd struct {
	logCmd
	connectedCmd
	UUID      string `long:"pool" required:"1" description:"UUID of DAOS pool"`
	Principal string `short:"p" long:"principal" required:"1" description:"Principal whose entry should be removed"`
}

// Execute is run when the PoolDeleteACLCmd subcommand is activated
func (d *PoolDeleteACLCmd) Execute(args []string) error {
	req := client.PoolDeleteACLReq{
		UUID:      d.UUID,
		Principal: d.Principal,
	}

	resp, err := d.conns.PoolDeleteACL(req)
	if err != nil {
		d.log.Infof("Pool-delete-ACL command failed: %s\n", err.Error())
		return err
	}

	d.log.Infof("Pool-delete-ACL command succeeded, UUID: %s\n", d.UUID)
	d.log.Info(formatACLDefault(resp.ACL))

	return nil
}
