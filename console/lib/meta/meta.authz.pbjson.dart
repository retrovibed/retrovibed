//
//  Generated code. Do not modify.
//  source: meta.authz.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use bearerDescriptor instead')
const Bearer$json = {
  '1': 'Bearer',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'jti'},
    {'1': 'issuer', '3': 2, '4': 1, '5': 9, '10': 'iss'},
    {'1': 'profile_id', '3': 3, '4': 1, '5': 9, '10': 'sub'},
    {'1': 'session_id', '3': 4, '4': 1, '5': 9, '10': 'sid'},
    {'1': 'issued', '3': 5, '4': 1, '5': 3, '10': 'iat'},
    {'1': 'expires', '3': 6, '4': 1, '5': 3, '10': 'exp'},
    {'1': 'not_before', '3': 7, '4': 1, '5': 3, '10': 'nbf'},
  ],
};

/// Descriptor for `Bearer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bearerDescriptor = $convert.base64Decode(
    'CgZCZWFyZXISDwoCaWQYASABKAlSA2p0aRITCgZpc3N1ZXIYAiABKAlSA2lzcxIXCgpwcm9maW'
    'xlX2lkGAMgASgJUgNzdWISFwoKc2Vzc2lvbl9pZBgEIAEoCVIDc2lkEhMKBmlzc3VlZBgFIAEo'
    'A1IDaWF0EhQKB2V4cGlyZXMYBiABKANSA2V4cBIXCgpub3RfYmVmb3JlGAcgASgDUgNuYmY=');

@$core.Deprecated('Use tokenDescriptor instead')
const Token$json = {
  '1': 'Token',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'jti'},
    {'1': 'account_id', '3': 2, '4': 1, '5': 9, '10': 'iss'},
    {'1': 'profile_id', '3': 3, '4': 1, '5': 9, '10': 'sub'},
    {'1': 'session_id', '3': 4, '4': 1, '5': 9, '10': 'sid'},
    {'1': 'issued', '3': 5, '4': 1, '5': 3, '10': 'iat'},
    {'1': 'expires', '3': 6, '4': 1, '5': 3, '10': 'exp'},
    {'1': 'not_before', '3': 7, '4': 1, '5': 3, '10': 'nbf'},
    {'1': 'usermanagement', '3': 1000, '4': 1, '5': 8, '10': 'usermanagement'},
  ],
  '9': [
    {'1': 9, '2': 1000},
  ],
};

/// Descriptor for `Token`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tokenDescriptor = $convert.base64Decode(
    'CgVUb2tlbhIPCgJpZBgBIAEoCVIDanRpEhcKCmFjY291bnRfaWQYAiABKAlSA2lzcxIXCgpwcm'
    '9maWxlX2lkGAMgASgJUgNzdWISFwoKc2Vzc2lvbl9pZBgEIAEoCVIDc2lkEhMKBmlzc3VlZBgF'
    'IAEoA1IDaWF0EhQKB2V4cGlyZXMYBiABKANSA2V4cBIXCgpub3RfYmVmb3JlGAcgASgDUgNuYm'
    'YSJwoOdXNlcm1hbmFnZW1lbnQY6AcgASgIUg51c2VybWFuYWdlbWVudEoFCAkQ6Ac=');

@$core.Deprecated('Use authzRequestDescriptor instead')
const AuthzRequest$json = {
  '1': 'AuthzRequest',
};

/// Descriptor for `AuthzRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authzRequestDescriptor = $convert.base64Decode(
    'CgxBdXRoelJlcXVlc3Q=');

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

@$core.Deprecated('Use grantRequestDescriptor instead')
const GrantRequest$json = {
  '1': 'GrantRequest',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `GrantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List grantRequestDescriptor = $convert.base64Decode(
    'CgxHcmFudFJlcXVlc3QSIQoFdG9rZW4YASABKAsyCy5tZXRhLlRva2VuUgV0b2tlbg==');

@$core.Deprecated('Use grantResponseDescriptor instead')
const GrantResponse$json = {
  '1': 'GrantResponse',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `GrantResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List grantResponseDescriptor = $convert.base64Decode(
    'Cg1HcmFudFJlc3BvbnNlEiEKBXRva2VuGAEgASgLMgsubWV0YS5Ub2tlblIFdG9rZW4=');

@$core.Deprecated('Use revokeRequestDescriptor instead')
const RevokeRequest$json = {
  '1': 'RevokeRequest',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `RevokeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revokeRequestDescriptor = $convert.base64Decode(
    'Cg1SZXZva2VSZXF1ZXN0EiEKBXRva2VuGAEgASgLMgsubWV0YS5Ub2tlblIFdG9rZW4=');

@$core.Deprecated('Use revokeResponseDescriptor instead')
const RevokeResponse$json = {
  '1': 'RevokeResponse',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `RevokeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revokeResponseDescriptor = $convert.base64Decode(
    'Cg5SZXZva2VSZXNwb25zZRIhCgV0b2tlbhgBIAEoCzILLm1ldGEuVG9rZW5SBXRva2Vu');

@$core.Deprecated('Use profileRequestDescriptor instead')
const ProfileRequest$json = {
  '1': 'ProfileRequest',
  '2': [
    {'1': 'profile_id', '3': 1, '4': 1, '5': 9, '10': 'profile_id'},
  ],
};

/// Descriptor for `ProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileRequestDescriptor = $convert.base64Decode(
    'Cg5Qcm9maWxlUmVxdWVzdBIeCgpwcm9maWxlX2lkGAEgASgJUgpwcm9maWxlX2lk');

@$core.Deprecated('Use profileResponseDescriptor instead')
const ProfileResponse$json = {
  '1': 'ProfileResponse',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 11, '6': '.meta.Token', '10': 'token'},
  ],
};

/// Descriptor for `ProfileResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileResponseDescriptor = $convert.base64Decode(
    'Cg9Qcm9maWxlUmVzcG9uc2USIQoFdG9rZW4YASABKAsyCy5tZXRhLlRva2VuUgV0b2tlbg==');

