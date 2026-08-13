import 'dart:convert';
import 'dart:typed_data';
import 'package:retrovibed/meta/meta.search.pb.dart';
import 'package:synchronized/synchronized.dart';
import 'package:lru/lru.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/design.kit/bytesx.dart';

export 'package:retrovibed/media.dart';

typedef FnKnownSearch =
    Future<KnownSearchResponse> Function(
      KnownSearchRequest req, {
      List<httpx.Option> options,
    });

typedef FnKnownLatest =
    Future<KnownLatestResponse> Function(
      KnownLatestRequest req, {
      List<httpx.Option> options,
    });

typedef FnRecommendations =
    Future<RecommendationSearchResponse> Function(
      RecommendationSearchRequest req, {
      List<httpx.Option> options,
    });

typedef FnRecent =
    Future<RecentSearchResponse> Function(
      RecentSearchRequest req, {
      List<httpx.Option> options,
    });

typedef FnRecentTombstone =
    Future<RecentDeleteResponse> Function(
      String id, {
      List<httpx.Option> options,
    });

typedef FnLibraryMetadataSync =
    Future<MediaUpdateResponse> Function(
      String id,
      Media media, {
      List<httpx.Option> options,
    });

typedef FnDiscoveredMetadataSync =
    Future<MetadataSyncResponse> Function(
      String id,
      Media media, {
      List<httpx.Option> options,
    });

abstract class known {
  static final Lock _lock = Lock();

  static LruTypedDataCache<String, Uint8List> cache = LruTypedDataCache<String, Uint8List>(
    capacity: 256,
    capacityInBytes: bytesx.MiB,
  );

  static KnownSearchRequest request({
    int limit = 0,
    String query = "",
    String mimetype = "",
    String? language,
    bool adult = false,
    timex.Range? released,
  }) {
    released = released ?? timex.Range.everything();
    return KnownSearchRequest(
      query: query,
      mimetype: mimetype,
      language: language,
      adult: adult,
      released: DateRange(
        oldest: timex.formatISO8601(released.begin),
        newest: timex.formatISO8601(released.end),
      ),
      limit: ds.Int64(limit),
    );
  }

  static KnownSearchResponse response({KnownSearchRequest? next}) =>
      KnownSearchResponse(next: next ?? request(limit: 100), items: []);

  static Future<Known> autodetect(
    Media m, {
    List<httpx.Option> options = const [],
  }) async {
    if (mimex.icon(m.mimetype) == mimex.icoimage) {
      return Future.value(
        Known(
          id: "",
          description: m.description,
          summary: "",
          rating: 0.0,
          image: m.image,
        ),
      );
    }

    if (!uuidx.isMinMax(uuidx.fromString(m.knownMediaId))) {
      return known
          .cached(
            m.knownMediaId,
            () => known.get(
              m.knownMediaId,
              options: options,
            ),
          )
          .then((v) => v.known);
    }

    return Future.value(
      Known(
        id: "",
        description: m.description,
        summary: "",
        rating: 0.0,
        image: "",
      ),
    );
  }

