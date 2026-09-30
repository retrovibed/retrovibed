// This is a generated file - do not edit.
//
// Generated from media/media.autoimport.proto.

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

@$core.Deprecated('Use autoimportDirectoryDescriptor instead')
const AutoimportDirectory$json = {
  '1': 'AutoimportDirectory',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'created_at', '3': 2, '4': 1, '5': 9, '10': 'created_at'},
    {'1': 'updated_at', '3': 3, '4': 1, '5': 9, '10': 'updated_at'},
    {'1': 'path', '3': 4, '4': 1, '5': 9, '10': 'path'},
    {'1': 'description', '3': 5, '4': 1, '5': 9, '10': 'description'},
    {'1': 'debounce', '3': 6, '4': 1, '5': 4, '10': 'debounce'},
    {'1': 'mode', '3': 7, '4': 1, '5': 13, '10': 'mode'},
    {
      '1': 'library_directory_id',
      '3': 8,
      '4': 1,
      '5': 9,
      '10': 'library_directory_id'
    },
    {'1': 'last_scanned_at', '3': 9, '4': 1, '5': 9, '10': 'last_scanned_at'},
    {'1': 'pending', '3': 10, '4': 1, '5': 4, '10': 'pending'},
    {'1': 'imported', '3': 11, '4': 1, '5': 4, '10': 'imported'},
  ],
};

/// Descriptor for `AutoimportDirectory`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoimportDirectoryDescriptor = $convert.base64Decode(
    'ChNBdXRvaW1wb3J0RGlyZWN0b3J5Eg4KAmlkGAEgASgJUgJpZBIeCgpjcmVhdGVkX2F0GAIgAS'
    'gJUgpjcmVhdGVkX2F0Eh4KCnVwZGF0ZWRfYXQYAyABKAlSCnVwZGF0ZWRfYXQSEgoEcGF0aBgE'
    'IAEoCVIEcGF0aBIgCgtkZXNjcmlwdGlvbhgFIAEoCVILZGVzY3JpcHRpb24SGgoIZGVib3VuY2'
    'UYBiABKARSCGRlYm91bmNlEhIKBG1vZGUYByABKA1SBG1vZGUSMgoUbGlicmFyeV9kaXJlY3Rv'
    'cnlfaWQYCCABKAlSFGxpYnJhcnlfZGlyZWN0b3J5X2lkEigKD2xhc3Rfc2Nhbm5lZF9hdBgJIA'
    'EoCVIPbGFzdF9zY2FubmVkX2F0EhgKB3BlbmRpbmcYCiABKARSB3BlbmRpbmcSGgoIaW1wb3J0'
    'ZWQYCyABKARSCGltcG9ydGVk');

@$core.Deprecated('Use autoimportDirectorySearchRequestDescriptor instead')
const AutoimportDirectorySearchRequest$json = {
  '1': 'AutoimportDirectorySearchRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 3, '5': 9, '10': 'id'},
    {'1': 'query', '3': 2, '4': 1, '5': 9, '10': 'query'},
    {'1': 'offset', '3': 900, '4': 1, '5': 4, '10': 'offset'},
    {'1': 'limit', '3': 901, '4': 1, '5': 4, '10': 'limit'},
  ],
  '9': [
    {'1': 3, '2': 900},
    {'1': 902, '2': 1000},
  ],
};

/// Descriptor for `AutoimportDirectorySearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoimportDirectorySearchRequestDescriptor =
    $convert.base64Decode(
        'CiBBdXRvaW1wb3J0RGlyZWN0b3J5U2VhcmNoUmVxdWVzdBIOCgJpZBgBIAMoCVICaWQSFAoFcX'
        'VlcnkYAiABKAlSBXF1ZXJ5EhcKBm9mZnNldBiEByABKARSBm9mZnNldBIVCgVsaW1pdBiFByAB'
        'KARSBWxpbWl0SgUIAxCEB0oGCIYHEOgH');

@$core.Deprecated('Use autoimportDirectorySearchResponseDescriptor instead')
const AutoimportDirectorySearchResponse$json = {
  '1': 'AutoimportDirectorySearchResponse',
  '2': [
    {
      '1': 'next',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.AutoimportDirectorySearchRequest',
      '10': 'next'
    },
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.media.AutoimportDirectory',
      '10': 'items'
    },
  ],
};

/// Descriptor for `AutoimportDirectorySearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoimportDirectorySearchResponseDescriptor =
    $convert.base64Decode(
        'CiFBdXRvaW1wb3J0RGlyZWN0b3J5U2VhcmNoUmVzcG9uc2USOwoEbmV4dBgBIAEoCzInLm1lZG'
        'lhLkF1dG9pbXBvcnREaXJlY3RvcnlTZWFyY2hSZXF1ZXN0UgRuZXh0EjAKBWl0ZW1zGAIgAygL'
        'MhoubWVkaWEuQXV0b2ltcG9ydERpcmVjdG9yeVIFaXRlbXM=');

