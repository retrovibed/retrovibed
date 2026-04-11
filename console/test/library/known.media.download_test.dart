import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/known.media.download.dart';
import 'package:retrovibed/library/known.media.card.dart';
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('KnownMediaDownload', () {
    testWidgets('renders without overflow', (tester) async {
      await tester.pumpApp(
        KnownMediaDownload(
          children: [api.Known(description: 'Test', summary: 'summary')],
          onTap: (k) async => k,
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow with onDoubleTap', (tester) async {
      await tester.pumpApp(
        KnownMediaDownload(
          children: [api.Known(description: 'Test', summary: 'summary')],
          onTap: (k) async => k,
          onDoubleTap: (k) async => k,
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('onTap is invoked when card is tapped', (tester) async {
      api.Known? tapped;
      final item = api.Known(description: 'Test', summary: 'summary');
      await tester.pumpApp(
        KnownMediaDownload(
          children: [item],
          onTap: (k) async {
            tapped = k;
            return k;
          },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaCard));
      await tester.pump();
      expect(tapped, equals(item));
    });

    testWidgets('onDoubleTap is null when not provided', (tester) async {
      await tester.pumpApp(
        KnownMediaDownload(
          children: [api.Known(description: 'Test', summary: 'summary')],
          onTap: (k) async => k,
        ),
      );
      await tester.pumpAndSettle();

      final card = tester.widget<KnownMediaCard>(find.byType(KnownMediaCard));
      expect(card.onDoubleTap, isNull);
    });

    testWidgets('renders empty children without overflow', (tester) async {
      await tester.pumpApp(
        KnownMediaDownload(
          onTap: (k) async => k,
        ),
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with leading widget', (tester) async {
      await tester.pumpApp(
        KnownMediaDownload(
          leading: const Text('leading'),
          children: [api.Known(description: 'Test', summary: 'summary')],
          onTap: (k) async => k,
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('leading'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
