import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/api.dart' as api;
import 'package:retrovibed/retrovibed.dart' as retro;

void main() {
  group('daemons.isLocalDevice', () {
    final localHost = retro.local_device().hostname.split(':').first;

    test('matches when hostname and port match the local device exactly', () {
      final library = api.Daemon(hostname: '$localHost:9998');
      expect(api.daemons.isLocalDevice(library), isTrue);
    });

    test('matches when hostname matches but the port differs from the local device', () {
      final library = api.Daemon(hostname: '$localHost:8443');
      expect(api.daemons.isLocalDevice(library), isTrue);
    });

    test('does not match a different hostname even with the local device port', () {
      final library = api.Daemon(hostname: 'not-$localHost:9998');
      expect(api.daemons.isLocalDevice(library), isFalse);
    });

    test('matches the localhost:9998 fallback regardless of the local device hostname', () {
      final library = api.Daemon(hostname: 'localhost:9998');
      expect(api.daemons.isLocalDevice(library), isTrue);
    });

    test('does not match a different hostname that shares the local device hostname as a prefix', () {
      final library = api.Daemon(hostname: '$localHost-nas:9998');
      expect(api.daemons.isLocalDevice(library), isFalse);
    });

    test('matches localhost on a non-default port', () {
      final library = api.Daemon(hostname: 'localhost:9999');
      expect(api.daemons.isLocalDevice(library), isTrue);
    });

    test('matches the ipv4 loopback address', () {
      final library = api.Daemon(hostname: '127.0.0.1:9998');
      expect(api.daemons.isLocalDevice(library), isTrue);
    });

    test('matches the ipv6 loopback address', () {
      final library = api.Daemon(hostname: '[::1]:9998');
      expect(api.daemons.isLocalDevice(library), isTrue);
    });

    test('matches the configured local daemon host', () {
      final library = api.Daemon(hostname: httpx.localhost());
      expect(api.daemons.isLocalDevice(library), isTrue);
    });
  });
}
