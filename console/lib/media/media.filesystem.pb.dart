// This is a generated file - do not edit.
//
// Generated from media/media.filesystem.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'media.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

/// the filesystem is its own surface. it shares the Media message with the library because
/// a folder and a file are both rows in library_metadata, but nothing else: reusing
/// MediaSearchRequest meant one message whose meaning changed with the presence of a field,
/// and one handler that had to serve two views.
class FilesystemSearchRequest extends $pb.GeneratedMessage {
  factory FilesystemSearchRequest({
    $core.String? query,
    $core.Iterable<$core.String>? mimetypes,
    $core.bool? hidden,
    $core.String? directoryId,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (mimetypes != null) result.mimetypes.addAll(mimetypes);
    if (hidden != null) result.hidden = hidden;
    if (directoryId != null) result.directoryId = directoryId;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  FilesystemSearchRequest._();

  factory FilesystemSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilesystemSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilesystemSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..pPS(2, _omitFieldNames ? '' : 'mimetypes')
    ..aOB(3, _omitFieldNames ? '' : 'hidden')
    ..aOS(4, _omitFieldNames ? '' : 'directory_id')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemSearchRequest copyWith(
          void Function(FilesystemSearchRequest) updates) =>
      super.copyWith((message) => updates(message as FilesystemSearchRequest))
          as FilesystemSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilesystemSearchRequest create() => FilesystemSearchRequest._();
  @$core.override
  FilesystemSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilesystemSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilesystemSearchRequest>(create);
  static FilesystemSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<$core.String> get mimetypes => $_getList(1);

  @$pb.TagNumber(3)
  $core.bool get hidden => $_getBF(2);
  @$pb.TagNumber(3)
  set hidden($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasHidden() => $_has(2);
  @$pb.TagNumber(3)
  void clearHidden() => $_clearField(3);

  /// the directory being listed. the zero uuid is the root of the library.
  @$pb.TagNumber(4)
  $core.String get directoryId => $_getSZ(3);
  @$pb.TagNumber(4)
  set directoryId($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDirectoryId() => $_has(3);
  @$pb.TagNumber(4)
  void clearDirectoryId() => $_clearField(4);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(4);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(4);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(5);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(5);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class FilesystemSearchResponse extends $pb.GeneratedMessage {
  factory FilesystemSearchResponse({
    FilesystemSearchRequest? next,
    $core.Iterable<$0.Media>? items,
    $core.Iterable<$0.Media>? breadcrumb,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    if (breadcrumb != null) result.breadcrumb.addAll(breadcrumb);
    return result;
  }

  FilesystemSearchResponse._();

  factory FilesystemSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilesystemSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilesystemSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<FilesystemSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: FilesystemSearchRequest.create)
    ..pPM<$0.Media>(2, _omitFieldNames ? '' : 'items',
        subBuilder: $0.Media.create)
    ..pPM<$0.Media>(3, _omitFieldNames ? '' : 'breadcrumb',
        subBuilder: $0.Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemSearchResponse copyWith(
          void Function(FilesystemSearchResponse) updates) =>
      super.copyWith((message) => updates(message as FilesystemSearchResponse))
          as FilesystemSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilesystemSearchResponse create() => FilesystemSearchResponse._();
  @$core.override
  FilesystemSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilesystemSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilesystemSearchResponse>(create);
  static FilesystemSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  FilesystemSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(FilesystemSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  FilesystemSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<$0.Media> get items => $_getList(1);

  /// the listed directory and its ancestors, root first. rides along with the listing so
  /// the path costs no second round trip.
  @$pb.TagNumber(3)
  $pb.PbList<$0.Media> get breadcrumb => $_getList(2);
}

class FilesystemCreateRequest extends $pb.GeneratedMessage {
  factory FilesystemCreateRequest({
    $core.String? directoryId,
    $core.String? name,
  }) {
    final result = create();
    if (directoryId != null) result.directoryId = directoryId;
    if (name != null) result.name = name;
    return result;
  }

  FilesystemCreateRequest._();

  factory FilesystemCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilesystemCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilesystemCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'directory_id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemCreateRequest copyWith(
          void Function(FilesystemCreateRequest) updates) =>
      super.copyWith((message) => updates(message as FilesystemCreateRequest))
          as FilesystemCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilesystemCreateRequest create() => FilesystemCreateRequest._();
  @$core.override
  FilesystemCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilesystemCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilesystemCreateRequest>(create);
  static FilesystemCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get directoryId => $_getSZ(0);
  @$pb.TagNumber(1)
  set directoryId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDirectoryId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDirectoryId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);
}

class FilesystemCreateResponse extends $pb.GeneratedMessage {
  factory FilesystemCreateResponse({
    $0.Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  FilesystemCreateResponse._();

  factory FilesystemCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilesystemCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilesystemCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<$0.Media>(1, _omitFieldNames ? '' : 'media',
        subBuilder: $0.Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemCreateResponse copyWith(
          void Function(FilesystemCreateResponse) updates) =>
      super.copyWith((message) => updates(message as FilesystemCreateResponse))
          as FilesystemCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilesystemCreateResponse create() => FilesystemCreateResponse._();
  @$core.override
  FilesystemCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilesystemCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilesystemCreateResponse>(create);
  static FilesystemCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $0.Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media($0.Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.Media ensureMedia() => $_ensure(0);
}

class FilesystemMoveRequest extends $pb.GeneratedMessage {
  factory FilesystemMoveRequest({
    $core.String? directoryId,
  }) {
    final result = create();
    if (directoryId != null) result.directoryId = directoryId;
    return result;
  }

  FilesystemMoveRequest._();

  factory FilesystemMoveRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilesystemMoveRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilesystemMoveRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'directory_id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemMoveRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemMoveRequest copyWith(
          void Function(FilesystemMoveRequest) updates) =>
      super.copyWith((message) => updates(message as FilesystemMoveRequest))
          as FilesystemMoveRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilesystemMoveRequest create() => FilesystemMoveRequest._();
  @$core.override
  FilesystemMoveRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilesystemMoveRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilesystemMoveRequest>(create);
  static FilesystemMoveRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get directoryId => $_getSZ(0);
  @$pb.TagNumber(1)
  set directoryId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDirectoryId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDirectoryId() => $_clearField(1);
}

class FilesystemMoveResponse extends $pb.GeneratedMessage {
  factory FilesystemMoveResponse({
    $0.Media? media,
  }) {
    final result = create();
    if (media != null) result.media = media;
    return result;
  }

  FilesystemMoveResponse._();

  factory FilesystemMoveResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilesystemMoveResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilesystemMoveResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<$0.Media>(1, _omitFieldNames ? '' : 'media',
        subBuilder: $0.Media.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemMoveResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemMoveResponse copyWith(
          void Function(FilesystemMoveResponse) updates) =>
      super.copyWith((message) => updates(message as FilesystemMoveResponse))
          as FilesystemMoveResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilesystemMoveResponse create() => FilesystemMoveResponse._();
  @$core.override
  FilesystemMoveResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilesystemMoveResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilesystemMoveResponse>(create);
  static FilesystemMoveResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $0.Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media($0.Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.Media ensureMedia() => $_ensure(0);
}

class FilesystemDeleteResponse extends $pb.GeneratedMessage {
  factory FilesystemDeleteResponse({
    $0.Media? media,
    $fixnum.Int64? removed,
  }) {
    final result = create();
    if (media != null) result.media = media;
    if (removed != null) result.removed = removed;
    return result;
  }

  FilesystemDeleteResponse._();

  factory FilesystemDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilesystemDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilesystemDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<$0.Media>(1, _omitFieldNames ? '' : 'media',
        subBuilder: $0.Media.create)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'removed', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilesystemDeleteResponse copyWith(
          void Function(FilesystemDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as FilesystemDeleteResponse))
          as FilesystemDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilesystemDeleteResponse create() => FilesystemDeleteResponse._();
  @$core.override
  FilesystemDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilesystemDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilesystemDeleteResponse>(create);
  static FilesystemDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $0.Media get media => $_getN(0);
  @$pb.TagNumber(1)
  set media($0.Media value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMedia() => $_has(0);
  @$pb.TagNumber(1)
  void clearMedia() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.Media ensureMedia() => $_ensure(0);

  /// everything below the deleted directory went with it.
  @$pb.TagNumber(2)
  $fixnum.Int64 get removed => $_getI64(1);
  @$pb.TagNumber(2)
  set removed($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRemoved() => $_has(1);
  @$pb.TagNumber(2)
  void clearRemoved() => $_clearField(2);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
