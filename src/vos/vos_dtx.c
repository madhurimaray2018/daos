/**
 * (C) Copyright 2019 Intel Corporation.
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
/**
 * This file is part of daos two-phase commit transaction.
 *
 * vos/vos_dtx.c
 */
#define D_LOGFAC	DD_FAC(vos)

#include <libpmem.h>
#include <daos_srv/vos.h>
#include "vos_layout.h"
#include "vos_internal.h"

/* Dummy offset for an aborted DTX address. */
#define DTX_UMOFF_ABORTED		1

/* 128 KB per SCM blob */
#define DTX_SCM_BLOB_SIZE		(1 << 17)

#define DTX_ACT_BLOB_MAGIC		0x14130a2b
#define DTX_CMT_BLOB_MAGIC		0x2502191c

static inline bool
dtx_is_aborted(umem_off_t umoff)
{
	return umem_off2flags(umoff) == DTX_UMOFF_ABORTED;
}

static void
dtx_set_aborted(umem_off_t *umoff)
{
	umem_off_set_null_flags(umoff, DTX_UMOFF_ABORTED);
}

static inline int
dtx_inprogress(struct vos_dtx_act_ent *dae, int pos)
{
	if (dae != NULL)
		D_DEBUG(DB_TRACE, "Hit uncommitted DTX "DF_DTI" at %d\n",
			DP_DTI(&DAE_XID(dae)), pos);
	else
		D_DEBUG(DB_TRACE, "Hit uncommitted (unknown) DTX at %d\n", pos);

	return -DER_INPROGRESS;
}

static struct dtx_batched_cleanup_blob *
dtx_bcb_alloc(struct vos_container *cont, struct vos_dtx_scm_blob *dsb)
{
	struct dtx_batched_cleanup_blob	*bcb;

	D_ALLOC(bcb, sizeof(*bcb) + sizeof(umem_off_t) * dsb->dsb_cap);
	if (bcb != NULL) {
		D_INIT_LIST_HEAD(&bcb->bcb_dce_list);
		bcb->bcb_dsb_off = umem_ptr2off(vos_cont2umm(cont), dsb);
		d_list_add_tail(&bcb->bcb_cont_link,
				&cont->vc_batched_cleanup_list);
	}

	return bcb;
}

static void
dtx_batched_cleanup(struct vos_container *cont,
		    struct dtx_batched_cleanup_blob *bcb, umem_off_t *next_dsb)
{
	struct umem_instance	*umm = vos_cont2umm(cont);
	struct vos_cont_df	*cont_df = cont->vc_cont_df;
	struct vos_dtx_scm_blob	*dsb;
	struct vos_dtx_scm_blob	*tmp;
	struct vos_dtx_cmt_ent	*dce;
	int			 i;

	dsb = umem_off2ptr(umm, bcb->bcb_dsb_off);
	for (i = 0; i < dsb->dsb_index; i++) {
		if (!dtx_is_null(bcb->bcb_recs[i]))
			umem_free(umm, bcb->bcb_recs[i]);
	}

	if (next_dsb != NULL)
		*next_dsb = dsb->dsb_next;

	tmp = umem_off2ptr(umm, dsb->dsb_prev);
	if (tmp != NULL) {
		umem_tx_add_ptr(umm, &tmp->dsb_next, sizeof(tmp->dsb_next));
		tmp->dsb_next = dsb->dsb_next;
	}

	tmp = umem_off2ptr(umm, dsb->dsb_next);
	if (tmp != NULL) {
		umem_tx_add_ptr(umm, &tmp->dsb_prev, sizeof(tmp->dsb_prev));
		tmp->dsb_prev = dsb->dsb_prev;
	}

	if (cont_df->cd_dtx_active_head == bcb->bcb_dsb_off) {
		umem_tx_add_ptr(umm, &cont_df->cd_dtx_active_head,
				sizeof(cont_df->cd_dtx_active_head));
		cont_df->cd_dtx_active_head = dsb->dsb_next;
	}

	if (cont_df->cd_dtx_active_tail == bcb->bcb_dsb_off) {
		umem_tx_add_ptr(umm, &cont_df->cd_dtx_active_tail,
				sizeof(cont_df->cd_dtx_active_tail));
		cont_df->cd_dtx_active_tail = dsb->dsb_prev;
	}

	while ((dce = d_list_pop_entry(&bcb->bcb_dce_list,
				       struct vos_dtx_cmt_ent,
				       dce_bcb_link)) != NULL)
		; /* NOP */

	d_list_del(&bcb->bcb_cont_link);
	umem_free(umm, bcb->bcb_dsb_off);
	D_FREE(bcb);
}

static int
dtx_hkey_size(void)
{
	return sizeof(struct dtx_id);
}

static void
dtx_hkey_gen(struct btr_instance *tins, d_iov_t *key_iov, void *hkey)
{
	D_ASSERT(key_iov->iov_len == sizeof(struct dtx_id));

	memcpy(hkey, key_iov->iov_buf, key_iov->iov_len);
}

static int
dtx_hkey_cmp(struct btr_instance *tins, struct btr_record *rec, void *hkey)
{
	struct dtx_id	*hkey1 = (struct dtx_id *)&rec->rec_hkey[0];
	struct dtx_id	*hkey2 = (struct dtx_id *)hkey;
	int		 rc;

	rc = memcmp(hkey1, hkey2, sizeof(struct dtx_id));

	return dbtree_key_cmp_rc(rc);
}

static int
dtx_act_ent_alloc(struct btr_instance *tins, d_iov_t *key_iov,
		  d_iov_t *val_iov, struct btr_record *rec)
{
	struct vos_dtx_act_ent	*dae = val_iov->iov_buf;

	rec->rec_off = umem_ptr2off(&tins->ti_umm, dae);

	return 0;
}

static int
dtx_act_ent_free(struct btr_instance *tins, struct btr_record *rec,
		 void *args)
{
	struct vos_dtx_act_ent	*dae;

	dae = umem_off2ptr(&tins->ti_umm, rec->rec_off);
	rec->rec_off = UMOFF_NULL;

	if (args != NULL) {
		/* Return the record addreass (offset in DRAM).
		 * The caller will release it after using.
		 */
		D_ASSERT(dae != NULL);
		*(struct vos_dtx_act_ent **)args = dae;
	} else if (dae != NULL) {
		D_FREE(dae->dae_records);
		D_FREE_PTR(dae);
	}

	return 0;
}

static int
dtx_act_ent_fetch(struct btr_instance *tins, struct btr_record *rec,
		  d_iov_t *key_iov, d_iov_t *val_iov)
{
	struct vos_dtx_act_ent	*dae;

	D_ASSERT(val_iov != NULL);

	dae = umem_off2ptr(&tins->ti_umm, rec->rec_off);
	d_iov_set(val_iov, dae, sizeof(*dae));

	return 0;
}

static int
dtx_act_ent_update(struct btr_instance *tins, struct btr_record *rec,
		   d_iov_t *key, d_iov_t *val)
{
	D_ASSERTF(0, "Should never been called\n");
	return 0;
}

static btr_ops_t dtx_active_btr_ops = {
	.to_hkey_size	= dtx_hkey_size,
	.to_hkey_gen	= dtx_hkey_gen,
	.to_hkey_cmp	= dtx_hkey_cmp,
	.to_rec_alloc	= dtx_act_ent_alloc,
	.to_rec_free	= dtx_act_ent_free,
	.to_rec_fetch	= dtx_act_ent_fetch,
	.to_rec_update	= dtx_act_ent_update,
};

static int
dtx_cmt_ent_alloc(struct btr_instance *tins, d_iov_t *key_iov,
		  d_iov_t *val_iov, struct btr_record *rec)
{
	struct vos_container	*cont = tins->ti_priv;
	struct vos_dtx_cmt_ent	*dce = val_iov->iov_buf;

	rec->rec_off = umem_ptr2off(&tins->ti_umm, dce);
	if (!cont->vc_reindex_cmt_dtx || dce->dce_reindex) {
		d_list_add_tail(&dce->dce_committed_link,
				&cont->vc_dtx_committed_list);
		cont->vc_dtx_committed_count++;
	} else {
		d_list_add_tail(&dce->dce_committed_link,
				&cont->vc_dtx_committed_tmp_list);
		cont->vc_dtx_committed_tmp_count++;
	}

	return 0;
}

