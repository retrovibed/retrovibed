// This is a generated file - do not edit.
//
// Generated from meta/meta.authz.proto.

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

@$core.Deprecated('Use bearerDescriptor instead')
const Bearer$json = {
  '1': 'Bearer',
  '2': [
    {'1': 'jti', '3': 1, '4': 1, '5': 9, '10': 'jti'},
    {'1': 'iss', '3': 2, '4': 1, '5': 9, '10': 'iss'},
    {'1': 'sub', '3': 3, '4': 1, '5': 9, '10': 'sub'},
    {'1': 'sid', '3': 4, '4': 1, '5': 9, '10': 'sid'},
    {'1': 'iat', '3': 5, '4': 1, '5': 3, '10': 'iat'},
    {'1': 'exp', '3': 6, '4': 1, '5': 3, '10': 'exp'},
    {'1': 'nbf', '3': 7, '4': 1, '5': 3, '10': 'nbf'},
  ],
};

/// Descriptor for `Bearer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bearerDescriptor = $convert.base64Decode(
    'CgZCZWFyZXISEAoDanRpGAEgASgJUgNqdGkSEAoDaXNzGAIgASgJUgNpc3MSEAoDc3ViGAMgAS'
    'gJUgNzdWISEAoDc2lkGAQgASgJUgNzaWQSEAoDaWF0GAUgASgDUgNpYXQSEAoDZXhwGAYgASgD'
    'UgNleHASEAoDbmJmGAcgASgDUgNuYmY=');

@$core.Deprecated('Use tokenDescriptor instead')
const Token$json = {
  '1': 'Token',
  '2': [
    {'1': 'jti', '3': 1, '4': 1, '5': 9, '10': 'jti'},
    {'1': 'iss', '3': 2, '4': 1, '5': 9, '10': 'iss'},
    {'1': 'sub', '3': 3, '4': 1, '5': 9, '10': 'sub'},
    {'1': 'sid', '3': 4, '4': 1, '5': 9, '10': 'sid'},
    {'1': 'iat', '3': 5, '4': 1, '5': 3, '10': 'iat'},
    {'1': 'exp', '3': 6, '4': 1, '5': 3, '10': 'exp'},
    {'1': 'nbf', '3': 7, '4': 1, '5': 3, '10': 'nbf'},
    {'1': 'usermanagement', '3': 1000, '4': 1, '5': 8, '10': 'usermanagement'},
    {'1': 'remote_control', '3': 1001, '4': 1, '5': 8, '10': 'remote_control'},
    {'1': 'billing_read', '3': 1002, '4': 1, '5': 8, '10': 'billing_read'},
    {'1': 'billing_modify', '3': 1003, '4': 1, '5': 8, '10': 'billing_modify'},
    {
      '1': 'community_modify',
      '3': 1004,
      '4': 1,
      '5': 8,
      '10': 'community_modify'
    },
    {'1': 'archive_upload', '3': 1005, '4': 1, '5': 4, '10': 'archive_upload'},
    {
      '1': 'archive_download',
      '3': 1006,
      '4': 1,
      '5': 4,
      '10': 'archive_download'
    },
    {'1': 'library_read', '3': 1007, '4': 1, '5': 8, '10': 'library_read'},
    {'1': 'library_modify', '3': 1008, '4': 1, '5': 8, '10': 'library_modify'},
    {'1': 'archive_sync', '3': 1009, '4': 1, '5': 8, '10': 'archive_sync'},
    {'1': 'community_sync', '3': 1010, '4': 1, '5': 8, '10': 'community_sync'},
    {'1': 'local_only', '3': 2000, '4': 1, '5': 8, '10': 'local_only'},
  ],
  '9': [
    {'1': 9, '2': 1000},
    {'1': 1012, '2': 2000},
  ],
};

