import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/community/community.pb.dart';

void main() {
  group('YouTubeStatus', () {
    test('mergeFromProto3Json parses linked true', () {
      final status = YouTubeStatus()..mergeFromProto3Json({'linked': true});
      expect(status.linked, isTrue);
    });

    test('mergeFromProto3Json parses linked false', () {
      final status = YouTubeStatus()..mergeFromProto3Json({'linked': false});
      expect(status.linked, isFalse);
    });

    test('default constructor defaults linked to false', () {
      final status = YouTubeStatus();
      expect(status.linked, isFalse);
    });
  });
}
