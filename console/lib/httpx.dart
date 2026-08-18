import 'dart:async';
import 'dart:io';
import 'dart:core';
import 'dart:convert';
import 'dart:math';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:http_parser/http_parser.dart';
import 'package:retrovibed/retrovibed.dart' as retro;

var _host = localhost();

String host() {
  return _host;
}

String? normalizeuri(String? s) {
  if (s == null) return s;
  final uri = Uri.parse(s);
  if (uri.host.isEmpty) return s;
  return "${uri.host}:${uri.port}";
}

String libraryhost({String? host}) {
  return normalizeuri(Platform.environment["RETROVIBED_DAEMON_SOCKET"]) ?? host ?? "localhost:9998";
}

String localhost() {
  return normalizeuri(Platform.environment["RETROVIBED_DAEMON_SOCKET"]) ?? "localhost:9998";
}

String metaendpoint() {
  return normalizeuri(Platform.environment["RETROVIBED_META_ENDPOINT"]) ?? "api.retrovibe.space";
}

String consoleendpoint() {
  return normalizeuri(Platform.environment["RETROVIBED_CONSOLE_ENDPOINT"]) ?? "console.retrovibe.space";
}

void set(String uri) {
  _host = uri;
}

// return an deeppool identity token for retrovibed api.
String deeppool_oauth2_bearer() {
  final token = retro.deeppool_oauth2_bearer();
  if (token.isEmpty) return "";
  return "bearer ${token}";
}

// return an identity token for the local device.
String auto_bearer() {
  final token = retro.bearer_token();
  if (token.isEmpty) return "";
  return "bearer ${token}";
}

// return a identity token from the currently connected host.
String auto_bearer_host({String? host}) {
  final token = retro.bearer_token_host("https://${host ?? _host}");
  return token.isEmpty ? "" : "bearer ${token}";
}

// formats an already-known token as a bearer credential, or "" if empty.
String bearer(String token) => token.isEmpty ? "" : "bearer $token";

abstract class mimetypes {
  static MediaType parse(String s) {
    try {
      return MediaType.parse(s);
    } catch (e) {
      print(
        "failed to parse mimetype ${s} ${e} returning application/octet-stream",
      );
      return MediaType("application", "octet-stream");
    }
  }

  static MediaType maybe(String? s) {
    if (s == null) return MediaType("application", "octet-stream");
    return parse(s);
  }
}

Future<http.MultipartFile> uploadable(
  String path,
  String name,
  String mimetype, {
  String field = 'content',
  ValueNotifier<int>? progress,
}) async {
  // Read file size
  final fileLength = await File(path).length();

  // Create a stream with progress tracking
  final fileStream = File(path).openRead().transform(
    StreamTransformer<List<int>, List<int>>.fromHandlers(
      handleData: (chunk, sink) {
        progress?.value += chunk.length;
        sink.add(chunk.cast<int>());
      },
    ),
  );

  // Create MultipartFile with progress-wrapped stream
  return http.MultipartFile(
    field,
    fileStream,
    fileLength,
    filename: name,
    contentType: mimetypes.parse(mimetype),
  );
}

Future<T> auto_error<T extends http.BaseResponse>(T v) {
  // if (v.statusCode >= 300) {
  //   print("failed ${v.request?.url.toString()} ${v.statusCode}");
  // }

  return v.statusCode >= 300 ? Future.error(v) : Future.value(v);
}

Future<HttpClientResponse> dart_io_auto_error(HttpClientRequest v) {
  return v.close().then((r) {
    // if (r.statusCode >= 300) {
    //   print("failed ${v.uri.toString()} ${r.statusCode}");
    // }

    return r.statusCode >= 300 ? Future.error(r) : Future.value(r);
  });
}

Future<HttpClientRequest> dart_io_request(HttpClientRequest v) {
  return v.close().then((r) {
    // if (r.statusCode >= 300) {
    //   print("failed ${v.uri.toString()} ${r.statusCode}");
    // }

    return r.statusCode >= 300 ? Future.error(v) : Future.value(v);
  });
}