/// Descriptor for `Token`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tokenDescriptor = $convert.base64Decode(
    'CgVUb2tlbhIQCgNqdGkYASABKAlSA2p0aRIQCgNpc3MYAiABKAlSA2lzcxIQCgNzdWIYAyABKA'
    'lSA3N1YhIQCgNzaWQYBCABKAlSA3NpZBIQCgNpYXQYBSABKANSA2lhdBIQCgNleHAYBiABKANS'
    'A2V4cBIQCgNuYmYYByABKANSA25iZhInCg51c2VybWFuYWdlbWVudBjoByABKAhSDnVzZXJtYW'
    '5hZ2VtZW50EicKDnJlbW90ZV9jb250cm9sGOkHIAEoCFIOcmVtb3RlX2NvbnRyb2wSIwoMYmls'
    'bGluZ19yZWFkGOoHIAEoCFIMYmlsbGluZ19yZWFkEicKDmJpbGxpbmdfbW9kaWZ5GOsHIAEoCF'
    'IOYmlsbGluZ19tb2RpZnkSKwoQY29tbXVuaXR5X21vZGlmeRjsByABKAhSEGNvbW11bml0eV9t'
    'b2RpZnkSJwoOYXJjaGl2ZV91cGxvYWQY7QcgASgEUg5hcmNoaXZlX3VwbG9hZBIrChBhcmNoaX'
    'ZlX2Rvd25sb2FkGO4HIAEoBFIQYXJjaGl2ZV9kb3dubG9hZBIjCgxsaWJyYXJ5X3JlYWQY7wcg'
    'ASgIUgxsaWJyYXJ5X3JlYWQSJwoObGlicmFyeV9tb2RpZnkY8AcgASgIUg5saWJyYXJ5X21vZG'
    'lmeRIjCgxhcmNoaXZlX3N5bmMY8QcgASgIUgxhcmNoaXZlX3N5bmMSJwoOY29tbXVuaXR5X3N5'
    'bmMY8gcgASgIUg5jb21tdW5pdHlfc3luYxIfCgpsb2NhbF9vbmx5GNAPIAEoCFIKbG9jYWxfb2'
    '5seUoFCAkQ6AdKBgj0BxDQDw==');

@$core.Deprecated('Use authzRequestDescriptor instead')
const AuthzRequest$json = {
  '1': 'AuthzRequest',
};

/// Descriptor for `AuthzRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authzRequestDescriptor =
    $convert.base64Decode('CgxBdXRoelJlcXVlc3Q=');

@$core.Deprecated('Use authzResponseDescriptor instead')
const AuthzResponse$json = {
  '1': 'AuthzResponse',
  '2': [
    {'1': 'bearer', '3': 1, '4': 1, '5': 9, '10': 'bearer'},
    {'1': 'token', '3': 2, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `AuthzResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authzResponseDescriptor = $convert.base64Decode(
    'Cg1BdXRoelJlc3BvbnNlEhYKBmJlYXJlchgBIAEoCVIGYmVhcmVyEiEKBXRva2VuGAIgASgLMg'
    'subWV0YS5Ub2tlblIFdG9rZW4=');

@$core.Deprecated('Use authzGrantRequestDescriptor instead')
const AuthzGrantRequest$json = {
  '1': 'AuthzGrantRequest',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `AuthzGrantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authzGrantRequestDescriptor = $convert.base64Decode(
    'ChFBdXRoekdyYW50UmVxdWVzdBIhCgV0b2tlbhgBIAEoCzILLm1ldGEuVG9rZW5SBXRva2Vu');

@$core.Deprecated('Use authzGrantResponseDescriptor instead')
const AuthzGrantResponse$json = {
  '1': 'AuthzGrantResponse',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `AuthzGrantResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authzGrantResponseDescriptor = $convert.base64Decode(
    'ChJBdXRoekdyYW50UmVzcG9uc2USIQoFdG9rZW4YASABKAsyCy5tZXRhLlRva2VuUgV0b2tlbg'
    '==');

@$core.Deprecated('Use authzRevokeRequestDescriptor instead')
const AuthzRevokeRequest$json = {
  '1': 'AuthzRevokeRequest',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `AuthzRevokeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authzRevokeRequestDescriptor = $convert.base64Decode(
    'ChJBdXRoelJldm9rZVJlcXVlc3QSIQoFdG9rZW4YASABKAsyCy5tZXRhLlRva2VuUgV0b2tlbg'
    '==');

@$core.Deprecated('Use authzRevokeResponseDescriptor instead')
const AuthzRevokeResponse$json = {
  '1': 'AuthzRevokeResponse',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `AuthzRevokeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authzRevokeResponseDescriptor = $convert.base64Decode(
    'ChNBdXRoelJldm9rZVJlc3BvbnNlEiEKBXRva2VuGAEgASgLMgsubWV0YS5Ub2tlblIFdG9rZW'
    '4=');

@$core.Deprecated('Use authzProfileRequestDescriptor instead')
const AuthzProfileRequest$json = {
  '1': 'AuthzProfileRequest',
  '2': [
    {'1': 'profile_id', '3': 1, '4': 1, '5': 9, '10': 'profile_id'},
  ],
};

/// Descriptor for `AuthzProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authzProfileRequestDescriptor = $convert.base64Decode(
    'ChNBdXRoelByb2ZpbGVSZXF1ZXN0Eh4KCnByb2ZpbGVfaWQYASABKAlSCnByb2ZpbGVfaWQ=');

@$core.Deprecated('Use authzProfileResponseDescriptor instead')
const AuthzProfileResponse$json = {
  '1': 'AuthzProfileResponse',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `AuthzProfileResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authzProfileResponseDescriptor = $convert.base64Decode(
    'ChRBdXRoelByb2ZpbGVSZXNwb25zZRIhCgV0b2tlbhgBIAEoCzILLm1ldGEuVG9rZW5SBXRva2'
    'Vu');
