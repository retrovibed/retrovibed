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

export 'meta.wireguard.pbenum.dart';

class Wireguard extends $pb.GeneratedMessage {
  factory Wireguard({
    $core.String? id,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? description,
    $core.int? nettype,
    $core.int? port,
    $core.int? rateLimitDns,
    $fixnum.Int64? maximumConnections,
    $core.int? rateLimitOutbound,
  }) {
    final result = Wireguard._();
    if (id != null) result.id = id;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (description != null) result.description = description;
    if (nettype != null) result.nettype = nettype;
    if (port != null) result.port = port;
    if (rateLimitDns != null) result.rateLimitDns = rateLimitDns;
    if (maximumConnections != null)
      result.maximumConnections = maximumConnections;
    if (rateLimitOutbound != null) result.rateLimitOutbound = rateLimitOutbound;
    return result;
  }

  Wireguard._();

  factory Wireguard.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Wireguard()..mergeFromBuffer(data, registry);
  factory Wireguard.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Wireguard()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Wireguard',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: Wireguard.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'created_at')
    ..aOS(3, _omitFieldNames ? '' : 'updated_at')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aI(5, _omitFieldNames ? '' : 'nettype', fieldType: $pb.PbFieldType.OU3)
    ..aI(6, _omitFieldNames ? '' : 'port', fieldType: $pb.PbFieldType.OU3)
    ..aI(7, _omitFieldNames ? '' : 'rate_limit_dns',
        fieldType: $pb.PbFieldType.OU3)
    ..a<$fixnum.Int64>(
        8, _omitFieldNames ? '' : 'maximum_connections', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aI(9, _omitFieldNames ? '' : 'rate_limit_outbound',
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
  @$core.Deprecated('Use Wireguard() / Wireguard.new instead')
  static Wireguard create() => Wireguard._();
  static $pb.GeneratedMessage $_createMessage() => Wireguard._();
  @$core.override
  Wireguard createEmptyInstance() => Wireguard._();
  @$core.pragma('dart2js:noInline')
  static Wireguard getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Wireguard>(Wireguard.$_createMessage);
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
  $core.int get nettype => $_getIZ(4);
  @$pb.TagNumber(5)
  set nettype($core.int value) => $_setUnsignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasNettype() => $_has(4);
  @$pb.TagNumber(5)
  void clearNettype() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.int get port => $_getIZ(5);
  @$pb.TagNumber(6)
  set port($core.int value) => $_setUnsignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasPort() => $_has(5);
  @$pb.TagNumber(6)
  void clearPort() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.int get rateLimitDns => $_getIZ(6);
  @$pb.TagNumber(7)
  set rateLimitDns($core.int value) => $_setUnsignedInt32(6, value);
  @$pb.TagNumber(7)
  $core.bool hasRateLimitDns() => $_has(6);
  @$pb.TagNumber(7)
  void clearRateLimitDns() => $_clearField(7);

  @$pb.TagNumber(8)
  $fixnum.Int64 get maximumConnections => $_getI64(7);
  @$pb.TagNumber(8)
  set maximumConnections($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(8)
  $core.bool hasMaximumConnections() => $_has(7);
  @$pb.TagNumber(8)
  void clearMaximumConnections() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.int get rateLimitOutbound => $_getIZ(8);
  @$pb.TagNumber(9)
  set rateLimitOutbound($core.int value) => $_setUnsignedInt32(8, value);
  @$pb.TagNumber(9)
  $core.bool hasRateLimitOutbound() => $_has(8);
  @$pb.TagNumber(9)
  void clearRateLimitOutbound() => $_clearField(9);
}

class WireguardSearchRequest extends $pb.GeneratedMessage {
  factory WireguardSearchRequest({
    $core.String? query,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = WireguardSearchRequest._();
    if (query != null) result.query = query;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  WireguardSearchRequest._();

  factory WireguardSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardSearchRequest()..mergeFromBuffer(data, registry);
  factory WireguardSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardSearchRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardSearchRequest.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardSearchRequest() / WireguardSearchRequest.new instead')
  static WireguardSearchRequest create() => WireguardSearchRequest._();
  static $pb.GeneratedMessage $_createMessage() => WireguardSearchRequest._();
  @$core.override
  WireguardSearchRequest createEmptyInstance() => WireguardSearchRequest._();
  @$core.pragma('dart2js:noInline')
  static WireguardSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardSearchRequest>(
          WireguardSearchRequest.$_createMessage);
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
    final result = WireguardSearchResponse._();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  WireguardSearchResponse._();

  factory WireguardSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardSearchResponse()..mergeFromBuffer(data, registry);
  factory WireguardSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardSearchResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardSearchResponse.$_createMessage)
    ..aOM<WireguardSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: WireguardSearchRequest.$_createMessage)
    ..pPM<Wireguard>(2, _omitFieldNames ? '' : 'items',
        subBuilder: Wireguard.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardSearchResponse() / WireguardSearchResponse.new instead')
  static WireguardSearchResponse create() => WireguardSearchResponse._();
  static $pb.GeneratedMessage $_createMessage() => WireguardSearchResponse._();
  @$core.override
  WireguardSearchResponse createEmptyInstance() => WireguardSearchResponse._();
  @$core.pragma('dart2js:noInline')
  static WireguardSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardSearchResponse>(
          WireguardSearchResponse.$_createMessage);
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
    final result = WireguardUpdateRequest._();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardUpdateRequest._();

  factory WireguardUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardUpdateRequest()..mergeFromBuffer(data, registry);
  factory WireguardUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardUpdateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardUpdateRequest.$_createMessage)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardUpdateRequest() / WireguardUpdateRequest.new instead')
  static WireguardUpdateRequest create() => WireguardUpdateRequest._();
  static $pb.GeneratedMessage $_createMessage() => WireguardUpdateRequest._();
  @$core.override
  WireguardUpdateRequest createEmptyInstance() => WireguardUpdateRequest._();
  @$core.pragma('dart2js:noInline')
  static WireguardUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardUpdateRequest>(
          WireguardUpdateRequest.$_createMessage);
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
    final result = WireguardUpdateResponse._();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardUpdateResponse._();

  factory WireguardUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardUpdateResponse()..mergeFromBuffer(data, registry);
  factory WireguardUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardUpdateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardUpdateResponse.$_createMessage)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardUpdateResponse() / WireguardUpdateResponse.new instead')
  static WireguardUpdateResponse create() => WireguardUpdateResponse._();
  static $pb.GeneratedMessage $_createMessage() => WireguardUpdateResponse._();
  @$core.override
  WireguardUpdateResponse createEmptyInstance() => WireguardUpdateResponse._();
  @$core.pragma('dart2js:noInline')
  static WireguardUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardUpdateResponse>(
          WireguardUpdateResponse.$_createMessage);
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
  factory WireguardTouchRequest({
    $core.int? nettype,
  }) {
    final result = WireguardTouchRequest._();
    if (nettype != null) result.nettype = nettype;
    return result;
  }

  WireguardTouchRequest._();

  factory WireguardTouchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardTouchRequest()..mergeFromBuffer(data, registry);
  factory WireguardTouchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardTouchRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardTouchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardTouchRequest.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'nettype', fieldType: $pb.PbFieldType.OU3)
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
  @$core.Deprecated(
      'Use WireguardTouchRequest() / WireguardTouchRequest.new instead')
  static WireguardTouchRequest create() => WireguardTouchRequest._();
  static $pb.GeneratedMessage $_createMessage() => WireguardTouchRequest._();
  @$core.override
  WireguardTouchRequest createEmptyInstance() => WireguardTouchRequest._();
  @$core.pragma('dart2js:noInline')
  static WireguardTouchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardTouchRequest>(
          WireguardTouchRequest.$_createMessage);
  static WireguardTouchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get nettype => $_getIZ(0);
  @$pb.TagNumber(1)
  set nettype($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasNettype() => $_has(0);
  @$pb.TagNumber(1)
  void clearNettype() => $_clearField(1);
}

class WireguardTouchResponse extends $pb.GeneratedMessage {
  factory WireguardTouchResponse({
    Wireguard? wireguard,
  }) {
    final result = WireguardTouchResponse._();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardTouchResponse._();

  factory WireguardTouchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardTouchResponse()..mergeFromBuffer(data, registry);
  factory WireguardTouchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardTouchResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardTouchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardTouchResponse.$_createMessage)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardTouchResponse() / WireguardTouchResponse.new instead')
  static WireguardTouchResponse create() => WireguardTouchResponse._();
  static $pb.GeneratedMessage $_createMessage() => WireguardTouchResponse._();
  @$core.override
  WireguardTouchResponse createEmptyInstance() => WireguardTouchResponse._();
  @$core.pragma('dart2js:noInline')
  static WireguardTouchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardTouchResponse>(
          WireguardTouchResponse.$_createMessage);
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
  factory WireguardUploadRequest() => WireguardUploadRequest._();

  WireguardUploadRequest._();

  factory WireguardUploadRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardUploadRequest()..mergeFromBuffer(data, registry);
  factory WireguardUploadRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardUploadRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardUploadRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardUploadRequest.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardUploadRequest() / WireguardUploadRequest.new instead')
  static WireguardUploadRequest create() => WireguardUploadRequest._();
  static $pb.GeneratedMessage $_createMessage() => WireguardUploadRequest._();
  @$core.override
  WireguardUploadRequest createEmptyInstance() => WireguardUploadRequest._();
  @$core.pragma('dart2js:noInline')
  static WireguardUploadRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardUploadRequest>(
          WireguardUploadRequest.$_createMessage);
  static WireguardUploadRequest? _defaultInstance;
}

class WireguardUploadResponse extends $pb.GeneratedMessage {
  factory WireguardUploadResponse({
    Wireguard? wireguard,
  }) {
    final result = WireguardUploadResponse._();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardUploadResponse._();

  factory WireguardUploadResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardUploadResponse()..mergeFromBuffer(data, registry);
  factory WireguardUploadResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardUploadResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardUploadResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardUploadResponse.$_createMessage)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardUploadResponse() / WireguardUploadResponse.new instead')
  static WireguardUploadResponse create() => WireguardUploadResponse._();
  static $pb.GeneratedMessage $_createMessage() => WireguardUploadResponse._();
  @$core.override
  WireguardUploadResponse createEmptyInstance() => WireguardUploadResponse._();
  @$core.pragma('dart2js:noInline')
  static WireguardUploadResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardUploadResponse>(
          WireguardUploadResponse.$_createMessage);
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
  factory WireguardCurrentRequest({
    $core.int? nettype,
  }) {
    final result = WireguardCurrentRequest._();
    if (nettype != null) result.nettype = nettype;
    return result;
  }

  WireguardCurrentRequest._();

  factory WireguardCurrentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardCurrentRequest()..mergeFromBuffer(data, registry);
  factory WireguardCurrentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardCurrentRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardCurrentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardCurrentRequest.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'nettype', fieldType: $pb.PbFieldType.OU3)
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
  @$core.Deprecated(
      'Use WireguardCurrentRequest() / WireguardCurrentRequest.new instead')
  static WireguardCurrentRequest create() => WireguardCurrentRequest._();
  static $pb.GeneratedMessage $_createMessage() => WireguardCurrentRequest._();
  @$core.override
  WireguardCurrentRequest createEmptyInstance() => WireguardCurrentRequest._();
  @$core.pragma('dart2js:noInline')
  static WireguardCurrentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardCurrentRequest>(
          WireguardCurrentRequest.$_createMessage);
  static WireguardCurrentRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get nettype => $_getIZ(0);
  @$pb.TagNumber(1)
  set nettype($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasNettype() => $_has(0);
  @$pb.TagNumber(1)
  void clearNettype() => $_clearField(1);
}

class WireguardCurrentResponse extends $pb.GeneratedMessage {
  factory WireguardCurrentResponse({
    Wireguard? wireguard,
    $core.bool? online,
    $core.String? ip,
    $core.String? ip4,
  }) {
    final result = WireguardCurrentResponse._();
    if (wireguard != null) result.wireguard = wireguard;
    if (online != null) result.online = online;
    if (ip != null) result.ip = ip;
    if (ip4 != null) result.ip4 = ip4;
    return result;
  }

  WireguardCurrentResponse._();

  factory WireguardCurrentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardCurrentResponse()..mergeFromBuffer(data, registry);
  factory WireguardCurrentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardCurrentResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardCurrentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardCurrentResponse.$_createMessage)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardCurrentResponse() / WireguardCurrentResponse.new instead')
  static WireguardCurrentResponse create() => WireguardCurrentResponse._();
  static $pb.GeneratedMessage $_createMessage() => WireguardCurrentResponse._();
  @$core.override
  WireguardCurrentResponse createEmptyInstance() =>
      WireguardCurrentResponse._();
  @$core.pragma('dart2js:noInline')
  static WireguardCurrentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardCurrentResponse>(
          WireguardCurrentResponse.$_createMessage);
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
  factory WireguardDeleteRequest() => WireguardDeleteRequest._();

  WireguardDeleteRequest._();

  factory WireguardDeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardDeleteRequest()..mergeFromBuffer(data, registry);
  factory WireguardDeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardDeleteRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardDeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardDeleteRequest.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardDeleteRequest() / WireguardDeleteRequest.new instead')
  static WireguardDeleteRequest create() => WireguardDeleteRequest._();
  static $pb.GeneratedMessage $_createMessage() => WireguardDeleteRequest._();
  @$core.override
  WireguardDeleteRequest createEmptyInstance() => WireguardDeleteRequest._();
  @$core.pragma('dart2js:noInline')
  static WireguardDeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardDeleteRequest>(
          WireguardDeleteRequest.$_createMessage);
  static WireguardDeleteRequest? _defaultInstance;
}

class WireguardDeleteResponse extends $pb.GeneratedMessage {
  factory WireguardDeleteResponse({
    Wireguard? wireguard,
  }) {
    final result = WireguardDeleteResponse._();
    if (wireguard != null) result.wireguard = wireguard;
    return result;
  }

  WireguardDeleteResponse._();

  factory WireguardDeleteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardDeleteResponse()..mergeFromBuffer(data, registry);
  factory WireguardDeleteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      WireguardDeleteResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WireguardDeleteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: WireguardDeleteResponse.$_createMessage)
    ..aOM<Wireguard>(1, _omitFieldNames ? '' : 'wireguard',
        protoName: 'Wireguard', subBuilder: Wireguard.$_createMessage)
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
  @$core.Deprecated(
      'Use WireguardDeleteResponse() / WireguardDeleteResponse.new instead')
  static WireguardDeleteResponse create() => WireguardDeleteResponse._();
  static $pb.GeneratedMessage $_createMessage() => WireguardDeleteResponse._();
  @$core.override
  WireguardDeleteResponse createEmptyInstance() => WireguardDeleteResponse._();
  @$core.pragma('dart2js:noInline')
  static WireguardDeleteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WireguardDeleteResponse>(
          WireguardDeleteResponse.$_createMessage);
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
