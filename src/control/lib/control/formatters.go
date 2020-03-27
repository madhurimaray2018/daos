//
// (C) Copyright 2020 Intel Corporation.
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
package control

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"

	"github.com/daos-stack/daos/src/control/lib/txtfmt"
	"github.com/daos-stack/daos/src/control/server/storage"
)

// Provide formatters for complex response types. The API should return structured
// responses so that callers can fully utilize them, but should also provide default
// human-readable representations so that callers aren't forced to implement formatters.

var (
	defaultFormatConfig = &FormatConfig{
		Verbose:       false,
		ShowHostPorts: false,
	}
)

type (
	// FormatConfig defines parameters for controlling formatter behavior.
	FormatConfig struct {
		// Verbose indicates that the output should include more detail.
		Verbose bool
		// ShowHostPorts indicates that the host output should include the network port.
		ShowHostPorts bool
	}

	// FormatConfigOption defines a config function.
	FormatConfigOption func(*FormatConfig)
)

// FmtWithVerboseOutput toggles verbose output from the formatter.
func FmtWithVerboseOutput(verbose bool) FormatConfigOption {
	return func(cfg *FormatConfig) {
		cfg.Verbose = verbose
	}
}

// FmtWithHostPorts enables display of host ports in output.
func FmtWithHostPorts() FormatConfigOption {
	return func(cfg *FormatConfig) {
		cfg.ShowHostPorts = true
	}
}

