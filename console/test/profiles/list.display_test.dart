import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/profiles/list.row.dart';
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
  String email = 'test@test.com',
}) {
  return meta.Profile(
    id: id,
    display: display,
    email: email,
    updatedAt: '2025-01-01T00:00:00Z',
  );
}

List<meta.Profile> _items() => [
  _profile(id: 'id-1', display: 'Alice'),
  _profile(id: 'id-2', display: 'Bob'),
  _profile(id: 'id-3', display: 'Charlie'),
];

List<meta.Profile> _longNameItems() => [
  _profile(
    id: 'very-long-identifier-that-could-potentially-overflow',
    display: 'A User With An Extremely Long Display Name',
  ),
  _profile(
    id: 'another-very-long-identifier-string',
    display: 'Another User With A Very Long Name That Overflows',
  ),
];

Future<meta.ProfileSearchResponse> _mockSearchWithItems(
  meta.ProfileSearchRequest req, {
  List<httpx.Option> options = const [],
}) {
  return Future.value(_response(items: _items()));
}

Future<meta.ProfileSearchResponse> _mockSearchWithLongNames(
  meta.ProfileSearchRequest req, {
  List<httpx.Option> options = const [],
}) {
  return Future.value(_response(items: _longNameItems()));
}

final _resolutions = Resolutions.variant();

void main() {
  group('profiles.ListDisplay', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        profiles.ListDisplay(search: _mockSearchWithItems),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders with long names without overflow', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        profiles.ListDisplay(search: _mockSearchWithLongNames),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });

  group('Profile ListDisplay onChange wiring', () {
    testWidgets('replaces the matching row in place when onChange fires', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(profiles.ListDisplay(search: _mockSearchWithItems));
      await tester.pumpAndSettle();

      final rows = tester.widgetList<ListRow>(find.byType(ListRow)).toList();
      expect(rows.length, 3);
      final target = rows.firstWhere((r) => r.current.id == 'id-2');

      target.onChange(_profile(id: 'id-2', display: 'Bob Updated'));
      await tester.pump();

      expect(find.text('Bob Updated'), findsOneWidget);
      expect(find.text('Bob'), findsNothing);
      expect(find.text('Alice'), findsOneWidget);
      expect(find.text('Charlie'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Profile ListDisplay Row Tests', () {
    group('Typography rows at default size', () {
      testWidgets('renders rows without overflow', (WidgetTester tester) async {
        await tester.pumpApp(
          profiles.ListDisplay(search: _mockSearchWithItems),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.text('Alice'), findsOneWidget);
        expect(find.text('Bob'), findsOneWidget);
        expect(find.text('Charlie'), findsOneWidget);
      });

      testWidgets('shows all columns at default size', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.ListDisplay(search: _mockSearchWithItems),
        );
        await tester.pumpAndSettle();

        expect(find.text('id'), findsOneWidget);
        expect(find.text('username'), findsOneWidget);
        expect(find.text('updated'), findsOneWidget);
        expect(find.text('id-1'), findsOneWidget);
      });
    });

    group('compact layout (< 300px)', () {
      testWidgets('hides id and updated columns below 300px', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          profiles.ListDisplay(search: _mockSearchWithItems),
          physicalSize: Size(280, 568),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        // header columns
        expect(find.text('id'), findsNothing);
        expect(find.text('updated'), findsNothing);
        expect(find.text('username'), findsOneWidget);
        // row data: id column hidden, username still visible
        expect(find.text('id-1'), findsNothing);
        expect(find.text('Alice'), findsOneWidget);
      });

      testWidgets('renders without overflow at 280x568 with items', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: Size(280, 568),
          profiles.ListDisplay(search: _mockSearchWithItems),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow at 280x568 with long names', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: Size(280, 568),
          profiles.ListDisplay(search: _mockSearchWithLongNames),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow at 240x400', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: Size(240, 400),
          profiles.ListDisplay(search: _mockSearchWithItems),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });

    group('shows all columns above breakpoint', () {
      testWidgets('shows all columns at 400px width', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: Size(400, 568),
          profiles.ListDisplay(search: _mockSearchWithItems),
        );
        await tester.pumpAndSettle();

        expect(find.text('id'), findsOneWidget);
        expect(find.text('username'), findsOneWidget);
        expect(find.text('updated'), findsOneWidget);
        expect(find.text('Alice'), findsOneWidget);
      });

      testWidgets('shows all columns at 420x568', (WidgetTester tester) async {
        await tester.pumpApp(
          physicalSize: Size(420, 568),
          profiles.ListDisplay(search: _mockSearchWithItems),
        );
        await tester.pumpAndSettle();

        expect(find.text('id'), findsOneWidget);
        expect(find.text('username'), findsOneWidget);
        expect(find.text('updated'), findsOneWidget);
      });
    });

    group('constrained containers', () {
      testWidgets('renders without overflow in SizedBox (250x400) with items', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: Size(250, 400),
          profiles.ListDisplay(search: _mockSearchWithItems),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        // compact: id and updated hidden
        expect(find.text('id'), findsNothing);
        expect(find.text('updated'), findsNothing);
      });

      testWidgets('renders without overflow in narrow Column (260px)', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: Size(260, 600),
          Column(
            children: [
              Expanded(
                child: profiles.ListDisplay(search: _mockSearchWithLongNames),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });
  });
}
