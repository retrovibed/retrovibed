// This is a generated file - do not edit.
//
// Generated from media/media.known.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import '../meta/meta.search.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class Known extends $pb.GeneratedMessage {
  factory Known({
    $core.String? id,
    $core.double? rating,
    $core.bool? adult,
    $core.String? description,
    $core.String? summary,
    $core.String? image,
    $core.String? released,
    $core.String? mimetype,
    $core.String? source,
    $core.String? uid,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (rating != null) result.rating = rating;
    if (adult != null) result.adult = adult;
    if (description != null) result.description = description;
    if (summary != null) result.summary = summary;
    if (image != null) result.image = image;
    if (released != null) result.released = released;
    if (mimetype != null) result.mimetype = mimetype;
    if (source != null) result.source = source;
    if (uid != null) result.uid = uid;
    return result;
  }

  Known._();

  factory Known.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Known.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Known',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aD(2, _omitFieldNames ? '' : 'rating', fieldType: $pb.PbFieldType.OF)
    ..aOB(3, _omitFieldNames ? '' : 'adult')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aOS(5, _omitFieldNames ? '' : 'summary')
    ..aOS(6, _omitFieldNames ? '' : 'image')
    ..aOS(7, _omitFieldNames ? '' : 'released')
    ..aOS(8, _omitFieldNames ? '' : 'mimetype')
    ..aOS(9, _omitFieldNames ? '' : 'source')
    ..aOS(10, _omitFieldNames ? '' : 'uid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Known clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Known copyWith(void Function(Known) updates) =>
      super.copyWith((message) => updates(message as Known)) as Known;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Known create() => Known._();
  @$core.override
  Known createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Known getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Known>(create);
  static Known? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.double get rating => $_getN(1);
  @$pb.TagNumber(2)
  set rating($core.double value) => $_setFloat(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRating() => $_has(1);
  @$pb.TagNumber(2)
  void clearRating() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.bool get adult => $_getBF(2);
  @$pb.TagNumber(3)
  set adult($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAdult() => $_has(2);
  @$pb.TagNumber(3)
  void clearAdult() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get summary => $_getSZ(4);
  @$pb.TagNumber(5)
  set summary($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasSummary() => $_has(4);
  @$pb.TagNumber(5)
  void clearSummary() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get image => $_getSZ(5);
  @$pb.TagNumber(6)
  set image($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasImage() => $_has(5);
  @$pb.TagNumber(6)
  void clearImage() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get released => $_getSZ(6);
  @$pb.TagNumber(7)
  set released($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasReleased() => $_has(6);
  @$pb.TagNumber(7)
  void clearReleased() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get mimetype => $_getSZ(7);
  @$pb.TagNumber(8)
  set mimetype($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasMimetype() => $_has(7);
  @$pb.TagNumber(8)
  void clearMimetype() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get source => $_getSZ(8);
  @$pb.TagNumber(9)
  set source($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasSource() => $_has(8);
  @$pb.TagNumber(9)
  void clearSource() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.String get uid => $_getSZ(9);
  @$pb.TagNumber(10)
  set uid($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasUid() => $_has(9);
  @$pb.TagNumber(10)
  void clearUid() => $_clearField(10);
}

class KnownSearchRequest extends $pb.GeneratedMessage {
  factory KnownSearchRequest({
    $core.String? query,
    $core.bool? adult,
    $core.String? language,
    $core.String? mimetype,
    $0.DateRange? released,
    $core.Iterable<$core.String>? source,
    $core.Iterable<$core.String>? id,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (adult != null) result.adult = adult;
    if (language != null) result.language = language;
    if (mimetype != null) result.mimetype = mimetype;
    if (released != null) result.released = released;
    if (source != null) result.source.addAll(source);
    if (id != null) result.id.addAll(id);
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  KnownSearchRequest._();

  factory KnownSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..aOB(2, _omitFieldNames ? '' : 'adult')
    ..aOS(3, _omitFieldNames ? '' : 'language')
    ..aOS(4, _omitFieldNames ? '' : 'mimetype')
    ..aOM<$0.DateRange>(5, _omitFieldNames ? '' : 'released',
        subBuilder: $0.DateRange.create)
    ..pPS(6, _omitFieldNames ? '' : 'source')
    ..pPS(7, _omitFieldNames ? '' : 'id')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownSearchRequest copyWith(void Function(KnownSearchRequest) updates) =>
      super.copyWith((message) => updates(message as KnownSearchRequest))
          as KnownSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownSearchRequest create() => KnownSearchRequest._();
  @$core.override
  KnownSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownSearchRequest>(create);
  static KnownSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get adult => $_getBF(1);
  @$pb.TagNumber(2)
  set adult($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAdult() => $_has(1);
  @$pb.TagNumber(2)
  void clearAdult() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get language => $_getSZ(2);
  @$pb.TagNumber(3)
  set language($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLanguage() => $_has(2);
  @$pb.TagNumber(3)
  void clearLanguage() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get mimetype => $_getSZ(3);
  @$pb.TagNumber(4)
  set mimetype($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasMimetype() => $_has(3);
  @$pb.TagNumber(4)
  void clearMimetype() => $_clearField(4);

  @$pb.TagNumber(5)
  $0.DateRange get released => $_getN(4);
  @$pb.TagNumber(5)
  set released($0.DateRange value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasReleased() => $_has(4);
  @$pb.TagNumber(5)
  void clearReleased() => $_clearField(5);
  @$pb.TagNumber(5)
  $0.DateRange ensureReleased() => $_ensure(4);

  @$pb.TagNumber(6)
  $pb.PbList<$core.String> get source => $_getList(5);

  @$pb.TagNumber(7)
  $pb.PbList<$core.String> get id => $_getList(6);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(7);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(7);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(8);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(8, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(8);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class KnownSearchResponse extends $pb.GeneratedMessage {
  factory KnownSearchResponse({
    KnownSearchRequest? next,
    $core.Iterable<Known>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  KnownSearchResponse._();

  factory KnownSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<KnownSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: KnownSearchRequest.create)
    ..pPM<Known>(2, _omitFieldNames ? '' : 'items', subBuilder: Known.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownSearchResponse copyWith(void Function(KnownSearchResponse) updates) =>
      super.copyWith((message) => updates(message as KnownSearchResponse))
          as KnownSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownSearchResponse create() => KnownSearchResponse._();
  @$core.override
  KnownSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownSearchResponse>(create);
  static KnownSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  KnownSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(KnownSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  KnownSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Known> get items => $_getList(1);
}

class KnownMatchRequest extends $pb.GeneratedMessage {
  factory KnownMatchRequest({
    $core.String? query,
  }) {
    final result = create();
    if (query != null) result.query = query;
    return result;
  }

  KnownMatchRequest._();

  factory KnownMatchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownMatchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownMatchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownMatchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownMatchRequest copyWith(void Function(KnownMatchRequest) updates) =>
      super.copyWith((message) => updates(message as KnownMatchRequest))
          as KnownMatchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownMatchRequest create() => KnownMatchRequest._();
  @$core.override
  KnownMatchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownMatchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownMatchRequest>(create);
  static KnownMatchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);
}

class KnownLookupRequest extends $pb.GeneratedMessage {
  factory KnownLookupRequest() => create();

  KnownLookupRequest._();

  factory KnownLookupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownLookupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownLookupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownLookupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownLookupRequest copyWith(void Function(KnownLookupRequest) updates) =>
      super.copyWith((message) => updates(message as KnownLookupRequest))
          as KnownLookupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownLookupRequest create() => KnownLookupRequest._();
  @$core.override
  KnownLookupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownLookupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownLookupRequest>(create);
  static KnownLookupRequest? _defaultInstance;
}

class KnownLookupResponse extends $pb.GeneratedMessage {
  factory KnownLookupResponse({
    Known? known,
  }) {
    final result = create();
    if (known != null) result.known = known;
    return result;
  }

  KnownLookupResponse._();

  factory KnownLookupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownLookupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownLookupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Known>(1, _omitFieldNames ? '' : 'known', subBuilder: Known.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownLookupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownLookupResponse copyWith(void Function(KnownLookupResponse) updates) =>
      super.copyWith((message) => updates(message as KnownLookupResponse))
          as KnownLookupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownLookupResponse create() => KnownLookupResponse._();
  @$core.override
  KnownLookupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownLookupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownLookupResponse>(create);
  static KnownLookupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Known get known => $_getN(0);
  @$pb.TagNumber(1)
  set known(Known value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasKnown() => $_has(0);
  @$pb.TagNumber(1)
  void clearKnown() => $_clearField(1);
  @$pb.TagNumber(1)
  Known ensureKnown() => $_ensure(0);
}

class KnownDownloadRequest extends $pb.GeneratedMessage {
  factory KnownDownloadRequest() => create();

  KnownDownloadRequest._();

  factory KnownDownloadRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownDownloadRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownDownloadRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownDownloadRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownDownloadRequest copyWith(void Function(KnownDownloadRequest) updates) =>
      super.copyWith((message) => updates(message as KnownDownloadRequest))
          as KnownDownloadRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownDownloadRequest create() => KnownDownloadRequest._();
  @$core.override
  KnownDownloadRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownDownloadRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownDownloadRequest>(create);
  static KnownDownloadRequest? _defaultInstance;
}

class KnownDownloadResponse extends $pb.GeneratedMessage {
  factory KnownDownloadResponse({
    Known? known,
  }) {
    final result = create();
    if (known != null) result.known = known;
    return result;
  }

  KnownDownloadResponse._();

  factory KnownDownloadResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownDownloadResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownDownloadResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Known>(1, _omitFieldNames ? '' : 'known', subBuilder: Known.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownDownloadResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownDownloadResponse copyWith(
          void Function(KnownDownloadResponse) updates) =>
      super.copyWith((message) => updates(message as KnownDownloadResponse))
          as KnownDownloadResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownDownloadResponse create() => KnownDownloadResponse._();
  @$core.override
  KnownDownloadResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownDownloadResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownDownloadResponse>(create);
  static KnownDownloadResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Known get known => $_getN(0);
  @$pb.TagNumber(1)
  set known(Known value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasKnown() => $_has(0);
  @$pb.TagNumber(1)
  void clearKnown() => $_clearField(1);
  @$pb.TagNumber(1)
  Known ensureKnown() => $_ensure(0);
}

class KnownCreateRequest extends $pb.GeneratedMessage {
  factory KnownCreateRequest({
    Known? known,
  }) {
    final result = create();
    if (known != null) result.known = known;
    return result;
  }

  KnownCreateRequest._();

  factory KnownCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Known>(1, _omitFieldNames ? '' : 'known', subBuilder: Known.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownCreateRequest copyWith(void Function(KnownCreateRequest) updates) =>
      super.copyWith((message) => updates(message as KnownCreateRequest))
          as KnownCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownCreateRequest create() => KnownCreateRequest._();
  @$core.override
  KnownCreateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownCreateRequest>(create);
  static KnownCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Known get known => $_getN(0);
  @$pb.TagNumber(1)
  set known(Known value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasKnown() => $_has(0);
  @$pb.TagNumber(1)
  void clearKnown() => $_clearField(1);
  @$pb.TagNumber(1)
  Known ensureKnown() => $_ensure(0);
}

class KnownCreateResponse extends $pb.GeneratedMessage {
  factory KnownCreateResponse({
    Known? known,
  }) {
    final result = create();
    if (known != null) result.known = known;
    return result;
  }

  KnownCreateResponse._();

  factory KnownCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Known>(1, _omitFieldNames ? '' : 'known', subBuilder: Known.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownCreateResponse copyWith(void Function(KnownCreateResponse) updates) =>
      super.copyWith((message) => updates(message as KnownCreateResponse))
          as KnownCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownCreateResponse create() => KnownCreateResponse._();
  @$core.override
  KnownCreateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownCreateResponse>(create);
  static KnownCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Known get known => $_getN(0);
  @$pb.TagNumber(1)
  set known(Known value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasKnown() => $_has(0);
  @$pb.TagNumber(1)
  void clearKnown() => $_clearField(1);
  @$pb.TagNumber(1)
  Known ensureKnown() => $_ensure(0);
}

class KnownLatestRequest extends $pb.GeneratedMessage {
  factory KnownLatestRequest({
    $0.DateRange? released,
    $core.bool? adult,
    $core.String? language,
    $core.String? mimetype,
    $core.Iterable<$core.String>? source,
    $core.Iterable<$core.String>? id,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (released != null) result.released = released;
    if (adult != null) result.adult = adult;
    if (language != null) result.language = language;
    if (mimetype != null) result.mimetype = mimetype;
    if (source != null) result.source.addAll(source);
    if (id != null) result.id.addAll(id);
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  KnownLatestRequest._();

  factory KnownLatestRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownLatestRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownLatestRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<$0.DateRange>(1, _omitFieldNames ? '' : 'released',
        subBuilder: $0.DateRange.create)
    ..aOB(2, _omitFieldNames ? '' : 'adult')
    ..aOS(3, _omitFieldNames ? '' : 'language')
    ..aOS(4, _omitFieldNames ? '' : 'mimetype')
    ..pPS(5, _omitFieldNames ? '' : 'source')
    ..pPS(6, _omitFieldNames ? '' : 'id')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownLatestRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownLatestRequest copyWith(void Function(KnownLatestRequest) updates) =>
      super.copyWith((message) => updates(message as KnownLatestRequest))
          as KnownLatestRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownLatestRequest create() => KnownLatestRequest._();
  @$core.override
  KnownLatestRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownLatestRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownLatestRequest>(create);
  static KnownLatestRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $0.DateRange get released => $_getN(0);
  @$pb.TagNumber(1)
  set released($0.DateRange value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasReleased() => $_has(0);
  @$pb.TagNumber(1)
  void clearReleased() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.DateRange ensureReleased() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.bool get adult => $_getBF(1);
  @$pb.TagNumber(2)
  set adult($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAdult() => $_has(1);
  @$pb.TagNumber(2)
  void clearAdult() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get language => $_getSZ(2);
  @$pb.TagNumber(3)
  set language($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLanguage() => $_has(2);
  @$pb.TagNumber(3)
  void clearLanguage() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get mimetype => $_getSZ(3);
  @$pb.TagNumber(4)
  set mimetype($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasMimetype() => $_has(3);
  @$pb.TagNumber(4)
  void clearMimetype() => $_clearField(4);

  @$pb.TagNumber(5)
  $pb.PbList<$core.String> get source => $_getList(4);

  @$pb.TagNumber(6)
  $pb.PbList<$core.String> get id => $_getList(5);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(6);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(6);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(7);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(7);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class KnownLatestResponse extends $pb.GeneratedMessage {
  factory KnownLatestResponse({
    KnownLatestRequest? next,
    $core.Iterable<Known>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  KnownLatestResponse._();

  factory KnownLatestResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KnownLatestResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KnownLatestResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<KnownLatestRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: KnownLatestRequest.create)
    ..pPM<Known>(2, _omitFieldNames ? '' : 'items', subBuilder: Known.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownLatestResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KnownLatestResponse copyWith(void Function(KnownLatestResponse) updates) =>
      super.copyWith((message) => updates(message as KnownLatestResponse))
          as KnownLatestResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KnownLatestResponse create() => KnownLatestResponse._();
  @$core.override
  KnownLatestResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KnownLatestResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KnownLatestResponse>(create);
  static KnownLatestResponse? _defaultInstance;

  @$pb.TagNumber(1)
  KnownLatestRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(KnownLatestRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  KnownLatestRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Known> get items => $_getList(1);
}

class RecommendationSearchRequest extends $pb.GeneratedMessage {
  factory RecommendationSearchRequest({
    $core.String? mimetype,
    $core.bool? adult,
    $core.String? language,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (mimetype != null) result.mimetype = mimetype;
    if (adult != null) result.adult = adult;
    if (language != null) result.language = language;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  RecommendationSearchRequest._();

  factory RecommendationSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecommendationSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecommendationSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'mimetype')
    ..aOB(2, _omitFieldNames ? '' : 'adult')
    ..aOS(3, _omitFieldNames ? '' : 'language')
    ..a<$fixnum.Int64>(
        900, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(901, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationSearchRequest copyWith(
          void Function(RecommendationSearchRequest) updates) =>
      super.copyWith(
              (message) => updates(message as RecommendationSearchRequest))
          as RecommendationSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecommendationSearchRequest create() =>
      RecommendationSearchRequest._();
  @$core.override
  RecommendationSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecommendationSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecommendationSearchRequest>(create);
  static RecommendationSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get mimetype => $_getSZ(0);
  @$pb.TagNumber(1)
  set mimetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasMimetype() => $_has(0);
  @$pb.TagNumber(1)
  void clearMimetype() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get adult => $_getBF(1);
  @$pb.TagNumber(2)
  set adult($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAdult() => $_has(1);
  @$pb.TagNumber(2)
  void clearAdult() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get language => $_getSZ(2);
  @$pb.TagNumber(3)
  set language($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLanguage() => $_has(2);
  @$pb.TagNumber(3)
  void clearLanguage() => $_clearField(3);

  @$pb.TagNumber(900)
  $fixnum.Int64 get offset => $_getI64(3);
  @$pb.TagNumber(900)
  set offset($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(900)
  $core.bool hasOffset() => $_has(3);
  @$pb.TagNumber(900)
  void clearOffset() => $_clearField(900);

  @$pb.TagNumber(901)
  $fixnum.Int64 get limit => $_getI64(4);
  @$pb.TagNumber(901)
  set limit($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(901)
  $core.bool hasLimit() => $_has(4);
  @$pb.TagNumber(901)
  void clearLimit() => $_clearField(901);
}

class RecommendationSearchResponse extends $pb.GeneratedMessage {
  factory RecommendationSearchResponse({
    RecommendationSearchRequest? next,
    $core.Iterable<Known>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  RecommendationSearchResponse._();

  factory RecommendationSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecommendationSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecommendationSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<RecommendationSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: RecommendationSearchRequest.create)
    ..pPM<Known>(2, _omitFieldNames ? '' : 'items', subBuilder: Known.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationSearchResponse copyWith(
          void Function(RecommendationSearchResponse) updates) =>
      super.copyWith(
              (message) => updates(message as RecommendationSearchResponse))
          as RecommendationSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecommendationSearchResponse create() =>
      RecommendationSearchResponse._();
  @$core.override
  RecommendationSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecommendationSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecommendationSearchResponse>(create);
  static RecommendationSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RecommendationSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(RecommendationSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  RecommendationSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Known> get items => $_getList(1);
}

class RecommendationFindRequest extends $pb.GeneratedMessage {
  factory RecommendationFindRequest() => create();

  RecommendationFindRequest._();

  factory RecommendationFindRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecommendationFindRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecommendationFindRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationFindRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationFindRequest copyWith(
          void Function(RecommendationFindRequest) updates) =>
      super.copyWith((message) => updates(message as RecommendationFindRequest))
          as RecommendationFindRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecommendationFindRequest create() => RecommendationFindRequest._();
  @$core.override
  RecommendationFindRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecommendationFindRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecommendationFindRequest>(create);
  static RecommendationFindRequest? _defaultInstance;
}

class RecommendationFindResponse extends $pb.GeneratedMessage {
  factory RecommendationFindResponse({
    Known? recommendation,
  }) {
    final result = create();
    if (recommendation != null) result.recommendation = recommendation;
    return result;
  }

  RecommendationFindResponse._();

  factory RecommendationFindResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecommendationFindResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecommendationFindResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Known>(1, _omitFieldNames ? '' : 'recomendation',
        protoName: 'recommendation', subBuilder: Known.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationFindResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationFindResponse copyWith(
          void Function(RecommendationFindResponse) updates) =>
      super.copyWith(
              (message) => updates(message as RecommendationFindResponse))
          as RecommendationFindResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecommendationFindResponse create() => RecommendationFindResponse._();
  @$core.override
  RecommendationFindResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecommendationFindResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecommendationFindResponse>(create);
  static RecommendationFindResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Known get recommendation => $_getN(0);
  @$pb.TagNumber(1)
  set recommendation(Known value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasRecommendation() => $_has(0);
  @$pb.TagNumber(1)
  void clearRecommendation() => $_clearField(1);
  @$pb.TagNumber(1)
  Known ensureRecommendation() => $_ensure(0);
}

class RecommendationDeleteRequest extends $pb.GeneratedMessage {
  factory RecommendationDeleteRequest() => create();

  RecommendationDeleteRequest._();

  factory RecommendationDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecommendationDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecommendationDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationDeleteRequest copyWith(
          void Function(RecommendationDeleteRequest) updates) =>
      super.copyWith(
              (message) => updates(message as RecommendationDeleteRequest))
          as RecommendationDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecommendationDeleteRequest create() =>
      RecommendationDeleteRequest._();
  @$core.override
  RecommendationDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecommendationDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecommendationDeleteRequest>(create);
  static RecommendationDeleteRequest? _defaultInstance;
}

class RecommendationDeleteResponse extends $pb.GeneratedMessage {
  factory RecommendationDeleteResponse({
    Known? recommendation,
  }) {
    final result = create();
    if (recommendation != null) result.recommendation = recommendation;
    return result;
  }

  RecommendationDeleteResponse._();

  factory RecommendationDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecommendationDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecommendationDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOM<Known>(1, _omitFieldNames ? '' : 'recomendation',
        protoName: 'recommendation', subBuilder: Known.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationDeleteResponse copyWith(
          void Function(RecommendationDeleteResponse) updates) =>
      super.copyWith(
              (message) => updates(message as RecommendationDeleteResponse))
          as RecommendationDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecommendationDeleteResponse create() =>
      RecommendationDeleteResponse._();
  @$core.override
  RecommendationDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecommendationDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecommendationDeleteResponse>(create);
  static RecommendationDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Known get recommendation => $_getN(0);
  @$pb.TagNumber(1)
  set recommendation(Known value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasRecommendation() => $_has(0);
  @$pb.TagNumber(1)
  void clearRecommendation() => $_clearField(1);
  @$pb.TagNumber(1)
  Known ensureRecommendation() => $_ensure(0);
}

class RecommendationRefreshRequest extends $pb.GeneratedMessage {
  factory RecommendationRefreshRequest({
    $core.String? profileId,
    $fixnum.Int64? limit,
    $core.String? mimetype,
    $core.bool? adult,
    $core.String? language,
  }) {
    final result = create();
    if (profileId != null) result.profileId = profileId;
    if (limit != null) result.limit = limit;
    if (mimetype != null) result.mimetype = mimetype;
    if (adult != null) result.adult = adult;
    if (language != null) result.language = language;
    return result;
  }

  RecommendationRefreshRequest._();

  factory RecommendationRefreshRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecommendationRefreshRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecommendationRefreshRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'profile_id')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(3, _omitFieldNames ? '' : 'mimetype')
    ..aOB(4, _omitFieldNames ? '' : 'adult')
    ..aOS(5, _omitFieldNames ? '' : 'language')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationRefreshRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationRefreshRequest copyWith(
          void Function(RecommendationRefreshRequest) updates) =>
      super.copyWith(
              (message) => updates(message as RecommendationRefreshRequest))
          as RecommendationRefreshRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecommendationRefreshRequest create() =>
      RecommendationRefreshRequest._();
  @$core.override
  RecommendationRefreshRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecommendationRefreshRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecommendationRefreshRequest>(create);
  static RecommendationRefreshRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get profileId => $_getSZ(0);
  @$pb.TagNumber(1)
  set profileId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasProfileId() => $_has(0);
  @$pb.TagNumber(1)
  void clearProfileId() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get limit => $_getI64(1);
  @$pb.TagNumber(2)
  set limit($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(2)
  void clearLimit() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get mimetype => $_getSZ(2);
  @$pb.TagNumber(3)
  set mimetype($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasMimetype() => $_has(2);
  @$pb.TagNumber(3)
  void clearMimetype() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get adult => $_getBF(3);
  @$pb.TagNumber(4)
  set adult($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasAdult() => $_has(3);
  @$pb.TagNumber(4)
  void clearAdult() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get language => $_getSZ(4);
  @$pb.TagNumber(5)
  set language($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasLanguage() => $_has(4);
  @$pb.TagNumber(5)
  void clearLanguage() => $_clearField(5);
}

class RecommendationRefreshResponse extends $pb.GeneratedMessage {
  factory RecommendationRefreshResponse() => create();

  RecommendationRefreshResponse._();

  factory RecommendationRefreshResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecommendationRefreshResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecommendationRefreshResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'media'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationRefreshResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecommendationRefreshResponse copyWith(
          void Function(RecommendationRefreshResponse) updates) =>
      super.copyWith(
              (message) => updates(message as RecommendationRefreshResponse))
          as RecommendationRefreshResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecommendationRefreshResponse create() =>
      RecommendationRefreshResponse._();
  @$core.override
  RecommendationRefreshResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecommendationRefreshResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecommendationRefreshResponse>(create);
  static RecommendationRefreshResponse? _defaultInstance;
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