// getFormatConfig is a helper that returns a format configuration
// for a format function.
func getFormatConfig(opts ...FormatConfigOption) *FormatConfig {
	cfg := &FormatConfig{}
	*cfg = *defaultFormatConfig
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// getFormatHosts is a helper that transforms the given list of
// host strings according to the format configuration.
func getFormatHosts(in string, opts ...FormatConfigOption) string {
	var out []string
	fc := getFormatConfig(opts...)

	for _, hostStr := range strings.Split(in, ",") {
		if fc.ShowHostPorts {
			out = append(out, hostStr)
			continue
		}

		hostPort := strings.Split(hostStr, ":")
		if len(hostPort) != 2 {
			out = append(out, hostStr)
			continue
		}
		out = append(out, hostPort[0])
	}

	return strings.Join(out, ",")
}

// FormatPoolQueryResponse generates a human-readable representation of the supplied
// PoolQueryResp struct and writes it to the supplied io.Writer.
func FormatPoolQueryResponse(pqr *PoolQueryResp, out io.Writer, opts ...FormatConfigOption) error {
	if pqr == nil {
		return errors.Errorf("nil %T", pqr)
	}
	w := txtfmt.NewErrWriter(out)

	// Maintain output compability with the `daos pool query` output.
	fmt.Fprintf(w, "Pool %s, ntarget=%d, disabled=%d\n",
		pqr.UUID, pqr.TotalTargets, pqr.DisabledTargets)
	fmt.Fprintln(w, "Pool space info:")
	fmt.Fprintf(w, "- Target(VOS) count:%d\n", pqr.ActiveTargets)
	if pqr.Scm != nil {
		fmt.Fprintln(w, "- SCM:")
		fmt.Fprintf(w, "  Total size: %s\n", humanize.Bytes(pqr.Scm.Total))
		fmt.Fprintf(w, "  Free: %s, min:%s, max:%s, mean:%s\n",
			humanize.Bytes(pqr.Scm.Free), humanize.Bytes(pqr.Scm.Min),
			humanize.Bytes(pqr.Scm.Max), humanize.Bytes(pqr.Scm.Mean))
	}
	if pqr.Nvme != nil {
		fmt.Fprintln(w, "- NVMe:")
		fmt.Fprintf(w, "  Total size: %s\n", humanize.Bytes(pqr.Nvme.Total))
		fmt.Fprintf(w, "  Free: %s, min:%s, max:%s, mean:%s\n",
			humanize.Bytes(pqr.Nvme.Free), humanize.Bytes(pqr.Nvme.Min),
			humanize.Bytes(pqr.Nvme.Max), humanize.Bytes(pqr.Nvme.Mean))
	}
	if pqr.Rebuild != nil {
		if pqr.Rebuild.Status == 0 {
			fmt.Fprintf(w, "Rebuild %s, %d objs, %d recs\n",
				pqr.Rebuild.State, pqr.Rebuild.Objects, pqr.Rebuild.Records)
		} else {
			fmt.Fprintf(w, "Rebuild failed, rc=%d, status=%d", pqr.Status, pqr.Rebuild.Status)
		}
	}

	return w.Err
}

func (pqr *PoolQueryResp) String() string {
	var bld strings.Builder
	if err := FormatPoolQueryResponse(pqr, &bld); err != nil {
		return fmt.Sprintf("failed to format %T", pqr)
	}
	return bld.String()
}

// FormatHostErrorsMap generates a human-readable representation of the supplied
// HostErrorsMap struct and writes it to the supplied io.Writer.
func FormatHostErrorsMap(hem HostErrorsMap, out io.Writer, opts ...FormatConfigOption) error {
	if len(hem) == 0 {
		return errors.Errorf("empty %T", hem)
	}

	setTitle := "Hosts"
	errTitle := "Error"

	tableFmt := txtfmt.NewTableFormatter(setTitle, errTitle)
	tableFmt.InitWriter(out)
	table := []txtfmt.TableRow{}

	for _, errStr := range hem.Keys() {
		errHosts := getFormatHosts(hem[errStr].RangedString(), opts...)
		row := txtfmt.TableRow{setTitle: errHosts}
		row[errTitle] = errStr

		table = append(table, row)
	}

	tableFmt.Format(table)
	return nil
}

func (hem HostErrorsMap) String() string {
	var bld strings.Builder
	if err := FormatHostErrorsMap(hem, &bld); err != nil {
		return fmt.Sprintf("failed to format %T", hem)
	}
	return bld.String()
}

func formatNvmeControllers(controllers storage.NvmeControllers, out io.Writer, opts ...FormatConfigOption) error {
	if len(controllers) == 0 {
		fmt.Fprintf(out, "\tNo NVMe devices found\n")
		return nil
	}

	pciTitle := "NVMe PCI"
	modelTitle := "Model"
	fwTitle := "FW Revision"
	socketTitle := "Socket ID"
	capacityTitle := "Capacity"

	formatter := txtfmt.NewTableFormatter(
		pciTitle, modelTitle, fwTitle, socketTitle, capacityTitle,
	)
	formatter.InitWriter(out)
	var table []txtfmt.TableRow

	sort.Slice(controllers, func(i, j int) bool { return controllers[i].PciAddr < controllers[j].PciAddr })

	for _, ctrlr := range controllers {
		row := txtfmt.TableRow{pciTitle: ctrlr.PciAddr}
		row[modelTitle] = ctrlr.Model
		row[fwTitle] = ctrlr.FwRev
		row[socketTitle] = fmt.Sprint(ctrlr.SocketID)
		row[capacityTitle] = humanize.Bytes(ctrlr.Capacity())

		table = append(table, row)
	}

	formatter.Format(table)
	return nil
}

func formatScmModules(modules storage.ScmModules, out io.Writer, opts ...FormatConfigOption) error {
	if len(modules) == 0 {
		fmt.Fprintf(out, "\tNo SCM modules found\n")
		return nil
	}

	physicalIdTitle := "SCM Module ID"
	socketTitle := "Socket ID"
	memCtrlrTitle := "Memory Ctrlr ID"
	channelTitle := "Channel ID"
	slotTitle := "Channel Slot"
	capacityTitle := "Capacity"

	formatter := txtfmt.NewTableFormatter(
		physicalIdTitle, socketTitle, memCtrlrTitle, channelTitle, slotTitle, capacityTitle,
	)
	formatter.InitWriter(out)
	var table []txtfmt.TableRow

	sort.Slice(modules, func(i, j int) bool { return modules[i].PhysicalID < modules[j].PhysicalID })

	for _, m := range modules {
		row := txtfmt.TableRow{physicalIdTitle: fmt.Sprint(m.PhysicalID)}
		row[socketTitle] = fmt.Sprint(m.SocketID)
		row[memCtrlrTitle] = fmt.Sprint(m.ControllerID)
		row[channelTitle] = fmt.Sprint(m.ChannelID)
		row[slotTitle] = fmt.Sprint(m.ChannelPosition)
		row[capacityTitle] = humanize.IBytes(m.Capacity)

		table = append(table, row)
	}

	formatter.Format(table)
	return nil
}

// formatHostStorageMapVerbose generates a human-readable representation of the supplied
// HostStorageMap struct and writes it to the supplied io.Writer.
func formatHostStorageMapVerbose(hsm HostStorageMap, out io.Writer, opts ...FormatConfigOption) error {
	for _, key := range hsm.Keys() {
		hss := hsm[key]
		hosts := getFormatHosts(hss.HostSet.RangedString(), opts...)
		lineBreak := strings.Repeat("-", len(hosts))
		fmt.Fprintf(out, "%s\n%s\n%s\n", lineBreak, hosts, lineBreak)
		if err := formatScmModules(hss.HostStorage.Scm, out, opts...); err != nil {
			return err
		}
		fmt.Fprintln(out)
		if err := formatNvmeControllers(hss.HostStorage.Nvme, out, opts...); err != nil {
			return err
		}
		fmt.Fprintln(out)
	}

	return nil
}

// FormatHostStorageMap generates a human-readable representation of the supplied
// HostStorageMap struct and writes it to the supplied io.Writer.
func FormatHostStorageMap(hsm HostStorageMap, out io.Writer, opts ...FormatConfigOption) error {
	if len(hsm) == 0 {
		return errors.Errorf("empty %T", hsm)
	}
	fc := getFormatConfig(opts...)

	if fc.Verbose {
		return formatHostStorageMapVerbose(hsm, out, opts...)
	}

	hostsTitle := "Hosts"
	scmTitle := "SCM Total"
	nvmeTitle := "NVMe Total"

	tableFmt := txtfmt.NewTableFormatter(hostsTitle, scmTitle, nvmeTitle)
	tableFmt.InitWriter(out)
	table := []txtfmt.TableRow{}

	for _, key := range hsm.Keys() {
		hss := hsm[key]
		hosts := getFormatHosts(hss.HostSet.RangedString(), opts...)
		row := txtfmt.TableRow{hostsTitle: hosts}
		row[scmTitle] = hss.HostStorage.Scm.Summary()
		row[nvmeTitle] = hss.HostStorage.Nvme.Summary()
		table = append(table, row)
	}

	tableFmt.Format(table)
	return nil
}

func (hsm HostStorageMap) String() string {
	var bld strings.Builder
	if err := FormatHostStorageMap(hsm, &bld); err != nil {
		return fmt.Sprintf("failed to format %T", hsm)
	}
	return bld.String()
}

// FormatStorageScanResponse generates a human-readable representation of the supplied
// StorageScanResp struct and writes it to the supplied io.Writer.
func FormatStorageScanResponse(ssr *StorageScanResp, out io.Writer, opts ...FormatConfigOption) error {
	if ssr == nil {
		return errors.Errorf("nil %T", ssr)
	}

	if len(ssr.HostErrors) > 0 {
		fmt.Fprintln(out, "Errors:")
		if err := FormatHostErrorsMap(ssr.HostErrors, txtfmt.NewIndentWriter(out), opts...); err != nil {
			return err
		}
		fmt.Fprintln(out)
	}

	if len(ssr.HostStorage) > 0 {
		return FormatHostStorageMap(ssr.HostStorage, out, opts...)
	}

	return nil
}

func (ssr *StorageScanResp) String() string {
	var bld strings.Builder
	if err := FormatStorageScanResponse(ssr, &bld); err != nil {
		return fmt.Sprintf("failed to format %T", ssr)
	}
	return bld.String()
}
