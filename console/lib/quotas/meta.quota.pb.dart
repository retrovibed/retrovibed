// This is a generated file - do not edit.
//
// Generated from meta.quota.proto.

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

class Quota extends $pb.GeneratedMessage {
  factory Quota({
    $core.String? sku,
    $core.String? accountId,
    $core.String? createdAt,
    $core.String? updatedAt,
    $core.String? renewedAt,
    $core.String? description,
    $core.bool? adjustable,
    $fixnum.Int64? maximum,
    $fixnum.Int64? credits,
    $fixnum.Int64? reserved,
    $fixnum.Int64? consumed,
    $fixnum.Int64? rollover,
    $fixnum.Int64? granted,
  }) {
    final result = create();
    if (sku != null) result.sku = sku;
    if (accountId != null) result.accountId = accountId;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (renewedAt != null) result.renewedAt = renewedAt;
    if (description != null) result.description = description;
    if (adjustable != null) result.adjustable = adjustable;
    if (maximum != null) result.maximum = maximum;
    if (credits != null) result.credits = credits;
    if (reserved != null) result.reserved = reserved;
    if (consumed != null) result.consumed = consumed;
    if (rollover != null) result.rollover = rollover;
    if (granted != null) result.granted = granted;
    return result;
  }

  Quota._();