class ErrorsTest {
  static int statusCode(Object obj) {
    if (obj is http.Response) return obj.statusCode;
    if (obj is HttpClientResponse) return obj.statusCode;
    return 0;
  }

  // test for common errors.
  static bool httpauto(Object obj) {
    final code = statusCode(obj);
    return const [401, 403, 404, 409, 429, 502].contains(code);
  }

  static bool httpnotimplemented(Object obj) {
    final code = statusCode(obj);
    return const [404, 501].contains(code);
  }

  static bool badrequest(Object obj) {
    const code = 400;
    return (obj is http.Response && obj.statusCode == code) || (obj is HttpClientResponse && obj.statusCode == code);
  }

  static bool unauthorized(Object obj) {
    const code = 401;
    return (obj is http.Response && obj.statusCode == code) || (obj is HttpClientResponse && obj.statusCode == code);
  }

  static bool forbidden(Object obj) {
    const code = 403;
    return (obj is http.Response && obj.statusCode == code) || (obj is HttpClientResponse && obj.statusCode == code);
  }

  static bool conflict(Object obj) {
    const code = 409;
    return (obj is http.Response && obj.statusCode == code) || (obj is HttpClientResponse && obj.statusCode == code);
  }

  static bool err404(Object obj) {
    const code = 404;
    return (obj is http.Response && obj.statusCode == code) || (obj is HttpClientResponse && obj.statusCode == code);
  }

  static bool notimplemented(Object obj) {
    const code = 501;
    return (obj is http.Response && obj.statusCode == code) || (obj is HttpClientResponse && obj.statusCode == code);
  }

  static bool unavailable(Object obj) {
    const code = 503;
    return (obj is http.Response && obj.statusCode == code) || (obj is HttpClientResponse && obj.statusCode == code);
  }
}

class Content {
  static Future<Request> json(Request request) {
    request.headers["Content-Type"] = "application/json";
    return Future.value(request);
  }

  static Future<Request> urlencoded(Request request) {
    request.headers["Content-Type"] = "application/x-www-form-urlencoded";
    return Future.value(request);
  }

  static Future<Request> formdata(Request request) {
    request.headers["Content-Type"] = "multipart/form-data";
    return Future.value(request);
  }
}

class Accept {
  static Future<Request> json(Request request) {
    request.headers["Accept"] = "application/json; charset=utf-8";
    return Future.value(request);
  }
}

typedef Endpoint<X, Y> = Future<Y> Function(X v, {List<Option> options});

typedef Option = Future<Request> Function(Request request);

class MissingTokenError implements Exception {
  final String message;
  const MissingTokenError([this.message = "bearer token is empty"]);

  @override
  String toString() => "MissingTokenError: $message";
}

class Request {
  Map<String, String> headers = {};

  static Future<Request> noop(Request request) {
    return Future.value(request);
  }

  static List<Option> empty() => [];

  static Option bearer(Future<String> Function() token) {
    return (Request request) {
      return token().then((v) {
        if (v.isEmpty) throw const MissingTokenError();
        request.headers["Authorization"] = v.toLowerCase().startsWith("bearer ") ? v : "Bearer ${v}";
        return Future.value(
          request,
        ); // Returns a completed Future with the modified request
      });
    };
  }

  static Option authorization(String token) {
    return (Request request) {
      if (token.isEmpty) {
        return Future.error(
          const MissingTokenError(),
          StackTrace.current,
        );
      }

      request.headers["Authorization"] = token;
      return Future.value(
        request,
      ); // Returns a completed Future with the modified request
    };
  }

  static Option header(String key, String value) {
    return (Request request) {
      if (value.isEmpty) return Future.value(request);

      request.headers[key] = value;

      return Future.value(request);
    };
  }
}

Future<Request> request(List<Option> options) {
  return options.fold(Future.value(Request()), (Future<Request> p, Option opt) {
    return p.then((r) {
      return opt(r);
    });
  });
}

