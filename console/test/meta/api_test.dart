import 'package:flutter_test/flutter_test.dart';
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
  });
}
