import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/qr.dart';

final _community = Community(
  url: 'https://testdomain.community.retrovibe.space',
  description: 'A test community',
  createdAt: '2024-01-15T14:30:00Z',
);

void main() {
  group('encodeQRPayload', () {
    test('wraps community in community key', () {
      final encoded = encodeQRPayload(_community);
      final decoded = jsonDecode(encoded) as Map<String, dynamic>;

      expect(decoded.containsKey('community'), isTrue);
      expect(decoded['community'], isA<Map<String, dynamic>>());
      expect((decoded['community'] as Map<String, dynamic>)['url'], 'https://testdomain.community.retrovibe.space');
    });

    test('includes attribution when provided', () {
      final encoded = encodeQRPayload(_community, attribution: 'eyJtoken');
      final decoded = jsonDecode(encoded) as Map<String, dynamic>;

      expect(decoded['attribution'], 'eyJtoken');
    });

    test('includes empty attribution when not provided', () {
      final encoded = encodeQRPayload(_community);
      final decoded = jsonDecode(encoded) as Map<String, dynamic>;

      expect(decoded['attribution'], isEmpty);
    });
  });

  group('decodeQRPayload', () {
    test('parses payload with attribution', () {
      final encoded = encodeQRPayload(_community, attribution: 'eyJtoken');
      final (community, attribution) = decodeQRPayload(encoded);

      expect(community, isNotNull);
      expect(community!.url, 'https://testdomain.community.retrovibe.space');
      expect(attribution, 'eyJtoken');
    });

    test('parses payload without attribution', () {
      final encoded = encodeQRPayload(_community);
      final (community, attribution) = decodeQRPayload(encoded);

      expect(community, isNotNull);
      expect(community!.url, 'https://testdomain.community.retrovibe.space');
      expect(attribution, isEmpty);
    });

    test('returns null for invalid data', () {
      final (community, attribution) = decodeQRPayload('not valid json');

      expect(community, isNull);
      expect(attribution, isEmpty);
    });

    test('roundtrips encode and decode', () {
      final encoded = encodeQRPayload(_community, attribution: 'attr-tok');
      final (community, attribution) = decodeQRPayload(encoded);

      expect(community!.url, _community.url);
      expect(community.description, _community.description);
      expect(attribution, 'attr-tok');
    });
  });
}