static int
dtx_cmt_ent_free(struct btr_instance *tins, struct btr_record *rec,
		 void *args)
{
	struct vos_container	*cont = tins->ti_priv;
	struct vos_dtx_cmt_ent	*dce;
	int			 rc = 0;

	dce = umem_off2ptr(&tins->ti_umm, rec->rec_off);
	D_ASSERT(dce != NULL);

	rec->rec_off = UMOFF_NULL;
	d_list_del(&dce->dce_committed_link);
	if (!cont->vc_reindex_cmt_dtx || dce->dce_reindex)
		cont->vc_dtx_committed_count--;
	else
		cont->vc_dtx_committed_tmp_count--;

	/* The committed DTX entry may be in aggregation now, if related
	 * active DTX entry is waitting for batched cleanup, then cleanup
	 * it by force before destroy the committed DTX entry (mainly for
	 * DTX aggregation).
	 */
	if (!d_list_empty(&dce->dce_bcb_link)) {
		struct umem_instance		*umm = vos_cont2umm(cont);
		struct dtx_batched_cleanup_blob	*bcb = dce->dce_bcb;
		struct vos_dtx_scm_blob		*dsb;
		int				 idx;

		D_ASSERT(bcb != NULL);

		idx = dce->dce_index;
		dsb = umem_off2ptr(umm, bcb->bcb_dsb_off);

		if (dsb->dsb_active_data[idx].dae_flags != DTX_EF_INVALID) {
			struct vos_dtx_act_ent_df	*dae_df;

			rc = vos_tx_begin(umm);
			if (rc != 0)
				return rc;

			dae_df = &dsb->dsb_active_data[idx];
			umem_tx_add_ptr(umm, &dae_df->dae_flags,
					sizeof(dae_df->dae_flags));
			dae_df->dae_flags = DTX_EF_INVALID;

			D_ASSERT(dsb->dsb_count > bcb->bcb_dae_count);
			umem_tx_add_ptr(umm, &dsb->dsb_count,
					sizeof(dsb->dsb_count));
			dsb->dsb_count--;

			if (!dtx_is_null(bcb->bcb_recs[idx])) {
				umem_free(umm, bcb->bcb_recs[idx]);
				bcb->bcb_recs[idx] = UMOFF_NULL;
			}

			rc = vos_tx_end(umm, 0);
		}

		d_list_del(&dce->dce_bcb_link);
	}

	D_FREE_PTR(dce);

	return rc;
}

static int
dtx_cmt_ent_fetch(struct btr_instance *tins, struct btr_record *rec,
		  d_iov_t *key_iov, d_iov_t *val_iov)
{
	if (val_iov != NULL) {
		struct vos_dtx_cmt_ent	*dce;

		dce = umem_off2ptr(&tins->ti_umm, rec->rec_off);
		d_iov_set(val_iov, dce, sizeof(*dce));
	}

	return 0;
}

static int
dtx_cmt_ent_update(struct btr_instance *tins, struct btr_record *rec,
		   d_iov_t *key, d_iov_t *val)
{
	struct vos_dtx_cmt_ent	*dce = val->iov_buf;

	dce->dce_exist = 1;

	return 0;
}

static btr_ops_t dtx_committed_btr_ops = {
	.to_hkey_size	= dtx_hkey_size,
	.to_hkey_gen	= dtx_hkey_gen,
	.to_hkey_cmp	= dtx_hkey_cmp,
	.to_rec_alloc	= dtx_cmt_ent_alloc,
	.to_rec_free	= dtx_cmt_ent_free,
	.to_rec_fetch	= dtx_cmt_ent_fetch,
	.to_rec_update	= dtx_cmt_ent_update,
};

int
vos_dtx_table_register(void)
{
	int	rc;

	rc = dbtree_class_register(VOS_BTR_DTX_ACT_TABLE, 0,
				   &dtx_active_btr_ops);
	if (rc != 0) {
		D_ERROR("Failed to register DTX active dbtree: %d\n", rc);
		return rc;
	}

	rc = dbtree_class_register(VOS_BTR_DTX_CMT_TABLE, 0,
				   &dtx_committed_btr_ops);
	if (rc != 0)
		D_ERROR("Failed to register DTX committed dbtree: %d\n", rc);

	return rc;
}

void
vos_dtx_table_destroy(struct umem_instance *umm, struct vos_cont_df *cont_df)
{
	struct vos_dtx_scm_blob		*dsb;
	umem_off_t			 dsb_off;

	while (!dtx_is_null(cont_df->cd_dtx_committed_head)) {
		dsb_off = cont_df->cd_dtx_committed_head;
		dsb = umem_off2ptr(umm, dsb_off);
		cont_df->cd_dtx_committed_head = dsb->dsb_next;
		umem_free(umm, dsb_off);
	}

	cont_df->cd_dtx_committed_head = UMOFF_NULL;
	cont_df->cd_dtx_committed_tail = UMOFF_NULL;

	while (!dtx_is_null(cont_df->cd_dtx_active_head)) {
		dsb_off = cont_df->cd_dtx_active_head;
		dsb = umem_off2ptr(umm, dsb_off);
		cont_df->cd_dtx_active_head = dsb->dsb_next;
		umem_free(umm, dsb_off);
	}

	cont_df->cd_dtx_active_head = UMOFF_NULL;
	cont_df->cd_dtx_active_tail = UMOFF_NULL;
}

static int
dtx_ilog_rec_release(struct umem_instance *umm, struct vos_container *cont,
		     struct vos_dtx_record_df *rec, struct vos_dtx_act_ent *dae,
		     bool abort)
{
	struct ilog_df		*ilog;
	daos_handle_t		 loh;
	struct ilog_desc_cbs	 cbs;
	struct ilog_id		 id;
	int			 rc;

	ilog = umem_off2ptr(umm, rec->dr_record);

	vos_ilog_desc_cbs_init(&cbs, vos_cont2hdl(cont));
	rc = ilog_open(umm, ilog, &cbs, &loh);
	if (rc != 0)
		return rc;

	id.id_epoch = DAE_EPOCH(dae);
	id.id_tx_id = dae->dae_df_off;

	if (abort)
		rc = ilog_abort(loh, &id);
	else
		rc = ilog_persist(loh, &id);

	ilog_close(loh);
	return rc;
}

static void
do_dtx_rec_release(struct umem_instance *umm, struct vos_container *cont,
		   struct vos_dtx_act_ent *dae, struct vos_dtx_record_df *rec,
		   struct vos_dtx_record_df *rec_df, int index, bool abort,
		   bool *sync)
{
	if (dtx_is_null(rec->dr_record))
		return;

	/* Has been deregistered. */
	if (rec_df != NULL && dtx_is_null(rec_df[index].dr_record))
		return;

	switch (rec->dr_type) {
	case DTX_RT_ILOG: {
		dtx_ilog_rec_release(umm, cont, rec, dae, abort);
		break;
	}
	case DTX_RT_SVT: {
		struct vos_irec_df	*svt;

		svt = umem_off2ptr(umm, rec->dr_record);
		if (abort) {
			if (DAE_INDEX(dae) != -1)
				umem_tx_add_ptr(umm, &svt->ir_dtx,
						sizeof(svt->ir_dtx));
			dtx_set_aborted(&svt->ir_dtx);
		} else {
			umem_tx_add_ptr(umm, &svt->ir_dtx, sizeof(svt->ir_dtx));
			svt->ir_dtx = UMOFF_NULL;
			*sync = false;
		}
		break;
	}
	case DTX_RT_EVT: {
		struct evt_desc		*evt;

		evt = umem_off2ptr(umm, rec->dr_record);
		if (abort) {
			if (DAE_INDEX(dae) != -1)
				umem_tx_add_ptr(umm, &evt->dc_dtx,
						sizeof(evt->dc_dtx));
			dtx_set_aborted(&evt->dc_dtx);
		} else {
			umem_tx_add_ptr(umm, &evt->dc_dtx, sizeof(evt->dc_dtx));
			evt->dc_dtx = UMOFF_NULL;
			*sync = false;
		}
		break;
	}
	default:
		D_ERROR(DF_UOID" unknown DTX "DF_DTI" type %u\n",
			DP_UOID(DAE_OID(dae)), DP_DTI(&DAE_XID(dae)),
			rec->dr_type);
		break;
	}
}