  factory Quota.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Quota.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Quota',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id', protoName: 'sku')
    ..aOS(2, _omitFieldNames ? '' : 'account_id')
    ..aOS(3, _omitFieldNames ? '' : 'created_at')
    ..aOS(4, _omitFieldNames ? '' : 'updated_at')
    ..aOS(5, _omitFieldNames ? '' : 'renewed_at')
    ..aOS(6, _omitFieldNames ? '' : 'description')
    ..aOB(7, _omitFieldNames ? '' : 'adjustable')
    ..aInt64(8, _omitFieldNames ? '' : 'maximum')
    ..aInt64(9, _omitFieldNames ? '' : 'credits')
    ..aInt64(10, _omitFieldNames ? '' : 'reserved')
    ..aInt64(11, _omitFieldNames ? '' : 'consumed')
    ..a<$fixnum.Int64>(
        12, _omitFieldNames ? '' : 'rollover', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aInt64(13, _omitFieldNames ? '' : 'granted')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Quota clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Quota copyWith(void Function(Quota) updates) =>
      super.copyWith((message) => updates(message as Quota)) as Quota;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Quota create() => Quota._();
  @$core.override
  Quota createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Quota getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Quota>(create);
  static Quota? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get sku => $_getSZ(0);
  @$pb.TagNumber(1)
  set sku($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSku() => $_has(0);
  @$pb.TagNumber(1)
  void clearSku() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get accountId => $_getSZ(1);
  @$pb.TagNumber(2)
  set accountId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAccountId() => $_has(1);
  @$pb.TagNumber(2)
  void clearAccountId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get createdAt => $_getSZ(2);
  @$pb.TagNumber(3)
  set createdAt($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCreatedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearCreatedAt() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get updatedAt => $_getSZ(3);
  @$pb.TagNumber(4)
  set updatedAt($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasUpdatedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearUpdatedAt() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get renewedAt => $_getSZ(4);
  @$pb.TagNumber(5)
  set renewedAt($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasRenewedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearRenewedAt() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get description => $_getSZ(5);
  @$pb.TagNumber(6)
  set description($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasDescription() => $_has(5);
  @$pb.TagNumber(6)
  void clearDescription() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.bool get adjustable => $_getBF(6);
  @$pb.TagNumber(7)
  set adjustable($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(7)
  $core.bool hasAdjustable() => $_has(6);
  @$pb.TagNumber(7)
  void clearAdjustable() => $_clearField(7);

  @$pb.TagNumber(8)
  $fixnum.Int64 get maximum => $_getI64(7);
  @$pb.TagNumber(8)
  set maximum($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(8)
  $core.bool hasMaximum() => $_has(7);
  @$pb.TagNumber(8)
  void clearMaximum() => $_clearField(8);

  @$pb.TagNumber(9)
  $fixnum.Int64 get credits => $_getI64(8);
  @$pb.TagNumber(9)
  set credits($fixnum.Int64 value) => $_setInt64(8, value);
  @$pb.TagNumber(9)
  $core.bool hasCredits() => $_has(8);
  @$pb.TagNumber(9)
  void clearCredits() => $_clearField(9);

  @$pb.TagNumber(10)
  $fixnum.Int64 get reserved => $_getI64(9);
  @$pb.TagNumber(10)
  set reserved($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(10)
  $core.bool hasReserved() => $_has(9);
  @$pb.TagNumber(10)
  void clearReserved() => $_clearField(10);

  @$pb.TagNumber(11)
  $fixnum.Int64 get consumed => $_getI64(10);
  @$pb.TagNumber(11)
  set consumed($fixnum.Int64 value) => $_setInt64(10, value);
  @$pb.TagNumber(11)
  $core.bool hasConsumed() => $_has(10);
  @$pb.TagNumber(11)
  void clearConsumed() => $_clearField(11);

  @$pb.TagNumber(12)
  $fixnum.Int64 get rollover => $_getI64(11);
  @$pb.TagNumber(12)
  set rollover($fixnum.Int64 value) => $_setInt64(11, value);
  @$pb.TagNumber(12)
  $core.bool hasRollover() => $_has(11);
  @$pb.TagNumber(12)
  void clearRollover() => $_clearField(12);

  @$pb.TagNumber(13)
  $fixnum.Int64 get granted => $_getI64(12);
  @$pb.TagNumber(13)
  set granted($fixnum.Int64 value) => $_setInt64(12, value);
  @$pb.TagNumber(13)
  $core.bool hasGranted() => $_has(12);
  @$pb.TagNumber(13)
  void clearGranted() => $_clearField(13);
}

class QuotaSearchRequest extends $pb.GeneratedMessage {
  factory QuotaSearchRequest({
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

  QuotaSearchRequest._();

  factory QuotaSearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuotaSearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuotaSearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaSearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaSearchRequest copyWith(void Function(QuotaSearchRequest) updates) =>
      super.copyWith((message) => updates(message as QuotaSearchRequest))
          as QuotaSearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuotaSearchRequest create() => QuotaSearchRequest._();
  @$core.override
  QuotaSearchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuotaSearchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuotaSearchRequest>(create);
  static QuotaSearchRequest? _defaultInstance;

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

class QuotaSearchResponse extends $pb.GeneratedMessage {
  factory QuotaSearchResponse({
    QuotaSearchRequest? next,
    $core.Iterable<Quota>? items,
  }) {
    final result = create();
    if (next != null) result.next = next;
    if (items != null) result.items.addAll(items);
    return result;
  }

  QuotaSearchResponse._();

  factory QuotaSearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuotaSearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuotaSearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<QuotaSearchRequest>(1, _omitFieldNames ? '' : 'next',
        subBuilder: QuotaSearchRequest.create)
    ..pPM<Quota>(2, _omitFieldNames ? '' : 'items', subBuilder: Quota.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaSearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaSearchResponse copyWith(void Function(QuotaSearchResponse) updates) =>
      super.copyWith((message) => updates(message as QuotaSearchResponse))
          as QuotaSearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuotaSearchResponse create() => QuotaSearchResponse._();
  @$core.override
  QuotaSearchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuotaSearchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuotaSearchResponse>(create);
  static QuotaSearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  QuotaSearchRequest get next => $_getN(0);
  @$pb.TagNumber(1)
  set next(QuotaSearchRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasNext() => $_has(0);
  @$pb.TagNumber(1)
  void clearNext() => $_clearField(1);
  @$pb.TagNumber(1)
  QuotaSearchRequest ensureNext() => $_ensure(0);

  @$pb.TagNumber(2)
  $pb.PbList<Quota> get items => $_getList(1);
}

class QuotaUpdateRequest extends $pb.GeneratedMessage {
  factory QuotaUpdateRequest({
    Quota? quota,
  }) {
    final result = create();
    if (quota != null) result.quota = quota;
    return result;
  }

  QuotaUpdateRequest._();

  factory QuotaUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuotaUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuotaUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Quota>(1, _omitFieldNames ? '' : 'quota', subBuilder: Quota.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaUpdateRequest copyWith(void Function(QuotaUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as QuotaUpdateRequest))
          as QuotaUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuotaUpdateRequest create() => QuotaUpdateRequest._();
  @$core.override
  QuotaUpdateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuotaUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuotaUpdateRequest>(create);
  static QuotaUpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Quota get quota => $_getN(0);
  @$pb.TagNumber(1)
  set quota(Quota value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasQuota() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuota() => $_clearField(1);
  @$pb.TagNumber(1)
  Quota ensureQuota() => $_ensure(0);
}

class QuotaUpdateResponse extends $pb.GeneratedMessage {
  factory QuotaUpdateResponse({
    Quota? quota,
  }) {
    final result = create();
    if (quota != null) result.quota = quota;
    return result;
  }

  QuotaUpdateResponse._();

  factory QuotaUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuotaUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuotaUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Quota>(1, _omitFieldNames ? '' : 'quota', subBuilder: Quota.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaUpdateResponse copyWith(void Function(QuotaUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as QuotaUpdateResponse))
          as QuotaUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuotaUpdateResponse create() => QuotaUpdateResponse._();
  @$core.override
  QuotaUpdateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuotaUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuotaUpdateResponse>(create);
  static QuotaUpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Quota get quota => $_getN(0);
  @$pb.TagNumber(1)
  set quota(Quota value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasQuota() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuota() => $_clearField(1);
  @$pb.TagNumber(1)
  Quota ensureQuota() => $_ensure(0);
}

class QuotaFindRequest extends $pb.GeneratedMessage {
  factory QuotaFindRequest() => create();

  QuotaFindRequest._();

  factory QuotaFindRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuotaFindRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuotaFindRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaFindRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaFindRequest copyWith(void Function(QuotaFindRequest) updates) =>
      super.copyWith((message) => updates(message as QuotaFindRequest))
          as QuotaFindRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuotaFindRequest create() => QuotaFindRequest._();
  @$core.override
  QuotaFindRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuotaFindRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuotaFindRequest>(create);
  static QuotaFindRequest? _defaultInstance;
}

class QuotaFindResponse extends $pb.GeneratedMessage {
  factory QuotaFindResponse({
    Quota? quota,
  }) {
    final result = create();
    if (quota != null) result.quota = quota;
    return result;
  }

  QuotaFindResponse._();

  factory QuotaFindResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuotaFindResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuotaFindResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOM<Quota>(1, _omitFieldNames ? '' : 'quota', subBuilder: Quota.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaFindResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaFindResponse copyWith(void Function(QuotaFindResponse) updates) =>
      super.copyWith((message) => updates(message as QuotaFindResponse))
          as QuotaFindResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuotaFindResponse create() => QuotaFindResponse._();
  @$core.override
  QuotaFindResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuotaFindResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuotaFindResponse>(create);
  static QuotaFindResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Quota get quota => $_getN(0);
  @$pb.TagNumber(1)
  set quota(Quota value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasQuota() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuota() => $_clearField(1);
  @$pb.TagNumber(1)
  Quota ensureQuota() => $_ensure(0);
}

class Adjustment extends $pb.GeneratedMessage {
  factory Adjustment({
    $core.String? sku,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (sku != null) result.sku = sku;
    if (limit != null) result.limit = limit;
    return result;
  }

  Adjustment._();

  factory Adjustment.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Adjustment.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Adjustment',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'sku')
    ..aInt64(2, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Adjustment clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Adjustment copyWith(void Function(Adjustment) updates) =>
      super.copyWith((message) => updates(message as Adjustment)) as Adjustment;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Adjustment create() => Adjustment._();
  @$core.override
  Adjustment createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Adjustment getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Adjustment>(create);
  static Adjustment? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get sku => $_getSZ(0);
  @$pb.TagNumber(1)
  set sku($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSku() => $_has(0);
  @$pb.TagNumber(1)
  void clearSku() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get limit => $_getI64(1);
  @$pb.TagNumber(2)
  set limit($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(2)
  void clearLimit() => $_clearField(2);
}

class QuotaAdjustmentRequest extends $pb.GeneratedMessage {
  factory QuotaAdjustmentRequest({
    $core.Iterable<Adjustment>? adjustments,
    $core.String? expiresAt,
  }) {
    final result = create();
    if (adjustments != null) result.adjustments.addAll(adjustments);
    if (expiresAt != null) result.expiresAt = expiresAt;
    return result;
  }

  QuotaAdjustmentRequest._();

  factory QuotaAdjustmentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuotaAdjustmentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuotaAdjustmentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..pPM<Adjustment>(1, _omitFieldNames ? '' : 'adjustments',
        subBuilder: Adjustment.create)
    ..aOS(2, _omitFieldNames ? '' : 'expires_at')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaAdjustmentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaAdjustmentRequest copyWith(
          void Function(QuotaAdjustmentRequest) updates) =>
      super.copyWith((message) => updates(message as QuotaAdjustmentRequest))
          as QuotaAdjustmentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuotaAdjustmentRequest create() => QuotaAdjustmentRequest._();
  @$core.override
  QuotaAdjustmentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuotaAdjustmentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuotaAdjustmentRequest>(create);
  static QuotaAdjustmentRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<Adjustment> get adjustments => $_getList(0);

  @$pb.TagNumber(2)
  $core.String get expiresAt => $_getSZ(1);
  @$pb.TagNumber(2)
  set expiresAt($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasExpiresAt() => $_has(1);
  @$pb.TagNumber(2)
  void clearExpiresAt() => $_clearField(2);
}

class QuotaAdjustmentResponse extends $pb.GeneratedMessage {
  factory QuotaAdjustmentResponse() => create();

  QuotaAdjustmentResponse._();

  factory QuotaAdjustmentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuotaAdjustmentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuotaAdjustmentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaAdjustmentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaAdjustmentResponse copyWith(
          void Function(QuotaAdjustmentResponse) updates) =>
      super.copyWith((message) => updates(message as QuotaAdjustmentResponse))
          as QuotaAdjustmentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuotaAdjustmentResponse create() => QuotaAdjustmentResponse._();
  @$core.override
  QuotaAdjustmentResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuotaAdjustmentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuotaAdjustmentResponse>(create);
  static QuotaAdjustmentResponse? _defaultInstance;
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
