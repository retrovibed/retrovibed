import 'dart:io';

import 'package:cross_file/cross_file.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/community/api.dart';
import 'package:retrovibed/community/socials.action.photo.dart';
import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:retrovibed/design.kit/theme.defaults.dart';
import 'package:retrovibed/media/media.pb.dart';

final _community = Community(
  id: 'c1',
  url: 'https://example-community.community.retrovibe.space',
  defaultPublishMode: PublishMode.SYNDICATED,
);

// the smallest thing the mime resolver will still identify as a png.
const _png = [0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A];

void main() {
  group('SocialActionPhoto', () {
    testWidgets('desktop renders a file drop well for images', (tester) async {
      await tester.pumpApp(
        SocialActionPhoto(
          _community,
          apiupload: (mkreq) async => throw UnimplementedError(),
          apipublish: (cid, req, {options = const []}) async => throw UnimplementedError(),
          capture: () async => throw UnimplementedError(),
        ),
        theme: ThemeData().copyWith(extensions: [Defaults(mobile: false)]),
      );
      await tester.pumpAndSettle();

      expect(find.byType(FileDropWell), findsOneWidget);
      expect(find.byIcon(Icons.add_a_photo_outlined), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('mobile uploads the captured photo then publishes it', (tester) async {
      final dir = Directory.systemTemp.createTempSync('retrovibed_photo_');
      addTearDown(() => dir.deleteSync(recursive: true));
      final photo = File('${dir.path}/snapshot.png')..writeAsBytesSync(_png);

      var uploaded = 0;
      String? publishedTo;
      PublishContentRequest? published;

      await tester.pumpApp(
        SocialActionPhoto(
          _community,
          apiupload: (mkreq) async {
            uploaded++;
            return MediaUploadResponse(media: Media(id: 'm1', knownMediaId: 'k1'));
          },
          apipublish: (cid, req, {options = const []}) async {
            publishedTo = cid;
            published = req;
            return PublishContentResponse();
          },
          capture: () async => XFile(photo.path),
        ),
        theme: ThemeData().copyWith(extensions: [Defaults(mobile: true)]),
      );
      await tester.pumpAndSettle();

      // reading the photo off disk is real i/o, which only progresses outside the fake clock.
      await tester.runAsync(() async {
        await tester.tap(find.byIcon(Icons.add_a_photo_outlined));
        await Future<void>.delayed(const Duration(milliseconds: 250));
      });
      await tester.pump();

      expect(uploaded, equals(1));
      expect(publishedTo, equals('c1'));
      expect(published?.publishMode, equals(PublishMode.SYNDICATED));
      expect(published?.publishedContent.communityId, equals('c1'));
      expect(published?.publishedContent.libraryId, equals('m1'));
      expect(published?.publishedContent.knownMediaId, equals('k1'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('mobile does nothing when the capture is cancelled', (tester) async {
      var uploaded = 0;
      var publishes = 0;

      await tester.pumpApp(
        SocialActionPhoto(
          _community,
          apiupload: (mkreq) async {
            uploaded++;
            return MediaUploadResponse(media: Media(id: 'm1'));
          },
          apipublish: (cid, req, {options = const []}) async {
            publishes++;
            return PublishContentResponse();
          },
          capture: () async => null,
        ),
        theme: ThemeData().copyWith(extensions: [Defaults(mobile: true)]),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.add_a_photo_outlined));
      await tester.pumpAndSettle();

      expect(uploaded, equals(0));
      expect(publishes, equals(0));
      expect(tester.takeException(), isNull);
    });

    testWidgets('mobile surfaces a publish failure instead of throwing', (tester) async {
      final dir = Directory.systemTemp.createTempSync('retrovibed_photo_');
      addTearDown(() => dir.deleteSync(recursive: true));
      final photo = File('${dir.path}/snapshot.png')..writeAsBytesSync(_png);

      var attempted = 0;

      // the failure is reported as a snackbar, which needs a Scaffold to land in.
      await tester.pumpApp(
        Scaffold(
          body: SocialActionPhoto(
            _community,
            apiupload: (mkreq) async => MediaUploadResponse(media: Media(id: 'm1')),
            apipublish: (cid, req, {options = const []}) async {
              attempted++;
              throw Exception('boom');
            },
            capture: () async => XFile(photo.path),
          ),
        ),
        theme: ThemeData().copyWith(extensions: [Defaults(mobile: true)]),
      );
      await tester.pumpAndSettle();

      await tester.runAsync(() async {
        await tester.tap(find.byIcon(Icons.add_a_photo_outlined));
        await Future<void>.delayed(const Duration(milliseconds: 250));
      });
      await tester.pump();
      await tester.pump(const Duration(milliseconds: 500));

      expect(attempted, equals(1));

      expect(find.byType(SnackBar), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