static void
dtx_rec_release(struct umem_instance *umm, struct vos_container *cont,
		struct vos_dtx_act_ent *dae, bool abort,
		struct dtx_batched_cleanup_blob **bcb_p)
{
	struct vos_dtx_act_ent_df	*dae_df;
	struct vos_dtx_record_df	*rec_df = NULL;
	struct vos_dtx_scm_blob		*dsb;
	struct dtx_batched_cleanup_blob	*bcb;
	int				 count;
	int				 i;
	bool				 sync = false;

	/* XXX: To avoid complex ilog related status check, let's
	 *	check next record directly. For pure ilog touched
	 *	DTX case, we do NOT handle it via batched cleanup.
	 *
	 *	For abort case, cleanup DTX entry synchronously.
	 */
	dae_df = umem_off2ptr(umm, dae->dae_df_off);
	if ((DAE_REC_CNT(dae) > 0 &&
	     DAE_REC_INLINE(dae)[0].dr_type == DTX_RT_ILOG) || abort)
		sync = true;

	if (dae->dae_records != NULL) {
		D_ASSERT(DAE_REC_CNT(dae) > DTX_INLINE_REC_CNT);

		if (dae_df != NULL &&
		    dae_df->dae_layout_gen != DAE_LAYOUT_GEN(dae))
			rec_df = umem_off2ptr(umm, dae_df->dae_rec_off);

		for (i = DAE_REC_CNT(dae) - DTX_INLINE_REC_CNT - 1; i >= 0; i--)
			do_dtx_rec_release(umm, cont, dae, &dae->dae_records[i],
					   rec_df, i, abort, &sync);

		D_FREE(dae->dae_records);
		dae->dae_records = NULL;
		dae->dae_rec_cap = 0;
	}

	if (dae_df != NULL && dae_df->dae_layout_gen != DAE_LAYOUT_GEN(dae))
		rec_df = dae_df->dae_rec_inline;

	if (DAE_REC_CNT(dae) > DTX_INLINE_REC_CNT)
		count = DTX_INLINE_REC_CNT;
	else
		count = DAE_REC_CNT(dae);

	for (i = count - 1; i >= 0; i--)
		do_dtx_rec_release(umm, cont, dae, &DAE_REC_INLINE(dae)[i],
				   rec_df, i, abort, &sync);

	/* Non-prepared case, not need to change on-disk things. */
	if (DAE_INDEX(dae) == -1) {
		D_ASSERT(abort);
		D_FREE_PTR(dae);
		return;
	}

	D_ASSERT(dae_df != NULL);

	dsb = dae->dae_dsb;
	D_ASSERT(dsb->dsb_magic == DTX_ACT_BLOB_MAGIC);

	bcb = dae->dae_bcb;
	D_ASSERT(bcb != NULL);
	D_ASSERT(bcb->bcb_dae_count > 0);

	if (sync) {
		if (DAE_INDEX(dae) != -1 &&
		    (bcb->bcb_dae_count > 1 || dsb->dsb_index < dsb->dsb_cap)) {
			umem_tx_add_ptr(umm, &dae_df->dae_flags,
					sizeof(dae_df->dae_flags));
			/* Mark the DTX entry as invalid in SCM. */
			dae_df->dae_flags = DTX_EF_INVALID;

			umem_tx_add_ptr(umm, &dsb->dsb_count,
					sizeof(dsb->dsb_count));
			dsb->dsb_count--;
		}

		if (!dtx_is_null(dae_df->dae_rec_off))
			umem_free(umm, dae_df->dae_rec_off);
	} else {
		bcb->bcb_recs[DAE_INDEX(dae)] = DAE_REC_OFF(dae);
	}

	bcb->bcb_dae_count--;
	if (bcb->bcb_dae_count == 0 && dsb->dsb_index >= dsb->dsb_cap) {
		dtx_batched_cleanup(cont, bcb, NULL);
	} else if (!sync) {
		D_ASSERT(bcb_p != NULL);

		*bcb_p = bcb;
	}

	D_FREE_PTR(dae);
}

static int
vos_dtx_commit_one(struct vos_container *cont, struct dtx_id *dti,
		   daos_epoch_t epoch, struct vos_dtx_cmt_ent **dce_p)
{
	struct umem_instance		*umm = vos_cont2umm(cont);
	struct vos_dtx_act_ent		*dae = NULL;
	struct vos_dtx_cmt_ent		*dce = NULL;
	d_iov_t				 kiov;
	d_iov_t				 riov;
	int				 rc = 0;

	d_iov_set(&kiov, dti, sizeof(*dti));
	if (epoch == 0) {
		rc = dbtree_delete(cont->vc_dtx_active_hdl, BTR_PROBE_EQ,
				   &kiov, &dae);
		if (rc == -DER_NONEXIST) {
			rc = dbtree_lookup(cont->vc_dtx_committed_hdl,
					   &kiov, NULL);
			goto out;
		}

		if (rc != 0)
			goto out;
	}

	D_ALLOC_PTR(dce);
	if (dce == NULL)
		D_GOTO(out, rc = -DER_NOMEM);

	D_INIT_LIST_HEAD(&dce->dce_bcb_link);
	DCE_XID(dce) = *dti;
	DCE_EPOCH(dce) = epoch != 0 ? epoch : DAE_EPOCH(dae);
	dce->dce_index = epoch != 0 ? -1 : DAE_INDEX(dae);
	dce->dce_reindex = 0;

	d_iov_set(&riov, dce, sizeof(*dce));
	rc = dbtree_upsert(cont->vc_dtx_committed_hdl, BTR_PROBE_EQ,
			   DAOS_INTENT_UPDATE, &kiov, &riov);
	if (rc != 0 || epoch != 0)
		goto out;

	vos_dtx_del_cc(cont, dti);

	/* XXX: Only mark the DTX as DTX_ST_COMMITTED (when commit) is not
	 *	enough. Otherwise, some subsequent modification may change
	 *	related data record's DTX reference or remove related data
	 *	record as to the current DTX will have invalid reference(s)
	 *	via its DTX record(s).
	 */
	dtx_rec_release(umm, cont, dae, false, &dce->dce_bcb);
	if (dce->dce_bcb != NULL)
		d_list_add_tail(&dce->dce_bcb_link,
				&dce->dce_bcb->bcb_dce_list);

out:
	D_DEBUG(DB_TRACE, "Commit the DTX "DF_DTI": rc = %d\n",
		DP_DTI(dti), rc);
	if (rc != 0)
		D_FREE_PTR(dce);
	else
		*dce_p = dce;

	return rc;
}

static int
vos_dtx_abort_one(struct vos_container *cont, daos_epoch_t epoch,
		  struct dtx_id *dti)
{
	struct vos_dtx_act_ent	*dae;
	d_iov_t			 riov;
	d_iov_t			 kiov;
	dbtree_probe_opc_t	 opc = BTR_PROBE_EQ;
	int			 rc;

	d_iov_set(&kiov, dti, sizeof(*dti));
	if (epoch != 0) {
		d_iov_set(&riov, NULL, 0);
		rc = dbtree_lookup(cont->vc_dtx_active_hdl, &kiov, &riov);
		if (rc != 0)
			goto out;

		dae = (struct vos_dtx_act_ent *)riov.iov_buf;
		if (DAE_EPOCH(dae) > epoch)
			D_GOTO(out, rc = -DER_NONEXIST);

		opc = BTR_PROBE_BYPASS;
	}

