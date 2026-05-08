import 'dart:convert';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:retrovibed/httpx.dart' as httpx;
import './community.pb.dart';

export './community.pb.dart';

typedef FnSubscribe = Future<CommunitySubscribeResponse> Function(String communityId, {List<httpx.Option> options});

class API {
  static Future<CommunitySearchResponse> search(
    CommunitySearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/c/",
            jsonDecode(jsonEncode(req.toProto3Json())),
          ),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            CommunitySearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<CommunityCreateResponse> create(
    CommunityCreateRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.metaendpoint(), "/c/"),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
          body: jsonEncode(req.toProto3Json()),
        )
        .then((v) {
          return Future.value(
            CommunityCreateResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<CommunityDeleteResponse> delete(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.metaendpoint(), "/c/$id"),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
          body: jsonEncode(CommunityDeleteRequest.create().toProto3Json()),
        )
        .then((v) {
          return Future.value(
            CommunityDeleteResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<CommunityUpdateResponse> update(
    String id,
    CommunityUpdateRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .put(
          Uri.https(httpx.metaendpoint(), "/c/$id"),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
          body: jsonEncode(req.toProto3Json()),
        )
        .then((v) {
          return Future.value(
            CommunityUpdateResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<PublishContentResponse> publish(
    String cid,
    PublishContentRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/c/$cid/publish"),
          body: jsonEncode(req.toProto3Json()),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
        )
        .then((v) {
          return Future.value(
            PublishContentResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<PublishedContentListResponse> published(
    String cid, {
    List<httpx.Option> options = const [],
    int offset = 0,
    int limit = 100,
  }) async {
    final req =
        PublishedContentListRequest()
          ..communityId = cid
          ..offset = fixnum.Int64(offset)
          ..limit = fixnum.Int64(limit);
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/c/$cid/published",
            jsonDecode(jsonEncode(req.toProto3Json())),
          ),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            PublishedContentListResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<CommunityMetricsResponse> metrics(
    String id, {
    required DateTime startDate,
    required DateTime endDate,
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/c/$id/metrics", {
            'start_date': startDate.toIso8601String(),
            'end_date': endDate.toIso8601String(),
          }),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            CommunityMetricsResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<MetricsSyncProgress> sync(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/c/$id/metrics/sync"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            MetricsSyncProgress.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<CommunitySubscribeResponse> subscribe(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/c/${id}/subscribe"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            CommunitySubscribeResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}
