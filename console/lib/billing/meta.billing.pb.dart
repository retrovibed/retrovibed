// This is a generated file - do not edit.
//
// Generated from billing/meta.billing.proto.

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

class Billing extends $pb.GeneratedMessage {
  factory Billing({
    $core.String? customerId,
    $core.String? planId,
    $core.String? subscriptionId,
    $core.String? subscriptionEndedAt,
  }) {
    final result = Billing._();
    if (customerId != null) result.customerId = customerId;
    if (planId != null) result.planId = planId;
    if (subscriptionId != null) result.subscriptionId = subscriptionId;
    if (subscriptionEndedAt != null)
      result.subscriptionEndedAt = subscriptionEndedAt;
    return result;
  }

  Billing._();

  factory Billing.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Billing()..mergeFromBuffer(data, registry);
  factory Billing.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Billing()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Billing',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: Billing.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'customer_id')
    ..aOS(2, _omitFieldNames ? '' : 'plan_id')
    ..aOS(3, _omitFieldNames ? '' : 'subscription_id')
    ..aOS(4, _omitFieldNames ? '' : 'subscription_ended_at')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Billing clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Billing copyWith(void Function(Billing) updates) =>
      super.copyWith((message) => updates(message as Billing)) as Billing;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Billing() / Billing.new instead')
  static Billing create() => Billing._();
  static $pb.GeneratedMessage $_createMessage() => Billing._();
  @$core.override
  Billing createEmptyInstance() => Billing._();
  @$core.pragma('dart2js:noInline')
  static Billing getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Billing>(Billing.$_createMessage);
  static Billing? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get customerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set customerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCustomerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearCustomerId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get planId => $_getSZ(1);
  @$pb.TagNumber(2)
  set planId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPlanId() => $_has(1);
  @$pb.TagNumber(2)
  void clearPlanId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get subscriptionId => $_getSZ(2);
  @$pb.TagNumber(3)
  set subscriptionId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSubscriptionId() => $_has(2);
  @$pb.TagNumber(3)
  void clearSubscriptionId() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get subscriptionEndedAt => $_getSZ(3);
  @$pb.TagNumber(4)
  set subscriptionEndedAt($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasSubscriptionEndedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearSubscriptionEndedAt() => $_clearField(4);
}

class Plan extends $pb.GeneratedMessage {
  factory Plan({
    $core.String? token,
    $core.String? id,
    $core.bool? legacy,
    $core.bool? hidden,
    $fixnum.Int64? profiles,
    $fixnum.Int64? storage,
    $fixnum.Int64? bandwidth,
    $core.bool? mobile,
    $core.String? stripeId,
  }) {
    final result = Plan._();
    if (token != null) result.token = token;
    if (id != null) result.id = id;
    if (legacy != null) result.legacy = legacy;
    if (hidden != null) result.hidden = hidden;
    if (profiles != null) result.profiles = profiles;
    if (storage != null) result.storage = storage;
    if (bandwidth != null) result.bandwidth = bandwidth;
    if (mobile != null) result.mobile = mobile;
    if (stripeId != null) result.stripeId = stripeId;
    return result;
  }

  Plan._();