	rc = dbtree_delete(cont->vc_dtx_active_hdl, opc, &kiov, &dae);
	if (rc == 0)
		dtx_rec_release(vos_cont2umm(cont), cont, dae, true, NULL);

out:
	D_DEBUG(DB_TRACE, "Abort the DTX "DF_DTI": rc = %d\n", DP_DTI(dti), rc);

	return rc;
}

static bool
vos_dtx_is_normal_entry(struct umem_instance *umm, umem_off_t entry)
{
	if (dtx_is_null(entry) || dtx_is_aborted(entry))
		return false;

	return true;
}

static int
vos_dtx_extend_act_table(struct vos_container *cont)
{
	struct umem_instance		*umm = vos_cont2umm(cont);
	struct vos_cont_df		*cont_df = cont->vc_cont_df;
	struct dtx_batched_cleanup_blob	*bcb;
	struct vos_dtx_scm_blob		*dsb;
	struct vos_dtx_scm_blob		*tmp;
	umem_off_t			 dsb_off;

	dsb_off = umem_zalloc(umm, DTX_SCM_BLOB_SIZE);
	if (dtx_is_null(dsb_off)) {
		D_ERROR("No space when create actvie DTX table.\n");
		return -DER_NOSPACE;
	}

	dsb = umem_off2ptr(umm, dsb_off);
	dsb->dsb_magic = DTX_ACT_BLOB_MAGIC;
	dsb->dsb_cap = (DTX_SCM_BLOB_SIZE - sizeof(struct vos_dtx_scm_blob)) /
			sizeof(struct vos_dtx_act_ent_df);

	bcb = dtx_bcb_alloc(cont, dsb);
	if (bcb == NULL) {
		umem_free(umm, dsb_off);
		return -DER_NOMEM;
	}

	tmp = umem_off2ptr(umm, cont_df->cd_dtx_active_tail);
	if (tmp == NULL) {
		D_ASSERT(dtx_is_null(cont_df->cd_dtx_active_head));

		/* cd_dtx_active_tail is next to cd_dtx_active_head */
		umem_tx_add_ptr(umm, &cont_df->cd_dtx_active_head,
				sizeof(cont_df->cd_dtx_active_head) +
				sizeof(cont_df->cd_dtx_active_tail));
		cont_df->cd_dtx_active_head = dsb_off;
	} else {
		umem_tx_add_ptr(umm, &tmp->dsb_next, sizeof(tmp->dsb_next));
		tmp->dsb_next = dsb_off;

		dsb->dsb_prev = cont_df->cd_dtx_active_tail;
		umem_tx_add_ptr(umm, &cont_df->cd_dtx_active_tail,
				sizeof(cont_df->cd_dtx_active_tail));
	}

	cont_df->cd_dtx_active_tail = dsb_off;

	return 0;
}

static int
vos_dtx_alloc(struct umem_instance *umm, struct dtx_handle *dth)
{
	struct vos_dtx_act_ent		*dae = NULL;
	struct dtx_batched_cleanup_blob	*bcb = NULL;
	struct vos_container		*cont;
	struct vos_cont_df		*cont_df;
	struct vos_dtx_scm_blob		*dsb;
	d_iov_t				 kiov;
	d_iov_t				 riov;
	int				 rc;

	cont = vos_hdl2cont(dth->dth_coh);
	D_ASSERT(cont != NULL);

	cont_df = cont->vc_cont_df;
	dth->dth_gen = cont->vc_dtx_resync_gen;

	D_ALLOC_PTR(dae);
	if (dae == NULL)
		return -DER_NOMEM;

	DAE_XID(dae) = dth->dth_xid;
	DAE_OID(dae) = dth->dth_oid;
	DAE_EPOCH(dae) = dth->dth_epoch;
	DAE_FLAGS(dae) = dth->dth_leader ? DTX_EF_LEADER : 0;
	DAE_SRV_GEN(dae) = dth->dth_gen;
	DAE_LAYOUT_GEN(dae) = dth->dth_gen;

	dsb = umem_off2ptr(umm, cont_df->cd_dtx_active_tail);
	if (dsb == NULL || dsb->dsb_index >= dsb->dsb_cap) {
		rc = vos_dtx_extend_act_table(cont);
		if (rc != 0)
			goto out;

		dsb = umem_off2ptr(umm, cont_df->cd_dtx_active_tail);
	}

	bcb = d_list_entry(cont->vc_batched_cleanup_list.prev,
			   struct dtx_batched_cleanup_blob, bcb_cont_link);

	/* Set it as dsb::dsb_index via vos_dtx_prepared(). */
	DAE_INDEX(dae) = -1;
	dae->dae_df_off = umem_ptr2off(umm, dsb) +
			offsetof(struct vos_dtx_scm_blob, dsb_active_data) +
			sizeof(struct vos_dtx_act_ent_df) * dsb->dsb_index;
	dae->dae_dsb = dsb;
	dae->dae_bcb = bcb;

	d_iov_set(&kiov, &DAE_XID(dae), sizeof(DAE_XID(dae)));
	d_iov_set(&riov, dae, sizeof(*dae));
	rc = dbtree_upsert(cont->vc_dtx_active_hdl, BTR_PROBE_EQ,
			   DAOS_INTENT_UPDATE, &kiov, &riov);
	if (rc == 0)
		dth->dth_ent = dae;

out:
	if (rc != 0) {
		D_FREE_PTR(dae);
		D_FREE_PTR(bcb);
	}

	return rc;
}

static int
vos_dtx_append(struct umem_instance *umm, struct dtx_handle *dth,
	       umem_off_t record, uint32_t type)
{
	struct vos_dtx_act_ent		*dae = dth->dth_ent;
	struct vos_dtx_record_df	*rec;

	D_ASSERT(dae != NULL);

	if (DAE_REC_CNT(dae) < DTX_INLINE_REC_CNT) {
		rec = &DAE_REC_INLINE(dae)[DAE_REC_CNT(dae)];
	} else {
		if (DAE_REC_CNT(dae) >= dae->dae_rec_cap + DTX_INLINE_REC_CNT) {
			int	count;

			if (dae->dae_rec_cap == 0)
				count = DTX_REC_CAP_DEFAULT;
			else
				count = dae->dae_rec_cap * 2;

			D_ALLOC(rec, sizeof(*rec) * count);
			if (rec == NULL)
				return -DER_NOMEM;

			if (dae->dae_records != NULL) {
				memcpy(rec, dae->dae_records,
				       sizeof(*rec) * dae->dae_rec_cap);
				D_FREE(dae->dae_records);
			}

			dae->dae_records = rec;
			dae->dae_rec_cap = count;
		}

		rec = &dae->dae_records[DAE_REC_CNT(dae) - DTX_INLINE_REC_CNT];
	}

	rec->dr_type = type;
	rec->dr_record = record;

	/* The rec_cnt on-disk value will be refreshed via vos_dtx_prepared() */
	DAE_REC_CNT(dae)++;

	return 0;
}

