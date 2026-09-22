import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

meta.ProfileSearchResponse _response({
  List<meta.Profile> items = const [],
  int limit = 32,
}) {
  return meta.ProfileSearchResponse(
    items: items,
    next: meta.profiles.request(limit: limit),
  );
}

meta.Profile _profile({
  String id = 'test-id-1',
  String display = 'Test User',
}) {
  return meta.Profile(id: id, display: display, updatedAt: '2025-01-01T00:00:00Z');
}

List<meta.Profile> _items() => [
  _profile(id: 'id-1', display: 'Alice'),
  _profile(id: 'id-2', display: 'Bob'),
];

Future<meta.ProfileSearchResponse> _mockSearchEmpty(
  meta.ProfileSearchRequest req, {
  List<httpx.Option> options = const [],
}) {
  return Future.value(_response(items: const []));
}

Future<meta.ProfileSearchResponse> _mockSearchWithItems(
  meta.ProfileSearchRequest req, {
  List<httpx.Option> options = const [],
}) {
  return Future.value(_response(items: _items()));
}

Future<meta.ProfileSearchResponse> _mockSearchFailure(
  meta.ProfileSearchRequest req, {
  List<httpx.Option> options = const [],
}) {
  return Future.error(Exception('search failed'));
}

final _resolutions = Resolutions.variant();

void main() {
  group('profiles.PendingListDisplay', () {
    testWidgets('renders nothing when there are no pending requests', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        profiles.PendingListDisplay(search: _mockSearchEmpty),
      );
      await tester.pumpAndSettle();

      expect(find.text('id'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders pending requests once loaded', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        profiles.PendingListDisplay(search: _mockSearchWithItems),
      );
      await tester.pumpAndSettle();

      expect(find.text('Alice'), findsOneWidget);
      expect(find.text('Bob'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets(
      'a failed search with no existing items renders empty without throwing',
      (WidgetTester tester) async {
        // The empty-state guard in build() hides everything, including the
        // error, when there are no items to show alongside it.
        await tester.pumpApp(
          profiles.PendingListDisplay(search: _mockSearchFailure),
        );
        await tester.pumpAndSettle();

        expect(find.text('an unexpected problem has occurred'), findsNothing);
        expect(tester.takeException(), isNull);
      },
    );

    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        profiles.PendingListDisplay(search: _mockSearchWithItems),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