// convert dynamic objects to a map<String, String>, flattening nested maps
// using bracket notation (e.g. {"created": {"newest": "..."}} → {"created[newest]": "..."})
Map<String, String> params(Object? m) {
  final decoded = jsonDecode(jsonEncode(m)) as Map<String, dynamic>;
  final result = <String, String>{};
  void flatten(Map<String, dynamic> map, String prefix) {
    for (final entry in map.entries) {
      final key = prefix.isEmpty ? entry.key : '$prefix[${entry.key}]';
      if (entry.value is Map<String, dynamic>) {
        flatten(entry.value as Map<String, dynamic>, key);
      } else {
        result[key] = entry.value.toString();
      }
    }
  }

  flatten(decoded, '');
  return result;
}

Future<http.StreamedResponse> send(
  Uri path, {
  List<Option> options = const [],
  dynamic query = const {},
}) {
  return request(options).then((r) {
    final req = http.Request("GET", path);
    req.headers.addAll(r.headers);

    return http.Client().send(req).then(auto_error);
  });
}

Future<http.Response> get(
  Uri path, {
  List<Option> options = const [],
  dynamic query = const {},
}) {
  return request(options).then((r) {
    return http.Client().get(path, headers: r.headers).then(auto_error);
  });
}

Future<WebSocket> websocket(
  Uri path, {
  List<Option> options = const [],
  dynamic query = const {},
}) {
  // https://stackoverflow.com/questions/53721745/dart-upgrade-client-socket-to-websocket#53727270
  return request([
    Request.header("Connection", "Upgrade"),
    Request.header("Upgrade", "websocket"),
    Request.header("sec-websocket-version", "13"),
    Request.header("sec-websocket-key", _generateWebSocketKey()),
    ...options,
  ]).then((r0) {
    return HttpClient()
        .getUrl(path)
        .then((r) {
          r0.headers.forEach((k, v) => r.headers.add(k, v));
          return r.close();
        })
        .then(
          (resp) => resp.statusCode == HttpStatus.switchingProtocols ? resp.detachSocket() : Future.error(resp),
        )
        .then((s) => WebSocket.fromUpgradedSocket(s, serverSide: false));
  });
}

Future<http.Response> post(
  Uri path, {
  List<Option> options = const [],
  Object? body,
}) {
  return request(options).then((r) {
    return http.Client().post(path, headers: r.headers, body: body).then(auto_error);
  });
}

Future<http.Response> patch(
  Uri path, {
  List<Option> options = const [],
  Object? body,
}) {
  return request(options).then((r) {
    return http.Client().patch(path, headers: r.headers, body: body).then(auto_error);
  });
}

Future<http.Response> put(
  Uri path, {
  List<Option> options = const [],
  Object? body,
}) {
  return request(options).then((r) {
    return http.Client().put(path, headers: r.headers, body: body).then(auto_error);
  });
}

Future<http.Response> delete(
  Uri path, {
  List<Option> options = const [],
  Object? body,
}) {
  return request(options).then((r) {
    return http.Client().delete(path, headers: r.headers, body: body).then(auto_error);
  });
}

String _generateWebSocketKey() {
  final random = Random.secure();
  final bytes = List<int>.generate(16, (i) => random.nextInt(256));
  return base64Encode(bytes);
}

/// Represents the result of checking whether an error should trigger a retry.
/// Returns null to indicate "don't retry", or a Duration to indicate the delay before retrying.
/// [backoff] is the computed backoff duration for this attempt; checks should return it (or an
/// override, e.g. respecting a retry-after header) rather than hardcoding a delay.
typedef _RetryCheck = Duration? Function(Object error, Duration backoff);

/// Configuration for automatic retry behavior.
class _RetryConfig {
  final int maximum;
  final Duration Function(int attempt) backoff;
  final List<_RetryCheck> checks;

  const _RetryConfig({
    this.maximum = 10,
    required this.backoff,
    this.checks = const [],
  });
}