int
vos_dtx_check_availability(struct umem_instance *umm, daos_handle_t coh,
			   umem_off_t entry, umem_off_t record, uint32_t intent,
			   uint32_t type)
{
	struct dtx_handle		*dth = vos_dth_get();
	struct vos_container		*cont;
	struct vos_dtx_act_ent		*dae = NULL;
	struct vos_dtx_act_ent_df	*dae_df = NULL;
	d_iov_t				 kiov;
	d_iov_t				 riov;
	int				 rc;

	switch (type) {
	case DTX_RT_SVT:
	case DTX_RT_EVT:
	case DTX_RT_ILOG:
		break;
	default:
		D_ERROR("Unexpected DTX type %u\n", type);
		/* Everything is available to PURGE, even if it belongs to some
		 * uncommitted DTX that may be garbage because of corruption.
		 */
		if (intent == DAOS_INTENT_PURGE)
			return ALB_AVAILABLE_DIRTY;

		return -DER_INVAL;
	}

	if (intent == DAOS_INTENT_CHECK || intent == DAOS_INTENT_COS) {
		if (dtx_is_aborted(entry))
			return ALB_UNAVAILABLE;

		return ALB_AVAILABLE_CLEAN;
	}

	/* Committed */
	if (dtx_is_null(entry))
		return ALB_AVAILABLE_CLEAN;

	if (intent == DAOS_INTENT_PURGE)
		return ALB_AVAILABLE_DIRTY;

	/* Aborted */
	if (dtx_is_aborted(entry))
		return ALB_UNAVAILABLE;

	/* The DTX owner can always see the DTX. */
	if (dth != NULL && dth->dth_ent != NULL) {
		dae = dth->dth_ent;
		if (dae->dae_df_off == entry)
			return ALB_AVAILABLE_CLEAN;
	}

	cont = vos_hdl2cont(coh);
	D_ASSERT(cont != NULL);

	dae_df = umem_off2ptr(umm, entry);
	d_iov_set(&kiov, &dae_df->dae_xid, sizeof(dae_df->dae_xid));
	d_iov_set(&riov, NULL, 0);
	rc = dbtree_lookup(cont->vc_dtx_active_hdl, &kiov, &riov);
	if (rc != 0) {
		if (rc == -DER_NONEXIST)
			/* Handle it as aborted one. */
			return ALB_UNAVAILABLE;

		D_WARN("Cannot find active DTX entry for "DF_DTI"\n",
		       DP_DTI(&dae_df->dae_xid));
		return rc;
	}

	dae = (struct vos_dtx_act_ent *)riov.iov_buf;

	rc = vos_dtx_lookup_cc(coh, &DAE_XID(dae));
	if (rc == 0)
		return ALB_AVAILABLE_CLEAN;

	if (rc != -DER_NONEXIST)
		return rc;

	/* The followings are for non-committable cases. */
	if (!(DAE_FLAGS(dae) & DTX_EF_LEADER) ||
	    DAOS_FAIL_CHECK(DAOS_VOS_NON_LEADER)) {
		/* Unavailable for rebuild case. */
		if (intent == DAOS_INTENT_REBUILD)
			return ALB_UNAVAILABLE;

		/* Non-leader and non-rebuild case, return
		 * -DER_INPROGRESS, then the caller will retry
		 * the RPC with leader replica.
		 */
		return dtx_inprogress(dae, 1);
	}

	/* For leader, non-committed DTX is unavailable. */
	return ALB_UNAVAILABLE;

}

umem_off_t
vos_dtx_get(void)
{
	struct dtx_handle	*dth = vos_dth_get();
	struct vos_dtx_act_ent	*dae;

	if (dth == NULL || dth->dth_solo)
		return UMOFF_NULL;

	dae = dth->dth_ent;

	return dae->dae_df_off;
}

/* The caller has started PMDK transaction. */
int
vos_dtx_register_record(struct umem_instance *umm, umem_off_t record,
			uint32_t type, umem_off_t *tx_id)
{
	struct dtx_handle	*dth = vos_dth_get();
	int			 rc = 0;

	/* For single participator case, we only need the DTX entry
	 * without DTX records for related targets to be modified.
	 */
	if (dth == NULL || dth->dth_solo) {
		*tx_id = UMOFF_NULL;
		return 0;
	}

	if (dth->dth_ent == NULL) {
		rc = vos_dtx_alloc(umm, dth);
		if (rc != 0)
			return rc;
	}

	rc = vos_dtx_append(umm, dth, record, type);
	if (rc == 0) {
		struct vos_dtx_act_ent	*dae = dth->dth_ent;

		/* Incarnation log entry implies a share */
		*tx_id = dae->dae_df_off;
	}

	return rc;
}

/* The caller has started PMDK transaction. */
void
vos_dtx_deregister_record(struct umem_instance *umm, daos_handle_t coh,
			  umem_off_t entry, umem_off_t record)
{
	struct vos_container		*cont;
	struct vos_dtx_act_ent_df	*dae_df;
	struct vos_dtx_record_df	*rec_df;
	int				 count;
	int				 rc;
	int				 i;

	if (!vos_dtx_is_normal_entry(umm, entry))
		return;

	dae_df = umem_off2ptr(umm, entry);
	if (daos_is_zero_dti(&dae_df->dae_xid) ||
	    dae_df->dae_flags & DTX_EF_INVALID)
		return;

	cont = vos_hdl2cont(coh);
	if (cont != NULL) {
		struct vos_dtx_act_ent		*dae;
		d_iov_t				 kiov;
		d_iov_t				 riov;

		d_iov_set(&kiov, &dae_df->dae_xid, sizeof(dae_df->dae_xid));
		d_iov_set(&riov, NULL, 0);
		rc = dbtree_lookup(cont->vc_dtx_active_hdl, &kiov, &riov);
		if (rc != 0) {
			D_WARN("NOT find active DTX entry when deregister for "
			       DF_DTI"\n", DP_DTI(&dae_df->dae_xid));
			return;
		}

		dae = (struct vos_dtx_act_ent *)riov.iov_buf;
		if (DAE_REC_CNT(dae) > DTX_INLINE_REC_CNT)
			count = DTX_INLINE_REC_CNT;
		else
			count = DAE_REC_CNT(dae);

		for (i = 0; i < count; i++) {
			if (record == DAE_REC_INLINE(dae)[i].dr_record) {
				DAE_REC_INLINE(dae)[i].dr_record = UMOFF_NULL;
				goto handle_df;
			}
		}

		for (i = 0; i < DAE_REC_CNT(dae) - DTX_INLINE_REC_CNT; i++) {
			if (record == dae->dae_records[i].dr_record) {
				dae->dae_records[i].dr_record = UMOFF_NULL;
				goto handle_df;
			}
		}

		/* Not found */
		return;
	}

handle_df:
	if (dae_df->dae_rec_cnt > DTX_INLINE_REC_CNT)
		count = DTX_INLINE_REC_CNT;
	else
		count = dae_df->dae_rec_cnt;

	rec_df = dae_df->dae_rec_inline;
	for (i = 0; i < count; i++) {
		if (rec_df[i].dr_record == record) {
			rec_df[i].dr_record = UMOFF_NULL;
			if (cont == NULL)
				dae_df->dae_layout_gen++;
			return;
		}
	}

	rec_df = umem_off2ptr(umm, dae_df->dae_rec_off);

	/* Not found */
	if (rec_df == NULL)
		return;

	for (i = 0; i < dae_df->dae_rec_cnt - DTX_INLINE_REC_CNT; i++) {
		if (rec_df[i].dr_record == record) {
			rec_df[i].dr_record = UMOFF_NULL;
			if (cont == NULL)
				dae_df->dae_layout_gen++;
			return;
		}
	}
}

