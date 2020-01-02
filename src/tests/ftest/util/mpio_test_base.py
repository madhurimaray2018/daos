#!/usr/bin/python
'''
  (C) Copyright 2019 Intel Corporation.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

  GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
  The Government's rights to use, modify, reproduce, release, perform, display,
  or disclose this software are subject to the terms of the Apache License as
  provided in Contract No. B609815.
  Any reproduction of computer software, computer software documentation, or
  portions thereof marked with this legend must also reproduce the markings.
'''
import os
from apricot import TestWithServers

from mpio_utils import MpioUtils, MpioFailed
from pydaos.raw import DaosPool, DaosApiError
from write_host_file import write_host_file


class LlnlMpi4pyHdf5(TestWithServers):
    """Runs LLNL, MPI4PY and HDF5 test suites.

    :avocado: recursive
    """

    def setUp(self):
        """Set up each LlnlMpi4pyHdf5 test."""
        super(LlnlMpi4pyHdf5, self).setUp()

        try:
            # parameters used in pool create
            createmode = self.params.get("mode", '/run/pool/createmode/*/')
            createuid = os.geteuid()
            creategid = os.getegid()
            createsetid = self.params.get("setname", '/run/pool/createset/')
            createsize = self.params.get("size", '/run/pool/createsize/')
            self.createsvc = self.params.get("svcn", '/run/pool/createsvc/')

            # initialize a python pool object then create the underlying
            # daos storage
            self.pool = DaosPool(self.context)
            self.pool.create(createmode, createuid, creategid,
                             createsize, createsetid, None, None,
                             self.createsvc)
        except (DaosApiError) as excep:
            self.fail("<Test Failed at pool create> \n{}".format(excep))

    def run_test(self, test_repo, test_name):
        """Execute the specified test.

        Args:
            test_repo (str): location of test repository
            test_name (str): name of the test to be run
        """
        # initialize MpioUtils
        mpio = MpioUtils()
        if not mpio.mpich_installed(self.hostlist_clients):
            self.fail("Exiting Test: Mpich not installed")

        try:
            # initialise test specific variables
            client_processes = self.params.get("np", '/run/client_processes/')

            # obtaining pool uuid and svc list
            pool_uuid = self.pool.get_uuid_str()
            svc_list = ""
            for i in range(self.createsvc):
                svc_list += str(int(self.pool.svc.rl_ranks[i])) + ":"
            svc_list = svc_list[:-1]

            # running tests
            hostfile = write_host_file(
                self.hostlist_clients, self.workdir, None)
            mpio.run_llnl_mpi4py_hdf5(
                hostfile, pool_uuid, test_repo, test_name, client_processes)

            # Parsing output to look for failures
            # stderr directed to stdout
            stdout = self.logdir + "/stdout"
            searchfile = open(stdout, "r")
            error_message = ["non-zero exit code", "MPI_Abort", "MPI_ABORT",
                             "ERROR"]

            for line in searchfile:
                # pylint: disable=C0200
                for i in range(len(error_message)):
                    if error_message[i] in line:
                        self.fail("Test Failed with error_message: {}"
                                  .format(error_message[i]))

        except (MpioFailed, DaosApiError) as excep:
            self.fail("<{0} Test Failed> \n{1}".format(test_name, excep))
