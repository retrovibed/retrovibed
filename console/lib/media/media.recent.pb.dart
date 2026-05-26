// This is a generated file - do not edit.
//
// Generated from media.recent.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'media.pb.dart' as $1;
import 'meta.search.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class RecentSearchRequest extends $pb.GeneratedMessage {
  factory RecentSearchRequest({
    $0.DateRange? created,
    $core.String? mimetype,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (created != null) result.created = created;
    if (mimetype != null) result.mimetype = mimetype;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  RecentSearchRequest._();

  factory RecentSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecentSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecentSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<$0.DateRange>(1, _omitFieldNames ? '' : 'created',
        subBuilder: $0.DateRange.create)
    ..aOS(2, _omitFieldNames ? '' : 'mimetype')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentSearchRequest copyWith(void Function(RecentSearchRequest) updates) =>
      super.copyWith((message) => updates(message as RecentSearchRequest))
          as RecentSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecentSearchRequest create() => RecentSearchRequest._();
  @$core.override
  RecentSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecentSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecentSearchRequest>(create);
  static RecentSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $0.DateRange get created => $_getN(0);
  @$pb.TagNumber(1)
  set created($0.DateRange value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCreated() => $_has(0);
  @$pb.TagNumber(1)
  void clearCreated() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.DateRange ensureCreated() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.String get mimetype => $_getSZ(1);
  @$pb.TagNumber(2)
  set mimetype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMimetype() => $_has(1);
  @$pb.TagNumber(2)
  void clearMimetype() => $_clearField(2);

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

class RecentSearchResponse extends $pb.GeneratedMessage {
  factory RecentSearchResponse({
    RecentSearchRequest? next,
    $core.Iterable<RecentRecordRequest>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  RecentSearchResponse._();

  factory RecentSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecentSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecentSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<RecentSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: RecentSearchRequest.create)
    ..pPM<RecentRecordRequest>(2, _omitFieldNames ? '' : 'items',
        subBuilder: RecentRecordRequest.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentSearchResponse copyWith(void Function(RecentSearchResponse) updates) =>
      super.copyWith((message) => updates(message as RecentSearchResponse))
          as RecentSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecentSearchResponse create() => RecentSearchResponse._();
  @$core.override
  RecentSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecentSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecentSearchResponse>(create);
  static RecentSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RecentSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(RecentSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  RecentSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<RecentRecordRequest> get items => $_getList(1);
}

class RecentRecordRequest extends $pb.GeneratedMessage {
  factory RecentRecordRequest({
    $core.String? id,
    $1.Media? media,
    $fixnum.Int64? duration,
    $fixnum.Int64? position,
    $1.MediaSearchRequest? query,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (media != null) result.media = media;
    if (duration != null) result.duration = duration;
    if (position != null) result.position = position;
    if (query != null) result.query = query;
    return result;
  }

  RecentRecordRequest._();

  factory RecentRecordRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecentRecordRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecentRecordRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOM<$1.Media>(2, _omitFieldNames ? '' : 'media',
        subBuilder: $1.Media.create)
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'duration', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        4, _omitFieldNames ? '' : 'position', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOM<$1.MediaSearchRequest>(5, _omitFieldNames ? '' : 'query',
        subBuilder: $1.MediaSearchRequest.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentRecordRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentRecordRequest copyWith(void Function(RecentRecordRequest) updates) =>
      super.copyWith((message) => updates(message as RecentRecordRequest))
          as RecentRecordRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecentRecordRequest create() => RecentRecordRequest._();
  @$core.override
  RecentRecordRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecentRecordRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecentRecordRequest>(create);
  static RecentRecordRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $1.Media get media => $_getN(1);
  @$pb.TagNumber(2)
  set media($1.Media value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasMedia() => $_has(1);
  @$pb.TagNumber(2)
  void clearMedia() => $_clearField(2);
  @$pb.TagNumber(2)
  $1.Media ensureMedia() => $_ensure(1);

  @$pb.TagNumber(3)
  $fixnum.Int64 get duration => $_getI64(2);
  @$pb.TagNumber(3)
  set duration($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasDuration() => $_has(2);
  @$pb.TagNumber(3)
  void clearDuration() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get position => $_getI64(3);
  @$pb.TagNumber(4)
  set position($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPosition() => $_has(3);
  @$pb.TagNumber(4)
  void clearPosition() => $_clearField(4);

  @$pb.TagNumber(5)
  $1.MediaSearchRequest get query => $_getN(4);
  @$pb.TagNumber(5)
  set query($1.MediaSearchRequest value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasQuery() => $_has(4);
  @$pb.TagNumber(5)
  void clearQuery() => $_clearField(5);
  @$pb.TagNumber(5)
  $1.MediaSearchRequest ensureQuery() => $_ensure(4);
}

class RecentRecordResponse extends $pb.GeneratedMessage {
  factory RecentRecordResponse() => create();

  RecentRecordResponse._();

  factory RecentRecordResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecentRecordResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecentRecordResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentRecordResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentRecordResponse copyWith(void Function(RecentRecordResponse) updates) =>
      super.copyWith((message) => updates(message as RecentRecordResponse))
          as RecentRecordResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecentRecordResponse create() => RecentRecordResponse._();
  @$core.override
  RecentRecordResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecentRecordResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecentRecordResponse>(create);
  static RecentRecordResponse? _defaultInstance;
}

class RecentDeleteRequest extends $pb.GeneratedMessage {
  factory RecentDeleteRequest() => create();

  RecentDeleteRequest._();

  factory RecentDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecentDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecentDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentDeleteRequest copyWith(void Function(RecentDeleteRequest) updates) =>
      super.copyWith((message) => updates(message as RecentDeleteRequest))
          as RecentDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecentDeleteRequest create() => RecentDeleteRequest._();
  @$core.override
  RecentDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecentDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecentDeleteRequest>(create);
  static RecentDeleteRequest? _defaultInstance;
}

class RecentDeleteResponse extends $pb.GeneratedMessage {
  factory RecentDeleteResponse() => create();

  RecentDeleteResponse._();

  factory RecentDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecentDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecentDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecentDeleteResponse copyWith(void Function(RecentDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as RecentDeleteResponse))
          as RecentDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecentDeleteResponse create() => RecentDeleteResponse._();
  @$core.override
  RecentDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecentDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecentDeleteResponse>(create);
  static RecentDeleteResponse? _defaultInstance;
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