  factory Plan.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Plan()..mergeFromBuffer(data, registry);
  factory Plan.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Plan()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Plan',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: Plan.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'token')
    ..aOS(2, _omitFieldNames ? '' : 'id')
    ..aOB(3, _omitFieldNames ? '' : 'legacy')
    ..aOB(4, _omitFieldNames ? '' : 'hidden')
    ..aInt64(1000, _omitFieldNames ? '' : 'profiles')
    ..aInt64(1001, _omitFieldNames ? '' : 'storage')
    ..aInt64(1002, _omitFieldNames ? '' : 'bandwidth')
    ..aOB(1003, _omitFieldNames ? '' : 'mobile')
    ..aOS(2000, _omitFieldNames ? '' : 'stripe_id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Plan clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Plan copyWith(void Function(Plan) updates) =>
      super.copyWith((message) => updates(message as Plan)) as Plan;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Plan() / Plan.new instead')
  static Plan create() => Plan._();
  static $pb.GeneratedMessage $_createMessage() => Plan._();
  @$core.override
  Plan createEmptyInstance() => Plan._();
  @$core.pragma('dart2js:noInline')
  static Plan getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Plan>(Plan.$_createMessage);
  static Plan? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get token => $_getSZ(0);
  @$pb.TagNumber(1)
  set token($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(2)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(2)
  void clearId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.bool get legacy => $_getBF(2);
  @$pb.TagNumber(3)
  set legacy($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLegacy() => $_has(2);
  @$pb.TagNumber(3)
  void clearLegacy() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get hidden => $_getBF(3);
  @$pb.TagNumber(4)
  set hidden($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasHidden() => $_has(3);
  @$pb.TagNumber(4)
  void clearHidden() => $_clearField(4);

  @$pb.TagNumber(1000)
  $fixnum.Int64 get profiles => $_getI64(4);
  @$pb.TagNumber(1000)
  set profiles($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(1000)
  $core.bool hasProfiles() => $_has(4);
  @$pb.TagNumber(1000)
  void clearProfiles() => $_clearField(1000);

  @$pb.TagNumber(1001)
  $fixnum.Int64 get storage => $_getI64(5);
  @$pb.TagNumber(1001)
  set storage($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(1001)
  $core.bool hasStorage() => $_has(5);
  @$pb.TagNumber(1001)
  void clearStorage() => $_clearField(1001);

  @$pb.TagNumber(1002)
  $fixnum.Int64 get bandwidth => $_getI64(6);
  @$pb.TagNumber(1002)
  set bandwidth($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(1002)
  $core.bool hasBandwidth() => $_has(6);
  @$pb.TagNumber(1002)
  void clearBandwidth() => $_clearField(1002);

  @$pb.TagNumber(1003)
  $core.bool get mobile => $_getBF(7);
  @$pb.TagNumber(1003)
  set mobile($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(1003)
  $core.bool hasMobile() => $_has(7);
  @$pb.TagNumber(1003)
  void clearMobile() => $_clearField(1003);

  @$pb.TagNumber(2000)
  $core.String get stripeId => $_getSZ(8);
  @$pb.TagNumber(2000)
  set stripeId($core.String value) => $_setString(8, value);
  @$pb.TagNumber(2000)
  $core.bool hasStripeId() => $_has(8);
  @$pb.TagNumber(2000)
  void clearStripeId() => $_clearField(2000);
}

class BillingCreateRequest extends $pb.GeneratedMessage {
  factory BillingCreateRequest() => BillingCreateRequest._();

  BillingCreateRequest._();

  factory BillingCreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingCreateRequest()..mergeFromBuffer(data, registry);
  factory BillingCreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingCreateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingCreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingCreateRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingCreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingCreateRequest copyWith(void Function(BillingCreateRequest) updates) =>
      super.copyWith((message) => updates(message as BillingCreateRequest))
          as BillingCreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use BillingCreateRequest() / BillingCreateRequest.new instead')
  static BillingCreateRequest create() => BillingCreateRequest._();
  static $pb.GeneratedMessage $_createMessage() => BillingCreateRequest._();
  @$core.override
  BillingCreateRequest createEmptyInstance() => BillingCreateRequest._();
  @$core.pragma('dart2js:noInline')
  static BillingCreateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingCreateRequest>(
          BillingCreateRequest.$_createMessage);
  static BillingCreateRequest? _defaultInstance;
}

class BillingCreateResponse extends $pb.GeneratedMessage {
  factory BillingCreateResponse({
    Billing? billing,
  }) {
    final result = BillingCreateResponse._();
    if (billing != null) result.billing = billing;
    return result;
  }

  BillingCreateResponse._();

  factory BillingCreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingCreateResponse()..mergeFromBuffer(data, registry);
  factory BillingCreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingCreateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingCreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingCreateResponse.$_createMessage)
    ..aOM<Billing>(1, _omitFieldNames ? '' : 'billing',
        subBuilder: Billing.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingCreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingCreateResponse copyWith(
          void Function(BillingCreateResponse) updates) =>
      super.copyWith((message) => updates(message as BillingCreateResponse))
          as BillingCreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use BillingCreateResponse() / BillingCreateResponse.new instead')
  static BillingCreateResponse create() => BillingCreateResponse._();
  static $pb.GeneratedMessage $_createMessage() => BillingCreateResponse._();
  @$core.override
  BillingCreateResponse createEmptyInstance() => BillingCreateResponse._();
  @$core.pragma('dart2js:noInline')
  static BillingCreateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingCreateResponse>(
          BillingCreateResponse.$_createMessage);
  static BillingCreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Billing get billing => $_getN(0);
  @$pb.TagNumber(1)
  set billing(Billing value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasBilling() => $_has(0);
  @$pb.TagNumber(1)
  void clearBilling() => $_clearField(1);
  @$pb.TagNumber(1)
  Billing ensureBilling() => $_ensure(0);
}

class BillingLookupRequest extends $pb.GeneratedMessage {
  factory BillingLookupRequest() => BillingLookupRequest._();

  BillingLookupRequest._();

  factory BillingLookupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingLookupRequest()..mergeFromBuffer(data, registry);
  factory BillingLookupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingLookupRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingLookupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingLookupRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingLookupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingLookupRequest copyWith(void Function(BillingLookupRequest) updates) =>
      super.copyWith((message) => updates(message as BillingLookupRequest))
          as BillingLookupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use BillingLookupRequest() / BillingLookupRequest.new instead')
  static BillingLookupRequest create() => BillingLookupRequest._();
  static $pb.GeneratedMessage $_createMessage() => BillingLookupRequest._();
  @$core.override
  BillingLookupRequest createEmptyInstance() => BillingLookupRequest._();
  @$core.pragma('dart2js:noInline')
  static BillingLookupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingLookupRequest>(
          BillingLookupRequest.$_createMessage);
  static BillingLookupRequest? _defaultInstance;
}

class BillingLookupResponse extends $pb.GeneratedMessage {
  factory BillingLookupResponse({
    Billing? billing,
    $fixnum.Int64? attributionCount,
    $core.int? attributionRate,
    Plan? plan,
  }) {
    final result = BillingLookupResponse._();
    if (billing != null) result.billing = billing;
    if (attributionCount != null) result.attributionCount = attributionCount;
    if (attributionRate != null) result.attributionRate = attributionRate;
    if (plan != null) result.plan = plan;
    return result;
  }

  BillingLookupResponse._();

  factory BillingLookupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingLookupResponse()..mergeFromBuffer(data, registry);
  factory BillingLookupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingLookupResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingLookupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingLookupResponse.$_createMessage)
    ..aOM<Billing>(1, _omitFieldNames ? '' : 'billing',
        subBuilder: Billing.$_createMessage)
    ..aInt64(2, _omitFieldNames ? '' : 'attribution_count')
    ..aI(3, _omitFieldNames ? '' : 'attribution_rate')
    ..aOM<Plan>(4, _omitFieldNames ? '' : 'plan',
        subBuilder: Plan.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingLookupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingLookupResponse copyWith(
          void Function(BillingLookupResponse) updates) =>
      super.copyWith((message) => updates(message as BillingLookupResponse))
          as BillingLookupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use BillingLookupResponse() / BillingLookupResponse.new instead')
  static BillingLookupResponse create() => BillingLookupResponse._();
  static $pb.GeneratedMessage $_createMessage() => BillingLookupResponse._();
  @$core.override
  BillingLookupResponse createEmptyInstance() => BillingLookupResponse._();
  @$core.pragma('dart2js:noInline')
  static BillingLookupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingLookupResponse>(
          BillingLookupResponse.$_createMessage);
  static BillingLookupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Billing get billing => $_getN(0);
  @$pb.TagNumber(1)
  set billing(Billing value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasBilling() => $_has(0);
  @$pb.TagNumber(1)
  void clearBilling() => $_clearField(1);
  @$pb.TagNumber(1)
  Billing ensureBilling() => $_ensure(0);

  @$pb.TagNumber(2)
  $fixnum.Int64 get attributionCount => $_getI64(1);
  @$pb.TagNumber(2)
  set attributionCount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAttributionCount() => $_has(1);
  @$pb.TagNumber(2)
  void clearAttributionCount() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get attributionRate => $_getIZ(2);
  @$pb.TagNumber(3)
  set attributionRate($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAttributionRate() => $_has(2);
  @$pb.TagNumber(3)
  void clearAttributionRate() => $_clearField(3);

  @$pb.TagNumber(4)
  Plan get plan => $_getN(3);
  @$pb.TagNumber(4)
  set plan(Plan value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasPlan() => $_has(3);
  @$pb.TagNumber(4)
  void clearPlan() => $_clearField(4);
  @$pb.TagNumber(4)
  Plan ensurePlan() => $_ensure(3);
}

class BillingSubscribeRequest extends $pb.GeneratedMessage {
  factory BillingSubscribeRequest({
    $core.String? plan,
  }) {
    final result = BillingSubscribeRequest._();
    if (plan != null) result.plan = plan;
    return result;
  }

  BillingSubscribeRequest._();

  factory BillingSubscribeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingSubscribeRequest()..mergeFromBuffer(data, registry);
  factory BillingSubscribeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingSubscribeRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingSubscribeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingSubscribeRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'plan')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingSubscribeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingSubscribeRequest copyWith(
          void Function(BillingSubscribeRequest) updates) =>
      super.copyWith((message) => updates(message as BillingSubscribeRequest))
          as BillingSubscribeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use BillingSubscribeRequest() / BillingSubscribeRequest.new instead')
  static BillingSubscribeRequest create() => BillingSubscribeRequest._();
  static $pb.GeneratedMessage $_createMessage() => BillingSubscribeRequest._();
  @$core.override
  BillingSubscribeRequest createEmptyInstance() => BillingSubscribeRequest._();
  @$core.pragma('dart2js:noInline')
  static BillingSubscribeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingSubscribeRequest>(
          BillingSubscribeRequest.$_createMessage);
  static BillingSubscribeRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get plan => $_getSZ(0);
  @$pb.TagNumber(1)
  set plan($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPlan() => $_has(0);
  @$pb.TagNumber(1)
  void clearPlan() => $_clearField(1);
}

class BillingSubscribeResponse extends $pb.GeneratedMessage {
  factory BillingSubscribeResponse({
    Billing? billing,
  }) {
    final result = BillingSubscribeResponse._();
    if (billing != null) result.billing = billing;
    return result;
  }

  BillingSubscribeResponse._();

  factory BillingSubscribeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingSubscribeResponse()..mergeFromBuffer(data, registry);
  factory BillingSubscribeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingSubscribeResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingSubscribeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingSubscribeResponse.$_createMessage)
    ..aOM<Billing>(1, _omitFieldNames ? '' : 'billing',
        subBuilder: Billing.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingSubscribeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingSubscribeResponse copyWith(
          void Function(BillingSubscribeResponse) updates) =>
      super.copyWith((message) => updates(message as BillingSubscribeResponse))
          as BillingSubscribeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use BillingSubscribeResponse() / BillingSubscribeResponse.new instead')
  static BillingSubscribeResponse create() => BillingSubscribeResponse._();
  static $pb.GeneratedMessage $_createMessage() => BillingSubscribeResponse._();
  @$core.override
  BillingSubscribeResponse createEmptyInstance() =>
      BillingSubscribeResponse._();
  @$core.pragma('dart2js:noInline')
  static BillingSubscribeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingSubscribeResponse>(
          BillingSubscribeResponse.$_createMessage);
  static BillingSubscribeResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Billing get billing => $_getN(0);
  @$pb.TagNumber(1)
  set billing(Billing value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasBilling() => $_has(0);
  @$pb.TagNumber(1)
  void clearBilling() => $_clearField(1);
  @$pb.TagNumber(1)
  Billing ensureBilling() => $_ensure(0);
}

class BillingSessionRequest extends $pb.GeneratedMessage {
  factory BillingSessionRequest({
    $core.String? plan,
  }) {
    final result = BillingSessionRequest._();
    if (plan != null) result.plan = plan;
    return result;
  }

  BillingSessionRequest._();

  factory BillingSessionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingSessionRequest()..mergeFromBuffer(data, registry);
  factory BillingSessionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingSessionRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingSessionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingSessionRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'plan')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingSessionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingSessionRequest copyWith(
          void Function(BillingSessionRequest) updates) =>
      super.copyWith((message) => updates(message as BillingSessionRequest))
          as BillingSessionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use BillingSessionRequest() / BillingSessionRequest.new instead')
  static BillingSessionRequest create() => BillingSessionRequest._();
  static $pb.GeneratedMessage $_createMessage() => BillingSessionRequest._();
  @$core.override
  BillingSessionRequest createEmptyInstance() => BillingSessionRequest._();
  @$core.pragma('dart2js:noInline')
  static BillingSessionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingSessionRequest>(
          BillingSessionRequest.$_createMessage);
  static BillingSessionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get plan => $_getSZ(0);
  @$pb.TagNumber(1)
  set plan($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPlan() => $_has(0);
  @$pb.TagNumber(1)
  void clearPlan() => $_clearField(1);
}

class BillingSessionResponse extends $pb.GeneratedMessage {
  factory BillingSessionResponse({
    $core.String? redirect,
  }) {
    final result = BillingSessionResponse._();
    if (redirect != null) result.redirect = redirect;
    return result;
  }

  BillingSessionResponse._();

  factory BillingSessionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingSessionResponse()..mergeFromBuffer(data, registry);
  factory BillingSessionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingSessionResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingSessionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingSessionResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'redirect')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingSessionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingSessionResponse copyWith(
          void Function(BillingSessionResponse) updates) =>
      super.copyWith((message) => updates(message as BillingSessionResponse))
          as BillingSessionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use BillingSessionResponse() / BillingSessionResponse.new instead')
  static BillingSessionResponse create() => BillingSessionResponse._();
  static $pb.GeneratedMessage $_createMessage() => BillingSessionResponse._();
  @$core.override
  BillingSessionResponse createEmptyInstance() => BillingSessionResponse._();
  @$core.pragma('dart2js:noInline')
  static BillingSessionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingSessionResponse>(
          BillingSessionResponse.$_createMessage);
  static BillingSessionResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get redirect => $_getSZ(0);
  @$pb.TagNumber(1)
  set redirect($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasRedirect() => $_has(0);
  @$pb.TagNumber(1)
  void clearRedirect() => $_clearField(1);
}

class BillingPlansRequest extends $pb.GeneratedMessage {
  factory BillingPlansRequest() => BillingPlansRequest._();

  BillingPlansRequest._();

  factory BillingPlansRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingPlansRequest()..mergeFromBuffer(data, registry);
  factory BillingPlansRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingPlansRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingPlansRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingPlansRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingPlansRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingPlansRequest copyWith(void Function(BillingPlansRequest) updates) =>
      super.copyWith((message) => updates(message as BillingPlansRequest))
          as BillingPlansRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use BillingPlansRequest() / BillingPlansRequest.new instead')
  static BillingPlansRequest create() => BillingPlansRequest._();
  static $pb.GeneratedMessage $_createMessage() => BillingPlansRequest._();
  @$core.override
  BillingPlansRequest createEmptyInstance() => BillingPlansRequest._();
  @$core.pragma('dart2js:noInline')
  static BillingPlansRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingPlansRequest>(
          BillingPlansRequest.$_createMessage);
  static BillingPlansRequest? _defaultInstance;
}

class BillingPlansResponse extends $pb.GeneratedMessage {
  factory BillingPlansResponse({
    $core.Iterable<Plan>? plans,
  }) {
    final result = BillingPlansResponse._();
    if (plans != null) result.plans.addAll(plans);
    return result;
  }

  BillingPlansResponse._();

  factory BillingPlansResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingPlansResponse()..mergeFromBuffer(data, registry);
  factory BillingPlansResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BillingPlansResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingPlansResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: BillingPlansResponse.$_createMessage)
    ..pPM<Plan>(1, _omitFieldNames ? '' : 'plans',
        subBuilder: Plan.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingPlansResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingPlansResponse copyWith(void Function(BillingPlansResponse) updates) =>
      super.copyWith((message) => updates(message as BillingPlansResponse))
          as BillingPlansResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use BillingPlansResponse() / BillingPlansResponse.new instead')
  static BillingPlansResponse create() => BillingPlansResponse._();
  static $pb.GeneratedMessage $_createMessage() => BillingPlansResponse._();
  @$core.override
  BillingPlansResponse createEmptyInstance() => BillingPlansResponse._();
  @$core.pragma('dart2js:noInline')
  static BillingPlansResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingPlansResponse>(
          BillingPlansResponse.$_createMessage);
  static BillingPlansResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<Plan> get plans => $_getList(0);
}

class AttributionTokenResponse extends $pb.GeneratedMessage {
  factory AttributionTokenResponse({
    $core.String? token,
  }) {
    final result = AttributionTokenResponse._();
    if (token != null) result.token = token;
    return result;
  }

  AttributionTokenResponse._();

  factory AttributionTokenResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AttributionTokenResponse()..mergeFromBuffer(data, registry);
  factory AttributionTokenResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AttributionTokenResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AttributionTokenResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: AttributionTokenResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'token')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributionTokenResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributionTokenResponse copyWith(
          void Function(AttributionTokenResponse) updates) =>
      super.copyWith((message) => updates(message as AttributionTokenResponse))
          as AttributionTokenResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AttributionTokenResponse() / AttributionTokenResponse.new instead')
  static AttributionTokenResponse create() => AttributionTokenResponse._();
  static $pb.GeneratedMessage $_createMessage() => AttributionTokenResponse._();
  @$core.override
  AttributionTokenResponse createEmptyInstance() =>
      AttributionTokenResponse._();
  @$core.pragma('dart2js:noInline')
  static AttributionTokenResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AttributionTokenResponse>(
          AttributionTokenResponse.$_createMessage);
  static AttributionTokenResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get token => $_getSZ(0);
  @$pb.TagNumber(1)
  set token($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);
}

class AttributionConsumeRequest extends $pb.GeneratedMessage {
  factory AttributionConsumeRequest({
    $core.String? token,
  }) {
    final result = AttributionConsumeRequest._();
    if (token != null) result.token = token;
    return result;
  }

  AttributionConsumeRequest._();

  factory AttributionConsumeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AttributionConsumeRequest()..mergeFromBuffer(data, registry);
  factory AttributionConsumeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AttributionConsumeRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AttributionConsumeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: AttributionConsumeRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'token')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributionConsumeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributionConsumeRequest copyWith(
          void Function(AttributionConsumeRequest) updates) =>
      super.copyWith((message) => updates(message as AttributionConsumeRequest))
          as AttributionConsumeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AttributionConsumeRequest() / AttributionConsumeRequest.new instead')
  static AttributionConsumeRequest create() => AttributionConsumeRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      AttributionConsumeRequest._();
  @$core.override
  AttributionConsumeRequest createEmptyInstance() =>
      AttributionConsumeRequest._();
  @$core.pragma('dart2js:noInline')
  static AttributionConsumeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AttributionConsumeRequest>(
          AttributionConsumeRequest.$_createMessage);
  static AttributionConsumeRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get token => $_getSZ(0);
  @$pb.TagNumber(1)
  set token($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);
}

class AttributionConsumeResponse extends $pb.GeneratedMessage {
  factory AttributionConsumeResponse({
    $core.String? attributionId,
  }) {
    final result = AttributionConsumeResponse._();
    if (attributionId != null) result.attributionId = attributionId;
    return result;
  }

  AttributionConsumeResponse._();

  factory AttributionConsumeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AttributionConsumeResponse()..mergeFromBuffer(data, registry);
  factory AttributionConsumeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AttributionConsumeResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AttributionConsumeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'meta'),
      createEmptyInstance: AttributionConsumeResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'attribution_id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributionConsumeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributionConsumeResponse copyWith(
          void Function(AttributionConsumeResponse) updates) =>
      super.copyWith(
              (message) => updates(message as AttributionConsumeResponse))
          as AttributionConsumeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AttributionConsumeResponse() / AttributionConsumeResponse.new instead')
  static AttributionConsumeResponse create() => AttributionConsumeResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      AttributionConsumeResponse._();
  @$core.override
  AttributionConsumeResponse createEmptyInstance() =>
      AttributionConsumeResponse._();
  @$core.pragma('dart2js:noInline')
  static AttributionConsumeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AttributionConsumeResponse>(
          AttributionConsumeResponse.$_createMessage);
  static AttributionConsumeResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get attributionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set attributionId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAttributionId() => $_has(0);
  @$pb.TagNumber(1)
  void clearAttributionId() => $_clearField(1);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
