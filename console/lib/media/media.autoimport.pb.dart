// This is a generated file - do not edit.
//
// Generated from media/media.autoimport.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class AutoimportDirectory extends $pb.GeneratedMessage {
  factory AutoimportDirectory({
    $core.String? id,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? path,
    $core.String? description,
    $fixnum.Int64? debounce,
    $core.int? mode,
    $core.String? libraryDirectoryId,
    $core.String? lastScannedAt,
    $fixnum.Int64? pending,
    $fixnum.Int64? imported,
  }) {
    final result = AutoimportDirectory._();
    if (id != null) result.id = id;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (path != null) result.path = path;
    if (description != null) result.description = description;
    if (debounce != null) result.debounce = debounce;
    if (mode != null) result.mode = mode;
    if (libraryDirectoryId != null)
      result.libraryDirectoryId = libraryDirectoryId;
    if (lastScannedAt != null) result.lastScannedAt = lastScannedAt;
    if (pending != null) result.pending = pending;
    if (imported != null) result.imported = imported;
    return result;
  }

  AutoimportDirectory._();

  factory AutoimportDirectory.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectory()..mergeFromBuffer(data, registry);
  factory AutoimportDirectory.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectory()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoimportDirectory',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: AutoimportDirectory.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'created_at')
    ..aOS(3, _omitFieldNames ? '' : 'updated_at')
    ..aOS(4, _omitFieldNames ? '' : 'path')
    ..aOS(5, _omitFieldNames ? '' : 'description')
    ..a<$fixnum.Int64>(
        6, _omitFieldNames ? '' : 'debounce', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aI(7, _omitFieldNames ? '' : 'mode', fieldType: $pb.PbFieldType.OU3)
    ..aOS(8, _omitFieldNames ? '' : 'library_directory_id')
    ..aOS(9, _omitFieldNames ? '' : 'last_scanned_at')
    ..a<$fixnum.Int64>(
        10, _omitFieldNames ? '' : 'pending', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        11, _omitFieldNames ? '' : 'imported', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectory clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectory copyWith(void Function(AutoimportDirectory) updates) =>
      super.copyWith((message) => updates(message as AutoimportDirectory))
          as AutoimportDirectory;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use AutoimportDirectory() / AutoimportDirectory.new instead')
  static AutoimportDirectory create() => AutoimportDirectory._();
  static $pb.GeneratedMessage $_createMessage() => AutoimportDirectory._();
  @$core.override
  AutoimportDirectory createEmptyInstance() => AutoimportDirectory._();
  @$core.pragma('dart2js:noInline')
  static AutoimportDirectory getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoimportDirectory>(
          AutoimportDirectory.$_createMessage);
  static AutoimportDirectory? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get createdAt => $_getSZ(1);
  @$pb.TagNumber(2)
  set createdAt($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCreatedAt() => $_has(1);
  @$pb.TagNumber(2)
  void clearCreatedAt() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get updatedAt => $_getSZ(2);
  @$pb.TagNumber(3)
  set updatedAt($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasUpdatedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearUpdatedAt() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get path => $_getSZ(3);
  @$pb.TagNumber(4)
  set path($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPath() => $_has(3);
  @$pb.TagNumber(4)
  void clearPath() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(5)
  set description($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(5)
  void clearDescription() => $_clearField(5);

  /// seconds a file must go unmodified before it is imported.
  @$pb.TagNumber(6)
  $fixnum.Int64 get debounce => $_getI64(5);
  @$pb.TagNumber(6)
  set debounce($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasDebounce() => $_has(5);
  @$pb.TagNumber(6)
  void clearDebounce() => $_clearField(6);

  /// 0 = copy (keep the original), 1 = move (remove the original after import).
  @$pb.TagNumber(7)
  $core.int get mode => $_getIZ(6);
  @$pb.TagNumber(7)
  set mode($core.int value) => $_setUnsignedInt32(6, value);
  @$pb.TagNumber(7)
  $core.bool hasMode() => $_has(6);
  @$pb.TagNumber(7)
  void clearMode() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get libraryDirectoryId => $_getSZ(7);
  @$pb.TagNumber(8)
  set libraryDirectoryId($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasLibraryDirectoryId() => $_has(7);
  @$pb.TagNumber(8)
  void clearLibraryDirectoryId() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get lastScannedAt => $_getSZ(8);
  @$pb.TagNumber(9)
  set lastScannedAt($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasLastScannedAt() => $_has(8);
  @$pb.TagNumber(9)
  void clearLastScannedAt() => $_clearField(9);

  @$pb.TagNumber(10)
  $fixnum.Int64 get pending => $_getI64(9);
  @$pb.TagNumber(10)
  set pending($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(10)
  $core.bool hasPending() => $_has(9);
  @$pb.TagNumber(10)
  void clearPending() => $_clearField(10);

  @$pb.TagNumber(11)
  $fixnum.Int64 get imported => $_getI64(10);
  @$pb.TagNumber(11)
  set imported($fixnum.Int64 value) => $_setInt64(10, value);
  @$pb.TagNumber(11)
  $core.bool hasImported() => $_has(10);
  @$pb.TagNumber(11)
  void clearImported() => $_clearField(11);
}

class AutoimportDirectorySearchRequest extends $pb.GeneratedMessage {
  factory AutoimportDirectorySearchRequest({
    $core.Iterable<$core.String>? id,
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = AutoimportDirectorySearchRequest._();
    if (id != null) result.id.addAll(id);
    if (query != null) result.query = query;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  AutoimportDirectorySearchRequest._();

  factory AutoimportDirectorySearchRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectorySearchRequest()..mergeFromBuffer(data, registry);
  factory AutoimportDirectorySearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectorySearchRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoimportDirectorySearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: AutoimportDirectorySearchRequest.$_createMessage)
    ..pPS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectorySearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectorySearchRequest copyWith(
          void Function(AutoimportDirectorySearchRequest) updates) =>
      super.copyWith(
              (message) => updates(message as AutoimportDirectorySearchRequest))
          as AutoimportDirectorySearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AutoimportDirectorySearchRequest() / AutoimportDirectorySearchRequest.new instead')
  static AutoimportDirectorySearchRequest create() =>
      AutoimportDirectorySearchRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      AutoimportDirectorySearchRequest._();
  @$core.override
  AutoimportDirectorySearchRequest createEmptyInstance() =>
      AutoimportDirectorySearchRequest._();
  @$core.pragma('dart2js:noInline')
  static AutoimportDirectorySearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoimportDirectorySearchRequest>(
          AutoimportDirectorySearchRequest.$_createMessage);
  static AutoimportDirectorySearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get id => $_getList(0);

  @$pb.TagNumber(2)
  $core.String get query => $_getSZ(1);
  @$pb.TagNumber(2)
  set query($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasQuery() => $_has(1);
  @$pb.TagNumber(2)
  void clearQuery() => $_clearField(2);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(2);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(2);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(3);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class AutoimportDirectorySearchResponse extends $pb.GeneratedMessage {
  factory AutoimportDirectorySearchResponse({
    AutoimportDirectorySearchRequest? next,
    $core.Iterable<AutoimportDirectory>? items,
  }) {
    final result = AutoimportDirectorySearchResponse._();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  AutoimportDirectorySearchResponse._();

  factory AutoimportDirectorySearchResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectorySearchResponse()..mergeFromBuffer(data, registry);
  factory AutoimportDirectorySearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectorySearchResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoimportDirectorySearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: AutoimportDirectorySearchResponse.$_createMessage)
    ..aOM<AutoimportDirectorySearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: AutoimportDirectorySearchRequest.$_createMessage)
    ..pPM<AutoimportDirectory>(2, _omitFieldNames ? '' : 'items',
        subBuilder: AutoimportDirectory.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectorySearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectorySearchResponse copyWith(
          void Function(AutoimportDirectorySearchResponse) updates) =>
      super.copyWith((message) =>
              updates(message as AutoimportDirectorySearchResponse))
          as AutoimportDirectorySearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AutoimportDirectorySearchResponse() / AutoimportDirectorySearchResponse.new instead')
  static AutoimportDirectorySearchResponse create() =>
      AutoimportDirectorySearchResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      AutoimportDirectorySearchResponse._();
  @$core.override
  AutoimportDirectorySearchResponse createEmptyInstance() =>
      AutoimportDirectorySearchResponse._();
  @$core.pragma('dart2js:noInline')
  static AutoimportDirectorySearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoimportDirectorySearchResponse>(
          AutoimportDirectorySearchResponse.$_createMessage);
  static AutoimportDirectorySearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AutoimportDirectorySearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(AutoimportDirectorySearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  AutoimportDirectorySearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<AutoimportDirectory> get items => $_getList(1);
}

class AutoimportDirectoryCreateRequest extends $pb.GeneratedMessage {
  factory AutoimportDirectoryCreateRequest({
    AutoimportDirectory? directory,
  }) {
    final result = AutoimportDirectoryCreateRequest._();
    if (directory != null) result.directory = directory;
    return result;
  }

  AutoimportDirectoryCreateRequest._();

  factory AutoimportDirectoryCreateRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryCreateRequest()..mergeFromBuffer(data, registry);
  factory AutoimportDirectoryCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryCreateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoimportDirectoryCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: AutoimportDirectoryCreateRequest.$_createMessage)
    ..aOM<AutoimportDirectory>(1, _omitFieldNames ? '' : 'directory',
        subBuilder: AutoimportDirectory.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryCreateRequest copyWith(
          void Function(AutoimportDirectoryCreateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as AutoimportDirectoryCreateRequest))
          as AutoimportDirectoryCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AutoimportDirectoryCreateRequest() / AutoimportDirectoryCreateRequest.new instead')
  static AutoimportDirectoryCreateRequest create() =>
      AutoimportDirectoryCreateRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      AutoimportDirectoryCreateRequest._();
  @$core.override
  AutoimportDirectoryCreateRequest createEmptyInstance() =>
      AutoimportDirectoryCreateRequest._();
  @$core.pragma('dart2js:noInline')
  static AutoimportDirectoryCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoimportDirectoryCreateRequest>(
          AutoimportDirectoryCreateRequest.$_createMessage);
  static AutoimportDirectoryCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  AutoimportDirectory get directory => $_getN(0);
  @$pb.TagNumber(1)
  set directory(AutoimportDirectory value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDirectory() => $_has(0);
  @$pb.TagNumber(1)
  void clearDirectory() => $_clearField(1);
  @$pb.TagNumber(1)
  AutoimportDirectory ensureDirectory() => $_ensure(0);
}

class AutoimportDirectoryCreateResponse extends $pb.GeneratedMessage {
  factory AutoimportDirectoryCreateResponse({
    AutoimportDirectory? directory,
  }) {
    final result = AutoimportDirectoryCreateResponse._();
    if (directory != null) result.directory = directory;
    return result;
  }

  AutoimportDirectoryCreateResponse._();

  factory AutoimportDirectoryCreateResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryCreateResponse()..mergeFromBuffer(data, registry);
  factory AutoimportDirectoryCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryCreateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoimportDirectoryCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: AutoimportDirectoryCreateResponse.$_createMessage)
    ..aOM<AutoimportDirectory>(1, _omitFieldNames ? '' : 'directory',
        subBuilder: AutoimportDirectory.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryCreateResponse copyWith(
          void Function(AutoimportDirectoryCreateResponse) updates) =>
      super.copyWith((message) =>
              updates(message as AutoimportDirectoryCreateResponse))
          as AutoimportDirectoryCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AutoimportDirectoryCreateResponse() / AutoimportDirectoryCreateResponse.new instead')
  static AutoimportDirectoryCreateResponse create() =>
      AutoimportDirectoryCreateResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      AutoimportDirectoryCreateResponse._();
  @$core.override
  AutoimportDirectoryCreateResponse createEmptyInstance() =>
      AutoimportDirectoryCreateResponse._();
  @$core.pragma('dart2js:noInline')
  static AutoimportDirectoryCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoimportDirectoryCreateResponse>(
          AutoimportDirectoryCreateResponse.$_createMessage);
  static AutoimportDirectoryCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AutoimportDirectory get directory => $_getN(0);
  @$pb.TagNumber(1)
  set directory(AutoimportDirectory value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDirectory() => $_has(0);
  @$pb.TagNumber(1)
  void clearDirectory() => $_clearField(1);
  @$pb.TagNumber(1)
  AutoimportDirectory ensureDirectory() => $_ensure(0);
}

class AutoimportDirectoryLookupResponse extends $pb.GeneratedMessage {
  factory AutoimportDirectoryLookupResponse({
    AutoimportDirectory? directory,
  }) {
    final result = AutoimportDirectoryLookupResponse._();
    if (directory != null) result.directory = directory;
    return result;
  }

  AutoimportDirectoryLookupResponse._();

  factory AutoimportDirectoryLookupResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryLookupResponse()..mergeFromBuffer(data, registry);
  factory AutoimportDirectoryLookupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryLookupResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoimportDirectoryLookupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: AutoimportDirectoryLookupResponse.$_createMessage)
    ..aOM<AutoimportDirectory>(1, _omitFieldNames ? '' : 'directory',
        subBuilder: AutoimportDirectory.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryLookupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryLookupResponse copyWith(
          void Function(AutoimportDirectoryLookupResponse) updates) =>
      super.copyWith((message) =>
              updates(message as AutoimportDirectoryLookupResponse))
          as AutoimportDirectoryLookupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AutoimportDirectoryLookupResponse() / AutoimportDirectoryLookupResponse.new instead')
  static AutoimportDirectoryLookupResponse create() =>
      AutoimportDirectoryLookupResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      AutoimportDirectoryLookupResponse._();
  @$core.override
  AutoimportDirectoryLookupResponse createEmptyInstance() =>
      AutoimportDirectoryLookupResponse._();
  @$core.pragma('dart2js:noInline')
  static AutoimportDirectoryLookupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoimportDirectoryLookupResponse>(
          AutoimportDirectoryLookupResponse.$_createMessage);
  static AutoimportDirectoryLookupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AutoimportDirectory get directory => $_getN(0);
  @$pb.TagNumber(1)
  set directory(AutoimportDirectory value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDirectory() => $_has(0);
  @$pb.TagNumber(1)
  void clearDirectory() => $_clearField(1);
  @$pb.TagNumber(1)
  AutoimportDirectory ensureDirectory() => $_ensure(0);
}

class AutoimportDirectoryUpdateRequest extends $pb.GeneratedMessage {
  factory AutoimportDirectoryUpdateRequest({
    AutoimportDirectory? directory,
  }) {
    final result = AutoimportDirectoryUpdateRequest._();
    if (directory != null) result.directory = directory;
    return result;
  }

  AutoimportDirectoryUpdateRequest._();

  factory AutoimportDirectoryUpdateRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryUpdateRequest()..mergeFromBuffer(data, registry);
  factory AutoimportDirectoryUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryUpdateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoimportDirectoryUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: AutoimportDirectoryUpdateRequest.$_createMessage)
    ..aOM<AutoimportDirectory>(1, _omitFieldNames ? '' : 'directory',
        subBuilder: AutoimportDirectory.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryUpdateRequest copyWith(
          void Function(AutoimportDirectoryUpdateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as AutoimportDirectoryUpdateRequest))
          as AutoimportDirectoryUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AutoimportDirectoryUpdateRequest() / AutoimportDirectoryUpdateRequest.new instead')
  static AutoimportDirectoryUpdateRequest create() =>
      AutoimportDirectoryUpdateRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      AutoimportDirectoryUpdateRequest._();
  @$core.override
  AutoimportDirectoryUpdateRequest createEmptyInstance() =>
      AutoimportDirectoryUpdateRequest._();
  @$core.pragma('dart2js:noInline')
  static AutoimportDirectoryUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoimportDirectoryUpdateRequest>(
          AutoimportDirectoryUpdateRequest.$_createMessage);
  static AutoimportDirectoryUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  AutoimportDirectory get directory => $_getN(0);
  @$pb.TagNumber(1)
  set directory(AutoimportDirectory value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDirectory() => $_has(0);
  @$pb.TagNumber(1)
  void clearDirectory() => $_clearField(1);
  @$pb.TagNumber(1)
  AutoimportDirectory ensureDirectory() => $_ensure(0);
}

class AutoimportDirectoryUpdateResponse extends $pb.GeneratedMessage {
  factory AutoimportDirectoryUpdateResponse({
    AutoimportDirectory? directory,
  }) {
    final result = AutoimportDirectoryUpdateResponse._();
    if (directory != null) result.directory = directory;
    return result;
  }

  AutoimportDirectoryUpdateResponse._();

  factory AutoimportDirectoryUpdateResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryUpdateResponse()..mergeFromBuffer(data, registry);
  factory AutoimportDirectoryUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryUpdateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoimportDirectoryUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: AutoimportDirectoryUpdateResponse.$_createMessage)
    ..aOM<AutoimportDirectory>(1, _omitFieldNames ? '' : 'directory',
        subBuilder: AutoimportDirectory.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryUpdateResponse copyWith(
          void Function(AutoimportDirectoryUpdateResponse) updates) =>
      super.copyWith((message) =>
              updates(message as AutoimportDirectoryUpdateResponse))
          as AutoimportDirectoryUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AutoimportDirectoryUpdateResponse() / AutoimportDirectoryUpdateResponse.new instead')
  static AutoimportDirectoryUpdateResponse create() =>
      AutoimportDirectoryUpdateResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      AutoimportDirectoryUpdateResponse._();
  @$core.override
  AutoimportDirectoryUpdateResponse createEmptyInstance() =>
      AutoimportDirectoryUpdateResponse._();
  @$core.pragma('dart2js:noInline')
  static AutoimportDirectoryUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoimportDirectoryUpdateResponse>(
          AutoimportDirectoryUpdateResponse.$_createMessage);
  static AutoimportDirectoryUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AutoimportDirectory get directory => $_getN(0);
  @$pb.TagNumber(1)
  set directory(AutoimportDirectory value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDirectory() => $_has(0);
  @$pb.TagNumber(1)
  void clearDirectory() => $_clearField(1);
  @$pb.TagNumber(1)
  AutoimportDirectory ensureDirectory() => $_ensure(0);
}

class AutoimportDirectoryDeleteResponse extends $pb.GeneratedMessage {
  factory AutoimportDirectoryDeleteResponse({
    AutoimportDirectory? directory,
  }) {
    final result = AutoimportDirectoryDeleteResponse._();
    if (directory != null) result.directory = directory;
    return result;
  }

  AutoimportDirectoryDeleteResponse._();

  factory AutoimportDirectoryDeleteResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryDeleteResponse()..mergeFromBuffer(data, registry);
  factory AutoimportDirectoryDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AutoimportDirectoryDeleteResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoimportDirectoryDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: AutoimportDirectoryDeleteResponse.$_createMessage)
    ..aOM<AutoimportDirectory>(1, _omitFieldNames ? '' : 'directory',
        subBuilder: AutoimportDirectory.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoimportDirectoryDeleteResponse copyWith(
          void Function(AutoimportDirectoryDeleteResponse) updates) =>
      super.copyWith((message) =>
              updates(message as AutoimportDirectoryDeleteResponse))
          as AutoimportDirectoryDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AutoimportDirectoryDeleteResponse() / AutoimportDirectoryDeleteResponse.new instead')
  static AutoimportDirectoryDeleteResponse create() =>
      AutoimportDirectoryDeleteResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      AutoimportDirectoryDeleteResponse._();
  @$core.override
  AutoimportDirectoryDeleteResponse createEmptyInstance() =>
      AutoimportDirectoryDeleteResponse._();
  @$core.pragma('dart2js:noInline')
  static AutoimportDirectoryDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoimportDirectoryDeleteResponse>(
          AutoimportDirectoryDeleteResponse.$_createMessage);
  static AutoimportDirectoryDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AutoimportDirectory get directory => $_getN(0);
  @$pb.TagNumber(1)
  set directory(AutoimportDirectory value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDirectory() => $_has(0);
  @$pb.TagNumber(1)
  void clearDirectory() => $_clearField(1);
  @$pb.TagNumber(1)
  AutoimportDirectory ensureDirectory() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
