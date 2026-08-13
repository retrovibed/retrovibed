// This is a generated file - do not edit.
//
// Generated from meta/meta.account.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use accountDescriptor instead')
const Account$json = {
  '1': 'Account',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'description', '3': 2, '4': 1, '5': 9, '10': 'description'},
    {'1': 'created_at', '3': 3, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 4, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'disabled_at', '3': 5, '4': 1, '5': 9, '10': 'disabled_at'},
    {'1': 'email', '3': 6, '4': 1, '5': 9, '10': 'email'},
    {'1': 'phone', '3': 7, '4': 1, '5': 9, '10': 'phone'},
  ],
};

/// Descriptor for `Account`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountDescriptor = $convert.base64Decode(
    'CgdBY2NvdW50Eg4KAmlkGAEgASgJUgJpZBIgCgtkZXNjcmlwdGlvbhgCIAEoCVILZGVzY3JpcH'
    'Rpb24SHgoKY3JlYXRlZF9hdBgDIAEoCVIKY3JlYXRlZF9hdBIeCgp1cGRhdGVkX2F0GAQgASgJ'
    'Ugp1cGRhdGVkX2F0EiAKC2Rpc2FibGVkX2F0GAUgASgJUgtkaXNhYmxlZF9hdBIUCgVlbWFpbB'
    'gGIAEoCVIFZW1haWwSFAoFcGhvbmUYByABKAlSBXBob25l');

@$core.Deprecated('Use accountLookupRequestDescriptor instead')
const AccountLookupRequest$json = {
  '1': 'AccountLookupRequest',
};

/// Descriptor for `AccountLookupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountLookupRequestDescriptor =
    $convert.base64Decode('ChRBY2NvdW50TG9va3VwUmVxdWVzdA==');

@$core.Deprecated('Use accountLookupResponseDescriptor instead')
const AccountLookupResponse$json = {
  '1': 'AccountLookupResponse',
  '2': [
    {
      '1': 'account',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Account',
      '10': 'account'
    },
  ],
};

/// Descriptor for `AccountLookupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountLookupResponseDescriptor = $convert.base64Decode(
    'ChVBY2NvdW50TG9va3VwUmVzcG9uc2USJwoHYWNjb3VudBgBIAEoCzINLm1ldGEuQWNjb3VudF'
    'IHYWNjb3VudA==');

@$core.Deprecated('Use accountUpdateRequestDescriptor instead')
const AccountUpdateRequest$json = {
  '1': 'AccountUpdateRequest',
  '2': [
    {
      '1': 'account',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Account',
      '10': 'account'
    },
  ],
};

/// Descriptor for `AccountUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountUpdateRequestDescriptor = $convert.base64Decode(
    'ChRBY2NvdW50VXBkYXRlUmVxdWVzdBInCgdhY2NvdW50GAEgASgLMg0ubWV0YS5BY2NvdW50Ug'
    'dhY2NvdW50');

@$core.Deprecated('Use accountUpdateResponseDescriptor instead')
const AccountUpdateResponse$json = {
  '1': 'AccountUpdateResponse',
  '2': [
    {
      '1': 'account',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.meta.Account',
      '10': 'account'
    },
  ],
};

/// Descriptor for `AccountUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountUpdateResponseDescriptor = $convert.base64Decode(
    'ChVBY2NvdW50VXBkYXRlUmVzcG9uc2USJwoHYWNjb3VudBgBIAEoCzINLm1ldGEuQWNjb3VudF'
    'IHYWNjb3VudA==');
