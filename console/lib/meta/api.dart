import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:qs_dart/qs_dart.dart' as qs;
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/retrovibed.dart' as retro;
import './meta.authn.pb.dart';
import './meta.daemon.pb.dart';
import './meta.authz.pb.dart';
import './meta.profile.pb.dart';

export './meta.daemon.pb.dart';
export './meta.authz.pb.dart';
export './meta.authn.pb.dart';
export './meta.profile.pb.dart';

Future<HttpClientResponse> healthz({String? host}) async {
  final hostport = host ?? httpx.host();
  final c = HttpClient();
  c.badCertificateCallback = (X509Certificate cert, String _host, int port) {
    final zhost = hostport.split(":").first;
    return zhost == _host;
  };

  return c
      .getUrl(Uri.https(hostport, "/healthz"))
      .then((req) {
        req.followRedirects = false;
        return req;
      })
      .then(httpx.dart_io_auto_error)
      .catchError((cause) {
        throw cause;
      });
}

Future<Authn> current({String? host}) {
  final hostport = host ?? httpx.host();
  return httpx
      .get(
        Uri.https(hostport, "/sso/"),
        options: [
          httpx.Request.authorization(httpx.auto_bearer_host(host: hostport)),
        ],
      )
      .then((v) => httpx.fromProto3JsonSafe(Authn.create(), jsonDecode(v.body)));
}

Future<Session> register(Identity iden, {String? host}) {
  print("registering ${httpx.libraryhost(host: host)} -> ${iden}");
  return httpx
      .post(
        Uri.https(
          httpx.libraryhost(host: host),
          "/sso/register",
          qs.decode(qs.encode(iden.toProto3Json())),
        ),
        options: [
          httpx.Request.bearer(
            () => Future.value(httpx.auto_bearer_host(host: host)),
          ),
        ],
      )
      .then((r) => httpx.fromProto3JsonSafe(Session.create(), jsonDecode(r.body)));
}

abstract class daemons {
  static DaemonSearchRequest request({int limit = 0}) => DaemonSearchRequest(limit: ds.Int64(limit));
  static DaemonSearchResponse response({DaemonSearchRequest? next}) =>
      DaemonSearchResponse(next: next ?? request(limit: 128), items: []);

  static bool isLocalDevice(Daemon library) =>
      library.hostname.startsWith(retro.local_device().hostname.split(":").first) ||
      library.hostname.startsWith("localhost:9998");
  static Future<DaemonSearchResponse> search(DaemonSearchRequest req) async {
    return http.Client()
        .get(
          Uri.https(
            httpx.localhost(),
            "/meta/d/",
            jsonDecode(jsonEncode(req.toProto3Json())),
          ),
          headers: {"Authorization": httpx.auto_bearer()},
        )
        .then(httpx.auto_error)
        .then((v) {
          return httpx.fromProto3JsonSafe(DaemonSearchResponse.create(), jsonDecode(v.body));
        });
  }

  static Future<DaemonCreateResponse> create(DaemonCreateRequest req) async {
    return http.Client()
        .post(
          Uri.https(httpx.localhost(), "/meta/d/"),
          headers: {"Authorization": httpx.auto_bearer()},
          body: jsonEncode(req.toProto3Json()),
        )
        .then(httpx.auto_error)
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(DaemonCreateResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<DaemonUpdateResponse> download(String id, Daemon upd) async {
    return http.Client()
        .post(
          Uri.https(httpx.localhost(), "/meta/d/${id}"),
          headers: {"Authorization": httpx.auto_bearer()},
          body: jsonEncode(
            DaemonUpdateRequest(daemon: upd..downloads = true).toProto3Json(),
          ),
        )
        .then(httpx.auto_error)
        .then((v) {
          print("SIH ${v.body}");
          return Future.value(
            httpx.fromProto3JsonSafe(DaemonUpdateResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  static Future<DaemonDisableResponse> touch(String id) async {
    final c = HttpClient();
    c.badCertificateCallback = (X509Certificate cert, String _host, int port) {
      return _host == "localhost" || _host == Platform.localHostname;
    };
    return c
        .putUrl(Uri.https(httpx.localhost(), "/meta/d/${id}"))
        .then((r) {
          r.headers.add("Authorization", httpx.auto_bearer());
          return r;
        })
        .then(httpx.dart_io_auto_error)
        .then((r) {
          return r.transform(utf8.decoder).join().then((body) {
            return httpx.fromProto3JsonSafe(DaemonDisableResponse.create(), jsonDecode(body));
          });
        });
  }

  // check if an already-persisted daemon is reachable, and touches it
  // server-side to update its heartbeat/last-seen state.
  static Future<Daemon> reachable(Daemon v) {
    return healthz(
      host: v.hostname,
    ).then((_) => authz.current(host: v.hostname)).then((_) => daemons.touch(v.id)).then((_) => v);
  }

  // check if a daemon is connectable, without touching any server-side
  // state. safe to call on a daemon that hasn't been created yet (no id).
  static Future<Daemon> connectable(Daemon v) {
    return healthz(
      host: v.hostname,
    ).then((_) => authz.current(host: v.hostname)).then((_) => v);
  }

  static Future<DaemonDisableResponse> delete(String id) async {
    return http.Client()
        .delete(
          Uri.https(httpx.localhost(), "/meta/d/${id}"),
          headers: {"Authorization": httpx.auto_bearer()},
        )
        .then(httpx.auto_error)
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(DaemonDisableResponse.create(), jsonDecode(v.body)),
          );
        });
  }

  // discover triggers a LAN mDNS scan on the connected daemon and streams
  // each discovered peer as it's found and persisted server-side.
  static Future<Stream<Daemon>> discover({
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .websocket(
          Uri.https(httpx.localhost(), "/meta/d/discover"),
          options: [httpx.Request.authorization(httpx.auto_bearer()), ...options],
        )
        .then((socket) {
          return socket.transform(
            StreamTransformer.fromHandlers(
              handleData: (data, sink) {
                if (data is List<int>) {
                  sink.add(httpx.fromProto3JsonSafe(Daemon.create(), jsonDecode(utf8.decode(data))));
                } else {
                  sink.addError('deserialization failed data: $data');
                }
              },
            ),
          );
        });
  }

  static Future<DaemonLookupResponse> latest() async {
    return http.Client()
        .get(
          Uri.https(httpx.localhost(), "/meta/d/latest"),
          headers: {"Authorization": httpx.auto_bearer()},
        )
        .then(httpx.auto_error)
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(DaemonLookupResponse.create(), jsonDecode(v.body)),
          );
        });
  }
}

abstract class profiles {
  static ProfileSearchRequest request({
    int limit = 0,
    ProfileStatus status = ProfileStatus.ENABLED,
  }) => ProfileSearchRequest(limit: ds.Int64(limit), status: status.value);
  static ProfileSearchResponse response({ProfileSearchRequest? next}) =>
      ProfileSearchResponse(next: next ?? request(limit: 128), items: []);
  static ProfileSearchResponse pending({int? limit}) => ProfileSearchResponse(
    next: request(limit: limit ?? 128, status: ProfileStatus.PENDING),
    items: [],
  );

  static Future<ProfileSearchResponse> search(
    ProfileSearchRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(
            httpx.host(),
            "/meta/u12t/",
            qs.decode(qs.encode(req.toProto3Json())),
          ),
          options: options,
        )
        .then(
          (v) => httpx.fromProto3JsonSafe(ProfileSearchResponse.create(), jsonDecode(v.body)),
        );
  }

  static Future<ProfileCreateResponse> create(
    ProfileCreateRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/meta/u12t/"),
          body: jsonEncode(req.toProto3Json()),
          options: options,
        )
        .then(
          (v) => httpx.fromProto3JsonSafe(ProfileCreateResponse.create(), jsonDecode(v.body)),
        );
  }