/// Backoff strategies for retries.
abstract class Backoff {
  /// Constant delay between retries.
  static Duration Function(int attempt) constant(Duration delay) {
    return (_) => delay;
  }

  /// Exponential backoff with optional max delay.
  static Duration Function(int attempt) exponential({
    Duration initial = const Duration(seconds: 1),
    Duration? max,
  }) {
    return (attempt) {
      final delay = initial * pow(2, attempt);
      if (max != null && delay > max) return max;
      return delay;
    };
  }
}

/// Error checks for determining retry behavior.
abstract class RetryChecks {
  /// Retry on rate limiting (429), respecting retry-after header.
  static Duration? ratelimited(Object error, Duration backoff) {
    const code = 429;
    if (error is http.Response && error.statusCode == code) {
      final retryAfter = error.headers['retry-after'];
      if (retryAfter != null) {
        final seconds = int.tryParse(retryAfter);
        if (seconds != null) return Duration(seconds: seconds);
      }
      return backoff;
    }
    if (error is HttpClientResponse && error.statusCode == code) {
      final retryAfter = error.headers.value('retry-after');
      if (retryAfter != null) {
        final seconds = int.tryParse(retryAfter);
        if (seconds != null) return Duration(seconds: seconds);
      }
      return backoff;
    }
    return null;
  }

  /// Retry on bad gateway (502).
  static Duration? badgateway(Object error, Duration backoff) {
    const code = 502;
    if (ErrorsTest.statusCode(error) == code) return backoff;
    return null;
  }

  /// Retry on service unavailable (503).
  static Duration? unavailable(Object error, Duration backoff) {
    const code = 503;
    if (ErrorsTest.statusCode(error) == code) return backoff;
    return null;
  }

  /// Retry on network issues (SocketException, ClientException, HandshakeException, etc).
  static Duration? networkissue(Object error, Duration backoff) {
    if (error is SocketException || error is HttpException) return backoff;
    if (error is HandshakeException) return backoff;
    if (error is http.ClientException) return backoff;
    return null;
  }

  /// Retry on unauthorized (401).
  static Duration? unauthorized(Object error, Duration backoff) {
    const code = 401;
    if (ErrorsTest.statusCode(error) == code) return backoff;
    return null;
  }

  /// Retry when a bearer token was not yet available.
  static Duration? missingtoken(Object error, Duration backoff) {
    if (error is MissingTokenError) return backoff;
    return null;
  }

  static const List<_RetryCheck> auto = [
    RetryChecks.ratelimited,
    RetryChecks.badgateway,
    RetryChecks.unavailable,
    RetryChecks.networkissue,
  ];
}

/// Creates a default auto retry configuration similar to the TypeScript version.
/// Retries on rate limiting, bad gateway, and network issues.
_RetryConfig _autoretry({
  int maxRetries = 10,
  List<_RetryCheck> checks = RetryChecks.auto,
  Duration Function(int attempt)? backoff,
}) {
  return _RetryConfig(
    maximum: maxRetries,
    backoff: backoff ?? Backoff.constant(const Duration(seconds: 1)),
    checks: checks,
  );
}

/// Wraps a future-returning function with automatic retry logic.
Future<T> withRetry<T>(
  Future<T> Function() operation, {
  int maxRetries = 10,
  Duration Function(int attempt)? backoff,
  List<_RetryCheck> checks = RetryChecks.auto,
}) async {
  final config = _autoretry(
    maxRetries: maxRetries,
    backoff: backoff,
    checks: checks,
  );

  for (var attempt = 0; true; attempt++) {
    try {
      return await operation();
    } catch (error) {
      if (attempt >= config.maximum) rethrow;
      final backoff = config.backoff(attempt);
      Duration? delay;
      for (final check in config.checks) {
        delay = check(error, backoff);
        if (delay != null) break;
      }

      if (delay == null) rethrow; // Error not retryable

      await Future.delayed(delay);
    }
  }
}
