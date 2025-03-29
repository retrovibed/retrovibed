//
//  Generated code. Do not modify.
//  source: meta.daemon.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

class Daemon extends $pb.GeneratedMessage {
  factory Daemon({
    $core.String? id,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? description,
    $core.String? hostname,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (createdAt != null) {
      $result.createdAt = createdAt;
    }
    if (updatedAt != null) {
      $result.updatedAt = updatedAt;
    }
    if (description != null) {
      $result.description = description;
    }
    if (hostname != null) {
      $result.hostname = hostname;
    }
    return $result;
  }
  Daemon._() : super();
  factory Daemon.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Daemon.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Daemon', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'created_at')
    ..aOS(3, _omitFieldNames ? '' : 'updated_at')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aOS(5, _omitFieldNames ? '' : 'hostname')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Daemon clone() => Daemon()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Daemon copyWith(void Function(Daemon) updates) => super.copyWith((message) => updates(message as Daemon)) as Daemon;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Daemon create() => Daemon._();
  Daemon createEmptyInstance() => create();
  static $pb.PbList<Daemon> createRepeated() => $pb.PbList<Daemon>();
  @$core.pragma('dart2js:noInline')
  static Daemon getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Daemon>(create);
  static Daemon? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get createdAt => $_getSZ(1);
  @$pb.TagNumber(2)
  set createdAt($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasCreatedAt() => $_has(1);
  @$pb.TagNumber(2)
  void clearCreatedAt() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get updatedAt => $_getSZ(2);
  @$pb.TagNumber(3)
  set updatedAt($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasUpdatedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearUpdatedAt() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get hostname => $_getSZ(4);
  @$pb.TagNumber(5)
  set hostname($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasHostname() => $_has(4);
  @$pb.TagNumber(5)
  void clearHostname() => clearField(5);
}

class DaemonSearchRequest extends $pb.GeneratedMessage {
  factory DaemonSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final $result = create();
    if (query != null) {
      $result.query = query;
    }
    if (offset != null) {
      $result.offset = offset;
    }
    if (limit != null) {
      $result.limit = limit;
    }
    return $result;
  }
  DaemonSearchRequest._() : super();
  factory DaemonSearchRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonSearchRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonSearchRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonSearchRequest clone() => DaemonSearchRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonSearchRequest copyWith(void Function(DaemonSearchRequest) updates) => super.copyWith((message) => updates(message as DaemonSearchRequest)) as DaemonSearchRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonSearchRequest create() => DaemonSearchRequest._();
  DaemonSearchRequest createEmptyInstance() => create();
  static $pb.PbList<DaemonSearchRequest> createRepeated() => $pb.PbList<DaemonSearchRequest>();
  @$core.pragma('dart2js:noInline')
  static DaemonSearchRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonSearchRequest>(create);
  static DaemonSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get offset => $_getI64(1);
  @$pb.TagNumber(2)
  set offset($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(2)
  void clearOffset() => clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get limit => $_getI64(2);
  @$pb.TagNumber(3)
  set limit($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => clearField(3);
}

class DaemonSearchResponse extends $pb.GeneratedMessage {
  factory DaemonSearchResponse({
    DaemonSearchRequest? next,
    $core.Iterable<Daemon>? items,
  }) {
    final $result = create();
    if (next != null) {
      $result.next = next;
    }
    if (items != null) {
      $result.items.addAll(items);
    }
    return $result;
  }
  DaemonSearchResponse._() : super();
  factory DaemonSearchResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonSearchResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonSearchResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<DaemonSearchRequest>(1, _omitFieldNames ? '' : 'next', subBuilder: DaemonSearchRequest.create)
    ..pc<Daemon>(2, _omitFieldNames ? '' : 'items', $pb.PbFieldType.PM, subBuilder: Daemon.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonSearchResponse clone() => DaemonSearchResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonSearchResponse copyWith(void Function(DaemonSearchResponse) updates) => super.copyWith((message) => updates(message as DaemonSearchResponse)) as DaemonSearchResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonSearchResponse create() => DaemonSearchResponse._();
  DaemonSearchResponse createEmptyInstance() => create();
  static $pb.PbList<DaemonSearchResponse> createRepeated() => $pb.PbList<DaemonSearchResponse>();
  @$core.pragma('dart2js:noInline')
  static DaemonSearchResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonSearchResponse>(create);
  static DaemonSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  DaemonSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(DaemonSearchRequest v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => clearField(1);
  @$pb.TagNumber(1)
  DaemonSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.List<Daemon> get items => $_getList(1);
}

class DaemonCreateRequest extends $pb.GeneratedMessage {
  factory DaemonCreateRequest({
    Daemon? daemon,
  }) {
    final $result = create();
    if (daemon != null) {
      $result.daemon = daemon;
    }
    return $result;
  }
  DaemonCreateRequest._() : super();
  factory DaemonCreateRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonCreateRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonCreateRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon', subBuilder: Daemon.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonCreateRequest clone() => DaemonCreateRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonCreateRequest copyWith(void Function(DaemonCreateRequest) updates) => super.copyWith((message) => updates(message as DaemonCreateRequest)) as DaemonCreateRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonCreateRequest create() => DaemonCreateRequest._();
  DaemonCreateRequest createEmptyInstance() => create();
  static $pb.PbList<DaemonCreateRequest> createRepeated() => $pb.PbList<DaemonCreateRequest>();
  @$core.pragma('dart2js:noInline')
  static DaemonCreateRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonCreateRequest>(create);
  static DaemonCreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonCreateResponse extends $pb.GeneratedMessage {
  factory DaemonCreateResponse({
    Daemon? daemon,
  }) {
    final $result = create();
    if (daemon != null) {
      $result.daemon = daemon;
    }
    return $result;
  }
  DaemonCreateResponse._() : super();
  factory DaemonCreateResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonCreateResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonCreateResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon', subBuilder: Daemon.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonCreateResponse clone() => DaemonCreateResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonCreateResponse copyWith(void Function(DaemonCreateResponse) updates) => super.copyWith((message) => updates(message as DaemonCreateResponse)) as DaemonCreateResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonCreateResponse create() => DaemonCreateResponse._();
  DaemonCreateResponse createEmptyInstance() => create();
  static $pb.PbList<DaemonCreateResponse> createRepeated() => $pb.PbList<DaemonCreateResponse>();
  @$core.pragma('dart2js:noInline')
  static DaemonCreateResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonCreateResponse>(create);
  static DaemonCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonLookupRequest extends $pb.GeneratedMessage {
  factory DaemonLookupRequest() => create();
  DaemonLookupRequest._() : super();
  factory DaemonLookupRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonLookupRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonLookupRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonLookupRequest clone() => DaemonLookupRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonLookupRequest copyWith(void Function(DaemonLookupRequest) updates) => super.copyWith((message) => updates(message as DaemonLookupRequest)) as DaemonLookupRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonLookupRequest create() => DaemonLookupRequest._();
  DaemonLookupRequest createEmptyInstance() => create();
  static $pb.PbList<DaemonLookupRequest> createRepeated() => $pb.PbList<DaemonLookupRequest>();
  @$core.pragma('dart2js:noInline')
  static DaemonLookupRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonLookupRequest>(create);
  static DaemonLookupRequest? _defaultInstance;
}

class DaemonLookupResponse extends $pb.GeneratedMessage {
  factory DaemonLookupResponse({
    Daemon? daemon,
  }) {
    final $result = create();
    if (daemon != null) {
      $result.daemon = daemon;
    }
    return $result;
  }
  DaemonLookupResponse._() : super();
  factory DaemonLookupResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonLookupResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonLookupResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon', subBuilder: Daemon.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonLookupResponse clone() => DaemonLookupResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonLookupResponse copyWith(void Function(DaemonLookupResponse) updates) => super.copyWith((message) => updates(message as DaemonLookupResponse)) as DaemonLookupResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonLookupResponse create() => DaemonLookupResponse._();
  DaemonLookupResponse createEmptyInstance() => create();
  static $pb.PbList<DaemonLookupResponse> createRepeated() => $pb.PbList<DaemonLookupResponse>();
  @$core.pragma('dart2js:noInline')
  static DaemonLookupResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonLookupResponse>(create);
  static DaemonLookupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonUpdateRequest extends $pb.GeneratedMessage {
  factory DaemonUpdateRequest({
    Daemon? daemon,
  }) {
    final $result = create();
    if (daemon != null) {
      $result.daemon = daemon;
    }
    return $result;
  }
  DaemonUpdateRequest._() : super();
  factory DaemonUpdateRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonUpdateRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonUpdateRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon', subBuilder: Daemon.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonUpdateRequest clone() => DaemonUpdateRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonUpdateRequest copyWith(void Function(DaemonUpdateRequest) updates) => super.copyWith((message) => updates(message as DaemonUpdateRequest)) as DaemonUpdateRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonUpdateRequest create() => DaemonUpdateRequest._();
  DaemonUpdateRequest createEmptyInstance() => create();
  static $pb.PbList<DaemonUpdateRequest> createRepeated() => $pb.PbList<DaemonUpdateRequest>();
  @$core.pragma('dart2js:noInline')
  static DaemonUpdateRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonUpdateRequest>(create);
  static DaemonUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonUpdateResponse extends $pb.GeneratedMessage {
  factory DaemonUpdateResponse({
    Daemon? daemon,
  }) {
    final $result = create();
    if (daemon != null) {
      $result.daemon = daemon;
    }
    return $result;
  }
  DaemonUpdateResponse._() : super();
  factory DaemonUpdateResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonUpdateResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonUpdateResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon', subBuilder: Daemon.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonUpdateResponse clone() => DaemonUpdateResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonUpdateResponse copyWith(void Function(DaemonUpdateResponse) updates) => super.copyWith((message) => updates(message as DaemonUpdateResponse)) as DaemonUpdateResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonUpdateResponse create() => DaemonUpdateResponse._();
  DaemonUpdateResponse createEmptyInstance() => create();
  static $pb.PbList<DaemonUpdateResponse> createRepeated() => $pb.PbList<DaemonUpdateResponse>();
  @$core.pragma('dart2js:noInline')
  static DaemonUpdateResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonUpdateResponse>(create);
  static DaemonUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonDisableRequest extends $pb.GeneratedMessage {
  factory DaemonDisableRequest({
    Daemon? daemon,
  }) {
    final $result = create();
    if (daemon != null) {
      $result.daemon = daemon;
    }
    return $result;
  }
  DaemonDisableRequest._() : super();
  factory DaemonDisableRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonDisableRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonDisableRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon', subBuilder: Daemon.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonDisableRequest clone() => DaemonDisableRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonDisableRequest copyWith(void Function(DaemonDisableRequest) updates) => super.copyWith((message) => updates(message as DaemonDisableRequest)) as DaemonDisableRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonDisableRequest create() => DaemonDisableRequest._();
  DaemonDisableRequest createEmptyInstance() => create();
  static $pb.PbList<DaemonDisableRequest> createRepeated() => $pb.PbList<DaemonDisableRequest>();
  @$core.pragma('dart2js:noInline')
  static DaemonDisableRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonDisableRequest>(create);
  static DaemonDisableRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}

class DaemonDisableResponse extends $pb.GeneratedMessage {
  factory DaemonDisableResponse({
    Daemon? daemon,
  }) {
    final $result = create();
    if (daemon != null) {
      $result.daemon = daemon;
    }
    return $result;
  }
  DaemonDisableResponse._() : super();
  factory DaemonDisableResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DaemonDisableResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DaemonDisableResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'), createEmptyInstance: create)
    ..aOM<Daemon>(1, _omitFieldNames ? '' : 'daemon', subBuilder: Daemon.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DaemonDisableResponse clone() => DaemonDisableResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DaemonDisableResponse copyWith(void Function(DaemonDisableResponse) updates) => super.copyWith((message) => updates(message as DaemonDisableResponse)) as DaemonDisableResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DaemonDisableResponse create() => DaemonDisableResponse._();
  DaemonDisableResponse createEmptyInstance() => create();
  static $pb.PbList<DaemonDisableResponse> createRepeated() => $pb.PbList<DaemonDisableResponse>();
  @$core.pragma('dart2js:noInline')
  static DaemonDisableResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DaemonDisableResponse>(create);
  static DaemonDisableResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Daemon get daemon => $_getN(0);
  @$pb.TagNumber(1)
  set daemon(Daemon v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasDaemon() => $_has(0);
  @$pb.TagNumber(1)
  void clearDaemon() => clearField(1);
  @$pb.TagNumber(1)
  Daemon ensureDaemon() => $_ensure(0);
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