int
vos_dtx_prepared(struct dtx_handle *dth)
{
	struct vos_container	*cont;

	cont = vos_hdl2cont(dth->dth_coh);
	D_ASSERT(cont != NULL);

	/* The caller has already started the PMDK transaction
	 * and add the DTX into the PMDK transaction.
	 */

	if (dth->dth_solo) {
		vos_dtx_commit_internal(cont, &dth->dth_xid, 1, dth->dth_epoch);
		dth->dth_sync = 1;
	} else {
		struct umem_instance		*umm = vos_cont2umm(cont);
		struct vos_dtx_act_ent		*dae = dth->dth_ent;
		struct vos_dtx_scm_blob		*dsb = dae->dae_dsb;

		/* If some DTXs share something (object/key) with others,
		 * or punch object that is quite possible affect subsequent
		 * operations, then synchronously commit the DTX when it
		 * becomes committable to avoid availability trouble.
		 */
		if (DAE_FLAGS(dae) & DTX_EF_SHARES || DAE_DKEY_HASH(dae) == 0)
			dth->dth_sync = 1;

		if (dae->dae_records != NULL) {
			struct vos_dtx_record_df	*rec_df;
			umem_off_t			 rec_off;
			int				 size;

			size = sizeof(struct vos_dtx_record_df) *
				(DAE_REC_CNT(dae) - DTX_INLINE_REC_CNT);
			rec_off = umem_zalloc(umm, size);
			if (dtx_is_null(rec_off)) {
				D_ERROR("No space to store active DTX (3) "
					DF_DTI"\n", DP_DTI(&DAE_XID(dae)));
				return -DER_NOSPACE;
			}

			rec_df = umem_off2ptr(umm, rec_off);
			memcpy(rec_df, dae->dae_records, size);
			DAE_REC_OFF(dae) = rec_off;
		}

		DAE_INDEX(dae) = dsb->dsb_index;
		if (DAE_INDEX(dae) > 0) {
			pmem_memcpy_nodrain(umem_off2ptr(umm, dae->dae_df_off),
					    &dae->dae_base,
					    sizeof(struct vos_dtx_act_ent_df));
			/* dsb_index is next to dsb_count */
			umem_tx_add_ptr(umm, &dsb->dsb_count,
					sizeof(dsb->dsb_count) +
					sizeof(dsb->dsb_index));
		} else {
			memcpy(umem_off2ptr(umm, dae->dae_df_off),
			       &dae->dae_base,
			       sizeof(struct vos_dtx_act_ent_df));
		}

		dsb->dsb_count++;
		dsb->dsb_index++;

		dae->dae_bcb->bcb_dae_count++;
	}

	return 0;
}

static int
do_vos_dtx_check(daos_handle_t coh, struct dtx_id *dti, daos_epoch_t *epoch)
{
	struct vos_container	*cont;
	struct vos_dtx_act_ent	*dae;
	d_iov_t			 kiov;
	d_iov_t			 riov;
	int			 rc;

	cont = vos_hdl2cont(coh);
	D_ASSERT(cont != NULL);

	d_iov_set(&kiov, dti, sizeof(*dti));
	d_iov_set(&riov, NULL, 0);
	rc = dbtree_lookup(cont->vc_dtx_active_hdl, &kiov, &riov);
	if (rc == 0) {
		dae = (struct vos_dtx_act_ent *)riov.iov_buf;
		if (epoch != NULL) {
			if (*epoch == 0)
				*epoch = DAE_EPOCH(dae);
			else if (*epoch != DAE_EPOCH(dae))
				return -DER_MISMATCH;
		}

		return DTX_ST_PREPARED;
	}

	if (rc == -DER_NONEXIST) {
		rc = dbtree_lookup(cont->vc_dtx_committed_hdl, &kiov, NULL);
		if (rc == 0)
			return DTX_ST_COMMITTED;
	}

	return rc;
}

int
vos_dtx_check_resend(daos_handle_t coh, daos_unit_oid_t *oid,
		     struct dtx_id *xid, daos_epoch_t *epoch)
{
	struct vos_container	*cont;
	int			 rc;

	rc = vos_dtx_lookup_cc(coh, xid);
	if (rc == 0)
		return DTX_ST_COMMITTED;

	if (rc != -DER_NONEXIST)
		return rc;

	rc = do_vos_dtx_check(coh, xid, epoch);
	if (rc != -DER_NONEXIST)
		return rc;

	cont = vos_hdl2cont(coh);
	D_ASSERT(cont != NULL);

	if (cont->vc_reindex_cmt_dtx)
		rc = -DER_AGAIN;

	return rc;
}

int
vos_dtx_check(daos_handle_t coh, struct dtx_id *dti)
{
	return do_vos_dtx_check(coh, dti, NULL);
}

int
vos_dtx_commit_internal(struct vos_container *cont, struct dtx_id *dtis,
			int count, daos_epoch_t epoch)
{
	struct vos_cont_df		*cont_df = cont->vc_cont_df;
	struct umem_instance		*umm = vos_cont2umm(cont);
	struct vos_dtx_scm_blob		*dsb;
	struct vos_dtx_scm_blob		*dsb_prev;
	umem_off_t			 dsb_off;
	struct vos_dtx_cmt_ent_df	*dce_df;
	int				 slots = 0;
	int				 cur = 0;
	int				 rc = 0;
	int				 rc1 = 0;
	int				 i;
	int				 j;

	dsb = umem_off2ptr(umm, cont_df->cd_dtx_committed_tail);
	if (dsb != NULL)
		slots = dsb->dsb_cap - dsb->dsb_count;

	if (slots == 0)
		goto new_blob;

	umem_tx_add_ptr(umm, &dsb->dsb_count, sizeof(dsb->dsb_count));

again:
	if (slots > count)
		slots = count;

	count -= slots;

	if (slots > 1) {
		D_ALLOC(dce_df, sizeof(*dce_df) * slots);
		if (dce_df == NULL)
			return -DER_NOMEM;
	} else {
		dce_df = &dsb->dsb_commmitted_data[dsb->dsb_count];
	}

	for (i = 0, j = 0; i < slots; i++, cur++) {
		struct vos_dtx_cmt_ent	*dce = NULL;

		rc = vos_dtx_commit_one(cont, &dtis[cur], epoch, &dce);
		if (rc1 == 0)
			rc1 = rc;

		if (dce != NULL) {
			if (slots == 1)
				pmem_memcpy_nodrain(dce_df, &dce->dce_base,
						    sizeof(*dce_df));
			else
				memcpy(&dce_df[j], &dce->dce_base,
				       sizeof(dce_df[j]));
			j++;
		}
	}

	if (dce_df != &dsb->dsb_commmitted_data[dsb->dsb_count]) {
		if (j > 0)
			pmem_memcpy_nodrain(
				&dsb->dsb_commmitted_data[dsb->dsb_count],
				dce_df, sizeof(*dce_df) * j);
		D_FREE(dce_df);
	}

	if (j > 0)
		dsb->dsb_count += j;

	if (count == 0)
		return rc != 0 ? rc : rc1;

	if (j < slots) {
		slots -= j;
		goto again;
	}

new_blob:
	dsb_prev = dsb;

	/* Need new @dsb */
	dsb_off = umem_zalloc(umm, DTX_SCM_BLOB_SIZE);
	if (dtx_is_null(dsb_off)) {
		D_ERROR("No space to store committed DTX %d "DF_DTI"\n",
			count, DP_DTI(&dtis[cur]));
		return -DER_NOSPACE;
	}

	dsb = umem_off2ptr(umm, dsb_off);
	dsb->dsb_magic = DTX_CMT_BLOB_MAGIC;
	dsb->dsb_cap = (DTX_SCM_BLOB_SIZE - sizeof(struct vos_dtx_scm_blob)) /
		       sizeof(struct vos_dtx_cmt_ent_df);
	dsb->dsb_prev = umem_ptr2off(umm, dsb_prev);

	/* Not allow to commit too many DTX together. */
	D_ASSERTF(count < dsb->dsb_cap, "Too many DTX: %d/%d\n",
		  count, dsb->dsb_cap);

	if (count > 1) {
		D_ALLOC(dce_df, sizeof(*dce_df) * count);
		if (dce_df == NULL) {
			umem_free(umm, dsb_off);
			return -DER_NOMEM;
		}
	} else {
		dce_df = &dsb->dsb_commmitted_data[0];
	}

	if (dsb_prev == NULL) {
		D_ASSERT(dtx_is_null(cont_df->cd_dtx_committed_head));
		D_ASSERT(dtx_is_null(cont_df->cd_dtx_committed_tail));

		/* cd_dtx_committed_tail is next to cd_dtx_committed_head */
		umem_tx_add_ptr(umm, &cont_df->cd_dtx_committed_head,
				sizeof(cont_df->cd_dtx_committed_head) +
				sizeof(cont_df->cd_dtx_committed_tail));
		cont_df->cd_dtx_committed_head = dsb_off;
	} else {
		umem_tx_add_ptr(umm, &dsb_prev->dsb_next,
				sizeof(dsb_prev->dsb_next));
		dsb_prev->dsb_next = dsb_off;

		umem_tx_add_ptr(umm, &cont_df->cd_dtx_committed_tail,
				sizeof(cont_df->cd_dtx_committed_tail));
	}

	cont_df->cd_dtx_committed_tail = dsb_off;

	for (i = 0, j = 0; i < count; i++, cur++) {
		struct vos_dtx_cmt_ent	*dce = NULL;

		rc = vos_dtx_commit_one(cont, &dtis[cur], epoch, &dce);
		if (rc1 == 0)
			rc1 = rc;

		if (dce != NULL) {
			memcpy(&dce_df[j], &dce->dce_base, sizeof(dce_df[j]));
			j++;
		}
	}

	if (dce_df != &dsb->dsb_commmitted_data[0]) {
		if (j > 0)
			memcpy(&dsb->dsb_commmitted_data[0], dce_df,
			       sizeof(*dce_df) * j);
		D_FREE(dce_df);
	}

	dsb->dsb_count = j;

	return rc != 0 ? rc : rc1;
}

