import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/discovery.dart' as disc;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

// Discovery is driven by a real streaming API (ddisc.api.locate) with no
// injection point, so these tests stick to the no-network path: a search
// with no mimetypes never resolves a category and never calls the network
// (see DiscoveryGrid.refresh's `category.isEmpty` early return).
ValueNotifier<media.MediaSearchState> _emptyCategorySearch() {
  return ValueNotifier(
    media.MediaSearchState(next: media.MediaSearchRequest()),
  );
}

void main() {
  group('discovery.DiscoveryGrid', () {
    testWidgets('with no mimetype category selected renders empty results without a network call', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        disc.DiscoveryGrid(search: _emptyCategorySearch()),
      );
      await tester.pumpAndSettle();

      expect(find.byType(disc.DiscoveryGrid), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders leading widgets', (WidgetTester tester) async {
      await tester.pumpApp(
        disc.DiscoveryGrid(
          search: _emptyCategorySearch(),
          leading: [Text('leading widget')],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('leading widget'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        disc.DiscoveryGrid(search: _emptyCategorySearch()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