@$core.Deprecated('Use autoimportDirectoryCreateRequestDescriptor instead')
const AutoimportDirectoryCreateRequest$json = {
  '1': 'AutoimportDirectoryCreateRequest',
  '2': [
    {
      '1': 'directory',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.AutoimportDirectory',
      '10': 'directory'
    },
  ],
};

/// Descriptor for `AutoimportDirectoryCreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoimportDirectoryCreateRequestDescriptor =
    $convert.base64Decode(
        'CiBBdXRvaW1wb3J0RGlyZWN0b3J5Q3JlYXRlUmVxdWVzdBI4CglkaXJlY3RvcnkYASABKAsyGi'
        '5tZWRpYS5BdXRvaW1wb3J0RGlyZWN0b3J5UglkaXJlY3Rvcnk=');

@$core.Deprecated('Use autoimportDirectoryCreateResponseDescriptor instead')
const AutoimportDirectoryCreateResponse$json = {
  '1': 'AutoimportDirectoryCreateResponse',
  '2': [
    {
      '1': 'directory',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.AutoimportDirectory',
      '10': 'directory'
    },
  ],
};

/// Descriptor for `AutoimportDirectoryCreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoimportDirectoryCreateResponseDescriptor =
    $convert.base64Decode(
        'CiFBdXRvaW1wb3J0RGlyZWN0b3J5Q3JlYXRlUmVzcG9uc2USOAoJZGlyZWN0b3J5GAEgASgLMh'
        'oubWVkaWEuQXV0b2ltcG9ydERpcmVjdG9yeVIJZGlyZWN0b3J5');

@$core.Deprecated('Use autoimportDirectoryLookupResponseDescriptor instead')
const AutoimportDirectoryLookupResponse$json = {
  '1': 'AutoimportDirectoryLookupResponse',
  '2': [
    {
      '1': 'directory',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.AutoimportDirectory',
      '10': 'directory'
    },
  ],
};

/// Descriptor for `AutoimportDirectoryLookupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoimportDirectoryLookupResponseDescriptor =
    $convert.base64Decode(
        'CiFBdXRvaW1wb3J0RGlyZWN0b3J5TG9va3VwUmVzcG9uc2USOAoJZGlyZWN0b3J5GAEgASgLMh'
        'oubWVkaWEuQXV0b2ltcG9ydERpcmVjdG9yeVIJZGlyZWN0b3J5');

@$core.Deprecated('Use autoimportDirectoryUpdateRequestDescriptor instead')
const AutoimportDirectoryUpdateRequest$json = {
  '1': 'AutoimportDirectoryUpdateRequest',
  '2': [
    {
      '1': 'directory',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.AutoimportDirectory',
      '10': 'directory'
    },
  ],
};

/// Descriptor for `AutoimportDirectoryUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoimportDirectoryUpdateRequestDescriptor =
    $convert.base64Decode(
        'CiBBdXRvaW1wb3J0RGlyZWN0b3J5VXBkYXRlUmVxdWVzdBI4CglkaXJlY3RvcnkYASABKAsyGi'
        '5tZWRpYS5BdXRvaW1wb3J0RGlyZWN0b3J5UglkaXJlY3Rvcnk=');

@$core.Deprecated('Use autoimportDirectoryUpdateResponseDescriptor instead')
const AutoimportDirectoryUpdateResponse$json = {
  '1': 'AutoimportDirectoryUpdateResponse',
  '2': [
    {
      '1': 'directory',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.AutoimportDirectory',
      '10': 'directory'
    },
  ],
};

/// Descriptor for `AutoimportDirectoryUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoimportDirectoryUpdateResponseDescriptor =
    $convert.base64Decode(
        'CiFBdXRvaW1wb3J0RGlyZWN0b3J5VXBkYXRlUmVzcG9uc2USOAoJZGlyZWN0b3J5GAEgASgLMh'
        'oubWVkaWEuQXV0b2ltcG9ydERpcmVjdG9yeVIJZGlyZWN0b3J5');

@$core.Deprecated('Use autoimportDirectoryDeleteResponseDescriptor instead')
const AutoimportDirectoryDeleteResponse$json = {
  '1': 'AutoimportDirectoryDeleteResponse',
  '2': [
    {
      '1': 'directory',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.media.AutoimportDirectory',
      '10': 'directory'
    },
  ],
};

/// Descriptor for `AutoimportDirectoryDeleteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoimportDirectoryDeleteResponseDescriptor =
    $convert.base64Decode(
        'CiFBdXRvaW1wb3J0RGlyZWN0b3J5RGVsZXRlUmVzcG9uc2USOAoJZGlyZWN0b3J5GAEgASgLMh'
        'oubWVkaWEuQXV0b2ltcG9ydERpcmVjdG9yeVIJZGlyZWN0b3J5');
