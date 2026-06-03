// This is a generated file - do not edit.
//
// Generated from community.metrics.proto.

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

class CommunityMetric extends $pb.GeneratedMessage {
  factory CommunityMetric({
    $core.String? id,
    $core.String? communityId,
    $core.String? periodStart,
    $core.String? periodEnd,
    $core.int? subscribers,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (communityId != null) result.communityId = communityId;
    if (periodStart != null) result.periodStart = periodStart;
    if (periodEnd != null) result.periodEnd = periodEnd;
    if (subscribers != null) result.subscribers = subscribers;
    return result;
  }

  CommunityMetric._();

  factory CommunityMetric.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityMetric.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityMetric',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'community_id')
    ..aOS(3, _omitFieldNames ? '' : 'period_start')
    ..aOS(4, _omitFieldNames ? '' : 'period_end')
    ..aI(5, _omitFieldNames ? '' : 'subscribers',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetric clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetric copyWith(void Function(CommunityMetric) updates) =>
      super.copyWith((message) => updates(message as CommunityMetric))
          as CommunityMetric;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityMetric create() => CommunityMetric._();
  @$core.override
  CommunityMetric createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityMetric getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityMetric>(create);
  static CommunityMetric? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get communityId => $_getSZ(1);
  @$pb.TagNumber(2)
  set communityId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCommunityId() => $_has(1);
  @$pb.TagNumber(2)
  void clearCommunityId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get periodStart => $_getSZ(2);
  @$pb.TagNumber(3)
  set periodStart($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPeriodStart() => $_has(2);
  @$pb.TagNumber(3)
  void clearPeriodStart() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get periodEnd => $_getSZ(3);
  @$pb.TagNumber(4)
  set periodEnd($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPeriodEnd() => $_has(3);
  @$pb.TagNumber(4)
  void clearPeriodEnd() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.int get subscribers => $_getIZ(4);
  @$pb.TagNumber(5)
  set subscribers($core.int value) => $_setUnsignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasSubscribers() => $_has(4);
  @$pb.TagNumber(5)
  void clearSubscribers() => $_clearField(5);
}

class PublishedContentMetric extends $pb.GeneratedMessage {
  factory PublishedContentMetric({
    $core.String? id,
    $core.String? publishedContentId,
    $core.String? periodStart,
    $core.String? periodEnd,
    $core.int? archivers,
    $fixnum.Int64? bytes,
    $fixnum.Int64? revenue,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (publishedContentId != null)
      result.publishedContentId = publishedContentId;
    if (periodStart != null) result.periodStart = periodStart;
    if (periodEnd != null) result.periodEnd = periodEnd;
    if (archivers != null) result.archivers = archivers;
    if (bytes != null) result.bytes = bytes;
    if (revenue != null) result.revenue = revenue;
    return result;
  }

  PublishedContentMetric._();

  factory PublishedContentMetric.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishedContentMetric.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishedContentMetric',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'published_content_id')
    ..aOS(3, _omitFieldNames ? '' : 'period_start')
    ..aOS(4, _omitFieldNames ? '' : 'period_end')
    ..aI(5, _omitFieldNames ? '' : 'archivers', fieldType: $pb.PbFieldType.OU3)
    ..aInt64(6, _omitFieldNames ? '' : 'bytes')
    ..aInt64(7, _omitFieldNames ? '' : 'revenue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentMetric clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishedContentMetric copyWith(
          void Function(PublishedContentMetric) updates) =>
      super.copyWith((message) => updates(message as PublishedContentMetric))
          as PublishedContentMetric;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishedContentMetric create() => PublishedContentMetric._();
  @$core.override
  PublishedContentMetric createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishedContentMetric getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishedContentMetric>(create);
  static PublishedContentMetric? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get publishedContentId => $_getSZ(1);
  @$pb.TagNumber(2)
  set publishedContentId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPublishedContentId() => $_has(1);
  @$pb.TagNumber(2)
  void clearPublishedContentId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get periodStart => $_getSZ(2);
  @$pb.TagNumber(3)
  set periodStart($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPeriodStart() => $_has(2);
  @$pb.TagNumber(3)
  void clearPeriodStart() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get periodEnd => $_getSZ(3);
  @$pb.TagNumber(4)
  set periodEnd($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPeriodEnd() => $_has(3);
  @$pb.TagNumber(4)
  void clearPeriodEnd() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.int get archivers => $_getIZ(4);
  @$pb.TagNumber(5)
  set archivers($core.int value) => $_setUnsignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasArchivers() => $_has(4);
  @$pb.TagNumber(5)
  void clearArchivers() => $_clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get bytes => $_getI64(5);
  @$pb.TagNumber(6)
  set bytes($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasBytes() => $_has(5);
  @$pb.TagNumber(6)
  void clearBytes() => $_clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get revenue => $_getI64(6);
  @$pb.TagNumber(7)
  set revenue($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasRevenue() => $_has(6);
  @$pb.TagNumber(7)
  void clearRevenue() => $_clearField(7);
}

class CommunityMetricsRequest extends $pb.GeneratedMessage {
  factory CommunityMetricsRequest({
    $core.String? communityId,
    $core.String? period,
    $core.String? startDate,
    $core.String? endDate,
  }) {
    final result = create();
    if (communityId != null) result.communityId = communityId;
    if (period != null) result.period = period;
    if (startDate != null) result.startDate = startDate;
    if (endDate != null) result.endDate = endDate;
    return result;
  }

  CommunityMetricsRequest._();

  factory CommunityMetricsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityMetricsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityMetricsRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'community_id')
    ..aOS(2, _omitFieldNames ? '' : 'period')
    ..aOS(3, _omitFieldNames ? '' : 'start_date')
    ..aOS(4, _omitFieldNames ? '' : 'end_date')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetricsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetricsRequest copyWith(
          void Function(CommunityMetricsRequest) updates) =>
      super.copyWith((message) => updates(message as CommunityMetricsRequest))
          as CommunityMetricsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityMetricsRequest create() => CommunityMetricsRequest._();
  @$core.override
  CommunityMetricsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityMetricsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityMetricsRequest>(create);
  static CommunityMetricsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get communityId => $_getSZ(0);
  @$pb.TagNumber(1)
  set communityId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunityId() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunityId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get period => $_getSZ(1);
  @$pb.TagNumber(2)
  set period($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPeriod() => $_has(1);
  @$pb.TagNumber(2)
  void clearPeriod() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get startDate => $_getSZ(2);
  @$pb.TagNumber(3)
  set startDate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasStartDate() => $_has(2);
  @$pb.TagNumber(3)
  void clearStartDate() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get endDate => $_getSZ(3);
  @$pb.TagNumber(4)
  set endDate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasEndDate() => $_has(3);
  @$pb.TagNumber(4)
  void clearEndDate() => $_clearField(4);
}

class CommunityMetricsResponse extends $pb.GeneratedMessage {
  factory CommunityMetricsResponse({
    CommunityMetric? summary,
    $core.int? totalArchivers,
    $core.Iterable<PublishedContentMetric>? items,
  }) {
    final result = create();
    if (summary != null) result.summary = summary;
    if (totalArchivers != null) result.totalArchivers = totalArchivers;
    if (items != null) result.items.addAll(items);
    return result;
  }

  CommunityMetricsResponse._();

  factory CommunityMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CommunityMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CommunityMetricsResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOM<CommunityMetric>(1, _omitFieldNames ? '' : 'summary',
        subBuilder: CommunityMetric.create)
    ..aI(2, _omitFieldNames ? '' : 'total_archivers')
    ..pPM<PublishedContentMetric>(3, _omitFieldNames ? '' : 'items',
        subBuilder: PublishedContentMetric.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CommunityMetricsResponse copyWith(
          void Function(CommunityMetricsResponse) updates) =>
      super.copyWith((message) => updates(message as CommunityMetricsResponse))
          as CommunityMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CommunityMetricsResponse create() => CommunityMetricsResponse._();
  @$core.override
  CommunityMetricsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CommunityMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CommunityMetricsResponse>(create);
  static CommunityMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  CommunityMetric get summary => $_getN(0);
  @$pb.TagNumber(1)
  set summary(CommunityMetric value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasSummary() => $_has(0);
  @$pb.TagNumber(1)
  void clearSummary() => $_clearField(1);
  @$pb.TagNumber(1)
  CommunityMetric ensureSummary() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.int get totalArchivers => $_getIZ(1);
  @$pb.TagNumber(2)
  set totalArchivers($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTotalArchivers() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalArchivers() => $_clearField(2);

  @$pb.TagNumber(3)
  $pb.PbList<PublishedContentMetric> get items => $_getList(2);
}

class MetricsSyncRequest extends $pb.GeneratedMessage {
  factory MetricsSyncRequest({
    $core.String? communityId,
    $core.String? since,
  }) {
    final result = create();
    if (communityId != null) result.communityId = communityId;
    if (since != null) result.since = since;
    return result;
  }

  MetricsSyncRequest._();

  factory MetricsSyncRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricsSyncRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricsSyncRequest',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'community_id')
    ..aOS(2, _omitFieldNames ? '' : 'since')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncRequest copyWith(void Function(MetricsSyncRequest) updates) =>
      super.copyWith((message) => updates(message as MetricsSyncRequest))
          as MetricsSyncRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricsSyncRequest create() => MetricsSyncRequest._();
  @$core.override
  MetricsSyncRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricsSyncRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricsSyncRequest>(create);
  static MetricsSyncRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get communityId => $_getSZ(0);
  @$pb.TagNumber(1)
  set communityId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCommunityId() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommunityId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get since => $_getSZ(1);
  @$pb.TagNumber(2)
  set since($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSince() => $_has(1);
  @$pb.TagNumber(2)
  void clearSince() => $_clearField(2);
}

class MetricsSyncResponse extends $pb.GeneratedMessage {
  factory MetricsSyncResponse({
    $core.Iterable<CommunityMetric>? communityMetrics,
    $core.Iterable<PublishedContentMetric>? contentMetrics,
    $core.String? syncedAt,
    $core.bool? complete,
  }) {
    final result = create();
    if (communityMetrics != null)
      result.communityMetrics.addAll(communityMetrics);
    if (contentMetrics != null) result.contentMetrics.addAll(contentMetrics);
    if (syncedAt != null) result.syncedAt = syncedAt;
    if (complete != null) result.complete = complete;
    return result;
  }

  MetricsSyncResponse._();

  factory MetricsSyncResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricsSyncResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricsSyncResponse',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..pPM<CommunityMetric>(1, _omitFieldNames ? '' : 'community_metrics',
        subBuilder: CommunityMetric.create)
    ..pPM<PublishedContentMetric>(2, _omitFieldNames ? '' : 'content_metrics',
        subBuilder: PublishedContentMetric.create)
    ..aOS(3, _omitFieldNames ? '' : 'synced_at')
    ..aOB(4, _omitFieldNames ? '' : 'complete')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncResponse copyWith(void Function(MetricsSyncResponse) updates) =>
      super.copyWith((message) => updates(message as MetricsSyncResponse))
          as MetricsSyncResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricsSyncResponse create() => MetricsSyncResponse._();
  @$core.override
  MetricsSyncResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricsSyncResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricsSyncResponse>(create);
  static MetricsSyncResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<CommunityMetric> get communityMetrics => $_getList(0);

  @$pb.TagNumber(2)
  $pb.PbList<PublishedContentMetric> get contentMetrics => $_getList(1);

  @$pb.TagNumber(3)
  $core.String get syncedAt => $_getSZ(2);
  @$pb.TagNumber(3)
  set syncedAt($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSyncedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearSyncedAt() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get complete => $_getBF(3);
  @$pb.TagNumber(4)
  set complete($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasComplete() => $_has(3);
  @$pb.TagNumber(4)
  void clearComplete() => $_clearField(4);
}

class MetricsSyncProgress extends $pb.GeneratedMessage {
  factory MetricsSyncProgress({
    $core.String? status,
    $core.int? communityMetricsCount,
    $core.int? contentMetricsCount,
    $core.String? error,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (communityMetricsCount != null)
      result.communityMetricsCount = communityMetricsCount;
    if (contentMetricsCount != null)
      result.contentMetricsCount = contentMetricsCount;
    if (error != null) result.error = error;
    return result;
  }

  MetricsSyncProgress._();

  factory MetricsSyncProgress.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricsSyncProgress.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricsSyncProgress',
      package: const $pb.PackageName(
          _omitMessageNames ? '' : 'retrovibed.community'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'status')
    ..aI(2, _omitFieldNames ? '' : 'community_metrics_count')
    ..aI(3, _omitFieldNames ? '' : 'content_metrics_count')
    ..aOS(4, _omitFieldNames ? '' : 'error')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncProgress clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricsSyncProgress copyWith(void Function(MetricsSyncProgress) updates) =>
      super.copyWith((message) => updates(message as MetricsSyncProgress))
          as MetricsSyncProgress;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricsSyncProgress create() => MetricsSyncProgress._();
  @$core.override
  MetricsSyncProgress createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricsSyncProgress getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricsSyncProgress>(create);
  static MetricsSyncProgress? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get status => $_getSZ(0);
  @$pb.TagNumber(1)
  set status($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(1)
  void clearStatus() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get communityMetricsCount => $_getIZ(1);
  @$pb.TagNumber(2)
  set communityMetricsCount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCommunityMetricsCount() => $_has(1);
  @$pb.TagNumber(2)
  void clearCommunityMetricsCount() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get contentMetricsCount => $_getIZ(2);
  @$pb.TagNumber(3)
  set contentMetricsCount($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasContentMetricsCount() => $_has(2);
  @$pb.TagNumber(3)
  void clearContentMetricsCount() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get error => $_getSZ(3);
  @$pb.TagNumber(4)
  set error($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasError() => $_has(3);
  @$pb.TagNumber(4)
  void clearError() => $_clearField(4);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