  static Future<ProfileLookupResponse> find(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(Uri.https(httpx.host(), "/meta/u12t/${id}"), options: options)
        .then(
          (v) => httpx.fromProto3JsonSafe(ProfileLookupResponse.create(), jsonDecode(v.body)),
        );
  }

  static Future<ProfileUpdateResponse> enable(
    Profile current, {
    List<httpx.Option> options = const [],
  }) async {
    final upd = current
      ..disabledAt = timex.inf.toIso8601String()
      ..disabledManuallyAt = timex.inf.toIso8601String()
      ..disabledPendingApprovalAt = timex.inf.toIso8601String();
    return httpx
        .patch(
          Uri.https(httpx.host(), "/meta/u12t/${current.id}"),
          body: jsonEncode(ProfileUpdateRequest(profile: upd).toProto3Json()),
          options: options,
        )
        .then(
          (v) => httpx.fromProto3JsonSafe(ProfileUpdateResponse.create(), jsonDecode(v.body)),
        );
  }

  static Future<ProfileUpdateResponse> update(
    ProfileUpdateRequest req, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .patch(
          Uri.https(httpx.host(), "/meta/u12t/${req.profile.id}"),
          body: jsonEncode(req.toProto3Json()),
          options: options,
        )
        .then(
          (v) => httpx.fromProto3JsonSafe(ProfileUpdateResponse.create(), jsonDecode(v.body)),
        );
  }

  static Future<ProfileDisableResponse> disable(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(Uri.https(httpx.host(), "/meta/u12t/${id}"), options: options)
        .then(
          (v) => httpx.fromProto3JsonSafe(ProfileDisableResponse.create(), jsonDecode(v.body)),
        );
  }
}

abstract class authz {
  static Future<AuthzResponse> current({String? host}) {
    final hostport = host ?? httpx.host();
    final c = HttpClient();
    c.badCertificateCallback = (X509Certificate cert, String _host, int port) {
      final zhost = hostport.split(":").first;
      return zhost == _host;
    };
    return c
        .getUrl(Uri.https(hostport, "/meta/authz/"))
        .then((r) {
          r.headers.add(
            "Authorization",
            httpx.auto_bearer_host(host: hostport),
          );
          return r;
        })
        .then(httpx.dart_io_auto_error)
        .then((r) {
          return r.transform(utf8.decoder).join().then((body) {
            return httpx.fromProto3JsonSafe(AuthzResponse.create(), jsonDecode(body));
          });
        });
  }

  static Future<AuthzProfileResponse> get(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(Uri.https(httpx.host(), "/meta/authz/${id}"), options: options)
        .then(
          (v) => httpx.fromProto3JsonSafe(AuthzProfileResponse.create(), jsonDecode(v.body)),
        );
  }

  static Future<AuthzGrantResponse> grant(
    String id,
    Token token, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/meta/authz/${id}"),
          body: jsonEncode(AuthzGrantRequest(token: token).toProto3Json()),
          options: options,
        )
        .then(
          (v) => httpx.fromProto3JsonSafe(AuthzGrantResponse.create(), jsonDecode(v.body)),
        );
  }

  static Future<AuthzRevokeResponse> revoke(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .delete(Uri.https(httpx.host(), "/meta/authz/${id}"), options: options)
        .then(
          (v) => httpx.fromProto3JsonSafe(AuthzRevokeResponse.create(), jsonDecode(v.body)),
        );
  }
}
