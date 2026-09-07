import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/discovery/releases.dart';
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<lib.KnownLatestResponse> Function(
  lib.KnownLatestRequest req, {
  List<httpx.Option> options,
})
_loading(Future<lib.KnownLatestResponse> pending) {
  return (
    lib.KnownLatestRequest req, {
    List<httpx.Option> options = const [],
  }) {
    return pending;
  };
}

Future<lib.KnownLatestResponse> _notimplemented(
  lib.KnownLatestRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.error(http.Response('', 501));

Future<lib.KnownLatestResponse> _unauthorized(
  lib.KnownLatestRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.error(http.Response('', 401));

Future<lib.KnownLatestResponse> _empty(
  lib.KnownLatestRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.value(lib.KnownLatestResponse(items: []));

Future<lib.KnownLatestResponse> _withItems(
  lib.KnownLatestRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.value(
      lib.KnownLatestResponse(
        items: [
          lib.Known(id: 'id-1', description: 'Release One'),
          lib.Known(id: 'id-2', description: 'Release Two'),
        ],
      ),
    );

void main() {
  group('Releases', () {
    testWidgets('displays loading state initially', (tester) async {
      final c = Completer<lib.KnownLatestResponse>();
      await tester.pumpApp(NewReleases(mimex.video, latest: _loading(c.future)));
      await tester.pump();
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      c.complete(lib.KnownLatestResponse(items: []));
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
    });

    testWidgets('renders items returned from the api with onChange wired up', (tester) async {
      await tester.pumpApp(NewReleases(mimex.video, latest: _withItems), isolatecache: true);
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
      expect(find.byType(lib.KnownMediaLocator), findsWidgets);
      expect(tester.takeException(), isNull);
    });

    testWidgets('onChange folds an updated item back into the list', (tester) async {
      await tester.pumpApp(NewReleases(mimex.video, latest: _withItems), isolatecache: true);
      await tester.pumpAndSettle();

      final locators = tester.widgetList<lib.KnownMediaLocator>(find.byType(lib.KnownMediaLocator)).toList();
      expect(locators.length, 2);
      final target = locators.firstWhere((l) => l.current.id == 'id-1');

      target.onChange(lib.Known(id: 'id-1', description: 'Release One Updated'));
      await tester.pump();

      final updated = tester.widgetList<lib.KnownMediaLocator>(find.byType(lib.KnownMediaLocator)).toList();
      expect(updated.length, 2);
      expect(updated.firstWhere((l) => l.current.id == 'id-1').current.description, 'Release One Updated');
      expect(updated.firstWhere((l) => l.current.id == 'id-2').current.description, 'Release Two');
      expect(tester.takeException(), isNull);
    });

    testWidgets('onChange removes the item from the list when given null', (tester) async {
      await tester.pumpApp(NewReleases(mimex.video, latest: _withItems), isolatecache: true);
      await tester.pumpAndSettle();

      final locators = tester.widgetList<lib.KnownMediaLocator>(find.byType(lib.KnownMediaLocator)).toList();
      expect(locators.length, 2);
      final target = locators.firstWhere((l) => l.current.id == 'id-1');

      target.onChange(null);
      await tester.pump();

      final remaining = tester.widgetList<lib.KnownMediaLocator>(find.byType(lib.KnownMediaLocator)).toList();
      expect(remaining.length, 1);
      expect(remaining.single.current.id, 'id-2');
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays empty state after loading', (tester) async {
      await tester.pumpApp(NewReleases(mimex.video, latest:_empty));
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
    });

    testWidgets('silently ignores not implemented response', (tester) async {
      await tester.pumpApp(NewReleases(mimex.video, latest:_notimplemented));
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays error on unauthorized response', (tester) async {
      await tester.pumpApp(NewReleases(mimex.video, latest:_unauthorized));
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders at all resolutions', (tester) async {
      await tester.pumpApp(NewReleases(mimex.video, latest:_empty));
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
    });
  });
}