int
vos_dtx_commit(daos_handle_t coh, struct dtx_id *dtis, int count)
{
	struct vos_container	*cont;
	int			 rc;

	cont = vos_hdl2cont(coh);
	D_ASSERT(cont != NULL);

	/* Commit multiple DTXs via single PMDK transaction. */
	rc = vos_tx_begin(vos_cont2umm(cont));
	if (rc == 0) {
		rc = vos_dtx_commit_internal(cont, dtis, count, 0);
		rc = vos_tx_end(vos_cont2umm(cont), rc);
	}

	return rc;
}

int
vos_dtx_abort(daos_handle_t coh, daos_epoch_t epoch, struct dtx_id *dtis,
	      int count)
{
	struct vos_container	*cont;
	int			 rc;
	int			 i;

	cont = vos_hdl2cont(coh);
	D_ASSERT(cont != NULL);

	/* Abort multiple DTXs via single PMDK transaction. */
	rc = vos_tx_begin(vos_cont2umm(cont));
	if (rc == 0) {
		for (i = 0; i < count; i++)
			vos_dtx_abort_one(cont, epoch, &dtis[i]);

		rc = vos_tx_end(vos_cont2umm(cont), 0);
	}

	return rc;
}

int
vos_dtx_aggregate(daos_handle_t coh)
{
	struct vos_container		*cont;
	struct vos_cont_df		*cont_df;
	struct umem_instance		*umm;
	struct vos_dtx_scm_blob		*dsb;
	struct vos_dtx_scm_blob		*tmp;
	umem_off_t			 dsb_off;
	int				 rc;
	int				 i;

	cont = vos_hdl2cont(coh);
	D_ASSERT(cont != NULL);

	umm = vos_cont2umm(cont);
	cont_df = cont->vc_cont_df;

	dsb_off = cont_df->cd_dtx_committed_head;
	dsb = umem_off2ptr(umm, dsb_off);
	if (dsb == NULL || dsb->dsb_count == 0)
		return 0;

	rc = vos_tx_begin(umm);
	if (rc != 0)
		return rc;

	for (i = 0; i < dsb->dsb_count &&
	     !d_list_empty(&cont->vc_dtx_committed_list); i++) {
		struct vos_dtx_cmt_ent	*dce;
		d_iov_t			 kiov;

		dce = d_list_entry(cont->vc_dtx_committed_list.next,
				   struct vos_dtx_cmt_ent, dce_committed_link);
		d_iov_set(&kiov, &DCE_XID(dce), sizeof(DCE_XID(dce)));
		dbtree_delete(cont->vc_dtx_committed_hdl, BTR_PROBE_EQ,
			      &kiov, NULL);
	}

	tmp = umem_off2ptr(umm, dsb->dsb_next);
	if (tmp == NULL) {
		/* The last blob for committed DTX blob. */
		D_ASSERT(cont_df->cd_dtx_committed_tail ==
			 cont_df->cd_dtx_committed_head);

		umem_tx_add_ptr(umm, &cont_df->cd_dtx_committed_tail,
				sizeof(cont_df->cd_dtx_committed_tail));
		cont_df->cd_dtx_committed_tail = UMOFF_NULL;
	} else {
		umem_tx_add_ptr(umm, &tmp->dsb_prev, sizeof(tmp->dsb_prev));
		tmp->dsb_prev = UMOFF_NULL;
	}

	umem_tx_add_ptr(umm, &cont_df->cd_dtx_committed_head,
			sizeof(cont_df->cd_dtx_committed_head));
	cont_df->cd_dtx_committed_head = dsb->dsb_next;

	umem_free(umm, dsb_off);

	return vos_tx_end(umm, 0);
}

void
vos_dtx_stat(daos_handle_t coh, struct dtx_stat *stat)
{
	struct vos_container	*cont;

	cont = vos_hdl2cont(coh);
	D_ASSERT(cont != NULL);

	stat->dtx_priority_count = cont->vc_dtx_priority_count;
	stat->dtx_committable_count = cont->vc_dtx_committable_count;
	stat->dtx_oldest_committable_time = vos_dtx_cc_oldest(cont);
	stat->dtx_committed_count = cont->vc_dtx_committed_count;
	if (d_list_empty(&cont->vc_dtx_committed_list)) {
		stat->dtx_oldest_committed_time = 0;
	} else {
		struct vos_dtx_cmt_ent	*dce;

		dce = d_list_entry(cont->vc_dtx_committed_list.next,
				   struct vos_dtx_cmt_ent, dce_committed_link);
		stat->dtx_oldest_committed_time = DCE_EPOCH(dce);
	}
}

int
vos_dtx_mark_sync(daos_handle_t coh, daos_unit_oid_t oid, daos_epoch_t epoch)
{
	struct vos_container	*cont;
	struct daos_lru_cache	*occ;
	struct vos_object	*obj;
	daos_epoch_range_t	 epr = {0, epoch};
	int	rc;

	cont = vos_hdl2cont(coh);
	occ = vos_obj_cache_current();
	rc = vos_obj_hold(occ, cont, oid, &epr, true,
			  DAOS_INTENT_DEFAULT, true, &obj);
	if (rc != 0) {
		D_ERROR(DF_UOID" fail to mark sync(1): rc = %d\n",
			DP_UOID(oid), rc);
		return rc;
	}

	if (obj->obj_df != NULL && obj->obj_df->vo_sync < epoch) {
		struct umem_instance	*umm = vos_cont2umm(cont);

		rc = vos_tx_begin(umm);
		if (rc == 0) {
			umem_tx_add_ptr(umm, &obj->obj_df->vo_sync,
					sizeof(obj->obj_df->vo_sync));
			obj->obj_df->vo_sync = epoch;
			rc = vos_tx_end(umm, 0);
		}

		if (rc == 0) {
			D_INFO("Update sync epoch "DF_U64" => "DF_U64
			       " for the obj "DF_UOID"\n",
			       obj->obj_sync_epoch, epoch, DP_UOID(oid));
			obj->obj_sync_epoch = epoch;
		} else {
			D_ERROR(DF_UOID" fail to mark sync(2): rc = %d\n",
				DP_UOID(oid), rc);
		}
	}

	vos_obj_release(occ, obj, false);
	return rc;
}

/*
 * \return	+2 if current DTX blob is batched cleanup.
 *		+1 if current DTX entry has ever been committed.
 *		Zero if the DTX entry is not committed.
 *		Negative value if error.
 */
