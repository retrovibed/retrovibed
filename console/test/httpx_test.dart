import 'dart:io';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;

void main() {
  test('normalizeuri', () {
    expect(httpx.normalizeuri(null), null);
    expect(httpx.normalizeuri("http://example.com"), "example.com:80");
    expect(httpx.normalizeuri("https://example.com"), "example.com:443");
    expect(httpx.normalizeuri("example.com"), "example.com");
    expect(httpx.normalizeuri("example.com:1000"), "example.com:1000");
  });

  group('params', () {
    test('flat object', () {
      expect(
        httpx.params({'limit': 100}),
        equals({'limit': '100'}),
      );
    });

    test('nested object uses bracket notation', () {
      expect(
        httpx.params({
          'created': {
            'newest': '2026-03-18T08:22:28.250090047Z',
            'oldest': '1970-01-01T00:00:00Z',
          },
          'limit': 100,
        }),
        equals({
          'created[newest]': '2026-03-18T08:22:28.250090047Z',
          'created[oldest]': '1970-01-01T00:00:00Z',
          'limit': '100',
        }),
      );
    });

    test('bool and null values', () {
      expect(
        httpx.params({'enabled': true, 'count': 0}),
        equals({'enabled': 'true', 'count': '0'}),
      );
    });
  });

  group('withRetry', () {
    test('returns value on immediate success', () async {
      final result = await httpx.withRetry(() async => 42);
      expect(result, equals(42));
    });

    test('retries on SocketException and succeeds', () async {
      int attempts = 0;
      final result = await httpx.withRetry(
        () async {
          attempts++;
          if (attempts < 3) {
            throw SocketException(
              'Connection refused',
              osError: OSError('Connection refused', 111),
            );
          }
          return 'ok';
        },
        maxRetries: 5,
        backoff: httpx.Backoff.constant(Duration.zero),
      );
      expect(result, equals('ok'));
      expect(attempts, equals(3));
    });

    test('retries on ClientException and succeeds', () async {
      int attempts = 0;
      final result = await httpx.withRetry(
        () async {
          attempts++;
          if (attempts < 3) {
            throw http.ClientException(
              'Connection refused',
              Uri.https('localhost:9998', '/meta/d/latest'),
            );
          }
          return 'ok';
        },
        maxRetries: 5,
        backoff: httpx.Backoff.constant(Duration.zero),
      );
      expect(result, equals('ok'));
      expect(attempts, equals(3));
    });

    test('stops retrying after maxRetries and rethrows', () async {
      int attempts = 0;
      await expectLater(
        httpx.withRetry(
          () async {
            attempts++;
            throw SocketException('Connection refused');
          },
          maxRetries: 2,
          backoff: httpx.Backoff.constant(Duration.zero),
        ),
        throwsA(isA<SocketException>()),
      );
      expect(attempts, equals(3)); // initial + 2 retries
    });

    test('does not retry on non-retryable errors', () async {
      int attempts = 0;
      await expectLater(
        httpx.withRetry(
          () async {
            attempts++;
            throw Exception('some unrelated error');
          },
          maxRetries: 5,
          backoff: httpx.Backoff.constant(Duration.zero),
        ),
        throwsA(isA<Exception>()),
      );
      expect(attempts, equals(1));
    });

    test('retries on 502 bad gateway', () async {
      int attempts = 0;
      await expectLater(
        httpx.withRetry(
          () async {
            attempts++;
            throw http.Response('bad gateway', 502);
          },
          maxRetries: 2,
          backoff: httpx.Backoff.constant(Duration.zero),
        ),
        throwsA(isA<http.Response>()),
      );
      expect(attempts, equals(3));
    });

    test('does not retry on 404', () async {
      int attempts = 0;
      await expectLater(
        httpx.withRetry(
          () async {
            attempts++;
            throw http.Response('not found', 404);
          },
          maxRetries: 5,
          backoff: httpx.Backoff.constant(Duration.zero),
        ),
        throwsA(isA<http.Response>()),
      );
      expect(attempts, equals(1));
    });
  });
}
