/* Generated by the protocol buffer compiler.  DO NOT EDIT! */
/* Generated from: srv.proto */

/* Do not generate deprecated warnings for self */
#ifndef PROTOBUF_C__NO_DEPRECATED
#define PROTOBUF_C__NO_DEPRECATED
#endif

#include "srv.pb-c.h"
void   mgmt__daos_resp__init
                     (Mgmt__DaosResp         *message)
{
  static const Mgmt__DaosResp init_value = MGMT__DAOS_RESP__INIT;
  *message = init_value;
}
size_t mgmt__daos_resp__get_packed_size
                     (const Mgmt__DaosResp *message)
{
  assert(message->base.descriptor == &mgmt__daos_resp__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__daos_resp__pack
                     (const Mgmt__DaosResp *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__daos_resp__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__daos_resp__pack_to_buffer
                     (const Mgmt__DaosResp *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__daos_resp__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__DaosResp *
       mgmt__daos_resp__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__DaosResp *)
     protobuf_c_message_unpack (&mgmt__daos_resp__descriptor,
                                allocator, len, data);
}
void   mgmt__daos_resp__free_unpacked
                     (Mgmt__DaosResp *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__daos_resp__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__join_req__init
                     (Mgmt__JoinReq         *message)
{
  static const Mgmt__JoinReq init_value = MGMT__JOIN_REQ__INIT;
  *message = init_value;
}
size_t mgmt__join_req__get_packed_size
                     (const Mgmt__JoinReq *message)
{
  assert(message->base.descriptor == &mgmt__join_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__join_req__pack
                     (const Mgmt__JoinReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__join_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__join_req__pack_to_buffer
                     (const Mgmt__JoinReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__join_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__JoinReq *
       mgmt__join_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__JoinReq *)
     protobuf_c_message_unpack (&mgmt__join_req__descriptor,
                                allocator, len, data);
}
void   mgmt__join_req__free_unpacked
                     (Mgmt__JoinReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__join_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__join_resp__init
                     (Mgmt__JoinResp         *message)
{
  static const Mgmt__JoinResp init_value = MGMT__JOIN_RESP__INIT;
  *message = init_value;
}
size_t mgmt__join_resp__get_packed_size
                     (const Mgmt__JoinResp *message)
{
  assert(message->base.descriptor == &mgmt__join_resp__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__join_resp__pack
                     (const Mgmt__JoinResp *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__join_resp__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__join_resp__pack_to_buffer
                     (const Mgmt__JoinResp *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__join_resp__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__JoinResp *
       mgmt__join_resp__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__JoinResp *)
     protobuf_c_message_unpack (&mgmt__join_resp__descriptor,
                                allocator, len, data);
}
void   mgmt__join_resp__free_unpacked
                     (Mgmt__JoinResp *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__join_resp__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__leader_query_req__init
                     (Mgmt__LeaderQueryReq         *message)
{
  static const Mgmt__LeaderQueryReq init_value = MGMT__LEADER_QUERY_REQ__INIT;
  *message = init_value;
}
size_t mgmt__leader_query_req__get_packed_size
                     (const Mgmt__LeaderQueryReq *message)
{
  assert(message->base.descriptor == &mgmt__leader_query_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__leader_query_req__pack
                     (const Mgmt__LeaderQueryReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__leader_query_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__leader_query_req__pack_to_buffer
                     (const Mgmt__LeaderQueryReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__leader_query_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__LeaderQueryReq *
       mgmt__leader_query_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__LeaderQueryReq *)
     protobuf_c_message_unpack (&mgmt__leader_query_req__descriptor,
                                allocator, len, data);
}
void   mgmt__leader_query_req__free_unpacked
                     (Mgmt__LeaderQueryReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__leader_query_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__leader_query_resp__init
                     (Mgmt__LeaderQueryResp         *message)
{
  static const Mgmt__LeaderQueryResp init_value = MGMT__LEADER_QUERY_RESP__INIT;
  *message = init_value;
}
size_t mgmt__leader_query_resp__get_packed_size
                     (const Mgmt__LeaderQueryResp *message)
{
  assert(message->base.descriptor == &mgmt__leader_query_resp__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__leader_query_resp__pack
                     (const Mgmt__LeaderQueryResp *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__leader_query_resp__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__leader_query_resp__pack_to_buffer
                     (const Mgmt__LeaderQueryResp *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__leader_query_resp__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__LeaderQueryResp *
       mgmt__leader_query_resp__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__LeaderQueryResp *)
     protobuf_c_message_unpack (&mgmt__leader_query_resp__descriptor,
                                allocator, len, data);
}
void   mgmt__leader_query_resp__free_unpacked
                     (Mgmt__LeaderQueryResp *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__leader_query_resp__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__get_attach_info_req__init
                     (Mgmt__GetAttachInfoReq         *message)
{
  static const Mgmt__GetAttachInfoReq init_value = MGMT__GET_ATTACH_INFO_REQ__INIT;
  *message = init_value;
}
size_t mgmt__get_attach_info_req__get_packed_size
                     (const Mgmt__GetAttachInfoReq *message)
{
  assert(message->base.descriptor == &mgmt__get_attach_info_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__get_attach_info_req__pack
                     (const Mgmt__GetAttachInfoReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__get_attach_info_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__get_attach_info_req__pack_to_buffer
                     (const Mgmt__GetAttachInfoReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__get_attach_info_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__GetAttachInfoReq *
       mgmt__get_attach_info_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__GetAttachInfoReq *)
     protobuf_c_message_unpack (&mgmt__get_attach_info_req__descriptor,
                                allocator, len, data);
}
void   mgmt__get_attach_info_req__free_unpacked
                     (Mgmt__GetAttachInfoReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__get_attach_info_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__get_attach_info_resp__psr__init
                     (Mgmt__GetAttachInfoResp__Psr         *message)
{
  static const Mgmt__GetAttachInfoResp__Psr init_value = MGMT__GET_ATTACH_INFO_RESP__PSR__INIT;
  *message = init_value;
}
void   mgmt__get_attach_info_resp__init
                     (Mgmt__GetAttachInfoResp         *message)
{
  static const Mgmt__GetAttachInfoResp init_value = MGMT__GET_ATTACH_INFO_RESP__INIT;
  *message = init_value;
}
size_t mgmt__get_attach_info_resp__get_packed_size
                     (const Mgmt__GetAttachInfoResp *message)
{
  assert(message->base.descriptor == &mgmt__get_attach_info_resp__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__get_attach_info_resp__pack
                     (const Mgmt__GetAttachInfoResp *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__get_attach_info_resp__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__get_attach_info_resp__pack_to_buffer
                     (const Mgmt__GetAttachInfoResp *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__get_attach_info_resp__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__GetAttachInfoResp *
       mgmt__get_attach_info_resp__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__GetAttachInfoResp *)
     protobuf_c_message_unpack (&mgmt__get_attach_info_resp__descriptor,
                                allocator, len, data);
}
void   mgmt__get_attach_info_resp__free_unpacked
                     (Mgmt__GetAttachInfoResp *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__get_attach_info_resp__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__prep_shutdown_req__init
                     (Mgmt__PrepShutdownReq         *message)
{
  static const Mgmt__PrepShutdownReq init_value = MGMT__PREP_SHUTDOWN_REQ__INIT;
  *message = init_value;
}
size_t mgmt__prep_shutdown_req__get_packed_size
                     (const Mgmt__PrepShutdownReq *message)
{
  assert(message->base.descriptor == &mgmt__prep_shutdown_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__prep_shutdown_req__pack
                     (const Mgmt__PrepShutdownReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__prep_shutdown_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__prep_shutdown_req__pack_to_buffer
                     (const Mgmt__PrepShutdownReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__prep_shutdown_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__PrepShutdownReq *
       mgmt__prep_shutdown_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__PrepShutdownReq *)
     protobuf_c_message_unpack (&mgmt__prep_shutdown_req__descriptor,
                                allocator, len, data);
}
void   mgmt__prep_shutdown_req__free_unpacked
                     (Mgmt__PrepShutdownReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__prep_shutdown_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__ping_rank_req__init
                     (Mgmt__PingRankReq         *message)
{
  static const Mgmt__PingRankReq init_value = MGMT__PING_RANK_REQ__INIT;
  *message = init_value;
}
size_t mgmt__ping_rank_req__get_packed_size
                     (const Mgmt__PingRankReq *message)
{
  assert(message->base.descriptor == &mgmt__ping_rank_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__ping_rank_req__pack
                     (const Mgmt__PingRankReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__ping_rank_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__ping_rank_req__pack_to_buffer
                     (const Mgmt__PingRankReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__ping_rank_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__PingRankReq *
       mgmt__ping_rank_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__PingRankReq *)
     protobuf_c_message_unpack (&mgmt__ping_rank_req__descriptor,
                                allocator, len, data);
}
void   mgmt__ping_rank_req__free_unpacked
                     (Mgmt__PingRankReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__ping_rank_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__set_rank_req__init
                     (Mgmt__SetRankReq         *message)
{
  static const Mgmt__SetRankReq init_value = MGMT__SET_RANK_REQ__INIT;
  *message = init_value;
}
size_t mgmt__set_rank_req__get_packed_size
                     (const Mgmt__SetRankReq *message)
{
  assert(message->base.descriptor == &mgmt__set_rank_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__set_rank_req__pack
                     (const Mgmt__SetRankReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__set_rank_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__set_rank_req__pack_to_buffer
                     (const Mgmt__SetRankReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__set_rank_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__SetRankReq *
       mgmt__set_rank_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__SetRankReq *)
     protobuf_c_message_unpack (&mgmt__set_rank_req__descriptor,
                                allocator, len, data);
}
void   mgmt__set_rank_req__free_unpacked
                     (Mgmt__SetRankReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__set_rank_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__create_ms_req__init
                     (Mgmt__CreateMsReq         *message)
{
  static const Mgmt__CreateMsReq init_value = MGMT__CREATE_MS_REQ__INIT;
  *message = init_value;
}
size_t mgmt__create_ms_req__get_packed_size
                     (const Mgmt__CreateMsReq *message)
{
  assert(message->base.descriptor == &mgmt__create_ms_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__create_ms_req__pack
                     (const Mgmt__CreateMsReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__create_ms_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__create_ms_req__pack_to_buffer
                     (const Mgmt__CreateMsReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__create_ms_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__CreateMsReq *
       mgmt__create_ms_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__CreateMsReq *)
     protobuf_c_message_unpack (&mgmt__create_ms_req__descriptor,
                                allocator, len, data);
}
void   mgmt__create_ms_req__free_unpacked
                     (Mgmt__CreateMsReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__create_ms_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
static const ProtobufCFieldDescriptor mgmt__daos_resp__field_descriptors[1] =
{
  {
    "status",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_INT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__DaosResp, status),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__daos_resp__field_indices_by_name[] = {
  0,   /* field[0] = status */
};
static const ProtobufCIntRange mgmt__daos_resp__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 1 }
};
const ProtobufCMessageDescriptor mgmt__daos_resp__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.DaosResp",
  "DaosResp",
  "Mgmt__DaosResp",
  "mgmt",
  sizeof(Mgmt__DaosResp),
  1,
  mgmt__daos_resp__field_descriptors,
  mgmt__daos_resp__field_indices_by_name,
  1,  mgmt__daos_resp__number_ranges,
  (ProtobufCMessageInit) mgmt__daos_resp__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__join_req__field_descriptors[5] =
{
  {
    "uuid",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__JoinReq, uuid),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "rank",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__JoinReq, rank),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "uri",
    3,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__JoinReq, uri),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "nctxs",
    4,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__JoinReq, nctxs),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "addr",
    5,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__JoinReq, addr),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__join_req__field_indices_by_name[] = {
  4,   /* field[4] = addr */
  3,   /* field[3] = nctxs */
  1,   /* field[1] = rank */
  2,   /* field[2] = uri */
  0,   /* field[0] = uuid */
};
static const ProtobufCIntRange mgmt__join_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 5 }
};
const ProtobufCMessageDescriptor mgmt__join_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.JoinReq",
  "JoinReq",
  "Mgmt__JoinReq",
  "mgmt",
  sizeof(Mgmt__JoinReq),
  5,
  mgmt__join_req__field_descriptors,
  mgmt__join_req__field_indices_by_name,
  1,  mgmt__join_req__number_ranges,
  (ProtobufCMessageInit) mgmt__join_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCEnumValue mgmt__join_resp__state__enum_values_by_number[2] =
{
  { "IN", "MGMT__JOIN_RESP__STATE__IN", 0 },
  { "OUT", "MGMT__JOIN_RESP__STATE__OUT", 1 },
};
static const ProtobufCIntRange mgmt__join_resp__state__value_ranges[] = {
{0, 0},{0, 2}
};
static const ProtobufCEnumValueIndex mgmt__join_resp__state__enum_values_by_name[2] =
{
  { "IN", 0 },
  { "OUT", 1 },
};
const ProtobufCEnumDescriptor mgmt__join_resp__state__descriptor =
{
  PROTOBUF_C__ENUM_DESCRIPTOR_MAGIC,
  "mgmt.JoinResp.State",
  "State",
  "Mgmt__JoinResp__State",
  "mgmt",
  2,
  mgmt__join_resp__state__enum_values_by_number,
  2,
  mgmt__join_resp__state__enum_values_by_name,
  1,
  mgmt__join_resp__state__value_ranges,
  NULL,NULL,NULL,NULL   /* reserved[1234] */
};
static const ProtobufCFieldDescriptor mgmt__join_resp__field_descriptors[3] =
{
  {
    "status",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_INT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__JoinResp, status),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "rank",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__JoinResp, rank),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "state",
    3,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_ENUM,
    0,   /* quantifier_offset */
    offsetof(Mgmt__JoinResp, state),
    &mgmt__join_resp__state__descriptor,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__join_resp__field_indices_by_name[] = {
  1,   /* field[1] = rank */
  2,   /* field[2] = state */
  0,   /* field[0] = status */
};
static const ProtobufCIntRange mgmt__join_resp__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 3 }
};
const ProtobufCMessageDescriptor mgmt__join_resp__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.JoinResp",
  "JoinResp",
  "Mgmt__JoinResp",
  "mgmt",
  sizeof(Mgmt__JoinResp),
  3,
  mgmt__join_resp__field_descriptors,
  mgmt__join_resp__field_indices_by_name,
  1,  mgmt__join_resp__number_ranges,
  (ProtobufCMessageInit) mgmt__join_resp__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__leader_query_req__field_descriptors[1] =
{
  {
    "system",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__LeaderQueryReq, system),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__leader_query_req__field_indices_by_name[] = {
  0,   /* field[0] = system */
};
static const ProtobufCIntRange mgmt__leader_query_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 1 }
};
const ProtobufCMessageDescriptor mgmt__leader_query_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.LeaderQueryReq",
  "LeaderQueryReq",
  "Mgmt__LeaderQueryReq",
  "mgmt",
  sizeof(Mgmt__LeaderQueryReq),
  1,
  mgmt__leader_query_req__field_descriptors,
  mgmt__leader_query_req__field_indices_by_name,
  1,  mgmt__leader_query_req__number_ranges,
  (ProtobufCMessageInit) mgmt__leader_query_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__leader_query_resp__field_descriptors[2] =
{
  {
    "currentLeader",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__LeaderQueryResp, currentleader),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "replicas",
    2,
    PROTOBUF_C_LABEL_REPEATED,
    PROTOBUF_C_TYPE_STRING,
    offsetof(Mgmt__LeaderQueryResp, n_replicas),
    offsetof(Mgmt__LeaderQueryResp, replicas),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__leader_query_resp__field_indices_by_name[] = {
  0,   /* field[0] = currentLeader */
  1,   /* field[1] = replicas */
};
static const ProtobufCIntRange mgmt__leader_query_resp__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 2 }
};
const ProtobufCMessageDescriptor mgmt__leader_query_resp__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.LeaderQueryResp",
  "LeaderQueryResp",
  "Mgmt__LeaderQueryResp",
  "mgmt",
  sizeof(Mgmt__LeaderQueryResp),
  2,
  mgmt__leader_query_resp__field_descriptors,
  mgmt__leader_query_resp__field_indices_by_name,
  1,  mgmt__leader_query_resp__number_ranges,
  (ProtobufCMessageInit) mgmt__leader_query_resp__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__get_attach_info_req__field_descriptors[2] =
{
  {
    "sys",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoReq, sys),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "numa",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoReq, numa),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__get_attach_info_req__field_indices_by_name[] = {
  1,   /* field[1] = numa */
  0,   /* field[0] = sys */
};
static const ProtobufCIntRange mgmt__get_attach_info_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 2 }
};
const ProtobufCMessageDescriptor mgmt__get_attach_info_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.GetAttachInfoReq",
  "GetAttachInfoReq",
  "Mgmt__GetAttachInfoReq",
  "mgmt",
  sizeof(Mgmt__GetAttachInfoReq),
  2,
  mgmt__get_attach_info_req__field_descriptors,
  mgmt__get_attach_info_req__field_indices_by_name,
  1,  mgmt__get_attach_info_req__number_ranges,
  (ProtobufCMessageInit) mgmt__get_attach_info_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__get_attach_info_resp__psr__field_descriptors[2] =
{
  {
    "rank",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoResp__Psr, rank),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "uri",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoResp__Psr, uri),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__get_attach_info_resp__psr__field_indices_by_name[] = {
  0,   /* field[0] = rank */
  1,   /* field[1] = uri */
};
static const ProtobufCIntRange mgmt__get_attach_info_resp__psr__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 2 }
};
const ProtobufCMessageDescriptor mgmt__get_attach_info_resp__psr__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.GetAttachInfoResp.Psr",
  "Psr",
  "Mgmt__GetAttachInfoResp__Psr",
  "mgmt",
  sizeof(Mgmt__GetAttachInfoResp__Psr),
  2,
  mgmt__get_attach_info_resp__psr__field_descriptors,
  mgmt__get_attach_info_resp__psr__field_indices_by_name,
  1,  mgmt__get_attach_info_resp__psr__number_ranges,
  (ProtobufCMessageInit) mgmt__get_attach_info_resp__psr__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__get_attach_info_resp__field_descriptors[7] =
{
  {
    "status",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_INT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoResp, status),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "psrs",
    2,
    PROTOBUF_C_LABEL_REPEATED,
    PROTOBUF_C_TYPE_MESSAGE,
    offsetof(Mgmt__GetAttachInfoResp, n_psrs),
    offsetof(Mgmt__GetAttachInfoResp, psrs),
    &mgmt__get_attach_info_resp__psr__descriptor,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "Provider",
    3,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoResp, provider),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "Interface",
    4,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoResp, interface),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "Domain",
    5,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoResp, domain),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "CrtCtxShareAddr",
    6,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoResp, crtctxshareaddr),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "CrtTimeout",
    7,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__GetAttachInfoResp, crttimeout),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__get_attach_info_resp__field_indices_by_name[] = {
  5,   /* field[5] = CrtCtxShareAddr */
  6,   /* field[6] = CrtTimeout */
  4,   /* field[4] = Domain */
  3,   /* field[3] = Interface */
  2,   /* field[2] = Provider */
  1,   /* field[1] = psrs */
  0,   /* field[0] = status */
};
static const ProtobufCIntRange mgmt__get_attach_info_resp__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 7 }
};
const ProtobufCMessageDescriptor mgmt__get_attach_info_resp__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.GetAttachInfoResp",
  "GetAttachInfoResp",
  "Mgmt__GetAttachInfoResp",
  "mgmt",
  sizeof(Mgmt__GetAttachInfoResp),
  7,
  mgmt__get_attach_info_resp__field_descriptors,
  mgmt__get_attach_info_resp__field_indices_by_name,
  1,  mgmt__get_attach_info_resp__number_ranges,
  (ProtobufCMessageInit) mgmt__get_attach_info_resp__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__prep_shutdown_req__field_descriptors[1] =
{
  {
    "rank",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__PrepShutdownReq, rank),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__prep_shutdown_req__field_indices_by_name[] = {
  0,   /* field[0] = rank */
};
static const ProtobufCIntRange mgmt__prep_shutdown_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 1 }
};
const ProtobufCMessageDescriptor mgmt__prep_shutdown_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.PrepShutdownReq",
  "PrepShutdownReq",
  "Mgmt__PrepShutdownReq",
  "mgmt",
  sizeof(Mgmt__PrepShutdownReq),
  1,
  mgmt__prep_shutdown_req__field_descriptors,
  mgmt__prep_shutdown_req__field_indices_by_name,
  1,  mgmt__prep_shutdown_req__number_ranges,
  (ProtobufCMessageInit) mgmt__prep_shutdown_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__ping_rank_req__field_descriptors[1] =
{
  {
    "rank",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__PingRankReq, rank),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__ping_rank_req__field_indices_by_name[] = {
  0,   /* field[0] = rank */
};
static const ProtobufCIntRange mgmt__ping_rank_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 1 }
};
const ProtobufCMessageDescriptor mgmt__ping_rank_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.PingRankReq",
  "PingRankReq",
  "Mgmt__PingRankReq",
  "mgmt",
  sizeof(Mgmt__PingRankReq),
  1,
  mgmt__ping_rank_req__field_descriptors,
  mgmt__ping_rank_req__field_indices_by_name,
  1,  mgmt__ping_rank_req__number_ranges,
  (ProtobufCMessageInit) mgmt__ping_rank_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__set_rank_req__field_descriptors[1] =
{
  {
    "rank",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__SetRankReq, rank),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__set_rank_req__field_indices_by_name[] = {
  0,   /* field[0] = rank */
};
static const ProtobufCIntRange mgmt__set_rank_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 1 }
};
const ProtobufCMessageDescriptor mgmt__set_rank_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.SetRankReq",
  "SetRankReq",
  "Mgmt__SetRankReq",
  "mgmt",
  sizeof(Mgmt__SetRankReq),
  1,
  mgmt__set_rank_req__field_descriptors,
  mgmt__set_rank_req__field_indices_by_name,
  1,  mgmt__set_rank_req__number_ranges,
  (ProtobufCMessageInit) mgmt__set_rank_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__create_ms_req__field_descriptors[3] =
{
  {
    "bootstrap",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_BOOL,
    0,   /* quantifier_offset */
    offsetof(Mgmt__CreateMsReq, bootstrap),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "uuid",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__CreateMsReq, uuid),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "addr",
    3,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__CreateMsReq, addr),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__create_ms_req__field_indices_by_name[] = {
  2,   /* field[2] = addr */
  0,   /* field[0] = bootstrap */
  1,   /* field[1] = uuid */
};
static const ProtobufCIntRange mgmt__create_ms_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 3 }
};
const ProtobufCMessageDescriptor mgmt__create_ms_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.CreateMsReq",
  "CreateMsReq",
  "Mgmt__CreateMsReq",
  "mgmt",
  sizeof(Mgmt__CreateMsReq),
  3,
  mgmt__create_ms_req__field_descriptors,
  mgmt__create_ms_req__field_indices_by_name,
  1,  mgmt__create_ms_req__number_ranges,
  (ProtobufCMessageInit) mgmt__create_ms_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
