/**
 * (C) Copyright 2017-2019 Intel Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
 * The Government's rights to use, modify, reproduce, release, perform, display,
 * or disclose this software are subject to the terms of the Apache License as
 * provided in Contract No. B609815.
 * Any reproduction of computer software, computer software documentation, or
 * portions thereof marked with this legend must also reproduce the markings.
 */

#include "dfuse_common.h"
#include "dfuse.h"

/* Lookup a unique inode for the specific dfs/oid combination */
int
dfuse_lookup_inode(struct dfuse_projection_info *fs_handle,
		   struct dfuse_dfs *dfs,
		   daos_obj_id_t *oid,
		   ino_t *_ino)
{
	struct dfuse_inode_record	*dfir;
	d_list_t			*rlink;
	int				rc = -DER_SUCCESS;

	D_ALLOC_PTR(dfir);
	if (!dfir) {
		D_GOTO(out, rc = -DER_NOMEM);
	}

	if (oid) {
		dfir->ir_id.irid_oid.lo = oid->lo;
		dfir->ir_id.irid_oid.hi = oid->hi;
	}

	dfir->ir_ino = atomic_fetch_add(&fs_handle->dfpi_ino_next, 1);
	dfir->ir_id.irid_dfs = dfs;

	rlink = d_hash_rec_find_insert(&fs_handle->dfpi_irt,
				       &dfir->ir_id,
				       sizeof(dfir->ir_id),
				       &dfir->ir_htl);

	if (rlink != &dfir->ir_htl) {
		D_FREE(dfir);
		dfir = container_of(rlink, struct dfuse_inode_record, ir_htl);
	}

	*_ino = dfir->ir_ino;

out:
	return rc;
};

int
find_inode(struct dfuse_request *request)
{
	struct dfuse_projection_info *fs_handle = request->fsh;
	struct dfuse_inode_entry *ie;
	d_list_t *rlink;

	rlink = d_hash_rec_find(&fs_handle->dfpi_iet,
				&request->ir_inode_num,
				sizeof(request->ir_inode_num));
	if (!rlink)
		return ENOENT;

	ie = container_of(rlink, struct dfuse_inode_entry, ie_htl);

	request->ir_inode = ie;
	return 0;
}

static void
drop_ino_ref(struct dfuse_projection_info *fs_handle, ino_t ino)
{
	d_list_t *rlink;

	rlink = d_hash_rec_find(&fs_handle->dfpi_iet, &ino, sizeof(ino));

	if (!rlink) {
		DFUSE_TRA_ERROR(fs_handle, "Could not find entry %lu", ino);
		return;
	}
	d_hash_rec_ndecref(&fs_handle->dfpi_iet, 2, rlink);
}

void
ie_close(struct dfuse_projection_info *fs_handle, struct dfuse_inode_entry *ie)
{
	int			rc;
	int			ref = atomic_load_consume(&ie->ie_ref);

	DFUSE_TRA_DEBUG(ie, "closing, inode %lu ref %u, name '%s', parent %lu",
			ie->ie_stat.st_ino, ref, ie->ie_name, ie->ie_parent);

	D_ASSERT(ref == 0);

	if (ie->ie_parent != 0) {
		drop_ino_ref(fs_handle, ie->ie_parent);
	}

	rc = dfs_release(ie->ie_obj);
	if (rc != -DER_SUCCESS) {
		DFUSE_TRA_ERROR(ie, "dfs_release failed: %d", rc);
	}
	D_FREE(ie);
}