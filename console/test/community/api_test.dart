import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/community/api.dart' as community;

void main() {
  group('communities.canonicaluri', () {
    test('builds a subdomain URL from a slug', () {
      expect(
        community.communities.canonicaluri('testdomain'),
        equals('https://testdomain.community.retrovibe.space'),
      );
    });

    test('defaults to "example" when the slug is empty', () {
      expect(
        community.communities.canonicaluri(''),
        equals('https://example.community.retrovibe.space'),
      );
    });

    test('returns an https URL unchanged', () {
      expect(
        community.communities.canonicaluri('https://custom.example.com'),
        equals('https://custom.example.com'),
      );
    });
  });

  group('communities.domain', () {
    test('extracts the subdomain from a community URL', () {
      expect(
        community.communities.domain('https://testdomain.community.retrovibe.space'),
        equals('testdomain'),
      );
    });

    test('extracts the subdomain ignoring path and query', () {
      expect(
        community.communities.domain('https://testdomain.community.retrovibe.space/some/path?x=1'),
        equals('testdomain'),
      );
    });

    test('returns the input unchanged when the host is not a community host', () {
      expect(
        community.communities.domain('https://example.com'),
        equals('https://example.com'),
      );
    });

    test('returns the input unchanged when it is not a valid URI', () {
      expect(
        community.communities.domain('not a uri'),
        equals('not a uri'),
      );
    });

    test('returns the input unchanged when empty', () {
      expect(
        community.communities.domain(''),
        equals(''),
      );
    });
  });
}
