// This is a generated file - do not edit.
//
// Generated from wireguard/meta.wireguard.proto.

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

class Wireguard extends $pb.GeneratedMessage {
  factory Wireguard({
    $core.String? id,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? description,
    $core.bool? default_5,
    $core.int? port,
    $core.int? dnsRateLimit,
    $fixnum.Int64? maximumConnections,
    $core.int? outboundRateLimit,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (description != null) result.description = description;
    if (default_5 != null) result.default_5 = default_5;
    if (port != null) result.port = port;
    if (dnsRateLimit != null) result.dnsRateLimit = dnsRateLimit;
    if (maximumConnections != null)
      result.maximumConnections = maximumConnections;
    if (outboundRateLimit != null) result.outboundRateLimit = outboundRateLimit;
    return result;
  }

  Wireguard._();

  factory Wireguard.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Wireguard.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Wireguard',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'created_at')
    ..aOS(3, _omitFieldNames ? '' : 'updated_at')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aOB(5, _omitFieldNames ? '' : 'default')
    ..aI(6, _omitFieldNames ? '' : 'port', fieldType: $pb.PbFieldType.OU3)
    ..aI(7, _omitFieldNames ? '' : 'dns_rate_limit',
        fieldType: $pb.PbFieldType.OU3)
    ..a<$fixnum.Int64>(
        8, _omitFieldNames ? '' : 'maximum_connections', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aI(9, _omitFieldNames ? '' : 'outbound_rate_limit',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Wireguard clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Wireguard copyWith(void Function(Wireguard) updates) =>
      super.copyWith((message) => updates(message as Wireguard)) as Wireguard;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Wireguard create() => Wireguard._();
  @$core.override
  Wireguard createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Wireguard getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Wireguard>(create);
  static Wireguard? _defaultInstance;

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
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.bool get default_5 => $_getBF(4);
  @$pb.TagNumber(5)
  set default_5($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(5)
  $core.bool hasDefault_5() => $_has(4);
  @$pb.TagNumber(5)
  void clearDefault_5() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.int get port => $_getIZ(5);
  @$pb.TagNumber(6)
  set port($core.int value) => $_setUnsignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasPort() => $_has(5);
  @$pb.TagNumber(6)
  void clearPort() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.int get dnsRateLimit => $_getIZ(6);
  @$pb.TagNumber(7)
  set dnsRateLimit($core.int value) => $_setUnsignedInt32(6, value);
  @$pb.TagNumber(7)
  $core.bool hasDnsRateLimit() => $_has(6);
  @$pb.TagNumber(7)
  void clearDnsRateLimit() => $_clearField(7);

  @$pb.TagNumber(8)
  $fixnum.Int64 get maximumConnections => $_getI64(7);
  @$pb.TagNumber(8)
  set maximumConnections($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(8)
  $core.bool hasMaximumConnections() => $_has(7);
  @$pb.TagNumber(8)
  void clearMaximumConnections() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.int get outboundRateLimit => $_getIZ(8);
  @$pb.TagNumber(9)
  set outboundRateLimit($core.int value) => $_setUnsignedInt32(8, value);
  @$pb.TagNumber(9)
  $core.bool hasOutboundRateLimit() => $_has(8);
  @$pb.TagNumber(9)
  void clearOutboundRateLimit() => $_clearField(9);
}

class WireguardSearchRequest extends $pb.GeneratedMessage {
  factory WireguardSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (query != null) result.query = query;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  WireguardSearchRequest._();

  factory WireguardSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardSearchRequest copyWith(
          void Function(WireguardSearchRequest) updates) =>
      super.copyWith((message) => updates(message as WireguardSearchRequest))
          as WireguardSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardSearchRequest create() => WireguardSearchRequest._();
  @$core.override
  WireguardSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardSearchRequest>(create);
  static WireguardSearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get offset => $_getI64(1);
  @$pb.TagNumber(2)
  set offset($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(2)
  void clearOffset() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get limit => $_getI64(2);
  @$pb.TagNumber(3)
  set limit($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => $_clearField(3);
}

class WireguardSearchResponse extends $pb.GeneratedMessage {
  factory WireguardSearchResponse({
    WireguardSearchRequest? next,
    $core.Iterable<Wireguard>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  WireguardSearchResponse._();

  factory WireguardSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<WireguardSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: WireguardSearchRequest.create)
    ..pPM<Wireguard>(2, _omitFieldNames ? '' : 'items',
        subBuilder: Wireguard.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardSearchResponse copyWith(
          void Function(WireguardSearchResponse) updates) =>
      super.copyWith((message) => updates(message as WireguardSearchResponse))
          as WireguardSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardSearchResponse create() => WireguardSearchResponse._();
  @$core.override
  WireguardSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardSearchResponse>(create);
  static WireguardSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  WireguardSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(WireguardSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  WireguardSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Wireguard> get items => $_getList(1);
}

class WireguardUpdateRequest extends $pb.GeneratedMessage {
  factory WireguardUpdateRequest({
    Wireguard? wireguard,
  }) {
    final result = create();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardUpdateRequest._();

  factory WireguardUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardUpdateRequest copyWith(
          void Function(WireguardUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as WireguardUpdateRequest))
          as WireguardUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardUpdateRequest create() => WireguardUpdateRequest._();
  @$core.override
  WireguardUpdateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardUpdateRequest>(create);
  static WireguardUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Wireguard get wireguard => $_getN(0);
  @$pb.TagNumber(1)
  set wireguard(Wireguard value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasWireguard() => $_has(0);
  @$pb.TagNumber(1)
  void clearWireguard() => $_clearField(1);
  @$pb.TagNumber(1)
  Wireguard ensureWireguard() => $_ensure(0);
}

class WireguardUpdateResponse extends $pb.GeneratedMessage {
  factory WireguardUpdateResponse({
    Wireguard? wireguard,
  }) {
    final result = create();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardUpdateResponse._();

  factory WireguardUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardUpdateResponse copyWith(
          void Function(WireguardUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as WireguardUpdateResponse))
          as WireguardUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardUpdateResponse create() => WireguardUpdateResponse._();
  @$core.override
  WireguardUpdateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardUpdateResponse>(create);
  static WireguardUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Wireguard get wireguard => $_getN(0);
  @$pb.TagNumber(1)
  set wireguard(Wireguard value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasWireguard() => $_has(0);
  @$pb.TagNumber(1)
  void clearWireguard() => $_clearField(1);
  @$pb.TagNumber(1)
  Wireguard ensureWireguard() => $_ensure(0);
}

class WireguardTouchRequest extends $pb.GeneratedMessage {
  factory WireguardTouchRequest() => create();

  WireguardTouchRequest._();

  factory WireguardTouchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardTouchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardTouchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardTouchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardTouchRequest copyWith(
          void Function(WireguardTouchRequest) updates) =>
      super.copyWith((message) => updates(message as WireguardTouchRequest))
          as WireguardTouchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardTouchRequest create() => WireguardTouchRequest._();
  @$core.override
  WireguardTouchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardTouchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardTouchRequest>(create);
  static WireguardTouchRequest? _defaultInstance;
}

class WireguardTouchResponse extends $pb.GeneratedMessage {
  factory WireguardTouchResponse({
    Wireguard? wireguard,
  }) {
    final result = create();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardTouchResponse._();

  factory WireguardTouchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardTouchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardTouchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardTouchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardTouchResponse copyWith(
          void Function(WireguardTouchResponse) updates) =>
      super.copyWith((message) => updates(message as WireguardTouchResponse))
          as WireguardTouchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardTouchResponse create() => WireguardTouchResponse._();
  @$core.override
  WireguardTouchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardTouchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardTouchResponse>(create);
  static WireguardTouchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Wireguard get wireguard => $_getN(0);
  @$pb.TagNumber(1)
  set wireguard(Wireguard value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasWireguard() => $_has(0);
  @$pb.TagNumber(1)
  void clearWireguard() => $_clearField(1);
  @$pb.TagNumber(1)
  Wireguard ensureWireguard() => $_ensure(0);
}

class WireguardUploadRequest extends $pb.GeneratedMessage {
  factory WireguardUploadRequest() => create();

  WireguardUploadRequest._();

  factory WireguardUploadRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardUploadRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardUploadRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardUploadRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardUploadRequest copyWith(
          void Function(WireguardUploadRequest) updates) =>
      super.copyWith((message) => updates(message as WireguardUploadRequest))
          as WireguardUploadRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardUploadRequest create() => WireguardUploadRequest._();
  @$core.override
  WireguardUploadRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardUploadRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardUploadRequest>(create);
  static WireguardUploadRequest? _defaultInstance;
}

class WireguardUploadResponse extends $pb.GeneratedMessage {
  factory WireguardUploadResponse({
    Wireguard? wireguard,
  }) {
    final result = create();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardUploadResponse._();

  factory WireguardUploadResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardUploadResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardUploadResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardUploadResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardUploadResponse copyWith(
          void Function(WireguardUploadResponse) updates) =>
      super.copyWith((message) => updates(message as WireguardUploadResponse))
          as WireguardUploadResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardUploadResponse create() => WireguardUploadResponse._();
  @$core.override
  WireguardUploadResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardUploadResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardUploadResponse>(create);
  static WireguardUploadResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Wireguard get wireguard => $_getN(0);
  @$pb.TagNumber(1)
  set wireguard(Wireguard value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasWireguard() => $_has(0);
  @$pb.TagNumber(1)
  void clearWireguard() => $_clearField(1);
  @$pb.TagNumber(1)
  Wireguard ensureWireguard() => $_ensure(0);
}

class WireguardCurrentRequest extends $pb.GeneratedMessage {
  factory WireguardCurrentRequest() => create();

  WireguardCurrentRequest._();

  factory WireguardCurrentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardCurrentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardCurrentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardCurrentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardCurrentRequest copyWith(
          void Function(WireguardCurrentRequest) updates) =>
      super.copyWith((message) => updates(message as WireguardCurrentRequest))
          as WireguardCurrentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardCurrentRequest create() => WireguardCurrentRequest._();
  @$core.override
  WireguardCurrentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardCurrentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardCurrentRequest>(create);
  static WireguardCurrentRequest? _defaultInstance;
}

class WireguardCurrentResponse extends $pb.GeneratedMessage {
  factory WireguardCurrentResponse({
    Wireguard? wireguard,
    $core.bool? online,
    $core.String? ip,
    $core.String? ip4,
  }) {
    final result = create();
    if (wireguard != null) result.wireguard = wireguard;
    if (online != null) result.online = online;
    if (ip != null) result.ip = ip;
    if (ip4 != null) result.ip4 = ip4;
    return result;
  }

  WireguardCurrentResponse._();

  factory WireguardCurrentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardCurrentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardCurrentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.create)
    ..aOB(2, _omitFieldNames ? '' : 'online')
    ..aOS(3, _omitFieldNames ? '' : 'ip')
    ..aOS(4, _omitFieldNames ? '' : 'ip4')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardCurrentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardCurrentResponse copyWith(
          void Function(WireguardCurrentResponse) updates) =>
      super.copyWith((message) => updates(message as WireguardCurrentResponse))
          as WireguardCurrentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardCurrentResponse create() => WireguardCurrentResponse._();
  @$core.override
  WireguardCurrentResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardCurrentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardCurrentResponse>(create);
  static WireguardCurrentResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Wireguard get wireguard => $_getN(0);
  @$pb.TagNumber(1)
  set wireguard(Wireguard value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasWireguard() => $_has(0);
  @$pb.TagNumber(1)
  void clearWireguard() => $_clearField(1);
  @$pb.TagNumber(1)
  Wireguard ensureWireguard() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.bool get online => $_getBF(1);
  @$pb.TagNumber(2)
  set online($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOnline() => $_has(1);
  @$pb.TagNumber(2)
  void clearOnline() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get ip => $_getSZ(2);
  @$pb.TagNumber(3)
  set ip($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasIp() => $_has(2);
  @$pb.TagNumber(3)
  void clearIp() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get ip4 => $_getSZ(3);
  @$pb.TagNumber(4)
  set ip4($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasIp4() => $_has(3);
  @$pb.TagNumber(4)
  void clearIp4() => $_clearField(4);
}

class WireguardDeleteRequest extends $pb.GeneratedMessage {
  factory WireguardDeleteRequest() => create();

  WireguardDeleteRequest._();

  factory WireguardDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardDeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardDeleteRequest copyWith(
          void Function(WireguardDeleteRequest) updates) =>
      super.copyWith((message) => updates(message as WireguardDeleteRequest))
          as WireguardDeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardDeleteRequest create() => WireguardDeleteRequest._();
  @$core.override
  WireguardDeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardDeleteRequest>(create);
  static WireguardDeleteRequest? _defaultInstance;
}

class WireguardDeleteResponse extends $pb.GeneratedMessage {
  factory WireguardDeleteResponse({
    Wireguard? wireguard,
  }) {
    final result = create();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardDeleteResponse._();

  factory WireguardDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WireguardDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardDeleteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WireguardDeleteResponse copyWith(
          void Function(WireguardDeleteResponse) updates) =>
      super.copyWith((message) => updates(message as WireguardDeleteResponse))
          as WireguardDeleteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WireguardDeleteResponse create() => WireguardDeleteResponse._();
  @$core.override
  WireguardDeleteResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WireguardDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardDeleteResponse>(create);
  static WireguardDeleteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Wireguard get wireguard => $_getN(0);
  @$pb.TagNumber(1)
  set wireguard(Wireguard value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasWireguard() => $_has(0);
  @$pb.TagNumber(1)
  void clearWireguard() => $_clearField(1);
  @$pb.TagNumber(1)
  Wireguard ensureWireguard() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
