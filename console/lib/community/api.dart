import 'dart:convert';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'community.pb.dart';
import 'community.metrics.pb.dart';
import 'community.publish.pb.dart';

export 'community.pb.dart';
export 'community.metrics.pb.dart';
export 'community.publish.pb.dart';

typedef FnSubscribe = Future<CommunitySubscribeResponse> Function(String communityId, {List<httpx.Option> options});

typedef FnPublishingSearch =
    Future<PublishedContentSearchResponse> Function(
      String id, {
      List<httpx.Option> options,
      PublishedContentSearchRequest? req,
    });

typedef FnPublishingTombstone =
    Future<PublishContentDeleteResponse> Function(
      String id, {
      List<httpx.Option> options,
    });

class publishing {
  static Future<PublishContentResponse> publish(
    String cid,
    PublishContentRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/c/p/$cid"),
          body: jsonEncode(req.toProto3Json()),
          options: [httpx.Accept.json, httpx.Content.json, ...options],
        )
        .then((v) {
          return Future.value(
            PublishContentResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<PublishedContentSearchResponse> search(
    String cid, {
    List<httpx.Option> options = const [],
    PublishedContentSearchRequest? req,
  }) async {
    req ??= PublishedContentSearchRequest();
    req.communityId = cid;
    if (req.limit.isZero) req.limit = fixnum.Int64(100);

    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/c/p/$cid",
            httpx.params(req.toProto3Json()),
          ),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            PublishedContentSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<PublishContentDeleteResponse> tombstone(
    String pid, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(
            httpx.host(),
            "/c/p/$pid",
          ),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            PublishContentDeleteResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}

class metrics {
  static Future<CommunityMetricsResponse> search(
    String id, {
    required DateTime startDate,
    required DateTime endDate,
    List<httpx.Option> options = const [],
  }) async {
    final req = CommunityMetricsRequest(
      communityId: id,
      startDate: startDate.toIso8601String(),
      endDate: endDate.toIso8601String(),
    );
    return httpx
        .get(
          Uri.https(httpx.host(), "/c/m/$id", httpx.params(req.toProto3Json())),
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
          Uri.https(httpx.host(), "/c/m/$id"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            MetricsSyncProgress.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}

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

  static Future<PublishedContentSearchResponse> published(
    String cid, {
    List<httpx.Option> options = const [],
    PublishedContentSearchRequest? req,
  }) async {
    req ??= PublishedContentSearchRequest();
    req.communityId = cid;
    req.sync = uuidx.min();
    if (req.limit.isZero) req.limit = fixnum.Int64(100);

    return httpx
        .get(
          Uri.https(
            httpx.metaendpoint(),
            "/p/sync",
            httpx.params(req.toProto3Json()),
          ),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            PublishedContentSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
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
