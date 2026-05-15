import 'dart:async';
import 'dart:convert';
import 'dart:typed_data';
import 'package:synchronized/synchronized.dart';
import 'package:lru/lru.dart';
import 'package:qs_dart/qs_dart.dart' as qs;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:protobuf/protobuf.dart' as pb;
import 'package:retrovibed/design.kit/bytesx.dart';
import './meta.billing.pb.dart';

export './meta.billing.pb.dart';

final _lock = Lock();

final _cache = LruTypedDataCache<String, Uint8List>(
  capacity: 1,
  capacityInBytes: bytesx.MiB,
);

Future<T> cached<T extends pb.GeneratedMessage>(
  String id,
  Future<T> Function() fetch,
  T Function(List<int> bytes) factory,
) {
  return _lock.synchronized(() {
    final c = _cache[id];
    return c == null
        ? fetch().then((v) {
          _cache[id] = v.writeToBuffer();
          return v;
        })
        : Future.value(factory(c));
  });
}

Future<BillingLookupResponse> lookup({List<httpx.Option> options = const []}) {
  return httpx.get(Uri.https(httpx.metaendpoint(), "/m/b/"), options: options).then((v) {
    return BillingLookupResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
  });
}

Future<BillingCreateResponse> create({List<httpx.Option> options = const []}) {
  return httpx.post(Uri.https(httpx.metaendpoint(), "/m/b/new"), options: options).then((v) {
    return BillingCreateResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
  });
}

Future<BillingSessionResponse> session(
  String plan, {
  List<httpx.Option> options = const [],
}) {
  return httpx
      .get(
        Uri.https(
          httpx.metaendpoint(),
          "/m/b/session",
          qs.decode(
            qs.encode(BillingSessionRequest(plan: plan).toProto3Json()),
          ),
        ),
        options: options,
      )
      .then((v) {
        return BillingSessionResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
      });
}

Future<BillingPlansResponse> plans({List<httpx.Option> options = const []}) {
  return httpx.get(Uri.https(httpx.metaendpoint(), "/m/b/plans"), options: options).then((v) {
    return BillingPlansResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
  });
}

Future<AttributionTokenResponse> attribution({List<httpx.Option> options = const []}) {
  return httpx.get(Uri.https(httpx.metaendpoint(), "/m/b/attribution"), options: options).then((v) {
    return AttributionTokenResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
  });
}

Future<AttributionConsumeResponse> consumeAttribution(
  String token, {
  List<httpx.Option> options = const [],
}) {
  if (token.isEmpty) return Future.value(AttributionConsumeResponse());
  return httpx
      .post(
        Uri.https(httpx.metaendpoint(), "/m/b/attribution"),
        body: jsonEncode({'token': token}),
        options: [httpx.Content.json, ...options],
      )
      .then((v) {
        return AttributionConsumeResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
      });
}

Future<void> delete({List<httpx.Option> options = const []}) {
  return httpx.delete(Uri.https(httpx.metaendpoint(), "/m/b/"), options: options).then((_) {});
}

Future<BillingSubscribeResponse> subscribe(String plan) {
  return httpx
      .get(
        Uri.https(
          httpx.metaendpoint(),
          "/m/b/subscribe",
          qs.decode(
            qs.encode(BillingSubscribeRequest(plan: plan).toProto3Json()),
          ),
        ),
      )
      .then((v) {
        return BillingSubscribeResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
      });
}