static int
dtx_act_handle_batched_cleanup(struct vos_container *cont,
			       struct dtx_batched_cleanup_blob *bcb,
			       struct vos_dtx_act_ent_df *dae_df,
			       umem_off_t *next_dsb)
{
	struct umem_instance		*umm;
	struct vos_dtx_scm_blob		*dsb;
	struct vos_dtx_record_df	*rec_df;
	int				 rc;

	if (daos_is_zero_dti(&dae_df->dae_xid) ||
	    dae_df->dae_flags & DTX_EF_INVALID)
		return 1;

	if (dae_df->dae_rec_cnt == 0)
		return 0;

	umm = vos_cont2umm(cont);
	rec_df = &dae_df->dae_rec_inline[0];
	switch (rec_df->dr_type) {
	case DTX_RT_ILOG:
		/* XXX: For pure ilog touched DTX case, we do NOT handle it
		 *	via batched cleanup.
		 */
		return 0;
	case DTX_RT_SVT: {
		struct vos_irec_df	*svt;

		svt = umem_off2ptr(umm, rec_df->dr_record);
		if (!dtx_is_null(svt->ir_dtx))
			return 0;
		break;
	}
	case DTX_RT_EVT: {
		struct evt_desc		*evt;

		evt = umem_off2ptr(umm, rec_df->dr_record);
		if (!dtx_is_null(evt->dc_dtx))
			return 0;
		break;
	}
	default:
		D_ERROR(DF_UOID" invalid DTX record type %d for DTX "DF_DTI"\n",
			DP_UOID(dae_df->dae_oid), rec_df->dr_type,
			DP_DTI(&dae_df->dae_xid));
		return -DER_IO;
	}

	rc = vos_tx_begin(umm);
	if (rc != 0)
		return rc;

	if (!dtx_is_null(dae_df->dae_rec_off))
		umem_free(umm, dae_df->dae_rec_off);

	bcb->bcb_dae_count--;
	dsb = umem_off2ptr(umm, bcb->bcb_dsb_off);
	if (bcb->bcb_dae_count > 0 || dsb->dsb_index < dsb->dsb_cap) {
		umem_tx_add_ptr(umm, &dae_df->dae_flags,
				sizeof(dae_df->dae_flags));
		dae_df->dae_flags = DTX_EF_INVALID;
		umem_tx_add_ptr(umm, &dsb->dsb_count, sizeof(dsb->dsb_count));
		dsb->dsb_count--;
		rc = 1;
	} else {
		dtx_batched_cleanup(cont, bcb, next_dsb);
		rc = 2;
	}

	return vos_tx_end(umm, 0);
}

int
vos_dtx_act_reindex(struct vos_container *cont)
{
	struct umem_instance		*umm = vos_cont2umm(cont);
	struct vos_cont_df		*cont_df = cont->vc_cont_df;
	struct vos_dtx_scm_blob		*dsb;
	umem_off_t			 dsb_off = cont_df->cd_dtx_active_head;
	d_iov_t				 kiov;
	d_iov_t				 riov;
	int				 rc = 0;
	int				 i;

	while (1) {
		struct dtx_batched_cleanup_blob	*bcb;

next:
		dsb = umem_off2ptr(umm, dsb_off);
		if (dsb == NULL)
			break;

		D_ASSERT(dsb->dsb_magic == DTX_ACT_BLOB_MAGIC);

		bcb = dtx_bcb_alloc(cont, dsb);
		if (bcb == NULL)
			D_GOTO(out, rc = -DER_NOMEM);

		bcb->bcb_dae_count = dsb->dsb_count;
		if (bcb->bcb_dae_count == 0) {
			if (cont_df->cd_dtx_active_tail != dsb_off) {
				D_ERROR("Invalid active DTX blob\n");
				rc = -DER_IO;
			}

			D_GOTO(out, rc);
		}

		for (i = 0; i < dsb->dsb_index; i++) {
			struct vos_dtx_act_ent_df	*dae_df;
			struct vos_dtx_act_ent		*dae;
			int				 count;

			dae_df = &dsb->dsb_active_data[i];
			rc = dtx_act_handle_batched_cleanup(cont, bcb, dae_df,
							    &dsb_off);
			if (rc > 1)
				goto next;

			if (rc == 1)
				continue;

			if (rc < 0)
				D_GOTO(out, rc);

			D_ALLOC_PTR(dae);
			if (dae == NULL)
				D_GOTO(out, rc = -DER_NOMEM);

			memcpy(&dae->dae_base, dae_df, sizeof(dae->dae_base));
			dae->dae_df_off = umem_ptr2off(umm, dae_df);
			dae->dae_dsb = dsb;
			dae->dae_bcb = bcb;

			if (DAE_REC_CNT(dae) <= DTX_INLINE_REC_CNT)
				goto insert;

			count = DAE_REC_CNT(dae) - DTX_INLINE_REC_CNT;
			D_ALLOC(dae->dae_records,
				sizeof(*dae->dae_records) * count);
			if (dae->dae_records == NULL) {
				D_FREE_PTR(dae);
				D_GOTO(out, rc = -DER_NOMEM);
			}

			memcpy(dae->dae_records,
			       umem_off2ptr(umm, dae_df->dae_rec_off),
			       sizeof(*dae->dae_records) * count);
			dae->dae_rec_cap = count;

insert:
			d_iov_set(&kiov, &DAE_XID(dae), sizeof(DAE_XID(dae)));
			d_iov_set(&riov, dae, sizeof(*dae));
			rc = dbtree_upsert(cont->vc_dtx_active_hdl,
					   BTR_PROBE_EQ, DAOS_INTENT_UPDATE,
					   &kiov, &riov);
			if (rc != 0) {
				D_FREE(dae->dae_records);
				D_FREE_PTR(dae);
				goto out;
			}
		}

		dsb_off = dsb->dsb_next;
	}

out:
	return rc > 0 ? 0 : rc;
}

int
vos_dtx_cmt_reindex(daos_handle_t coh, void *hint)
{
	struct umem_instance		*umm;
	struct vos_container		*cont;
	struct vos_cont_df		*cont_df;
	struct vos_dtx_cmt_ent		*dce;
	struct vos_dtx_scm_blob		*dsb;
	umem_off_t			*dsb_off = hint;
	d_iov_t				 kiov;
	d_iov_t				 riov;
	int				 rc = 0;
	int				 i;

	cont = vos_hdl2cont(coh);
	D_ASSERT(cont != NULL);

	umm = vos_cont2umm(cont);
	cont_df = cont->vc_cont_df;

	if (dtx_is_null(*dsb_off))
		dsb = umem_off2ptr(umm, cont_df->cd_dtx_committed_head);
	else
		dsb = umem_off2ptr(umm, *dsb_off);

	if (dsb == NULL)
		D_GOTO(out, rc = 1);

	D_ASSERT(dsb->dsb_magic == DTX_CMT_BLOB_MAGIC);

	cont->vc_reindex_cmt_dtx = 1;

	for (i = 0; i < dsb->dsb_count; i++) {
		if (dsb->dsb_commmitted_data[i].dce_epoch == 0)
			continue;

		D_ALLOC_PTR(dce);
		if (dce == NULL)
			D_GOTO(out, rc = -DER_NOMEM);

		D_INIT_LIST_HEAD(&dce->dce_bcb_link);
		memcpy(&dce->dce_base, &dsb->dsb_commmitted_data[i],
		       sizeof(dce->dce_base));
		dce->dce_reindex = 1;

		d_iov_set(&kiov, &DCE_XID(dce), sizeof(DCE_XID(dce)));
		d_iov_set(&riov, dce, sizeof(*dce));
		rc = dbtree_upsert(cont->vc_dtx_committed_hdl, BTR_PROBE_EQ,
				   DAOS_INTENT_UPDATE, &kiov, &riov);
		if (rc != 0)
			goto out;

		/* The committed DTX entry is already in the index.
		 * Related re-index logic can stop.
		 */
		if (dce->dce_exist)
			D_GOTO(out, rc = 1);
	}

	if (dsb->dsb_count < dsb->dsb_cap || dtx_is_null(dsb->dsb_next))
		D_GOTO(out, rc = 1);

	*dsb_off = dsb->dsb_next;

out:
	if (rc > 0) {
		d_list_splice_init(&cont->vc_dtx_committed_tmp_list,
				   &cont->vc_dtx_committed_list);
		cont->vc_dtx_committed_count +=
				cont->vc_dtx_committed_tmp_count;
		cont->vc_dtx_committed_tmp_count = 0;
		cont->vc_reindex_cmt_dtx = 0;
	}

	return rc;
}
