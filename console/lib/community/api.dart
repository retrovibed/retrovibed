import 'dart:convert';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'community.pb.dart';
import 'community.metrics.pb.dart';
import 'community.publish.pb.dart';
import 'community.social.pb.dart';

export 'community.pb.dart';
export 'community.metrics.pb.dart';
export 'community.publish.pb.dart';
export 'community.social.pb.dart';

bool isStale(Community c, {Duration threshold = const Duration(hours: 1)}) {
  return timex.now().difference(timex.iso8601(c.lastSyncAt)) > threshold;
}

typedef FnSubscribe = Future<CommunitySubscribeResponse> Function(String communityId, {List<httpx.Option> options});

typedef FnCommunitySearch =
    Future<CommunitySearchResponse> Function(
      CommunitySearchRequest req, {
      List<httpx.Option> options,
    });

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
            httpx.fromProto3JsonSafe(PublishContentResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(PublishedContentSearchResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(PublishContentDeleteResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(CommunityMetricsResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(MetricsSyncProgress.create(), jsonDecode(v.body)),
          );
        });
  }
}

typedef FnSocialsSearch =
    Future<SocialsSearchResponse> Function(
      SocialsSearchRequest req, {
      List<httpx.Option> options,
    });

typedef FnSocialsEnable =
    Future<CommunityPublisherEnableResponse> Function(
      String communityId,
      String publisherId, {
      List<httpx.Option> options,
    });

typedef FnSocialsDisable =
    Future<CommunityPublisherDisableResponse> Function(
      String communityId,
      String publisherId, {
      List<httpx.Option> options,
    });

class socials {
  static Future<SocialsSearchResponse> search(
    SocialsSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    if (req.limit.isZero) req.limit = fixnum.Int64(100);

    return httpx
        .get(
          Uri.https(httpx.host(), "/c/social/", httpx.params(req.toProto3Json())),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(SocialsSearchResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<CommunityPublisherEnableResponse> enable(
    String communityId,
    String publisherId, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/c/social/$communityId/publishers/$publisherId"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(CommunityPublisherEnableResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<CommunityPublisherDisableResponse> disable(
    String communityId,
    String publisherId, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(httpx.host(), "/c/social/$communityId/publishers/$publisherId"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(CommunityPublisherDisableResponse.create(), jsonDecode(v.body)),
          );
        });
  }
}

class communities {
  static String canonicaluri(String v) {
    if (v.startsWith("https")) {
      return v;
    }

    return "https://${v.isEmpty ? 'example' : v}.community.retrovibe.space";
  }

  static String domain(String uri) {
    final host = Uri.tryParse(uri.trim())?.host ?? "";
    if (host.isEmpty) return uri;
    if (!host.endsWith("community.retrovibe.space")) return uri;
    return host.split(".").first;
  }

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
            httpx.fromProto3JsonSafe(CommunitySearchResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(CommunityCreateResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(CommunityDeleteResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(CommunityUpdateResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(PublishedContentSearchResponse.create(), jsonDecode(v.body)),
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
            httpx.fromProto3JsonSafe(CommunitySubscribeResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<PublishedContentSearchResponse> resync(
    String id, {
    List<httpx.Option> options = const [],
    PublishedContentSearchRequest? req,
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/c/$id/resync"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(PublishedContentSearchResponse.create(), jsonDecode(v.body)),
          );
        });
  }
}

/// Configuration for one installed publisher plugin, proxied as the raw
/// bytes of its .env sidecar - the same convention the search plugin
/// endpoints use, and for the same reason: parsing and re-serializing is a
/// client concern (see envfile.dart), the server never interprets it.
///
/// A GET returns the variables the plugin itself declares, with whatever
/// has been configured filled in over the top, so the editor can render a
/// form for a plugin the console knows nothing about.
abstract class publisherenvironment {
  static Future<String> get(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/c/publishers/environment/${id}"), options: options).then((v) => v.body);
  }

  static Future<String> update(
    String id,
    String content, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/c/publishers/environment/${id}"),
          options: options,
          body: content,
        )
        .then((v) => v.body);
  }

  static Future<void> delete(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.delete(Uri.https(httpx.host(), "/c/publishers/environment/${id}"), options: options).then((_) {});
  }
}
