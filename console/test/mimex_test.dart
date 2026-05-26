import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/mimex.dart' as mimex;

void main() {
  group('category', () {
    test('returns video for video mime types', () {
      expect(mimex.category(mimex.videos), equals('video'));
    });

    test('returns audio for audio mime types', () {
      expect(mimex.category(mimex.audios), equals('audio'));
    });

    test('returns image for image mime types', () {
      expect(mimex.category(mimex.images), equals('image'));
    });

    test('returns empty string for unknown mime types', () {
      expect(mimex.category(['application/octet-stream']), equals(''));
    });

    test('returns empty string for empty list', () {
      expect(mimex.category([]), equals(''));
    });
  });

  group('category consistency with predicates', () {
    test('all video mimes satisfy isVideo', () {
      expect(mimex.videos.every(mimex.isVideo), isTrue);
    });

    test('video mimes produce category video', () {
      final videoMimes = mimex.videos.where(mimex.isVideo).toList();
      expect(mimex.category(videoMimes), equals('video'));
    });

    test('all audio mimes satisfy isAudio', () {
      expect(mimex.audios.every(mimex.isAudio), isTrue);
    });

    test('audio mimes produce category audio', () {
      final audioMimes = mimex.audios.where(mimex.isAudio).toList();
      expect(mimex.category(audioMimes), equals('audio'));
    });

    test('all image mimes satisfy isImage', () {
      expect(mimex.images.every(mimex.isImage), isTrue);
    });

    test('image mimes produce category image', () {
      final imageMimes = mimex.images.where(mimex.isImage).toList();
      expect(mimex.category(imageMimes), equals('image'));
    });
  });
}
