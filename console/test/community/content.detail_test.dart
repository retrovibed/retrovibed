import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/community.dart';
import 'package:retrovibed/community/content.detail.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

PublishedContent _item({
  String title = 'Movie Title',
  String description = '',
  String mimetype = 'video/mp4',
  String magnetUri = 'magnet:?xt=urn:btih:abc123',
  String publishedAt = '2026-01-15T10:00:00Z',
  Int64? bytes,
}) => PublishedContent(
  id: 'pc-1',
  communityId: 'community-1',
  knownMediaId: 'magnet:?xt=urn:btih:abc123',
  title: title,
  description: description,
  mimetype: mimetype,
  magnetUri: magnetUri,
  publishedAt: publishedAt,
  bytes: bytes ?? Int64(1500000000),
);

void main() {
  group('PublishedContentDetail', () {
    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          PublishedContentDetail(item: _item(description: 'A great film')),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('description', () {
      testWidgets('shown when non-empty', (tester) async {
        await tester.pumpApp(
          PublishedContentDetail(item: _item(description: 'A compelling description')),
        );
        await tester.pumpAndSettle();

        expect(find.text('A compelling description'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('hidden when empty', (tester) async {
        const sentinel = 'should not appear';
        await tester.pumpApp(
          PublishedContentDetail(item: _item(description: '')),
        );
        await tester.pumpAndSettle();

        expect(find.text(sentinel), findsNothing);
        expect(find.text('video/mp4'), findsOneWidget); // metadata still renders
        expect(tester.takeException(), isNull);
      });
    });

    group('metadata', () {
      testWidgets('shows published timestamp', (tester) async {
        await tester.pumpApp(
          PublishedContentDetail(item: _item(publishedAt: '2026-01-15T10:00:00Z')),
        );
        await tester.pumpAndSettle();

        expect(find.text('published '), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows mimetype', (tester) async {
        await tester.pumpApp(
          PublishedContentDetail(item: _item(mimetype: 'video/mp4')),
        );
        await tester.pumpAndSettle();

        expect(find.text('video/mp4'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows magnetUri as selectable text', (tester) async {
        await tester.pumpApp(
          PublishedContentDetail(item: _item(magnetUri: 'magnet:?xt=urn:btih:deadbeef')),
        );
        await tester.pumpAndSettle();

        expect(find.text('magnet:?xt=urn:btih:deadbeef'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('overflow', () {
      testWidgets('renders inside Expanded within Row without overflow', (tester) async {
        await tester.pumpApp(
          Row(
            children: [
              Expanded(
                child: PublishedContentDetail(
                  item: _item(
                    description: 'A description that is quite long and verbose',
                    magnetUri: 'magnet:?xt=urn:btih:${'a' * 80}&dn=Very+Long+Title',
                  ),
                ),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders inside SingleChildScrollView without overflow', (tester) async {
        await tester.pumpApp(
          SingleChildScrollView(
            child: PublishedContentDetail(item: _item(description: 'Some description')),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });
  });
}