  static Future<KnownSearchResponse> search(
    KnownSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/k/",
            httpx.params(req.toProto3Json()),
          ),
          options: options,
        )
        .then((v) {
          final resp = KnownSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body));
          return Future.value(resp);
        });
  }

  static Future<KnownLookupResponse> cached(
    String id,
    Future<KnownLookupResponse> Function() fetch,
  ) {
    return _lock.synchronized(() {
      final c = cache[id];
      return c == null
          ? fetch().then((v) {
              cache[id] = v.known.writeToBuffer();
              return v;
            })
          : Future.value(KnownLookupResponse(known: Known.fromBuffer(c)));
    });
  }

  static Future<KnownLookupResponse> get(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/k/${id}", {}), options: options).then((v) {
      return Future.value(
        KnownLookupResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
      );
    });
  }

  static Future<KnownCreateResponse> create(
    KnownCreateRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/k/", {}),
          options: options,
          body: jsonEncode(req.toProto3Json()),
        )
        .then((v) {
          return Future.value(
            KnownCreateResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static KnownLatestRequest latestRequest({
    String mimetype = "",
    String language = "",
    bool adult = false,
    timex.Range? released,
    int limit = 100,
  }) {
    released = released ?? timex.Range.latest(Duration(days: 360));
    return KnownLatestRequest(
      mimetype: mimetype,
      language: language,
      adult: adult,
      limit: ds.Int64(limit),
      released: DateRange(
        oldest: timex.formatISO8601(released.begin),
        newest: timex.formatISO8601(released.end),
      ),
    );
  }

  static Future<KnownLatestResponse> latest(
    KnownLatestRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/k/latest",
            httpx.params(req.toProto3Json()),
          ),
          options: [httpx.Content.urlencoded, httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            KnownLatestResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}

abstract class recommendations {
  static RecommendationSearchRequest request({
    String mimetype = "",
    String language = "",
    bool adult = false,
    int limit = 100,
    timex.Range? created,
  }) {
    created = created ?? timex.Range.latest(Duration(days: 30));
    return RecommendationSearchRequest(mimetype: mimetype, language: language, adult: adult, limit: ds.Int64(limit));
  }

  static RecommendationSearchResponse response({RecommendationSearchRequest? next}) =>
      RecommendationSearchResponse(items: []);

  static Future<RecommendationSearchResponse> latest(
    RecommendationSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/r/",
            httpx.params(req.toProto3Json()),
          ),
          options: [httpx.Content.urlencoded, httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            RecommendationSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<void> random(
    RecommendationSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(
            httpx.host(),
            "/r/random",
          ),
          options: [httpx.Content.json, httpx.Accept.json, ...options],
          body: jsonEncode(req.toProto3Json()),
        )
        .then((v) {
          return Future.value();
        });
  }

  static Future<RecommendationFindResponse> find(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/r/${id}", {}), options: options).then((v) {
      return Future.value(
        RecommendationFindResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
      );
    });
  }

  static Future<RecommendationFindResponse> content(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/r/content/${id}", {}), options: options).then((v) {
      return Future.value(
        RecommendationFindResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
      );
    });
  }

  static Future<RecommendationDeleteResponse> delete(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(
          Uri.https(
            httpx.host(),
            "/r/${id}",
          ),
          options: [httpx.Content.json, httpx.Accept.json, ...options],
          body: jsonEncode(RecommendationDeleteRequest().toProto3Json()),
        )
        .then((v) {
          return Future.value(
            RecommendationDeleteResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<RecommendationSearchResponse> refresh({
    List<httpx.Option> options = const [],
  }) {
    return httpx
        .post(
          Uri.https(httpx.host(), "/r/"),
          options: options,
          body: jsonEncode(RecommendationSearchResponse.create().toProto3Json()),
        )
        .then((v) {
          return Future.value(
            RecommendationSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}

abstract class recent {
  static RecentSearchRequest request({String mimetype = "", int limit = 100, timex.Range? created}) {
    created = created ?? timex.Range.latest(Duration(days: 30));
    return RecentSearchRequest(
      mimetype: mimetype,
      limit: ds.Int64(limit),
      created: DateRange(
        oldest: timex.formatISO8601(created.begin),
        newest: timex.formatISO8601(created.end),
      ),
    );
  }

  static RecentSearchResponse response({RecentSearchRequest? next}) =>
      RecentSearchResponse(next: next ?? request(), items: []);

  static Future<RecentSearchResponse> latest(
    RecentSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.host(), "/w/", httpx.params(req.toProto3Json())),
          options: [httpx.Content.urlencoded, httpx.Accept.json, ...options],
        )
        .then((v) {
          return Future.value(
            RecentSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<RecentRecordResponse> record(
    RecentRecordRequest req, {
    List<httpx.Option> options = const [],
  }) {
    return httpx
        .post(
          Uri.https(httpx.host(), "/w/"),
          body: jsonEncode(req.toProto3Json()),
          options: options,
        )
        .then((v) {
          return Future.value(
            RecentRecordResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<RecentDeleteResponse> delete(
    String id, {
    List<httpx.Option> options = const [],
  }) {
    return httpx
        .delete(
          Uri.https(httpx.host(), "/w/$id"),
          options: options,
        )
        .then((v) {
          return Future.value(
            RecentDeleteResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}

abstract class releases {
  static KnownLatestRequest request({int limit = 0}) => KnownLatestRequest(limit: ds.Int64(limit));
  static KnownLatestResponse response({List<Known>? items}) => KnownLatestResponse(items: items ?? <Known>[]);

  static Future<KnownLatestResponse> get(
    KnownLatestResponse req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/k/latest",
            jsonDecode(jsonEncode(req.toProto3Json())),
          ),
          options: options,
        )
        .then((v) {
          return Future.value(
            KnownLatestResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}

abstract class locate {
  static LocateSearchRequest request({int limit = 0, String query = ""}) =>
      LocateSearchRequest(limit: ds.Int64(limit), query: query);
  static LocateSearchResponse response({LocateSearchRequest? next}) =>
      LocateSearchResponse(next: next ?? request(limit: 100), items: []);

  static Future<LocateSearchResponse> search(
    LocateSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/l/",
            jsonDecode(jsonEncode(req.toProto3Json())),
          ),
          options: options,
        )
        .then((v) {
          return Future.value(
            LocateSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }

  static Future<LocateLookupResponse> get(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/l/${id}", {}), options: options).then((v) {
      return Future.value(
        LocateLookupResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
      );
    });
  }

  static Future<LocateCreateResponse> create(
    Locate req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/l/", {}),
          options: options,
          body: jsonEncode(LocateCreateRequest(locate: req).toProto3Json()),
        )
        .then((v) {
          return Future.value(
            LocateCreateResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}
